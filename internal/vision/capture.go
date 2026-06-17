package vision

import (
	"bytes"
	"fmt"
	"image"
	"image/png"

	"github.com/kbinani/screenshot"
	"github.com/yc446833448/VisuTask/internal/model"
)

// ScreenCapturer captures screen using kbinani/screenshot
type ScreenCapturer struct{}

func NewScreenCapturer() *ScreenCapturer {
	return &ScreenCapturer{}
}

// Capture takes a full screen screenshot, returns PNG bytes
func (c *ScreenCapturer) Capture() ([]byte, error) {
	n := screenshot.NumActiveDisplays()
	if n == 0 {
		return nil, fmt.Errorf("no active displays found")
	}

	// Capture the primary display (index 0)
	bounds := screenshot.GetDisplayBounds(0)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		return nil, fmt.Errorf("capture screen: %w", err)
	}

	return imageToPNG(img)
}

// CaptureRegion captures a specific screen region
func (c *ScreenCapturer) CaptureRegion(rect model.Rect) ([]byte, error) {
	bounds := image.Rect(rect.X, rect.Y, rect.X+rect.Width, rect.Y+rect.Height)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		return nil, fmt.Errorf("capture region: %w", err)
	}
	return imageToPNG(img)
}

// CaptureWindow captures a specific window by handle
// Note: kbinani/screenshot doesn't support window-specific capture directly,
// so we capture the full screen and crop to window bounds if available.
func (c *ScreenCapturer) CaptureWindow(handle string) ([]byte, error) {
	// Fallback: capture full screen
	// A more advanced implementation would use Win32 API to get window rect
	return c.Capture()
}

// imageToPNG converts an image.Image to PNG bytes
func imageToPNG(img *image.RGBA) ([]byte, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("encode png: %w", err)
	}
	return buf.Bytes(), nil
}
