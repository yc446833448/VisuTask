package monitor

import (
	"fmt"
	"strings"
	"time"

	"github.com/yc446833448/VisuTask/internal/vision"
)

// Checker verifies step execution results
type Checker struct {
	vision *vision.Engine
}

func NewChecker(v *vision.Engine) *Checker {
	return &Checker{vision: v}
}

// VerifyResult represents the outcome of a verification
type VerifyResult struct {
	Success    bool    `json:"success"`
	Method     string  `json:"method"`
	Confidence float64 `json:"confidence"`
	Message    string  `json:"message"`
	Duration   float64 `json:"duration"` // seconds
}

// VerifyByOCR checks if target text appears on screen via OCR
func (c *Checker) VerifyByOCR(targetText string, timeout int) (*VerifyResult, error) {
	start := time.Now()
	deadline := time.Now().Add(time.Duration(timeout) * time.Millisecond)
	if timeout == 0 {
		deadline = time.Now().Add(5 * time.Second)
	}

	for time.Now().Before(deadline) {
		img, err := c.vision.CaptureScreen()
		if err != nil {
			return nil, fmt.Errorf("capture screen: %w", err)
		}

		results, err := c.vision.Recognize(img)
		if err != nil {
			return nil, fmt.Errorf("ocr recognize: %w", err)
		}

		for _, r := range results {
			if strings.Contains(r.Text, targetText) {
				return &VerifyResult{
					Success:    true,
					Method:     "ocr",
					Confidence: r.Confidence,
					Message:    fmt.Sprintf("found '%s' in OCR results", targetText),
					Duration:   time.Since(start).Seconds(),
				}, nil
			}
		}

		time.Sleep(500 * time.Millisecond)
	}

	return &VerifyResult{
		Success:  false,
		Method:   "ocr",
		Message:  fmt.Sprintf("'%s' not found within timeout", targetText),
		Duration: time.Since(start).Seconds(),
	}, nil
}
