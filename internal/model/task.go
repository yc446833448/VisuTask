package model

import "time"

type TaskStatus string

const (
	TaskStatusIdle      TaskStatus = "idle"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusPaused    TaskStatus = "paused"
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusFailed    TaskStatus = "failed"
)

type Task struct {
	ID           string            `gorm:"primaryKey" json:"id"`
	Name         string            `gorm:"not null" json:"name"`
	ScriptID     string            `gorm:"index" json:"scriptId"`
	ScriptName   string            `json:"scriptName"`
	WindowHandle string            `json:"windowHandle"`
	WindowTitle  string            `json:"windowTitle"`
	ProcessName  string            `json:"processName"`
	Parameters   map[string]string `gorm:"serializer:json" json:"parameters"`
	Trigger      Trigger           `gorm:"serializer:json" json:"trigger"`
	Status       TaskStatus        `gorm:"default:idle" json:"status"`
	CreatedAt    time.Time         `json:"createdAt"`
	UpdatedAt    time.Time         `json:"updatedAt"`
}

type Trigger struct {
	Type   string `json:"type"`   // manual / cron / hotkey
	Cron   string `json:"cron,omitempty"`
	Hotkey string `json:"hotkey,omitempty"`
}
