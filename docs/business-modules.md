# 业务模块设计

> VisuTask Go 后端核心业务模块规划

---

## 核心概念

```
脚本 (Script)                     任务 (Task)                       执行 (Execution)
────────────                     ─────────────                     ─────────────
Agent 生成的自动化模板              脚本 + 窗口句柄 + 参数的实例         一次任务的运行过程
可复用、可编辑                     一个脚本可创建多个任务               并发管理，最多 N 个
保存步骤序列                      绑定不同窗口可同时跑               记录成功/失败/截图
```

---

## 用户等级系统

| 等级 | 并发任务上限 | 说明 |
|------|------------|------|
| 普通用户 | 3 | 默认 |
| VIP 1 | 5 | |
| VIP 2 | 6 | |
| VIP 3 | 7 | |
| VIP 4 | 8 | |
| VIP 5 | 10 | 最高 |

并发控制器在启动任务时检查当前运行数是否 < 用户上限。

---

## 模块总览

```
internal/
├── agent/          # Agent 调度核心 —— 脚本创建 + 执行
├── vision/         # 视觉感知层 —— 屏幕理解能力
├── action/         # 动作执行层 —— 系统操作能力
├── monitor/        # 监控验证层 —— 执行结果校验
├── llm/            # LLM 网关 —— 多模型统一管理
├── scheduler/      # 定时调度 —— 任务触发
├── hotkey/         # 全局快捷键 —— 快捷触发
├── concurrency/    # 并发管理 —— 任务调度 + 槽位控制
├── model/          # 数据模型 —— 实体定义
└── store/          # 持久化层 —— 数据存储
```

---

## 1. Agent 调度核心 (`internal/agent/`)

Agent 是系统的"大脑"，负责脚本的创建（解析、规划、模拟）和任务的执行。
架构参考 OpenCode 的 Agent Loop 模式，适配为 GUI 自动化场景。

### 1.1 模块总览

```
internal/agent/
├── session.go       # Session 管理 —— Agent 运行实例
├── loop.go          # Agent Loop —— while(true) 核心循环
├── processor.go     # Stream Processor —— LLM 流式事件处理
├── tool.go          # Tool 定义 + 注册表
├── system_prompt.go # System Prompt 构建
├── memory.go        # 上下文管理 + 压缩
├── retry.go         # 重试策略 + 退避
├── safety.go        # Doom Loop 检测 + 安全机制
├── service.go       # Wails Binding 入口
│
├── tools/           # VisuTask 专用工具集
│   ├── capture.go   # 截图工具
│   ├── ocr.go       # OCR 识别工具
│   ├── click.go     # 鼠标点击工具
│   ├── type.go      # 键盘输入工具
│   ├── hotkey.go    # 快捷键工具
│   ├── scroll.go    # 滚动工具
│   ├── wait.go      # 等待工具
│   ├── verify.go    # 步骤验证工具
│   ├── simulate.go  # 模拟演示工具（仅标注，不执行）
│   ├── window.go    # 窗口管理工具
│   └── locate.go    # 目标定位工具（OCR + 图像匹配）
│
└── agents/          # 内置 Agent 角色定义
    ├── planner.go   # 规划 Agent —— 自然语言 → 步骤计划
    ├── executor.go  # 执行 Agent —— 按脚本逐步执行
    └── reviewer.go  # 审查 Agent —— 执行结果校验
```

---

### 1.2 Agent Loop 核心循环

参考 OpenCode 的 `runLoop` 模式，核心是一个 `for` 循环驱动的多轮对话。

```go
// loop.go

type Loop struct {
    session     *Session
    processor   *Processor
    registry    *ToolRegistry
    memory      *Memory
    safety      *Safety
    llm         *llm.Gateway
    events      *EventBus
}

// Run Agent 主循环 —— 每轮: 加载上下文 → 调用 LLM → 处理工具 → 评估结果
func (l *Loop) Run(ctx context.Context) error {
    step := 0
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
        }

        // 1. 加载消息上下文（过滤已压缩的旧消息）
        messages := l.memory.LoadMessages()

        // 2. 检查上下文是否溢出 → 触发压缩
        if l.memory.IsOverflow(messages) {
            if err := l.memory.Compact(ctx, l.llm); err != nil {
                return fmt.Errorf("compaction failed: %w", err)
            }
            continue
        }

        // 3. 解析当前 Agent 角色 + 可用工具集
        agent := l.session.CurrentAgent()
        tools := l.registry.ToolsForAgent(agent.Name)
        system := l.buildSystemPrompt(agent)

        // 4. 通知前端：新一轮开始
        l.events.Emit(EventStepStarted, StepEvent{Step: step})

        // 5. 流式调用 LLM，处理事件流
        result, err := l.processor.Process(ctx, ProcessInput{
            Messages: messages,
            Tools:    tools,
            System:   system,
            Model:    agent.Model,
        })
        if err != nil {
            // 重试判定
            if l.retry.ShouldRetry(err) {
                l.retry.Wait(ctx, err)
                continue
            }
            return err
        }

        // 6. 评估结果，决定下一步
        switch result.Outcome {
        case OutcomeStop:
            // Agent 完成，退出循环
            l.events.Emit(EventCompleted, nil)
            return nil

        case OutcomeCompact:
            // 上下文溢出，压缩后继续
            l.memory.Compact(ctx, l.llm)
            continue

        case OutcomeContinue:
            // 有工具调用，继续下一轮
            step++

            // Doom Loop 检测
            if l.safety.IsDoomLoop(result.ToolCalls) {
                l.events.Emit(EventDoomLoop, DoomLoopEvent{
                    Tool:   result.ToolCalls[0].Name,
                    Count:  l.safety.RepeatCount(),
                })
                // 等待用户确认后继续，或中止
                if err := l.safety.WaitForConfirmation(ctx); err != nil {
                    return err
                }
            }
            continue
        }
    }
}
```

