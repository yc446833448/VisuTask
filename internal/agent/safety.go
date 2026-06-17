package agent

import "sync"

// Safety detects doom loops and other problematic agent behaviors
type Safety struct {
	doomThreshold int
	history       []toolCallRecord
	mu            sync.Mutex
}

type toolCallRecord struct {
	Name string
	Args string // serialized args for comparison
}

func NewSafety(threshold int) *Safety {
	if threshold == 0 {
		threshold = 3
	}
	return &Safety{
		doomThreshold: threshold,
	}
}

// RecordToolCall adds a tool call to the history
func (s *Safety) RecordToolCall(name string, args string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.history = append(s.history, toolCallRecord{Name: name, Args: args})

	// Keep only last 20 records
	if len(s.history) > 20 {
		s.history = s.history[len(s.history)-20:]
	}
}

// IsDoomLoop detects if the same tool+args has been called repeatedly
func (s *Safety) IsDoomLoop() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.history) < s.doomThreshold {
		return false
	}

	last := s.history[len(s.history)-1]
	count := 1
	for i := len(s.history) - 2; i >= 0; i-- {
		if s.history[i].Name == last.Name && s.history[i].Args == last.Args {
			count++
		} else {
			break
		}
	}

	return count >= s.doomThreshold
}

// RepeatCount returns how many times the last call has been repeated
func (s *Safety) RepeatCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.history) == 0 {
		return 0
	}

	last := s.history[len(s.history)-1]
	count := 1
	for i := len(s.history) - 2; i >= 0; i-- {
		if s.history[i].Name == last.Name && s.history[i].Args == last.Args {
			count++
		} else {
			break
		}
	}
	return count
}

// Reset clears the history
func (s *Safety) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.history = nil
}
