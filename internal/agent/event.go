package agent

import "sync"

// EventType represents the type of event emitted by the Agent
type EventType string

const (
	EventStepStarted  EventType = "step:started"
	EventTextDelta    EventType = "text:delta"
	EventToolCalled   EventType = "tool:called"
	EventToolResult   EventType = "tool:result"
	EventStepProgress EventType = "step:progress"
	EventScreenshot   EventType = "screenshot"
	EventOCRResult    EventType = "ocr:result"
	EventCompleted    EventType = "completed"
	EventError        EventType = "error"
	EventDoomLoop     EventType = "doom_loop"
)

// EventHandler is a callback function for events
type EventHandler func(event EventType, data interface{})

// EventBus manages event subscriptions and emissions
type EventBus struct {
	handlers []EventHandler
	mu       sync.RWMutex
}

func NewEventBus() *EventBus {
	return &EventBus{}
}

// On registers an event handler
func (e *EventBus) On(handler EventHandler) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.handlers = append(e.handlers, handler)
}

// Emit sends an event to all registered handlers
func (e *EventBus) Emit(event EventType, data interface{}) {
	e.mu.RLock()
	handlers := make([]EventHandler, len(e.handlers))
	copy(handlers, e.handlers)
	e.mu.RUnlock()

	for _, h := range handlers {
		h(event, data)
	}
}

// ─── Event Data Types ───

type StepEvent struct {
	Step   int    `json:"step"`
	Action string `json:"action,omitempty"`
	Target string `json:"target,omitempty"`
}

type ToolCalledEvent struct {
	Name string      `json:"name"`
	Args interface{} `json:"args"`
}

type ToolResultEvent struct {
	Name    string      `json:"name"`
	Success bool        `json:"success"`
	Result  interface{} `json:"result"`
}

type ProgressEvent struct {
	Index  int    `json:"index"`
	Total  int    `json:"total"`
	Action string `json:"action"`
	Target string `json:"target"`
}

type ScreenshotEvent struct {
	Data       string `json:"data"`       // base64
	Annotation string `json:"annotation"` // base64 overlay
}

type OCRResultEvent struct {
	Text       string  `json:"text"`
	X          int     `json:"x"`
	Y          int     `json:"y"`
	Width      int     `json:"width"`
	Height     int     `json:"height"`
	Confidence float64 `json:"confidence"`
}

type CompletedEvent struct {
	Success  bool   `json:"success"`
	Duration string `json:"duration"`
}

type ErrorEvent struct {
	Message string `json:"message"`
	Step    int    `json:"step,omitempty"`
}

type DoomLoopEvent struct {
	Tool  string `json:"tool"`
	Count int    `json:"count"`
}