### 循环退出条件

| 条件 | 说明 |
|------|------|
| `OutcomeStop` | LLM 返回最终结果，无工具调用 |
| `ctx.Done()` | 用户取消 / 超时 |
| 最大步数 | 防止无限循环（默认 50 步） |
| 不可重试错误 | 认证失败、参数错误等 |

---

### 1.3 Session 管理

```go
// session.go

type Session struct {
    ID          string
    Type        SessionType       // ScriptCreation / TaskExecution
    Agent       string            // 当前 Agent 角色
    Model       string            // 当前 LLM 模型
    Status      SessionStatus     // idle / busy / paused / completed / failed
    CreatedAt   time.Time

    // 运行时状态
    cancel      context.CancelFunc
    runState    *RunState         // 防止重复执行
}

type SessionType string
const (
    SessionScriptCreation SessionType = "script_creation"  // 创建脚本阶段
    SessionTaskExecution SessionType = "task_execution"   // 执行任务阶段
)

// RunState 防止同一 Session 并发执行 Loop
type RunState struct {
    mu      sync.Mutex
    running bool
}

func (r *RunState) EnsureRunning() (func(), error) {
    r.mu.Lock()
    defer r.mu.Unlock()
    if r.running {
        return nil, ErrSessionAlreadyRunning
    }
    r.running = true
    return func() {
        r.mu.Lock()
        r.running = false
        r.mu.Unlock()
    }, nil
}
```

---

### 1.4 Tool 系统

参考 OpenCode 的声明式 Tool 定义 + 注册表模式。

#### Tool 定义

```go
// tool.go

type Tool struct {
    Name        string
    Description string
    Parameters  json.RawMessage   // JSON Schema
    Execute     func(ctx context.Context, args json.RawMessage, tc *ToolContext) (*ToolResult, error)
}

type ToolContext struct {
    SessionID   string
    WindowHandle uintptr          // 绑定的目标窗口
    Screenshot  []byte            // 最新截图（缓存）
    OCRResults  []vision.OCRResult // 最新 OCR 结果（缓存）
    Events      *EventBus
}

type ToolResult struct {
    Content     string            // 文本结果（返回给 LLM）
    Screenshot  []byte            // 截图（如有）
    Annotation  []byte            // 标注图（如有）
    Metadata    map[string]any    // 额外元数据
}
```

#### Tool 注册表

```go
// tool.go

type ToolRegistry struct {
    tools    map[string]*Tool
    agentMap map[string][]string   // agent name → tool names
}

func NewToolRegistry(v *vision.Engine, a *action.Engine, m *monitor.Checker) *ToolRegistry {
    r := &ToolRegistry{
        tools:    make(map[string]*Tool),
        agentMap: make(map[string][]string),
    }

    // 注册所有工具
    r.Register(NewCaptureTool(v))
    r.Register(NewOCRTool(v))
    r.Register(NewClickTool(a))
    r.Register(NewTypeTool(a))
    r.Register(NewHotkeyTool(a))
    r.Register(NewScrollTool(a))
    r.Register(NewWaitTool())
    r.Register(NewVerifyTool(m))
    r.Register(NewSimulateTool(v))
    r.Register(NewWindowTool(a))
    r.Register(NewLocateTool(v))

    // 按 Agent 角色分配工具
    r.agentMap["planner"]  = []string{"capture", "ocr", "locate", "window"}
    r.agentMap["executor"] = []string{"capture", "ocr", "click", "type", "hotkey", "scroll", "wait", "verify", "locate", "window"}
    r.agentMap["reviewer"] = []string{"capture", "ocr", "verify"}

    return r
}

func (r *ToolRegistry) ToolsForAgent(agent string) []*Tool { ... }
func (r *ToolRegistry) Get(name string) (*Tool, bool) { ... }
```

#### VisuTask 专用工具集

| 工具 | 说明 | Agent 可见性 |
|------|------|-------------|
| `capture` | 截取当前屏幕/指定窗口 | planner, executor, reviewer |
| `ocr` | 对截图执行 OCR 识别 | planner, executor, reviewer |
| `locate` | 定位目标元素（OCR + 图像匹配），返回坐标 | planner, executor |
| `click` | 在指定坐标执行鼠标点击 | executor |
| `type` | 在指定输入框输入文字 | executor |
| `hotkey` | 执行快捷键组合 | executor |
| `scroll` | 在指定区域滚动 | executor |
| `wait` | 等待指定时间或条件满足 | executor |
| `verify` | 验证步骤执行结果（OCR/图像对比） | executor, reviewer |
| `simulate` | 模拟演示（截图标注目标，不真正操作） | planner |
| `window` | 窗口操作（聚焦/移动/列出） | planner, executor |

#### Tool 执行流程

```
LLM 返回 tool_call
    │
    ▼
Permission Check ── deny → 返回拒绝结果给 LLM
    │
    ▼ allow
创建超时 context (默认 30s，可配置)
    │
    ▼
Tool.Execute(ctx, args, toolCtx)
    │
    ├── 成功 → ToolResult{Content: "操作成功..."} → 追加到消息 → 继续循环
    ├── 超时 → ToolResult{Content: "操作超时(30s)"} → 追加到消息 → LLM 决策
    └── 失败 → ToolResult{Content: "操作失败: ..."} → 追加到消息 → LLM 决策重试或调整
```

