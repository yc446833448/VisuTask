package vision

import "github.com/yc446833448/VisuTask/internal/model"

// Engine provides screen capture and OCR capabilities
type Engine struct {
	capture   Capturer
	ocr       Recognizer
	detector  Detector
}

func NewEngine(c Capturer, o Recognizer, d Detector) *Engine {
	return &Engine{capture: c, ocr: o, detector: d}
}

// CaptureScreen takes a full screen screenshot
func (e *Engine) CaptureScreen() ([]byte, error) {
	return e.capture.Capture()
}

// CaptureRegion takes a screenshot of a specific region
func (e *Engine) CaptureRegion(rect model.Rect) ([]byte, error) {
	return e.capture.CaptureRegion(rect)
}

// CaptureWindow takes a screenshot of a specific window
func (e *Engine) CaptureWindow(handle string) ([]byte, error) {
	return e.capture.CaptureWindow(handle)
}

// Recognize performs OCR on an image
func (e *Engine) Recognize(image []byte) ([]OCRResult, error) {
	return e.ocr.Recognize(image)
}

// DetectControls finds interactive UI elements in an image
func (e *Engine) DetectControls(image []byte) ([]Control, error) {
	return e.detector.Detect(image)
}

// ─── Interfaces ───

type Capturer interface {
	Capture() ([]byte, error)
	CaptureRegion(rect model.Rect) ([]byte, error)
	CaptureWindow(handle string) ([]byte, error)
}

type Recognizer interface {
	Recognize(image []byte) ([]OCRResult, error)
}

type Detector interface {
	Detect(image []byte) ([]Control, error)
}

// ─── Data Types ───

type OCRResult struct {
	Text       string    `json:"text"`
	Rect       model.Rect `json:"rect"`
	Confidence float64   `json:"confidence"`
}

type Control struct {
	Type       string    `json:"type"` // button / input / dropdown / checkbox / label
	Rect       model.Rect `json:"rect"`
	Text       string    `json:"text"`
	Confidence float64   `json:"confidence"`
}
