package vision

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/yc446833448/VisuTask/internal/model"
)

// RemoteOCRRecognizer calls a remote OCR service via HTTP API
type RemoteOCRRecognizer struct {
	endpoint string
	apiKey   string
	client   *http.Client
}

func NewRemoteOCRRecognizer() (*RemoteOCRRecognizer, error) {
	endpoint := os.Getenv("OCR_ENDPOINT")
	if endpoint == "" {
		endpoint = "https://api.deepseek.com/v1/ocr"
	}

	apiKey := os.Getenv("OCR_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OCR_API_KEY not set")
	}

	return &RemoteOCRRecognizer{
		endpoint: endpoint,
		apiKey:   apiKey,
		client:   &http.Client{Timeout: 30 * time.Second},
	}, nil
}

// NewRemoteOCRRecognizerWithConfig creates an OCR recognizer from explicit config values
func NewRemoteOCRRecognizerWithConfig(endpoint, apiKey string) (*RemoteOCRRecognizer, error) {
	if endpoint == "" {
		endpoint = os.Getenv("OCR_ENDPOINT")
	}
	if apiKey == "" {
		apiKey = os.Getenv("OCR_API_KEY")
	}
	if endpoint == "" || apiKey == "" {
		return nil, fmt.Errorf("OCR endpoint or API key not configured")
	}
	return &RemoteOCRRecognizer{
		endpoint: endpoint,
		apiKey:   apiKey,
		client:   &http.Client{Timeout: 30 * time.Second},
	}, nil
}

// ocrRequest is the request body sent to the remote OCR service
type ocrRequest struct {
	Image string `json:"image"` // base64 encoded PNG
}

// ocrResponse is the expected response from the OCR service
type ocrResponse struct {
	Results []ocrItem `json:"results"`
	Error   string    `json:"error,omitempty"`
}

type ocrItem struct {
	Text       string  `json:"text"`
	X          int     `json:"x"`
	Y          int     `json:"y"`
	Width      int     `json:"width"`
	Height     int     `json:"height"`
	Confidence float64 `json:"confidence"` // 0.0 - 1.0
}

// Recognize sends the image to the remote OCR service and parses the result
func (r *RemoteOCRRecognizer) Recognize(imageData []byte) ([]OCRResult, error) {
	// Encode image to base64
	b64Image := base64.StdEncoding.EncodeToString(imageData)

	reqBody := ocrRequest{Image: b64Image}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	// Build HTTP request
	req, err := http.NewRequest("POST", r.endpoint, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+r.apiKey)

	// Send request
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ocr request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ocr API error %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var ocrResp ocrResponse
	if err := json.Unmarshal(body, &ocrResp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	if ocrResp.Error != "" {
		return nil, fmt.Errorf("ocr service error: %s", ocrResp.Error)
	}

	// Convert to internal format
	results := make([]OCRResult, 0, len(ocrResp.Results))
	for _, item := range ocrResp.Results {
		if item.Confidence < 0.3 {
			continue // skip low confidence
		}
		results = append(results, OCRResult{
			Text: item.Text,
			Rect: model.Rect{
				X:      item.X,
				Y:      item.Y,
				Width:  item.Width,
				Height: item.Height,
			},
			Confidence: item.Confidence,
		})
	}

	return results, nil
}

func (r *RemoteOCRRecognizer) Close() error { return nil }
