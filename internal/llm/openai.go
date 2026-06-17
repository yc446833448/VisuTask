package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const openaiAPIURL = "https://api.openai.com/v1/chat/completions"

// OpenAIAdapter wraps OpenAI-compatible APIs and converts to Anthropic-style events
type OpenAIAdapter struct {
	apiKey  string
	baseURL string
	client  *http.Client
	name    string
}

func NewOpenAIAdapter() *OpenAIAdapter {
	return &OpenAIAdapter{
		apiKey:  os.Getenv("OPENAI_API_KEY"),
		baseURL: getEnvDefault("OPENAI_BASE_URL", openaiAPIURL),
		client:  &http.Client{},
		name:    "openai",
	}
}

// NewOllamaAdapter creates an adapter for Ollama's OpenAI-compatible API
func NewOllamaAdapter() *OpenAIAdapter {
	return &OpenAIAdapter{
		apiKey:  "ollama",
		baseURL: getEnvDefault("OLLAMA_BASE_URL", "http://localhost:11434/v1/chat/completions"),
		client:  &http.Client{},
		name:    "ollama",
	}
}

func (a *OpenAIAdapter) Name() string   { return a.name }
func (a *OpenAIAdapter) Available() bool { return a.apiKey != "" }

// SetConfig updates the OpenAI adapter configuration from config file
func (a *OpenAIAdapter) SetConfig(apiKey, baseURL, name string) {
	if apiKey != "" {
		a.apiKey = apiKey
	}
	if baseURL != "" {
		a.baseURL = baseURL
	}
	if name != "" {
		a.name = name
	}
}

// SetOllamaConfig updates the Ollama adapter configuration
func (a *OpenAIAdapter) SetOllamaConfig(baseURL string) {
	if baseURL != "" {
		// Ensure the URL includes the chat completions path
		if !strings.HasSuffix(baseURL, "/chat/completions") {
			baseURL = strings.TrimSuffix(baseURL, "/") + "/v1/chat/completions"
		}
		a.baseURL = baseURL
	}
}

// ─── OpenAI API Types ───

type openaiRequest struct {
	Model    string           `json:"model"`
	Messages []openaiMessage  `json:"messages"`
	Tools    []openaiTool     `json:"tools,omitempty"`
	Stream   bool             `json:"stream"`
}

type openaiMessage struct {
	Role       string          `json:"role"`
	Content    interface{}     `json:"content,omitempty"`
	ToolCalls  []openaiToolCall `json:"tool_calls,omitempty"`
	ToolCallID string          `json:"tool_call_id,omitempty"`
}

type openaiTool struct {
	Type     string         `json:"type"`
	Function openaiFunction `json:"function"`
}

type openaiFunction struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Parameters  json.RawMessage `json:"parameters"`
}

type openaiToolCall struct {
	ID       string           `json:"id"`
	Type     string           `json:"type"`
	Function openaiFuncCall   `json:"function"`
}

type openaiFuncCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

func (a *OpenAIAdapter) Stream(ctx context.Context, req StreamRequest) (<-chan StreamEvent, error) {
	// Convert Anthropic-style request to OpenAI format
	openaiMsgs := toOpenAIMessages(req.Messages)

	tools := make([]openaiTool, len(req.Tools))
	for i, t := range req.Tools {
		tools[i] = openaiTool{
			Type: "function",
			Function: openaiFunction{
				Name:        t.Name,
				Description: t.Description,
				Parameters:  t.InputSchema,
			},
		}
	}

	// Prepend system message
	if req.System != "" {
		openaiMsgs = append([]openaiMessage{{Role: "system", Content: req.System}}, openaiMsgs...)
	}

	body := openaiRequest{
		Model:    req.Model,
		Messages: openaiMsgs,
		Tools:    tools,
		Stream:   true,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", a.baseURL, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if a.name != "ollama" {
		httpReq.Header.Set("Authorization", "Bearer "+a.apiKey)
	}

	resp, err := a.client.Do(httpReq)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("openai API error %d: %s", resp.StatusCode, string(respBody))
	}

	ch := make(chan StreamEvent, 64)
	go a.readSSEStream(resp, ch)
	return ch, nil
}