---

### 1.5 Stream Processor

处理 LLM 流式响应的事件状态机。

```go
// processor.go

type ProcessInput struct {
    Messages []Message
    Tools    []*Tool
    System   string
    Model    string
}

type ProcessResult struct {
    Outcome   Outcome          // stop / compact / continue
    ToolCalls []ToolCallInfo   // 本轮工具调用
    TextDelta string           // 文本增量
}

type Outcome string
const (
    OutcomeStop    Outcome = "stop"     // Agent 完成
    OutcomeCompact Outcome = "compact"  // 需要压缩上下文
    OutcomeContinue Outcome = "continue" // 有工具调用，继续
)

type Processor struct {
    llm       *llm.Gateway
    registry  *ToolRegistry
    events    *EventBus
    safety    *Safety
}

func (p *Processor) Process(ctx context.Context, input ProcessInput) (*ProcessResult, error) {
    // 1. 调用 LLM Gateway 获取 Anthropic 风格流式响应
    stream, err := p.llm.Stream(ctx, llm.StreamRequest{
        Model:    input.Model,
        System:   input.System,
        Messages: input.Messages,
        Tools:    toLLMTools(input.Tools),
    })
    if err != nil {
        return nil, err
    }

    // 2. 处理 Anthropic SSE 事件流
    var result ProcessResult
    var toolCalls []ToolCallInfo
    var textBuf strings.Builder
    var currentToolInput strings.Builder

    for event := range stream {
        switch e := event.(type) {

        case llm.ContentBlockDeltaEvent:
            switch d := e.Delta.(type) {
            case llm.TextDelta:
                textBuf.WriteString(d.Text)
                p.events.Emit(EventTextDelta, d.Text)

            case llm.InputJSONDelta:
                currentToolInput.WriteString(d.PartialJSON)
            }

        case llm.ContentBlockStartEvent:
            if tb, ok := e.Block.(*llm.ToolUseBlock); ok {
                p.events.Emit(EventToolCalled, ToolCalledEvent{Name: tb.Name})
            }

        case llm.ContentBlockStopEvent:
            // ToolUse block 结束 → 执行工具
            if toolUse := p.getToolUse(e.Index); toolUse != nil {
                tool, ok := p.registry.Get(toolUse.Name)
                if !ok {
                    toolCalls = append(toolCalls, ToolCallInfo{
                        Name:   toolUse.Name,
                        Result: &ToolResult{Content: fmt.Sprintf("unknown tool: %s", toolUse.Name)},
                    })
                    continue
                }

                tc := &ToolContext{SessionID: "...", Events: p.events}
                toolResult, err := tool.Execute(ctx, toolUse.Input, tc)
                if err != nil {
                    toolResult = &ToolResult{Content: fmt.Sprintf("error: %v", err)}
                }

                p.events.Emit(EventToolResult, ToolResultEvent{Name: toolUse.Name, Result: toolResult})
                toolCalls = append(toolCalls, ToolCallInfo{
                    Name:   toolUse.Name,
                    Args:   toolUse.Input,
                    Result: toolResult,
                })
            }

        case llm.MessageDeltaEvent:
            if e.StopReason == "tool_use" && len(toolCalls) > 0 {
                result.Outcome = OutcomeContinue
            } else {
                result.Outcome = OutcomeStop
            }

        case llm.ErrorEvent:
            return nil, e.Err
        }
    }

    result.ToolCalls = toolCalls
    result.TextDelta = textBuf.String()
    return &result, nil
}
```

---

### 1.6 LLM 网关（流式，Anthropic 协议）

```go
// internal/llm/gateway.go

// StreamEvent 统一使用 Anthropic SSE 事件类型
type StreamEvent interface{}

type MessageStartEvent struct{ Message Message }
type ContentBlockStartEvent struct{ Index int; Block ContentBlock }
type ContentBlockDeltaEvent struct {
    Index int
    Delta Delta  // TextDelta / InputJSONDelta / ThinkingDelta
}
type ContentBlockStopEvent struct{ Index int }
type MessageDeltaEvent struct{ StopReason string; Usage TokenUsage }
type ErrorEvent struct{ Err error }

type Provider interface {
    Name() string
    Available() bool
    Stream(ctx context.Context, req StreamRequest) (<-chan StreamEvent, error)
}

type Gateway struct {
    providers []Provider        // 按环境变量配置优先级排序
    retry     *RetryPolicy
}

// Stream 自动选择可用 provider，失败时 fallback
func (g *Gateway) Stream(ctx context.Context, req StreamRequest) (<-chan StreamEvent, error) {
    for _, p := range g.providers {
        if !p.Available() {
            continue
        }
        stream, err := p.Stream(ctx, req)
        if err != nil {
            if g.retry.ShouldRetry(err) {
                continue  // 尝试下一个 provider
            }
            return nil, err
        }
        return stream, nil
    }
    return nil, ErrNoAvailableProvider
}
```

#### Provider 适配器

| 文件 | 提供方 | 协议 | 角色 |
|------|--------|------|------|
| `anthropic.go` | Anthropic Claude | Anthropic Messages API (SSE) | **主协议** |
| `openai.go` | OpenAI / 兼容 API | OpenAI Chat Completions (SSE) | 适配层，转换为 Anthropic 内部格式 |
| `ollama.go` | Ollama 本地模型 | OpenAI-compatible (SSE) | 复用 openai.go 适配器 |

---

### 1.7 上下文管理 (Memory)

参考 OpenCode 的两阶段压缩策略。

