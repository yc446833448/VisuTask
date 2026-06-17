package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/yc446833448/VisuTask/internal/llm"
)

// Outcome represents the result of processing one LLM turn
type Outcome string

const (
	OutcomeStop     Outcome = "stop"     // Agent completed
	OutcomeCompact  Outcome = "compact"  // Context overflow, needs compaction
	OutcomeContinue Outcome = "continue" // Tool calls made, loop continues
)

// ProcessInput contains everything needed for one LLM call
type ProcessInput struct {
	Messages []llm.Message
	Tools    []*Tool
	System   string
	Model    string
}

// ProcessResult contains the outcome of processing one LLM response
type ProcessResult struct {
	Outcome   Outcome
	ToolCalls []ToolCallInfo
	TextDelta string
}

// Processor handles LLM streaming responses and tool execution
type Processor struct {
	registry *ToolRegistry
	events   *EventBus
	safety   *Safety
}

func NewProcessor(registry *ToolRegistry, events *EventBus, safety *Safety) *Processor {
	return &Processor{
		registry: registry,
		events:   events,
		safety:   safety,
	}
}

// Process calls the LLM and handles the streaming response
func (p *Processor) Process(ctx context.Context, gateway *llm.Gateway, input ProcessInput) (*ProcessResult, error) {
	// Convert tools to LLM schema
	toolSchemas := make([]llm.ToolSchema, len(input.Tools))
	for i, t := range input.Tools {
		toolSchemas[i] = llm.ToolSchema{
			Name:        t.Name,
			Description: t.Description,
			InputSchema: t.InputSchema,
		}
	}

	// Call LLM Gateway
	stream, err := gateway.Stream(ctx, llm.StreamRequest{
		Model:     input.Model,
		System:    input.System,
		Messages:  input.Messages,
		Tools:     toolSchemas,
		MaxTokens: 8192,
	})
	if err != nil {
		return nil, fmt.Errorf("llm stream: %w", err)
	}

	// Process stream events
	var result ProcessResult
	var textBuf strings.Builder

	// Track pending tool calls by index
	type pendingTool struct {
		id        string
		name      string
		inputBuf  strings.Builder
	}
	pendingTools := make(map[int]*pendingTool)

	for event := range stream {
		switch e := event.(type) {
		case *llm.MessageStartEvent:
			// Session started, nothing to do

		case *llm.ContentBlockStartEvent:
			if tb, ok := e.Block.(*llm.ToolUseBlock); ok {
				pendingTools[e.Index] = &pendingTool{
					id:   tb.ID,
					name: tb.Name,
				}
				p.events.Emit(EventToolCalled, ToolCalledEvent{
					Name: tb.Name,
				})
			}

		case *llm.ContentBlockDeltaEvent:
			switch d := e.Delta.(type) {
			case llm.TextDelta:
				textBuf.WriteString(d.Text)
				p.events.Emit(EventTextDelta, d.Text)

			case llm.InputJSONDelta:
				if pt, ok := pendingTools[e.Index]; ok {
					pt.inputBuf.WriteString(d.PartialJSON)
				}
			}

		case *llm.ContentBlockStopEvent:
			if pt, ok := pendingTools[e.Index]; ok {
				// Execute the tool
				toolResult := p.executeTool(ctx, pt.name, pt.inputBuf.String())

				p.events.Emit(EventToolResult, ToolResultEvent{
					Name:    pt.name,
					Success: !toolResult.IsError,
					Result:  toolResult,
				})

				result.ToolCalls = append(result.ToolCalls, ToolCallInfo{
					ID:   pt.id,
					Name: pt.name,
					Args: json.RawMessage(pt.inputBuf.String()),
					Result: &ToolResultInfo{
						Content:    toolResult.Content,
						Screenshot: toolResult.Screenshot,
						IsError:    toolResult.IsError,
					},
				})

				// Record for doom loop detection
				p.safety.RecordToolCall(pt.name, pt.inputBuf.String())

				delete(pendingTools, e.Index)
			}

		case *llm.MessageDeltaEvent:
			if e.StopReason == "tool_use" && len(result.ToolCalls) > 0 {
				result.Outcome = OutcomeContinue
			} else {
				result.Outcome = OutcomeStop
			}

		case *llm.LLMErrorEvent:
			return nil, e.Err
		}
	}

	result.TextDelta = textBuf.String()
	return &result, nil
}

// executeTool runs a single tool and returns the result
func (p *Processor) executeTool(ctx context.Context, name string, argsJSON string) *ToolResult {
	tool, ok := p.registry.Get(name)
	if !ok {
		return &ToolResult{
			Content: fmt.Sprintf("error: unknown tool '%s'", name),
			IsError: true,
		}
	}

	// Apply timeout
	toolCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	tc := &ToolContext{
		SessionID: uuid.New().String(),
		Events:    p.events,
	}

	result, err := tool.Execute(toolCtx, json.RawMessage(argsJSON), tc)
	if err != nil {
		return &ToolResult{
			Content: fmt.Sprintf("error: %v", err),
			IsError: true,
		}
	}

	return result
}
