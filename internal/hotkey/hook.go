package hotkey

import (
	"log"
	"sync"
)

// Manager handles global hotkey registration
type Manager struct {
	hooks map[string]func() // hotkey → callback
	mu    sync.Mutex
}

func New() *Manager {
	return &Manager{
		hooks: make(map[string]func()),
	}
}

// Register adds a global hotkey
func (m *Manager) Register(hotkey string, callback func()) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.hooks[hotkey] = callback
	log.Printf("registered hotkey: %s", hotkey)
	// TODO: implement with gohook
	return nil
}

// Unregister removes a global hotkey
func (m *Manager) Unregister(hotkey string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.hooks, hotkey)
	log.Printf("unregistered hotkey: %s", hotkey)
}

// List returns all registered hotkeys
func (m *Manager) List() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	keys := make([]string, 0, len(m.hooks))
	for k := range m.hooks {
		keys = append(keys, k)
	}
	return keys
}

// Trigger simulates a hotkey press (for testing)
func (m *Manager) Trigger(hotkey string) {
	m.mu.Lock()
	cb, ok := m.hooks[hotkey]
	m.mu.Unlock()

	if ok {
		cb()
	}
}