```go
// memory.go

type Memory struct {
    messages []Message
    mu       sync.RWMutex

    // 配置
    contextWindow int     // 模型上下文窗口大小 (tokens)
    protectRecent int     // 保护最近的 N tokens 不被修剪
    pruneMin      int     // 最少修剪 N tokens 才值得执行
}

// LoadMessages 加载消息（过滤已压缩的）
func (m *Memory) LoadMessages() []Message

// AppendMessage 追加消息
func (m *Memory) AppendMessage(msg Message)

// IsOverflow 检查上下文是否超出模型窗口
func (m *Memory) IsOverflow(messages []Message) bool

// Compact 两阶段压缩
func (m *Memory) Compact(ctx context.Context, llm *Gateway) error {
    // Phase 1: 丢弃所有截图类消息（只保留文本）
    m.discardScreenshots()

    // Phase 2: 用 LLM 总结旧消息
    if m.IsOverflow(m.messages) {
        summary, err := m.summarizeWithLLM(ctx, llm)
        if err != nil {
            return err
        }
        m.replaceOldMessages(summary)
    }
    return nil
}

// discardScreenshots 移除所有截图数据，仅保留文本描述
func (m *Memory) discardScreenshots()
```

---

### 1.8 事件系统

通过 Wails Events 实时推送到前端。

```go
// event.go

type EventBus struct {
    app *application.Application  // Wails app instance
}

type EventType string
const (
    EventStepStarted   EventType = "step:started"    // 新一轮开始
    EventTextDelta     EventType = "text:delta"      // LLM 文本流
    EventToolCalled    EventType = "tool:called"     // 工具调用开始
    EventToolResult    EventType = "tool:result"     // 工具执行结果
    EventStepProgress  EventType = "step:progress"   // 步骤执行进度
    EventScreenshot    EventType = "screenshot"      // 截图帧
    EventOCRResult     EventType = "ocr:result"      // OCR 识别结果
    EventCompleted     EventType = "completed"       // Agent 完成
    EventError         EventType = "error"           // 错误
    EventDoomLoop      EventType = "doom_loop"       // Doom Loop 检测
)

func (e *EventBus) Emit(event EventType, data any) {
    e.app.EmitEvent(string(event), data)
}
```

前端通过 `EventsOn` 订阅：

```typescript
// 前端订阅示例
EventsOn("step:progress", (data) => { ... })
EventsOn("screenshot", (base64) => { ... })
EventsOn("tool:called", (data) => { ... })
```

---

### 1.9 安全机制 (Safety)

```go
// safety.go

type Safety struct {
    doomThreshold int                    // 连续重复调用阈值 (默认 3)
    history       []ToolCallInfo         // 最近的工具调用记录
}

// IsDoomLoop 检测是否陷入死循环（同一工具+参数连续调用 N 次）
func (s *Safety) IsDoomLoop(calls []ToolCallInfo) bool {
    if len(calls) == 0 {
        return false
    }
    last := calls[len(calls)-1]
    count := 1
    for i := len(s.history) - 1; i >= 0; i-- {
        if s.history[i].Name == last.Name &&
           bytes.Equal(s.history[i].Args, last.Args) {
            count++
        } else {
            break
        }
    }
    return count >= s.doomThreshold
}

// WaitForConfirmation 通过 Wails Event 请求用户确认
func (s *Safety) WaitForConfirmation(ctx context.Context) error
```

---

### 1.10 重试策略 (Retry)

参考 OpenCode 的指数退避 + retry-after header。

```go
// retry.go

type RetryPolicy struct {
    initialDelay time.Duration   // 2s
    maxDelay     time.Duration   // 30s
    backoffFactor float64        // 2.0
    maxRetries   int             // 5
}

func (r *RetryPolicy) ShouldRetry(err error) bool {
    // 可重试: 5xx 服务端错误, 429 限流, 网络超时
    // 不可重试: 4xx 客户端错误, 认证失败, 上下文溢出
}

func (r *RetryPolicy) Wait(ctx context.Context, err error) error {
    delay := r.calculateDelay(err)
    select {
    case <-time.After(delay):
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}

func (r *RetryPolicy) calculateDelay(err error) time.Duration {
    // 检查 retry-after header
    // 否则使用指数退避: initialDelay * backoffFactor^attempt
}
```

---

### 1.11 System Prompt 构建

```go
// system_prompt.go

func (l *Loop) buildSystemPrompt(agent AgentConfig) string {
    var parts []string

    // 1. 基础指令（按 Agent 角色）
    parts = append(parts, agent.BasePrompt)

    // 2. 环境信息
    parts = append(parts, fmt.Sprintf(`
## 环境信息
- 操作系统: %s
- 屏幕分辨率: %s
- 当前日期: %s
- 绑定窗口: %s
`, runtime.GOOS, screenResolution(), time.Now().Format("2006-01-02"), l.session.WindowTitle))

    // 3. 可用工具说明（自动从注册表生成）
    parts = append(parts, l.generateToolDescriptions(agent.Name))

    // 4. 输出格式要求
    parts = append(parts, `
## 输出要求
- 每次操作前先截图确认当前屏幕状态
- 使用工具执行操作，不要凭空猜测坐标
- 每步操作后验证结果
- 如果连续失败 2 次，停下来分析原因
`)

    return strings.Join(parts, "\n\n")
}
```

---

### 1.12 内置 Agent 角色

