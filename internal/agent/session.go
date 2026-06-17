package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// SessionType distinguishes script creation from task execution
type SessionType string

const (
	SessionScriptCreation SessionType = "script_creation"
	SessionTaskExecution  SessionType = "task_execution"
)

// SessionStatus tracks the session lifecycle
type SessionStatus string

const (
	SessionIdle      SessionStatus = "idle"
	SessionBusy      SessionStatus = "busy"
	SessionPaused    SessionStatus = "paused"
	SessionCompleted SessionStatus = "completed"
	SessionFailed    SessionStatus = "failed"
)

// Session represents a running Agent instance
type Session struct {
	ID           string
	Type         SessionType
	Agent        string        // planner / executor / reviewer
	Model        string        // LLM model ID
	Status       SessionStatus
	WindowHandle string
	WindowTitle  string
	CreatedAt    time.Time

	cancel   context.CancelFunc
	runState *RunState
}

// RunState prevents concurrent loop execution on the same session
type RunState struct {
	mu      sync.Mutex
	running bool
}

func (r *RunState) EnsureRunning() (func(), error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.running {
		return nil, fmt.Errorf("session already running")
	}
	r.running = true
	return func() {
		r.mu.Lock()
		r.running = false
		r.mu.Unlock()
	}, nil
}

// NewSession creates a new agent session
func NewSession(sessionType SessionType, agent, model string) *Session {
	return &Session{
		ID:        uuid.New().String(),
		Type:      sessionType,
		Agent:     agent,
		Model:     model,
		Status:    SessionIdle,
		CreatedAt: time.Now(),
		runState:  &RunState{},
	}
}

// Stop cancels the session's context
func (s *Session) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
}

// Manager tracks active sessions
type SessionManager struct {
	sessions map[string]*Session
	mu       sync.RWMutex
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*Session),
	}
}

func (m *SessionManager) Create(sessionType SessionType, agent, model string) *Session {
	s := NewSession(sessionType, agent, model)
	m.mu.Lock()
	m.sessions[s.ID] = s
	m.mu.Unlock()
	return s
}

func (m *SessionManager) Get(id string) (*Session, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.sessions[id]
	return s, ok
}

func (m *SessionManager) Remove(id string) {
	m.mu.Lock()
	if s, ok := m.sessions[id]; ok {
		s.Stop()
	}
	delete(m.sessions, id)
	m.mu.Unlock()
}

func (m *SessionManager) List() []*Session {
	m.mu.RLock()
	defer m.mu.RUnlock()
	list := make([]*Session, 0, len(m.sessions))
	for _, s := range m.sessions {
		list = append(list, s)
	}
	return list
}

// ─── Message Types for Agent Memory ───

// AgentMessage represents a message in the agent conversation
type AgentMessage struct {
	Role      string          `json:"role"` // "user" / "assistant" / "system"
	Content   string          `json:"content"`
	ToolCalls []ToolCallInfo  `json:"tool_calls,omitempty"`
	ToolResult *ToolResultInfo `json:"tool_result,omitempty"`
	Timestamp time.Time       `json:"timestamp"`
	Tokens    int             `json:"tokens,omitempty"`
}

// ToolCallInfo records a tool invocation
type ToolCallInfo struct {
	ID     string          `json:"id"`
	Name   string          `json:"name"`
	Args   json.RawMessage `json:"args"`
	Result *ToolResultInfo `json:"result,omitempty"`
}

// ToolResultInfo records a tool execution result
type ToolResultInfo struct {
	Content    string `json:"content"`
	Screenshot string `json:"screenshot,omitempty"` // base64
	IsError    bool   `json:"is_error"`
}
