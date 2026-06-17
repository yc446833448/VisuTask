package agent

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yc446833448/VisuTask/internal/action"
	"github.com/yc446833448/VisuTask/internal/concurrency"
	"github.com/yc446833448/VisuTask/internal/llm"
	"github.com/yc446833448/VisuTask/internal/model"
	"github.com/yc446833448/VisuTask/internal/monitor"
	"github.com/yc446833448/VisuTask/internal/store"
	"github.com/yc446833448/VisuTask/internal/vision"
)

// Service is the main entry point for Wails bindings.
// It orchestrates all agent subsystems.
type Service struct {
	store       *store.DB
	vision      *vision.Engine
	action      *action.Engine
	monitor     *monitor.Checker
	gateway     *llm.Gateway
	concurrency *concurrency.Manager
	sessions    *SessionManager
	registry    *ToolRegistry
	events      *EventBus
	safety      *Safety
}

// NewService creates a fully wired agent service
func NewService(
	db *store.DB,
	v *vision.Engine,
	a *action.Engine,
	m *monitor.Checker,
	gw *llm.Gateway,
	cm *concurrency.Manager,
	registry *ToolRegistry,
) *Service {
	return &Service{
		store:       db,
		vision:      v,
		action:      a,
		monitor:     m,
		gateway:     gw,
		concurrency: cm,
		sessions:    NewSessionManager(),
		registry:    registry,
		events:      NewEventBus(),
		safety:      NewSafety(3),
	}
}

// Events returns the event bus for frontend subscription
func (s *Service) Events() *EventBus {
	return s.events
}

// ─── Script Creation ───

// CreateScriptPlan generates a script plan from natural language description
func (s *Service) CreateScriptPlan(ctx context.Context, description string) (*Session, error) {
	session := s.sessions.Create(SessionScriptCreation, "planner", "claude-sonnet-4-20250514")

	memory := NewMemory(128000)
	memory.Append(AgentMessage{
		Role:      "user",
		Content:   description,
		Timestamp: time.Now(),
	})

	processor := NewProcessor(s.registry, s.events, s.safety)

	loop := NewLoop(session, AgentConfig{
		Name:        "planner",
		Model:       "claude-sonnet-4-20250514",
		MaxSteps:    20,
		BasePrompt:  "You are a GUI automation planner. Help the user create an automation script by observing the screen and planning steps.",
		ToolNames:   []string{"capture", "ocr", "locate", "window", "simulate"},
	}, processor, memory, s.registry, s.safety, s.gateway, s.events)

	go func() {
		if err := loop.Run(ctx); err != nil {
			s.events.Emit(EventError, ErrorEvent{Message: err.Error()})
		}
	}()

	return session, nil
}

// ─── Task Execution ───

// StartTask begins executing a task
func (s *Service) StartTask(ctx context.Context, taskID string) error {
	task, err := s.store.GetTask(taskID)
	if err != nil {
		return fmt.Errorf("get task: %w", err)
	}

	// Acquire concurrency slot + window lock
	execCtx, err := s.concurrency.Acquire(taskID, task.WindowHandle)
	if err != nil {
		return fmt.Errorf("acquire slot: %w", err)
	}

	// Update status
	s.store.UpdateTaskStatus(taskID, model.TaskStatusRunning)

	session := s.sessions.Create(SessionTaskExecution, "executor", "claude-sonnet-4-20250514")
	session.WindowHandle = task.WindowHandle
	session.WindowTitle = task.WindowTitle

	// Build execution prompt from script steps
	script, err := s.store.GetScript(task.ScriptID)
	if err != nil {
		s.concurrency.Release(taskID, task.WindowHandle)
		return fmt.Errorf("get script: %w", err)
	}

	memory := NewMemory(128000)
	memory.Append(AgentMessage{
		Role:      "user",
		Content:   buildExecutionPrompt(task, script),
		Timestamp: time.Now(),
	})

	processor := NewProcessor(s.registry, s.events, s.safety)

	loop := NewLoop(session, AgentConfig{
		Name:        "executor",
		Model:       "claude-sonnet-4-20250514",
		MaxSteps:    100,
		BasePrompt:  "You are a GUI automation executor. Execute the given steps on the target window.",
		ToolNames:   []string{"capture", "ocr", "click", "type", "hotkey", "scroll", "wait", "verify", "locate", "window"},
	}, processor, memory, s.registry, s.safety, s.gateway, s.events)

	go func() {
		defer func() {
			s.concurrency.Release(taskID, task.WindowHandle)
			s.sessions.Remove(session.ID)
		}()

		err := loop.Run(execCtx)
		if err != nil {
			s.store.UpdateTaskStatus(taskID, model.TaskStatusFailed)
			s.events.Emit(EventError, ErrorEvent{Message: err.Error()})
		} else {
			s.store.UpdateTaskStatus(taskID, model.TaskStatusCompleted)

			// Save execution record
			execution := &model.Execution{
				ID:         uuid.New().String(),
				TaskID:     taskID,
				TaskName:   task.Name,
				Status:     "success",
				StartedAt:  session.CreatedAt,
				FinishedAt: time.Now(),
				Duration:   int64(time.Since(session.CreatedAt).Seconds()),
			}
			s.store.CreateExecution(execution)
		}
	}()

	return nil
}

// StopTask stops a running task
func (s *Service) StopTask(taskID string) error {
	sessions := s.sessions.List()
	for _, sess := range sessions {
		// Find the session associated with this task
		sess.Stop()
	}
	s.store.UpdateTaskStatus(taskID, model.TaskStatusFailed)
	return nil
}

// ─── Helpers ───

func buildExecutionPrompt(task *model.Task, script *model.Script) string {
	prompt := fmt.Sprintf("Execute the following automation script on window: %s\n\n", task.WindowTitle)
	prompt += fmt.Sprintf("Script: %s\n", script.Name)
	prompt += fmt.Sprintf("Description: %s\n\n", script.Description)
	prompt += "Steps:\n"
	for i, step := range script.Steps {
		prompt += fmt.Sprintf("%d. %s → \"%s\"", i+1, step.Action, step.Target)
		if step.Value != "" {
			prompt += fmt.Sprintf(" (value: %s)", step.Value)
		}
		prompt += "\n"
	}

	if len(task.Parameters) > 0 {
		prompt += "\nParameters:\n"
		for k, v := range task.Parameters {
			prompt += fmt.Sprintf("- %s = %s\n", k, v)
		}
	}

	return prompt
}