| Agent | 用途 | 可用工具 | 模型建议 |
|-------|------|---------|---------|
| `planner` | 自然语言 → 步骤计划 + 模拟演示 | capture, ocr, locate, window, simulate | GPT-4o / Claude (带视觉) |
| `executor` | 按脚本步骤在目标窗口执行 | 全部执行工具 + capture + ocr + verify | GPT-4o / Claude |
| `reviewer` | 执行结果校验 + 异常分析 | capture, ocr, verify | 轻量模型即可 |

```go
// agents/planner.go
var PlannerAgent = AgentConfig{
    Name:        "planner",
    Description: "规划 Agent：将自然语言描述转化为可执行的自动化步骤",
    Model:       "gpt-4o",
    MaxSteps:    20,
    BasePrompt: `你是一个 GUI 自动化规划专家。用户会描述一个要自动化的操作任务。
你的职责是：
1. 截图观察当前屏幕状态
2. 用 OCR 识别界面元素
3. 规划操作步骤序列
4. 逐步模拟演示（只标注目标，不真正操作）
5. 等待用户确认后保存为脚本`,
}

// agents/executor.go
var ExecutorAgent = AgentConfig{
    Name:        "executor",
    Description: "执行 Agent：按脚本步骤在目标窗口上执行操作",
    Model:       "gpt-4o",
    MaxSteps:    100,
    BasePrompt: `你是一个 GUI 自动化执行专家。按照给定的步骤脚本，在指定窗口上逐步执行操作。
每步执行流程：
1. 截图确认当前屏幕状态
2. 用 OCR/locate 定位目标元素
3. 执行操作（click/type/hotkey 等）
4. 截图验证操作结果
5. 如果失败，分析原因并重试`,
}
```

---

### 1.13 脚本创建流程（使用 Agent Loop）

```
前端: 用户输入 "把 Excel 数据录入到 CRM"
    │
    ▼ Wails IPC
Session.Create(ScriptCreation)
    │
    ▼ 创建 planner session
Loop.Run(ctx)
    │
    ├── Step 1: LLM 返回 tool_call: capture()
    │           → 截图 → 返回 base64 给 LLM
    │
    ├── Step 2: LLM 返回 tool_call: ocr(screenshot)
    │           → OCR 识别 → 返回文字列表+坐标
    │
    ├── Step 3: LLM 返回 tool_call: simulate(step_plan)
    │           → 截图标注目标 → 推送给前端展示
    │
    ├── Step 4-8: 更多 simulate 调用...
    │
    └── Step N: LLM 返回最终步骤计划（无工具调用）
                → OutcomeStop → 返回 ScriptPlan 给前端
```

### 1.14 任务执行流程（使用 Agent Loop）

```
前端: 用户点击 "启动任务"
    │
    ▼ Wails IPC
Concurrency.Acquire() → 检查槽位
    │
    ▼ 创建 executor session
Session.Create(TaskExecution, task)
    │
    ▼ 注入脚本步骤到 system prompt
Loop.Run(ctx)
    │
    ├── Step 1: LLM 读取步骤 1 "click 新建"
    │           → tool_call: locate("新建") → 返回坐标
    │           → tool_call: click(x, y) → 执行点击
    │           → tool_call: verify() → 验证结果
    │
    ├── Step 2: LLM 读取步骤 2 "input 姓名"
    │           → tool_call: locate("姓名") → 定位输入框
    │           → tool_call: type("张三") → 输入文字
    │           → tool_call: verify() → 验证
    │
    ├── ... 逐步执行 ...
    │
    └── 全部步骤完成 → OutcomeStop
        → Concurrency.Release()
        → 保存 Execution 记录
```

---

### 1.15 前端订阅的事件清单

| 事件 | 数据 | 前端处理 |
|------|------|---------|
| `step:started` | `{step: number}` | 更新步骤计数 |
| `text:delta` | `string` | 追加到 AI 对话区 |
| `tool:called` | `{name, args}` | 显示"正在执行: click..." |
| `tool:result` | `{name, success, screenshot}` | 更新步骤状态 |
| `step:progress` | `{index, total, action, target}` | 更新进度条 |
| `screenshot` | `base64 string` | 更新实时画面 |
| `ocr:result` | `{text, rect, confidence}[]` | 更新 OCR 结果列表 |
| `completed` | `{success, duration}` | 显示完成提示 |
| `error` | `{message, step}` | Toast 错误通知 |
| `doom_loop` | `{tool, count}` | 弹窗询问用户 |

---

### 脚本创建流程（简化版）

```
用户描述 → Parser.ParseIntent()
              │
              ▼
         Planner.GeneratePlan()  ← 截图 + OCR 结果
              │
              ▼
      Simulator.SimulateStep()   ← 逐步：截图 + 标注目标
              │
              ▼ (用户确认/修改)
         Store.SaveScript()      ← 保存为脚本模板
              │
              ▼ (用户可从此脚本创建多个任务)
```

### 任务执行流程（简化版）

```
用户创建任务 → 选择脚本 + 绑定窗口 + 配置参数
              │
              ▼ (用户启动任务)
      Concurrency.Acquire()      ← 检查并发槽位
              │
              ├── 槽位不足 → 返回错误，提示用户等待或升级
              │
              ▼
       Executor.ExecuteTask()    ← 在绑定窗口上逐步执行
              │                    每步: 定位窗口 → OCR → 操作 → 验证
              ├── 成功 → 记录日志 → Concurrency.Release()
              └── 失败 → Recovery 介入 → Concurrency.Release()
```

### 数据模型

