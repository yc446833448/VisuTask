package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/yc446833448/VisuTask/internal/action"
	"github.com/yc446833448/VisuTask/internal/agent"
	"github.com/yc446833448/VisuTask/internal/monitor"
	"github.com/yc446833448/VisuTask/internal/vision"
)

// RegisterAll adds all GUI automation tools to the registry
func RegisterAll(r *agent.ToolRegistry, v *vision.Engine, a *action.Engine, m *monitor.Checker) {
	r.Register(newCaptureTool(v))
	r.Register(newOCRTool(v))
	r.Register(newLocateTool(v))
	r.Register(newClickTool(a))
	r.Register(newTypeTool(a))
	r.Register(newHotkeyTool(a))
	r.Register(newScrollTool(a))
	r.Register(newWaitTool())
	r.Register(newVerifyTool(m))
	r.Register(newSimulateTool(v))
	r.Register(newWindowTool(a))
}

// ─── Capture Tool ───

func newCaptureTool(v *vision.Engine) *agent.Tool {
	return &agent.Tool{
		Name:        "capture",
		Description: "Take a screenshot of the current screen or a specific window",
		InputSchema: json.RawMessage(`{"type":"object","properties":{"window":{"type":"string","description":"window handle (optional, defaults to full screen)"}}}`),
		Execute: func(ctx context.Context, args json.RawMessage, tc *agent.ToolContext) (*agent.ToolResult, error) {
			var input struct {
				Window string `json:"window"`
			}
			json.Unmarshal(args, &input)

			var img []byte
			var err error
			if input.Window != "" {
				img, err = v.CaptureWindow(input.Window)
			} else if tc.WindowHandle != "" {
				img, err = v.CaptureWindow(tc.WindowHandle)
			} else {
				img, err = v.CaptureScreen()
			}
			if err != nil {
				return nil, fmt.Errorf("capture: %w", err)
			}
			return &agent.ToolResult{
				Content:    "screenshot captured successfully",
				Screenshot: string(img),
			}, nil
		},
	}
}

// ─── OCR Tool ───

func newOCRTool(v *vision.Engine) *agent.Tool {
	return &agent.Tool{
		Name:        "ocr",
		Description: "Perform OCR text recognition on a screenshot, returns text with coordinates",
		InputSchema: json.RawMessage(`{"type":"object","properties":{"image":{"type":"string","description":"base64 image data (optional, uses latest screenshot if omitted)"}}}`),
		Execute: func(ctx context.Context, args json.RawMessage, tc *agent.ToolContext) (*agent.ToolResult, error) {
			var input struct {
				Image string `json:"image"`
			}
			json.Unmarshal(args, &input)

			imageData := []byte(input.Image)
			if len(imageData) == 0 {
				// Use latest screenshot
				var err error
				imageData, err = v.CaptureScreen()
				if err != nil {
					return nil, err
				}
			}

			results, err := v.Recognize(imageData)
			if err != nil {
				return nil, fmt.Errorf("ocr: %w", err)
			}

			// Format results
			var lines []string
			for _, r := range results {
				lines = append(lines, fmt.Sprintf("\"%s\" at (%d,%d) conf=%.2f", r.Text, r.Rect.X, r.Rect.Y, r.Confidence))
			}

			content := fmt.Sprintf("OCR found %d elements:\n%s", len(results), joinLines(lines))
			return &agent.ToolResult{Content: content}, nil
		},
	}
}

// ─── Locate Tool ───

func newLocateTool(v *vision.Engine) *agent.Tool {
	return &agent.Tool{
		Name:        "locate",
		Description: "Locate a UI element by text using OCR, returns coordinates",
		InputSchema: json.RawMessage(`{"type":"object","properties":{"text":{"type":"string","description":"text to search for on screen"}},"required":["text"]}`),
		Execute: func(ctx context.Context, args json.RawMessage, tc *agent.ToolContext) (*agent.ToolResult, error) {
			var input struct {
				Text string `json:"text"`
			}
			if err := json.Unmarshal(args, &input); err != nil {
				return nil, err
			}

			img, err := v.CaptureScreen()
			if err != nil {
				return nil, err
			}

			results, err := v.Recognize(img)
			if err != nil {
				return nil, err
			}

			// Find best match
			for _, r := range results {
				if containsIgnoreCase(r.Text, input.Text) {
					cx := r.Rect.X + r.Rect.Width/2
					cy := r.Rect.Y + r.Rect.Height/2
					return &agent.ToolResult{
						Content: fmt.Sprintf("found \"%s\" at (%d, %d), rect=(%d,%d,%d,%d), confidence=%.2f",
							r.Text, cx, cy, r.Rect.X, r.Rect.Y, r.Rect.Width, r.Rect.Height, r.Confidence),
					}, nil
				}
			}

			return &agent.ToolResult{
				Content: fmt.Sprintf("target \"%s\" not found on screen", input.Text),
				IsError: true,
			}, nil
		},
	}
}

