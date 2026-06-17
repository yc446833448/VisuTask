package scheduler

import (
	"log"
	"sync"

	"github.com/robfig/cron/v3"
)

// Scheduler manages cron-based task triggers
type Scheduler struct {
	cron *cron.Cron
	jobs map[string]cron.EntryID // taskID → entryID
	mu   sync.Mutex
}

func New() *Scheduler {
	return &Scheduler{
		cron: cron.New(cron.WithSeconds()),
		jobs: make(map[string]cron.EntryID),
	}
}

// Start begins the cron scheduler
func (s *Scheduler) Start() {
	s.cron.Start()
	log.Println("scheduler started")
}

// Stop halts the cron scheduler
func (s *Scheduler) Stop() {
	s.cron.Stop()
	log.Println("scheduler stopped")
}

// AddJob registers a task with a cron expression
func (s *Scheduler) AddJob(taskID, cronExpr string, fn func()) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entryID, err := s.cron.AddFunc(cronExpr, fn)
	if err != nil {
		return err
	}

	s.jobs[taskID] = entryID
	log.Printf("scheduled task %s with cron: %s", taskID, cronExpr)
	return nil
}

// RemoveJob unregisters a task from the scheduler
func (s *Scheduler) RemoveJob(taskID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if entryID, ok := s.jobs[taskID]; ok {
		s.cron.Remove(entryID)
		delete(s.jobs, taskID)
		log.Printf("unscheduled task %s", taskID)
	}
}

// ListJobs returns all scheduled task IDs
func (s *Scheduler) ListJobs() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	ids := make([]string, 0, len(s.jobs))
	for id := range s.jobs {
		ids = append(ids, id)
	}
	return ids
}