func (a *OpenAIAdapter) readSSEStream(resp *http.Response, ch chan<- StreamEvent) {
	defer resp.Body.Close()
	defer close(ch)

	scanner := bufio.NewScanner(resp.Body)
	toolCallIndex := -1
	var currentToolName string
	var toolArgBuf strings.Builder

	for scanner.Scan() {
		line := scanner.Text()

		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			// Flush any pending tool call
			if currentToolName != "" {
				ch <- &ContentBlockStopEvent{Index: toolCallIndex}
			}
			return
		}

		var chunk struct {
			Choices []struct {
				Delta struct {
					Role      string           `json:"role"`
					Content   string           `json:"content"`
					ToolCalls []openaiToolCall `json:"tool_calls"`
				} `json:"delta"`
				FinishReason *string `json:"finish_reason"`
			} `json:"choices"`
		}

		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue
		}

		if len(chunk.Choices) == 0 {
			continue
		}

		choice := chunk.Choices[0]

		// Text content
		if choice.Delta.Content != "" {
			ch <- &ContentBlockDeltaEvent{
				Index: 0,
				Delta: TextDelta{Text: choice.Delta.Content},
			}
		}

		// Tool calls
		for _, tc := range choice.Delta.ToolCalls {
			if tc.ID != "" {
				// New tool call starts
				if currentToolName != "" {
					ch <- &ContentBlockStopEvent{Index: toolCallIndex}
				}
				toolCallIndex++
				currentToolName = tc.Function.Name
				toolArgBuf.Reset()

				ch <- &ContentBlockStartEvent{
					Index: toolCallIndex,
					Block: &ToolUseBlock{
						Type: "tool_use",
						ID:   tc.ID,
						Name: tc.Function.Name,
					},
				}
			}

			if tc.Function.Arguments != "" {
				toolArgBuf.WriteString(tc.Function.Arguments)
				ch <- &ContentBlockDeltaEvent{
					Index: toolCallIndex,
					Delta: InputJSONDelta{PartialJSON: tc.Function.Arguments},
				}
			}
		}

		// Finish
		if choice.FinishReason != nil {
			reason := *choice.FinishReason
			stopReason := "end_turn"
			if reason == "tool_calls" {
				stopReason = "tool_use"
			}
			if currentToolName != "" {
				ch <- &ContentBlockStopEvent{Index: toolCallIndex}
				currentToolName = ""
			}
			ch <- &MessageDeltaEvent{StopReason: stopReason}
		}
	}
}

// toOpenAIMessages converts Anthropic-style messages to OpenAI format
func toOpenAIMessages(messages []Message) []openaiMessage {
	var result []openaiMessage

	for _, msg := range messages {
		var contentParts []string
		var toolCalls []openaiToolCall
		var toolResultContent string

		for _, block := range msg.Content {
			switch b := block.(type) {
			case *TextBlock:
				contentParts = append(contentParts, b.Text)
			case *ToolUseBlock:
				toolCalls = append(toolCalls, openaiToolCall{
					ID:   b.ID,
					Type: "function",
					Function: openaiFuncCall{
						Name:      b.Name,
						Arguments: string(b.Input),
					},
				})
			case *ToolResultBlock:
				toolResultContent = b.Content
			case *ImageBlock:
				// OpenAI vision: would need multipart content
				// For now, skip images in OpenAI adapter
			}
		}

		om := openaiMessage{Role: msg.Role}

		if msg.Role == "assistant" && len(toolCalls) > 0 {
			om.Content = strings.Join(contentParts, "\n")
			om.ToolCalls = toolCalls
		} else if msg.Role == "user" && toolResultContent != "" {
			om.Role = "tool"
			om.Content = toolResultContent
			// ToolCallID would need to be tracked
		} else {
			om.Content = strings.Join(contentParts, "\n")
		}

		result = append(result, om)
	}

	return result
}

func getEnvDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
