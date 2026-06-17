package agent

import (
	"fmt"
	"sync"
)

// Memory manages the agent's conversation context with compaction support
type Memory struct {
	messages      []AgentMessage
	mu            sync.RWMutex
	contextWindow int // model context window in tokens
}

func NewMemory(contextWindow int) *Memory {
	if contextWindow == 0 {
		contextWindow = 128000 // default
	}
	return &Memory{
		contextWindow: contextWindow,
	}
}

// Append adds a message to the conversation
func (m *Memory) Append(msg AgentMessage) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.messages = append(m.messages, msg)
}

// Messages returns all messages (read-only copy)
func (m *Memory) Messages() []AgentMessage {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]AgentMessage, len(m.messages))
	copy(result, m.messages)
	return result
}

// Clear resets the conversation
func (m *Memory) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.messages = nil
}

// TokenEstimate provides a rough token count (1 token ≈ 4 chars)
func (m *Memory) TokenEstimate() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	total := 0
	for _, msg := range m.messages {
		total += len(msg.Content) / 4
		for _, tc := range msg.ToolCalls {
			total += len(tc.Name) / 4
			total += len(tc.Args) / 4
			if tc.Result != nil {
				total += len(tc.Result.Content) / 4
			}
		}
	}
	return total
}

// IsOverflow checks if the context exceeds the model's window
func (m *Memory) IsOverflow() bool {
	return m.TokenEstimate() > m.contextWindow
}

// Compact performs two-phase context compression:
// Phase 1: Discard all screenshot/image data, keep only text
// Phase 2: Summarize old messages if still overflowing
func (m *Memory) Compact() {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Phase 1: Discard screenshots from tool results
	for i := range m.messages {
		for j := range m.messages[i].ToolCalls {
			if m.messages[i].ToolCalls[j].Result != nil {
				m.messages[i].ToolCalls[j].Result.Screenshot = ""
			}
		}
	}

	// Phase 2: If still overflowing, keep only the last N messages
	if m.isOverflowLocked() {
		keepRecent := 10
		if len(m.messages) > keepRecent {
			// Replace old messages with a summary placeholder
			oldCount := len(m.messages) - keepRecent
			summary := AgentMessage{
				Role:    "system",
				Content: fmt.Sprintf("[Compacted: %d earlier messages removed to fit context window. Only text data retained.]", oldCount),
			}
			m.messages = append([]AgentMessage{summary}, m.messages[oldCount:]...)
		}
	}
}

func (m *Memory) isOverflowLocked() bool {
	total := 0
	for _, msg := range m.messages {
		total += len(msg.Content) / 4
	}
	return total > m.contextWindow
}
