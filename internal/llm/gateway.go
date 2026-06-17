package llm

import (
	"context"
	"encoding/json"
	"fmt"
)

// ─── Anthropic-style Message Types (internal unified model) ───

type Message struct {
	Role    string         `json:"role"` // "user" / "assistant"
	Content []ContentBlock `json:"content"`
}

type ContentBlock interface {
	contentType() string
}

type TextBlock struct {
	Type string `json:"type"` // "text"
	Text string `json:"text"`
}

func (t *TextBlock) contentType() string { return "text" }

type ImageBlock struct {
	Type   string      `json:"type"` // "image"
	Source ImageSource `json:"source"`
}

type ImageSource struct {
	Type      string `json:"type"` // "base64"
	MediaType string `json:"media_type"`
	Data      string `json:"data"`
}

func (i *ImageBlock) contentType() string { return "image" }

type ToolUseBlock struct {
	Type  string          `json:"type"` // "tool_use"
	ID    string          `json:"id"`
	Name  string          `json:"name"`
	Input json.RawMessage `json:"input"`
}

func (t *ToolUseBlock) contentType() string { return "tool_use" }

type ToolResultBlock struct {
	Type      string `json:"type"` // "tool_result"
	ToolUseID string `json:"tool_use_id"`
	Content   string `json:"content"`
	IsError   bool   `json:"is_error"`
}

func (t *ToolResultBlock) contentType() string { return "tool_result" }

// ─── Tool Schema ───

type ToolSchema struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputSchema json.RawMessage `json:"input_schema"`
}

// ─── Stream Request ───

type StreamRequest struct {
	Model     string       `json:"model"`
	System    string       `json:"system"`
	Messages  []Message    `json:"messages"`
	Tools     []ToolSchema `json:"tools,omitempty"`
	MaxTokens int          `json:"max_tokens"`
}

// ─── Stream Events (Anthropic SSE style) ───

type StreamEvent interface {
	eventType() string
}

type MessageStartEvent struct {
	Message Message
}

func (e *MessageStartEvent) eventType() string { return "message_start" }

type ContentBlockStartEvent struct {
	Index int
	Block ContentBlock
}

func (e *ContentBlockStartEvent) eventType() string { return "content_block_start" }

type ContentBlockDeltaEvent struct {
	Index int
	Delta Delta
}

func (e *ContentBlockDeltaEvent) eventType() string { return "content_block_delta" }

type ContentBlockStopEvent struct {
	Index int
}

func (e *ContentBlockStopEvent) eventType() string { return "content_block_stop" }

type MessageDeltaEvent struct {
	StopReason string
	Usage      TokenUsage
}

func (e *MessageDeltaEvent) eventType() string { return "message_delta" }

type LLMErrorEvent struct {
	Err error
}

func (e *LLMErrorEvent) eventType() string { return "error" }

// ─── Delta Types ───

type Delta interface {
	deltaType() string
}

type TextDelta struct {
	Text string
}

func (d TextDelta) deltaType() string { return "text_delta" }

type InputJSONDelta struct {
	PartialJSON string
}

func (d InputJSONDelta) deltaType() string { return "input_json_delta" }

type ThinkingDelta struct {
	Thinking string
}

func (d ThinkingDelta) deltaType() string { return "thinking_delta" }

// ─── Token Usage ───

type TokenUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// ─── Provider Interface ───

type Provider interface {
	Name() string
	Available() bool
	Stream(ctx context.Context, req StreamRequest) (<-chan StreamEvent, error)
}

// ─── Gateway ───

type Gateway struct {
	providers []Provider
}

func NewGateway(providers ...Provider) *Gateway {
	return &Gateway{providers: providers}
}

// Stream selects the first available provider and streams the response.
// Falls back to the next provider on failure.
func (g *Gateway) Stream(ctx context.Context, req StreamRequest) (<-chan StreamEvent, error) {
	for _, p := range g.providers {
		if !p.Available() {
			continue
		}
		stream, err := p.Stream(ctx, req)
		if err != nil {
			continue
		}
		return stream, nil
	}
	return nil, fmt.Errorf("no available LLM provider")
}

// AvailableProviders returns the list of provider names and their availability
func (g *Gateway) AvailableProviders() []ProviderStatus {
	var result []ProviderStatus
	for _, p := range g.providers {
		result = append(result, ProviderStatus{
			Name:      p.Name(),
			Available: p.Available(),
		})
	}
	return result
}

type ProviderStatus struct {
	Name      string `json:"name"`
	Available bool   `json:"available"`
}