// ─── Click Tool ───

func newClickTool(a *action.Engine) *agent.Tool {
	return &agent.Tool{
		Name:        "click",
		Description: "Click at specific screen coordinates",
		InputSchema: json.RawMessage(`{"type":"object","properties":{"x":{"type":"integer"},"y":{"type":"integer"},"button":{"type":"string","enum":["left","right","double"]}},"required":["x","y"]}`),
		Execute: func(ctx context.Context, args json.RawMessage, tc *agent.ToolContext) (*agent.ToolResult, error) {
			var input struct {
				X      int    `json:"x"`
				Y      int    `json:"y"`
				Button string `json:"button"`
			}
			json.Unmarshal(args, &input)

			var err error
			switch input.Button {
			case "right":
				err = a.RightClick(input.X, input.Y)
			case "double":
				err = a.DoubleClick(input.X, input.Y)
			default:
				err = a.Click(input.X, input.Y)
			}
			if err != nil {
				return nil, err
			}
			return &agent.ToolResult{Content: fmt.Sprintf("clicked at (%d, %d)", input.X, input.Y)}, nil
		},
	}
}

// ─── Type Tool ───

func newTypeTool(a *action.Engine) *agent.Tool {
	return &agent.Tool{
		Name:        "type",
		Description: "Type text at the current cursor position",
		InputSchema: json.RawMessage(`{"type":"object","properties":{"text":{"type":"string","description":"text to type"}},"required":["text"]}`),
		Execute: func(ctx context.Context, args json.RawMessage, tc *agent.ToolContext) (*agent.ToolResult, error) {
			var input struct {
				Text string `json:"text"`
			}
			json.Unmarshal(args, &input)

			if err := a.Type(input.Text); err != nil {
				return nil, err
			}
			return &agent.ToolResult{Content: fmt.Sprintf("typed: %s", input.Text)}, nil
		},
	}
}

// ─── Hotkey Tool ───

func newHotkeyTool(a *action.Engine) *agent.Tool {
	return &agent.Tool{
		Name:        "hotkey",
		Description: "Press a keyboard shortcut (e.g. ctrl+c, alt+tab)",
		InputSchema: json.RawMessage(`{"type":"object","properties":{"keys":{"type":"array","items":{"type":"string"},"description":"key combination, e.g. [\"ctrl\",\"c\"]"}},"required":["keys"]}`),
		Execute: func(ctx context.Context, args json.RawMessage, tc *agent.ToolContext) (*agent.ToolResult, error) {
			var input struct {
				Keys []string `json:"keys"`
			}
			json.Unmarshal(args, &input)

			if err := a.HotKey(input.Keys...); err != nil {
				return nil, err
			}
			return &agent.ToolResult{Content: fmt.Sprintf("pressed hotkey: %v", input.Keys)}, nil
		},
	}
}

// ─── Scroll Tool ───

func newScrollTool(a *action.Engine) *agent.Tool {
	return &agent.Tool{
		Name:        "scroll",
		Description: "Scroll at a position (positive=up, negative=down)",
		InputSchema: json.RawMessage(`{"type":"object","properties":{"x":{"type":"integer"},"y":{"type":"integer"},"delta":{"type":"integer"}},"required":["x","y","delta"]}`),
		Execute: func(ctx context.Context, args json.RawMessage, tc *agent.ToolContext) (*agent.ToolResult, error) {
			var input struct {
				X     int `json:"x"`
				Y     int `json:"y"`
				Delta int `json:"delta"`
			}
			json.Unmarshal(args, &input)

			if err := a.Scroll(input.X, input.Y, input.Delta); err != nil {
				return nil, err
			}
			return &agent.ToolResult{Content: fmt.Sprintf("scrolled %d at (%d,%d)", input.Delta, input.X, input.Y)}, nil
		},
	}
}

// ─── Wait Tool ───

