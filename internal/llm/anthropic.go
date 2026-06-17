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

const anthropicAPIURL = "https://api.anthropic.com/v1/messages"

type AnthropicProvider struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func NewAnthropicProvider() *AnthropicProvider {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	baseURL := os.Getenv("ANTHROPIC_BASE_URL")
	if baseURL == "" {
		baseURL = anthropicAPIURL
	}
	return &AnthropicProvider{
		apiKey:  apiKey,
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (p *AnthropicProvider) Name() string      { return "anthropic" }
func (p *AnthropicProvider) Available() bool    { return p.apiKey != "" }

// SetConfig updates the provider configuration from config file
func (p *AnthropicProvider) SetConfig(apiKey, baseURL string) {
	if apiKey != "" {
		p.apiKey = apiKey
	}
	if baseURL != "" {
		p.baseURL = baseURL
	}
}

type anthropicRequest struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	System    string             `json:"system,omitempty"`
	Messages  []anthropicMessage `json:"messages"`
	Tools     []anthropicTool    `json:"tools,omitempty"`
	Stream    bool               `json:"stream"`
}

type anthropicMessage struct {
	Role    string          `json:"role"`
	Content json.RawMessage `json:"content"`
}

type anthropicTool struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputSchema json.RawMessage `json:"input_schema"`
}

func (p *AnthropicProvider) Stream(ctx context.Context, req StreamRequest) (<-chan StreamEvent, error) {
	// Convert internal messages to Anthropic format
	messages, err := toAnthropicMessages(req.Messages)
	if err != nil {
		return nil, fmt.Errorf("convert messages: %w", err)
	}

	tools := make([]anthropicTool, len(req.Tools))
	for i, t := range req.Tools {
		tools[i] = anthropicTool{
			Name:        t.Name,
			Description: t.Description,
			InputSchema: t.InputSchema,
		}
	}

	maxTokens := req.MaxTokens
	if maxTokens == 0 {
		maxTokens = 4096
	}

	body := anthropicRequest{
		Model:     req.Model,
		MaxTokens: maxTokens,
		System:    req.System,
		Messages:  messages,
		Tools:     tools,
		Stream:    true,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.baseURL, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", p.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("anthropic API error %d: %s", resp.StatusCode, string(respBody))
	}

	ch := make(chan StreamEvent, 64)
	go p.readSSEStream(resp, ch)
	return ch, nil
}

func (p *AnthropicProvider) readSSEStream(resp *http.Response, ch chan<- StreamEvent) {
	defer resp.Body.Close()
	defer close(ch)

	scanner := bufio.NewScanner(resp.Body)
	var currentToolID string
	var currentToolName string
	var toolInputBuf strings.Builder
	blockIndex := -1

	for scanner.Scan() {
		line := scanner.Text()

		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			return
		}

		var event map[string]interface{}
		if err := json.Unmarshal([]byte(data), &event); err != nil {
			continue
		}

		eventType, _ := event["type"].(string)

		switch eventType {
		case "message_start":
			ch <- &MessageStartEvent{}

		case "content_block_start":
			idx, _ := event["index"].(float64)
			blockIndex = int(idx)

			block, _ := event["content_block"].(map[string]interface{})
			blockType, _ := block["type"].(string)

			switch blockType {
			case "text":
				ch <- &ContentBlockStartEvent{
					Index: blockIndex,
					Block: &TextBlock{Type: "text"},
				}
			case "tool_use":
				currentToolID, _ = block["id"].(string)
				currentToolName, _ = block["name"].(string)
				toolInputBuf.Reset()
				ch <- &ContentBlockStartEvent{
					Index: blockIndex,
					Block: &ToolUseBlock{
						Type: "tool_use",
						ID:   currentToolID,
						Name: currentToolName,
					},
				}
			}

		case "content_block_delta":
			delta, _ := event["delta"].(map[string]interface{})
			deltaType, _ := delta["type"].(string)

			switch deltaType {
			case "text_delta":
				text, _ := delta["text"].(string)
				ch <- &ContentBlockDeltaEvent{
					Index: blockIndex,
					Delta: TextDelta{Text: text},
				}
			case "input_json_delta":
				partial, _ := delta["partial_json"].(string)
				toolInputBuf.WriteString(partial)
				ch <- &ContentBlockDeltaEvent{
					Index: blockIndex,
					Delta: InputJSONDelta{PartialJSON: partial},
				}
			case "thinking_delta":
				thinking, _ := delta["thinking"].(string)
				ch <- &ContentBlockDeltaEvent{
					Index: blockIndex,
					Delta: ThinkingDelta{Thinking: thinking},
				}
			}

		case "content_block_stop":
			if currentToolName != "" {
				// Emit the tool use block with accumulated input
				input := json.RawMessage(toolInputBuf.String())
				if len(input) == 0 {
					input = json.RawMessage("{}")
				}
				ch <- &ContentBlockDeltaEvent{
					Index: blockIndex,
					Delta: InputJSONDelta{PartialJSON: toolInputBuf.String()},
				}
				// Reset for next tool
				currentToolID = ""
				currentToolName = ""
				toolInputBuf.Reset()
			}
			ch <- &ContentBlockStopEvent{Index: blockIndex}

		case "message_delta":
			delta, _ := event["delta"].(map[string]interface{})
			stopReason, _ := delta["stop_reason"].(string)

			usage, _ := event["usage"].(map[string]interface{})
			var tokenUsage TokenUsage
			if usage != nil {
				if v, ok := usage["output_tokens"].(float64); ok {
					tokenUsage.OutputTokens = int(v)
				}
			}

			ch <- &MessageDeltaEvent{
				StopReason: stopReason,
				Usage:      tokenUsage,
			}

		case "error":
			errMsg, _ := event["error"].(map[string]interface{})
			msg, _ := errMsg["message"].(string)
			ch <- &LLMErrorEvent{Err: fmt.Errorf("anthropic stream error: %s", msg)}
			return
		}
	}
}

// toAnthropicMessages converts internal Messages to Anthropic API format
func toAnthropicMessages(messages []Message) ([]anthropicMessage, error) {
	result := make([]anthropicMessage, len(messages))
	for i, msg := range messages {
		content, err := json.Marshal(msg.Content)
		if err != nil {
			return nil, err
		}
		result[i] = anthropicMessage{
			Role:    msg.Role,
			Content: content,
		}
	}
	return result, nil
}
