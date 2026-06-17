package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

// Tool defines an executable action available to the Agent
type Tool struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputSchema json.RawMessage `json:"input_schema"`
	Execute     func(ctx context.Context, args json.RawMessage, tc *ToolContext) (*ToolResult, error)
}

// ToolContext provides execution context to tools
type ToolContext struct {
	SessionID    string
	WindowHandle string
	WindowTitle  string
	Events       *EventBus
}

// ToolResult is the output of a tool execution
type ToolResult struct {
	Content    string `json:"content"`
	Screenshot string `json:"screenshot,omitempty"` // base64
	Annotation string `json:"annotation,omitempty"` // base64 overlay
	IsError    bool   `json:"is_error"`
}

// ToolRegistry manages all available tools and their visibility per agent
type ToolRegistry struct {
	tools    map[string]*Tool
	agentMap map[string][]string // agent name → tool names
	mu       sync.RWMutex
}

func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools:    make(map[string]*Tool),
		agentMap: make(map[string][]string),
	}
}

// Register adds a tool to the registry
func (r *ToolRegistry) Register(tool *Tool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tools[tool.Name] = tool
}

// AssignAgentTools sets the visible tool list for an agent
func (r *ToolRegistry) AssignAgentTools(agent string, toolNames []string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.agentMap[agent] = toolNames
}

// ToolsForAgent returns the tools visible to a specific agent
func (r *ToolRegistry) ToolsForAgent(agent string) []*Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names, ok := r.agentMap[agent]
	if !ok {
		// Return all tools if agent has no specific assignment
		result := make([]*Tool, 0, len(r.tools))
		for _, t := range r.tools {
			result = append(result, t)
		}
		return result
	}

	result := make([]*Tool, 0, len(names))
	for _, name := range names {
		if t, ok := r.tools[name]; ok {
			result = append(result, t)
		}
	}
	return result
}

// Get returns a tool by name
func (r *ToolRegistry) Get(name string) (*Tool, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.tools[name]
	return t, ok
}

// All returns all registered tools
func (r *ToolRegistry) All() []*Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*Tool, 0, len(r.tools))
	for _, t := range r.tools {
		result = append(result, t)
	}
	return result
}

// ExecuteTool runs a tool by name with the given arguments
func (r *ToolRegistry) ExecuteTool(ctx context.Context, name string, args json.RawMessage, tc *ToolContext) (*ToolResult, error) {
	tool, ok := r.Get(name)
	if !ok {
		return nil, fmt.Errorf("unknown tool: %s", name)
	}
	return tool.Execute(ctx, args, tc)
}