func newWaitTool() *agent.Tool {
	return &agent.Tool{
		Name:        "wait",
		Description: "Wait for a specified duration in milliseconds",
		InputSchema: json.RawMessage(`{"type":"object","properties":{"ms":{"type":"integer","description":"milliseconds to wait"}},"required":["ms"]}`),
		Execute: func(ctx context.Context, args json.RawMessage, tc *agent.ToolContext) (*agent.ToolResult, error) {
			var input struct {
				MS int `json:"ms"`
			}
			json.Unmarshal(args, &input)

			if input.MS <= 0 {
				input.MS = 1000
			}
			if input.MS > 30000 {
				input.MS = 30000
			}

			select {
			case <-time.After(time.Duration(input.MS) * time.Millisecond):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
			return &agent.ToolResult{Content: fmt.Sprintf("waited %dms", input.MS)}, nil
		},
	}
}

// ─── Verify Tool ───

func newVerifyTool(m *monitor.Checker) *agent.Tool {
	return &agent.Tool{
		Name:        "verify",
		Description: "Verify that target text appears on screen within a timeout",
		InputSchema: json.RawMessage(`{"type":"object","properties":{"text":{"type":"string","description":"text to verify on screen"},"timeout_ms":{"type":"integer","description":"timeout in ms, default 5000"}},"required":["text"]}`),
		Execute: func(ctx context.Context, args json.RawMessage, tc *agent.ToolContext) (*agent.ToolResult, error) {
			var input struct {
				Text      string `json:"text"`
				TimeoutMS int    `json:"timeout_ms"`
			}
			json.Unmarshal(args, &input)

			if input.TimeoutMS == 0 {
				input.TimeoutMS = 5000
			}

			result, err := m.VerifyByOCR(input.Text, input.TimeoutMS)
			if err != nil {
				return nil, err
			}

			return &agent.ToolResult{
				Content: fmt.Sprintf("verify \"%s\": success=%v, method=%s, confidence=%.2f, %s",
					input.Text, result.Success, result.Method, result.Confidence, result.Message),
				IsError: !result.Success,
			}, nil
		},
	}
}

// ─── Simulate Tool ───

func newSimulateTool(v *vision.Engine) *agent.Tool {
	return &agent.Tool{
		Name:        "simulate",
		Description: "Simulate a step by annotating the target on a screenshot (does NOT execute the action)",
		InputSchema: json.RawMessage(`{"type":"object","properties":{"action":{"type":"string","description":"action type (click/input/verify)"},"target":{"type":"string","description":"target description"}},"required":["action","target"]}`),
		Execute: func(ctx context.Context, args json.RawMessage, tc *agent.ToolContext) (*agent.ToolResult, error) {
			var input struct {
				Action string `json:"action"`
				Target string `json:"target"`
			}
			json.Unmarshal(args, &input)

			img, err := v.CaptureScreen()
			if err != nil {
				return nil, err
			}

			results, err := v.Recognize(img)
			if err != nil {
				return nil, err
			}

			// Find target in OCR results
			for _, r := range results {
				if containsIgnoreCase(r.Text, input.Target) {
					return &agent.ToolResult{
						Content:    fmt.Sprintf("simulated %s on \"%s\" at (%d,%d)", input.Action, r.Text, r.Rect.X+r.Rect.Width/2, r.Rect.Y+r.Rect.Height/2),
						Screenshot: string(img),
					}, nil
				}
			}

			return &agent.ToolResult{
				Content:    fmt.Sprintf("simulated %s on \"%s\" — target not found, needs manual review", input.Action, input.Target),
				Screenshot: string(img),
				IsError:    true,
			}, nil
		},
	}
}

// ─── Window Tool ───

func newWindowTool(a *action.Engine) *agent.Tool {
	return &agent.Tool{
		Name:        "window",
		Description: "List visible windows or focus a specific window",
		InputSchema: json.RawMessage(`{"type":"object","properties":{"action":{"type":"string","enum":["list","focus"]},"handle":{"type":"string","description":"window handle for focus action"}},"required":["action"]}`),
		Execute: func(ctx context.Context, args json.RawMessage, tc *agent.ToolContext) (*agent.ToolResult, error) {
			var input struct {
				Action string `json:"action"`
				Handle string `json:"handle"`
			}
			json.Unmarshal(args, &input)

			switch input.Action {
			case "list":
				windows, err := a.ListWindows()
				if err != nil {
					return nil, err
				}
				var lines []string
				for _, w := range windows {
					lines = append(lines, fmt.Sprintf("[%s] %s (%s)", w.Handle, w.Title, w.Process))
				}
				return &agent.ToolResult{
					Content: fmt.Sprintf("found %d windows:\n%s", len(windows), joinLines(lines)),
				}, nil

			case "focus":
				if err := a.FocusWindow(input.Handle); err != nil {
					return nil, err
				}
				return &agent.ToolResult{Content: fmt.Sprintf("focused window %s", input.Handle)}, nil
			}

			return &agent.ToolResult{Content: "unknown window action", IsError: true}, nil
		},
	}
}

// ─── Helpers ───

func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) > 0 && containsLower(lowerStr(s), lowerStr(substr)))
}

func containsLower(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func lowerStr(s string) string {
	b := make([]byte, len(s))
	for i := range s {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		b[i] = c
	}
	return string(b)
}

func joinLines(lines []string) string {
	result := ""
	for i, l := range lines {
		if i > 0 {
			result += "\n"
		}
		result += l
	}
	return result
}
