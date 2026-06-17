package agent

import (
	"context"
	"fmt"
	"time"

	"github.com/yc446833448/VisuTask/internal/llm"
)

const (
	DefaultMaxSteps      = 50
	DefaultToolTimeout   = 30 * time.Second
)

// AgentConfig defines a built-in agent role
type AgentConfig struct {
	Name        string
	Description string
	Model       string
	MaxSteps    int
	BasePrompt  string
	ToolNames   []string
}

// Loop is the core agent loop: load context → call LLM → process tools → evaluate
type Loop struct {
	session    *Session
	config     AgentConfig
	processor  *Processor
	memory     *Memory
	registry   *ToolRegistry
	safety     *Safety
	retry      *RetryPolicy
	gateway    *llm.Gateway
	events     *EventBus
	maxSteps   int
}

// NewLoop creates a new agent loop
func NewLoop(
	session *Session,
	config AgentConfig,
	processor *Processor,
	memory *Memory,
	registry *ToolRegistry,
	safety *Safety,
	gateway *llm.Gateway,
	events *EventBus,
) *Loop {
	maxSteps := config.MaxSteps
	if maxSteps == 0 {
		maxSteps = DefaultMaxSteps
	}
	return &Loop{
		session:   session,
		config:    config,
		processor: processor,
		memory:    memory,
		registry:  registry,
		safety:    safety,
		retry:     DefaultRetryPolicy(),
		gateway:   gateway,
		events:    events,
		maxSteps:  maxSteps,
	}
}

// Run executes the agent loop until completion, error, or context cancellation
func (l *Loop) Run(ctx context.Context) error {
	// Prevent concurrent execution
	done, err := l.session.runState.EnsureRunning()
	if err != nil {
		return err
	}
	defer done()

	l.session.Status = SessionBusy
	step := 0
	attempt := 0

	for {
		// Check context cancellation
		select {
		case <-ctx.Done():
			l.session.Status = SessionFailed
			return ctx.Err()
		default:
		}

		// Max steps guard
		if step >= l.maxSteps {
			l.events.Emit(EventError, ErrorEvent{
				Message: fmt.Sprintf("max steps (%d) reached", l.maxSteps),
				Step:    step,
			})
			l.session.Status = SessionFailed
			return fmt.Errorf("max steps reached: %d", l.maxSteps)
		}

		// Phase 1: Check context overflow → compact
		if l.memory.IsOverflow() {
			l.memory.Compact()
		}

		// Phase 2: Build system prompt
		system := l.buildSystemPrompt()

		// Phase 3: Get available tools for this agent
		tools := l.registry.ToolsForAgent(l.config.Name)

		// Phase 4: Convert messages to LLM format
		messages := l.toLLMMessages()

		// Phase 5: Notify step started
		l.events.Emit(EventStepStarted, StepEvent{Step: step})

		// Phase 6: Call LLM via processor
		result, err := l.processor.Process(ctx, l.gateway, ProcessInput{
			Messages: messages,
			Tools:    tools,
			System:   system,
			Model:    l.config.Model,
		})

		if err != nil {
			// Retry logic
			if l.retry.ShouldRetry(err) && attempt < l.retry.MaxRetries {
				attempt++
				l.events.Emit(EventError, ErrorEvent{
					Message: fmt.Sprintf("retrying (attempt %d): %v", attempt, err),
					Step:    step,
				})
				if waitErr := l.retry.Wait(ctx, attempt, err); waitErr != nil {
					return waitErr
				}
				continue
			}
			l.session.Status = SessionFailed
			l.events.Emit(EventError, ErrorEvent{Message: err.Error(), Step: step})
			return fmt.Errorf("process failed: %w", err)
		}

		// Reset retry counter on success
		attempt = 0

		// Phase 7: Save assistant response to memory
		assistantMsg := AgentMessage{
			Role:      "assistant",
			Content:   result.TextDelta,
			ToolCalls: result.ToolCalls,
			Timestamp: time.Now(),
		}
		l.memory.Append(assistantMsg)

		// Phase 8: Append tool results as user messages (for next LLM turn)
		if len(result.ToolCalls) > 0 {
			for _, tc := range result.ToolCalls {
				if tc.Result != nil {
					toolMsg := AgentMessage{
						Role:      "user",
						Timestamp: time.Now(),
						ToolResult: &ToolResultInfo{
							Content:    tc.Result.Content,
							Screenshot: tc.Result.Screenshot,
							IsError:    tc.Result.IsError,
						},
					}
					l.memory.Append(toolMsg)
				}
			}
		}

		// Phase 9: Evaluate outcome
		switch result.Outcome {
		case OutcomeStop:
			l.session.Status = SessionCompleted
			l.events.Emit(EventCompleted, CompletedEvent{Success: true})
			return nil

		case OutcomeCompact:
			l.memory.Compact()
			continue

		case OutcomeContinue:
			step++

			// Doom loop detection
			if l.safety.IsDoomLoop() {
				l.events.Emit(EventDoomLoop, DoomLoopEvent{
					Tool:  result.ToolCalls[len(result.ToolCalls)-1].Name,
					Count: l.safety.RepeatCount(),
				})
				// Continue anyway — LLM should adapt based on tool error feedback
			}
			continue
		}
	}
}

// buildSystemPrompt constructs the system prompt for the current agent
func (l *Loop) buildSystemPrompt() string {
	var parts []string

	// Base agent prompt
	parts = append(parts, l.config.BasePrompt)

	// Environment info
	parts = append(parts, fmt.Sprintf(`## Environment
- Date: %s
- Window: %s
- OS: windows`, time.Now().Format("2006-01-02 15:04"), l.session.WindowTitle))

	// Tool descriptions
	tools := l.registry.ToolsForAgent(l.config.Name)
	if len(tools) > 0 {
		var toolDescs []string
		for _, t := range tools {
			toolDescs = append(toolDescs, fmt.Sprintf("- **%s**: %s", t.Name, t.Description))
		}
		parts = append(parts, "## Available Tools\n"+joinStrings(toolDescs, "\n"))
	}

	// Output rules
	parts = append(parts, `## Rules
- Always take a screenshot before acting
- Use OCR to locate targets, don't guess coordinates
- Verify results after each action
- If failing twice in a row, stop and explain why`)

	return joinStrings(parts, "\n\n")
}

// toLLMMessages converts agent messages to LLM message format
func (l *Loop) toLLMMessages() []llm.Message {
	msgs := l.memory.Messages()
	result := make([]llm.Message, 0, len(msgs))

	for _, msg := range msgs {
		lm := llm.Message{Role: msg.Role}

		if msg.Content != "" {
			lm.Content = append(lm.Content, &llm.TextBlock{Type: "text", Text: msg.Content})
		}

		// Tool calls from assistant
		for _, tc := range msg.ToolCalls {
			lm.Content = append(lm.Content, &llm.ToolUseBlock{
				Type:  "tool_use",
				ID:    tc.ID,
				Name:  tc.Name,
				Input: tc.Args,
			})
		}

		// Tool results from user
		if msg.ToolResult != nil {
			lm.Content = append(lm.Content, &llm.ToolResultBlock{
				Type:      "tool_result",
				ToolUseID: "", // would need to track which tool this responds to
				Content:   msg.ToolResult.Content,
				IsError:   msg.ToolResult.IsError,
			})
		}

		if len(lm.Content) > 0 {
			result = append(result, lm)
		}
	}

	return result
}

func joinStrings(parts []string, sep string) string {
	result := ""
	for i, p := range parts {
		if i > 0 {
			result += sep
		}
		result += p
	}
	return result
}