```go
// Intent 解析后的用户意图
type Intent struct {
    Action      string            // 高层动作：录入、安装、测试、查询
    Target      string            // 目标应用/系统
    Parameters  map[string]string // 关键参数
    LoopCount   int               // 循环次数（0=单次）
}

// ScriptPlan Agent 生成的脚本计划（创建阶段）
type ScriptPlan struct {
    Steps             []*PlannedStep
    EstimatedDuration time.Duration
}

// PlannedStep 计划中的单步
type PlannedStep struct {
    Index       int
    Action      string  // click / input / verify / scroll / hotkey / wait
    Target      string  // 目标描述（自然语言）
    TargetOCR   string  // OCR 匹配文字
    Value       string  // 输入值（input 时）
    Confidence  float64 // Agent 对此步骤的置信度
}

// StepSimulation 模拟演示结果
type StepSimulation struct {
    StepIndex   int
    Action      string
    Target      string
    Screenshot  []byte  // 当前屏幕截图
    Annotation  []byte  // 标注覆盖图（目标区域高亮）
    Confidence  float64 // OCR 匹配置信度
    TargetRect  Rect    // 目标区域坐标
}

// ─── 脚本（可复用模板）───

// Script 保存的脚本模板
type Script struct {
    ID          string
    Name        string
    Description string
    Steps       []*Step
    Variables   map[string]string   // 可配置参数
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// Step 脚本中的已确认步骤
type Step struct {
    ID          int
    Action      string
    Target      string
    TargetOCR   string
    Value       string
    Timeout     int
    Confidence  float64
    Confirmed   bool
}

// ─── 任务（执行实例）───

// Task 一个执行任务 = 脚本 + 窗口 + 参数
type Task struct {
    ID          string
    Name        string
    ScriptID    string              // 关联的脚本
    ScriptName  string              // 冗余，列表显示用
    WindowHandle uintptr            // 绑定的窗口句柄 (HWND)
    WindowTitle string              // 窗口标题（显示用）
    Parameters  map[string]string   // 脚本参数的实际值
    Trigger     *Trigger            // 触发方式
    Status      TaskStatus          // 状态
    CreatedAt   time.Time
}

type TaskStatus string
const (
    TaskStatusIdle      TaskStatus = "idle"      // 已配置，未运行
    TaskStatusRunning   TaskStatus = "running"   // 运行中
    TaskStatusPaused    TaskStatus = "paused"    // 暂停
    TaskStatusCompleted TaskStatus = "completed" // 已完成
    TaskStatusFailed    TaskStatus = "failed"    // 失败
)

// Trigger 触发方式
type Trigger struct {
    Type    string  // manual / cron / hotkey
    Cron    string  // cron 表达式
    Hotkey  string  // 快捷键
}

// ─── 执行记录 ───

// Execution 一次任务执行的记录
type Execution struct {
    ID          string
    TaskID      string
    TaskName    string
    Status      string            // success / failed / cancelled
    StartedAt   time.Time
    FinishedAt  time.Time
    Duration    time.Duration
    StepResults []*StepResult
}

// StepResult 单步执行结果
type StepResult struct {
    StepID      int
    Action      string
    Target      string
    Success     bool
    Duration    time.Duration
    Screenshot  []byte
    Error       string
}

// ─── 用户 ───

type User struct {
    ID          string
    VIPLevel    int               // 0=普通用户, 1-5=VIP等级
    MaxConcurrent int             // 并发上限：普通=3, VIP1=5, ..., VIP5=10
}
```

---

## 2. 视觉感知层 (`internal/vision/`)

负责"看懂"屏幕——截图、OCR、控件检测。

| 文件 | 职责 | 关键接口 |
|------|------|----------|
| `capture.go` | 屏幕截图（全屏/区域） | `CaptureScreen() ([]byte, error)` / `CaptureRegion(rect Rect) ([]byte, error)` |
| `ocr.go` | 文字识别 + 坐标定位 | `Recognize(image []byte) ([]OCRResult, error)` |
| `detection.go` | UI 控件检测（按钮、输入框等） | `DetectControls(image []byte) ([]Control, error)` |
| `matching.go` | 图像模板匹配 | `MatchTemplate(screen, template []byte) (*MatchResult, error)` |
| `layout.go` | 布局结构分析 | `AnalyzeLayout(image []byte) (*LayoutTree, error)` |

### 数据模型

```go
type OCRResult struct {
    Text       string
    Rect       Rect     // 文字区域坐标
    Confidence float64
}

type Control struct {
    Type       string   // button / input / dropdown / checkbox / label
    Rect       Rect
    Text       string   // OCR 识别到的文字
    Confidence float64
}

type Rect struct {
    X, Y, Width, Height int
}
```

---

## 3. 动作执行层 (`internal/action/`)

负责"动手操作"——鼠标、键盘、窗口、剪贴板。

| 文件 | 职责 | 关键接口 |
|------|------|----------|
| `mouse.go` | 鼠标操作 | `Click(x, y int)` / `DoubleClick()` / `RightClick()` / `Drag(from, to Point)` / `Scroll(x, y, delta int)` |
| `keyboard.go` | 键盘操作 | `Type(text string)` / `HotKey(keys ...string)` / `KeyDown()` / `KeyUp()` |
| `window.go` | 窗口管理 | `FindWindow(title string) (*Window, error)` / `Focus()` / `Move()` / `Resize()` / `ListWindows()` |
| `clipboard.go` | 剪贴板 | `GetText() (string, error)` / `SetText(text string) error` |

---

## 4. 监控验证层 (`internal/monitor/`)

负责"检查结果"——每步操作后验证是否成功。

