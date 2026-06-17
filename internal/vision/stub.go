package vision

import "github.com/yc446833448/VisuTask/internal/model"

// StubCapturer is a placeholder implementation for development
type StubCapturer struct{}

func NewStubCapturer() *StubCapturer { return &StubCapturer{} }

func (s *StubCapturer) Capture() ([]byte, error) {
	return nil, nil // TODO: implement with robotgo/screenshot
}

func (s *StubCapturer) CaptureRegion(rect model.Rect) ([]byte, error) {
	return nil, nil // TODO: implement
}

func (s *StubCapturer) CaptureWindow(handle string) ([]byte, error) {
	return nil, nil // TODO: implement
}

// StubRecognizer is a placeholder implementation for development
type StubRecognizer struct{}

func NewStubRecognizer() *StubRecognizer { return &StubRecognizer{} }

func (s *StubRecognizer) Recognize(image []byte) ([]OCRResult, error) {
	return nil, nil // TODO: implement with gosseract
}

// StubDetector is a placeholder implementation for development
type StubDetector struct{}

func NewStubDetector() *StubDetector { return &StubDetector{} }

func (s *StubDetector) Detect(image []byte) ([]Control, error) {
	return nil, nil // TODO: implement
}
