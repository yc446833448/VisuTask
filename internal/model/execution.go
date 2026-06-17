package model

import "time"

type Execution struct {
	ID          string       `gorm:"primaryKey" json:"id"`
	TaskID      string       `gorm:"index" json:"taskId"`
	TaskName    string       `json:"taskName"`
	Status      string       `json:"status"` // success / failed / cancelled
	StartedAt   time.Time    `json:"startedAt"`
	FinishedAt  time.Time    `json:"finishedAt"`
	Duration    int64        `json:"duration"` // seconds
	StepResults []StepResult `gorm:"serializer:json" json:"stepResults"`
}

type StepResult struct {
	StepID     int     `json:"stepId"`
	Action     string  `json:"action"`
	Target     string  `json:"target"`
	Success    bool    `json:"success"`
	Duration   float64 `json:"duration"` // seconds
	Screenshot string  `json:"screenshot,omitempty"` // base64
	Error      string  `json:"error,omitempty"`
}