| 文件 | 职责 | 关键接口 |
|------|------|----------|
| `checker.go` | 步骤结果校验 | `Verify(step *Step) (*VerifyResult, error)` |
| `logger.go` | 执行日志 + 截图序列记录 | `LogStep(result *StepResult)` / `GetExecutionLog(taskID string)` |

### 验证策略

| 验证类型 | 方法 | 适用场景 |
|---------|------|---------|
| OCR 文字验证 | 截图 OCR 检查目标文字是否出现 | 提交成功提示、错误信息 |
| 图像对比 | 操作前后截图直方图对比 | 界面是否发生变化 |
| 控件状态 | 检测控件是否变为预期状态 | 按钮变灰、输入框有值 |
| 超时判定 | 等待指定时间内条件满足 | 页面加载、弹窗出现 |

---

## 5. LLM 网关 (`internal/llm/`)

以 **Anthropic Messages API** 为主协议，兼容 OpenAI 格式。内部统一使用 Anthropic 数据模型。

| 文件 | 职责 |
|------|------|
| `gateway.go` | 统一调用入口，路由 + 重试 + fallback |
| `anthropic.go` | **Anthropic Messages API**（主协议，支持 streaming / vision / thinking） |
| `openai.go` | OpenAI Chat Completions 适配（转换为内部 Anthropic 格式） |
| `ollama.go` | Ollama 本地模型（通过 OpenAI-compatible 接口） |

### 设计原则

- 内部数据模型统一使用 Anthropic 格式（`Message`, `Content`, `ToolUse`, `ToolResult`）
- OpenAI 适配器负责将 OpenAI 请求/响应转换为 Anthropic 格式
- Ollama 通过 OpenAI-compatible API 接入，复用 OpenAI 适配器

### 核心接口

```go
// Anthropic 风格的消息格式（内部统一模型）
type Message struct {
    Role    string      // "user" / "assistant"
    Content []ContentBlock
}

type ContentBlock interface {
    contentType() string
}

type TextBlock struct {
    Type string `json:"type"` // "text"
    Text string `json:"text"`
}

type ImageBlock struct {
    Type   string    `json:"type"` // "image"
    Source ImageSource `json:"source"`
}

type ToolUseBlock struct {
    Type  string          `json:"type"` // "tool_use"
    ID    string          `json:"id"`
    Name  string          `json:"name"`
    Input json.RawMessage `json:"input"`
}

type ToolResultBlock struct {
    Type      string `json:"type"` // "tool_result"
    ToolUseID string `json:"tool_use_id"`
    Content   string `json:"content"`
    IsError   bool   `json:"is_error"`
}

// Provider 统一接口
type Provider interface {
    Name() string
    Available() bool
    Stream(ctx context.Context, req StreamRequest) (<-chan StreamEvent, error)
}

// StreamRequest 统一使用 Anthropic Messages API 格式
type StreamRequest struct {
    Model     string
    System    string
    Messages  []Message
    Tools     []ToolSchema
    MaxTokens int
}

type Gateway struct {
    providers []Provider  // 按环境变量配置的优先级排序
    retry     *RetryPolicy
}

func (g *Gateway) Stream(ctx context.Context, req StreamRequest) (<-chan StreamEvent, error)
```

### OpenAI 适配层

```go
// openai.go — 将 Anthropic 内部格式 ↔ OpenAI 格式双向转换

type OpenAIAdapter struct {
    apiKey  string
    baseURL string
}

// Stream 接收 Anthropic 格式请求 → 转换为 OpenAI 格式 → 调用 → 转回 Anthropic 事件流
func (a *OpenAIAdapter) Stream(ctx context.Context, req StreamRequest) (<-chan StreamEvent, error) {
    openaiReq := toOpenAIRequest(req)     // Anthropic → OpenAI
    stream, err := a.callOpenAI(ctx, openaiReq)
    return fromOpenAIStream(stream), err  // OpenAI events → Anthropic events
}
```

---

## 6. 定时调度 (`internal/scheduler/`)

| 文件 | 职责 | 关键接口 |
|------|------|----------|
| `cron.go` | Cron 定时任务管理 | `AddJob(taskID, cronExpr string)` / `RemoveJob(taskID)` / `ListJobs()` |

---

## 7. 全局快捷键 (`internal/hotkey/`)

| 文件 | 职责 | 关键接口 |
|------|------|----------|
| `hook.go` | 系统级热键注册 | `Register(hotkey string, callback func())` / `Unregister(hotkey string)` |

---

## 8. 并发管理 (`internal/concurrency/`)

控制任务并发执行数量，管理执行槽位。

| 文件 | 职责 |
|------|------|
| `manager.go` | 槽位分配/释放 + 并发上限控制 |

```go
type Manager struct {
    userMax  int                // 用户等级对应的上限
    running  map[string]context.CancelFunc  // taskID → cancel
    windows  map[uintptr]string            // windowHandle → taskID（窗口占用表）
    mu       sync.Mutex
}

// Acquire 尝试获取执行槽位
// 检查: 1) 并发上限  2) 窗口冲突
func (m *Manager) Acquire(taskID string, windowHandle uintptr) (context.Context, error) {
    m.mu.Lock()
    defer m.mu.Unlock()

    // 检查并发上限
    if len(m.running) >= m.userMax {
        return nil, ErrConcurrencyLimit
    }

    // 检查窗口冲突
    if occupiedBy, ok := m.windows[windowHandle]; ok {
        return nil, fmt.Errorf("窗口已被任务 %s 占用", occupiedBy)
    }

    ctx, cancel := context.WithCancel(context.Background())
    m.running[taskID] = cancel
    m.windows[windowHandle] = taskID
    return ctx, nil
}

// Release 释放槽位 + 窗口占用
func (m *Manager) Release(taskID string, windowHandle uintptr) {
    m.mu.Lock()
    defer m.mu.Unlock()
    if cancel, ok := m.running[taskID]; ok {
        cancel()
        delete(m.running, taskID)
    }
    delete(m.windows, windowHandle)
}
```

