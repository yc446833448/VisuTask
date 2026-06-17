package concurrency

import (
	"context"
	"fmt"
	"sync"
)

// Manager controls concurrent task execution with slot and window conflict management
type Manager struct {
	userMax int
	running map[string]context.CancelFunc // taskID → cancel
	windows map[string]string             // windowHandle → taskID
	mu      sync.Mutex
}

func NewManager(maxConcurrent int) *Manager {
	return &Manager{
		userMax: maxConcurrent,
		running: make(map[string]context.CancelFunc),
		windows: make(map[string]string),
	}
}

// Acquire attempts to get an execution slot and lock the window.
// Returns error if concurrency limit reached or window is occupied.
func (m *Manager) Acquire(taskID, windowHandle string) (context.Context, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check concurrency limit
	if len(m.running) >= m.userMax {
		return nil, fmt.Errorf("concurrency limit reached: %d/%d", len(m.running), m.userMax)
	}

	// Check window conflict
	if windowHandle != "" {
		if occupiedBy, ok := m.windows[windowHandle]; ok {
			return nil, fmt.Errorf("window is occupied by task %s", occupiedBy)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	m.running[taskID] = cancel
	if windowHandle != "" {
		m.windows[windowHandle] = taskID
	}
	return ctx, nil
}

// Release frees the execution slot and window lock
func (m *Manager) Release(taskID, windowHandle string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if cancel, ok := m.running[taskID]; ok {
		cancel()
		delete(m.running, taskID)
	}
	if windowHandle != "" {
		delete(m.windows, windowHandle)
	}
}

// RunningCount returns the number of currently running tasks
func (m *Manager) RunningCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.running)
}

// RunningTasks returns the IDs of running tasks
func (m *Manager) RunningTasks() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	ids := make([]string, 0, len(m.running))
	for id := range m.running {
		ids = append(ids, id)
	}
	return ids
}

// MaxConcurrent returns the user's concurrency limit
func (m *Manager) MaxConcurrent() int {
	return m.userMax
}

// SetMaxConcurrent updates the concurrency limit (e.g., after VIP upgrade)
func (m *Manager) SetMaxConcurrent(max int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.userMax = max
}

// Status returns current concurrency status
type Status struct {
	Running int    `json:"running"`
	Max     int    `json:"max"`
	Tasks   []string `json:"tasks"`
}

func (m *Manager) Status() Status {
	m.mu.Lock()
	defer m.mu.Unlock()
	tasks := make([]string, 0, len(m.running))
	for id := range m.running {
		tasks = append(tasks, id)
	}
	return Status{
		Running: len(m.running),
		Max:     m.userMax,
		Tasks:   tasks,
	}
}
