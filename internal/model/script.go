package model

import "time"

type Script struct {
	ID          string            `gorm:"primaryKey" json:"id"`
	Name        string            `gorm:"not null" json:"name"`
	Description string            `json:"description"`
	Steps       []Step            `gorm:"serializer:json" json:"steps"`
	Variables   map[string]string `gorm:"serializer:json" json:"variables"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
}

type Step struct {
	ID         int     `json:"id"`
	Action     string  `json:"action"`     // click / input / verify / scroll / hotkey / wait
	Target     string  `json:"target"`     // 目标描述
	TargetOCR  string  `json:"targetOCR"`  // OCR 匹配文字
	Value      string  `json:"value,omitempty"`
	Timeout    int     `json:"timeout,omitempty"`
	Confidence float64 `json:"confidence"`
	Confirmed  bool    `json:"confirmed"`
}