---

## 9. 数据模型 (`internal/model/`)

| 文件 | 实体 |
|------|------|
| `script.go` | Script, Step, Variable |
| `task.go` | Task, Trigger, TaskStatus |
| `execution.go` | Execution, StepResult |
| `user.go` | User, VIPLevel |
| `config.go` | HotkeyConfig, ThemeConfig |

---

## 10. 持久化层 (`internal/store/`)

| 文件 | 职责 |
|------|------|
| `sqlite.go` | SQLite 初始化 + 自动迁移 |
| `repository.go` | 仓储接口（Script CRUD, Task CRUD, Execution CRUD, User） |

---

## Wails 绑定 API 汇总

### 脚本管理

| 方法 | 说明 |
|------|------|
| `ScriptService.CreatePlan(description string) (*ScriptPlan, error)` | 自然语言 → 生成脚本计划 |
| `ScriptService.SimulateStep(plan *ScriptPlan, index int) (*StepSimulation, error)` | 模拟单步（截图标注） |
| `ScriptService.SaveScript(plan *ScriptPlan, name, desc string) (*Script, error)` | 保存为脚本模板 |
| `ScriptService.ListScripts() ([]*Script, error)` | 脚本列表 |
| `ScriptService.GetScript(id string) (*Script, error)` | 脚本详情 |
| `ScriptService.UpdateScript(script *Script) error` | 更新脚本 |
| `ScriptService.DeleteScript(id string) error` | 删除脚本 |
| `ScriptService.CopyScript(id string) (*Script, error)` | 复制脚本 |

### 任务管理

| 方法 | 说明 |
|------|------|
| `TaskService.CreateTask(req *CreateTaskRequest) (*Task, error)` | 创建任务（绑定脚本+窗口+参数） |
| `TaskService.ListTasks() ([]*Task, error)` | 任务列表 |
| `TaskService.GetTask(id string) (*Task, error)` | 任务详情 |
| `TaskService.UpdateTask(task *Task) error` | 更新任务配置 |
| `TaskService.DeleteTask(id string) error` | 删除任务 |
| `TaskService.StartTask(id string) error` | 启动任务（获取并发槽位） |
| `TaskService.PauseTask(id string) error` | 暂停任务 |
| `TaskService.ResumeTask(id string) error` | 恢复任务 |
| `TaskService.StopTask(id string) error` | 停止任务（释放槽位） |
| `TaskService.ListWindows() ([]*WindowInfo, error)` | 列出系统窗口（供绑定选择） |

### 并发 & 监控

| 方法 | 说明 |
|------|------|
| `ConcurrencyService.GetStatus() (*ConcurrencyStatus, error)` | 当前并发状态（运行数/上限） |
| `ConcurrencyService.GetRunningTasks() ([]*RunningTaskInfo, error)` | 运行中任务列表+进度 |

### 其他

| 方法 | 说明 |
|------|------|
| `VisionService.CaptureScreen() (string, error)` | 截图（返回 base64） |
| `VisionService.StartOCRStream() error` | 开始 OCR 实时识别 |
| `VisionService.StopOCRStream()` | 停止 OCR |
| `LogService.ListExecutions(filter *Filter) ([]*Execution, error)` | 执行日志 |
| `LogService.GetExecutionDetail(id string) (*Execution, error)` | 日志详情 |
| `ConfigService.GetConfig() (*Config, error)` | 获取配置（快捷键、主题等） |
| `ConfigService.SaveConfig(config *Config) error` | 保存配置 |
| `ConfigService.CheckUpdate() (*UpdateInfo, error)` | 检查版本更新 |
| `UserService.GetUser() (*User, error)` | 获取当前用户信息 |
| `SchedulerService.AddSchedule(taskID, cronExpr string) error` | 添加定时任务 |
| `SchedulerService.RemoveSchedule(taskID string) error` | 移除定时任务 |

---

## 待讨论

- [x] ~~脚本删除时，已关联的任务如何处理？~~ → **禁止删除**，提示用户先删除关联任务
- [x] ~~任务绑定的窗口句柄在系统重启后失效，如何处理？~~ → **按窗口标题/进程名重新匹配**，找不到则终止任务并通知用户原因
- [x] ~~并发执行时，多个任务操作同一个窗口怎么办？~~ → **启动时检测窗口冲突**，若同一窗口已有任务运行，拒绝执行并提示用户
- [ ] **TODO（后期完成）**：登录、VIP 等级验证、用户信息体系
- [x] ~~任务暂停后恢复时，是从当前步骤继续还是从头开始？~~ → **从当前步骤继续**
- [x] ~~Agent Loop 的最大步数默认值？~~ → **可在设置中配置**，默认 planner: 20, executor: 100
- [x] ~~上下文压缩时，截图类消息如何处理？~~ → **完全丢弃截图**，只保留对话文本数据
- [x] ~~Tool 执行超时时间如何配置？~~ → **全局默认 30s**，可在设置中配置
- [x] ~~LLM Gateway 的 fallback 顺序？~~ → **当前通过环境变量配置**（构建时注入），后期改为远端获取 (TODO)
- [x] ~~是否需要支持 Agent 角色自定义？~~ → **不需要**，所有 Agent 角色均为内置，用户无需额外配置
