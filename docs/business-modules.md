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

### 模块划分

| 文件 | 职责 | 对外接口 (Wails Binding) |
|------|------|--------------------------|
| `parser.go` | 自然语言 → 结构化意图 | `ParseIntent(description string) (*Intent, error)` |
| `planner.go` | 意图 + 屏幕状态 → 操作计划 | `GeneratePlan(intent *Intent) (*ScriptPlan, error)` |
| `simulator.go` | 单步模拟（截图标注，不执行） | `SimulateStep(plan *ScriptPlan, stepIndex int) (*StepSimulation, error)` |
| `executor.go` | 按脚本步骤在指定窗口执行 | `ExecuteTask(taskID string) (*ExecutionResult, error)` |
| `recovery.go` | 执行失败时的重试/回退/人工介入 | 内部使用 |
| `memory.go` | 短期记忆 + 长期记忆 | 内部使用 |
| `service.go` | Wails 绑定入口 | 前端调用的 API |

### 脚本创建流程

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

### 任务执行流程

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

统一管理多个 LLM 提供方，自动 fallback。

| 文件 | 职责 |
|------|------|
| `gateway.go` | 统一调用入口，路由 + 重试 + fallback |
| `openai.go` | OpenAI API 适配器 |
| `claude.go` | Claude API 适配器 |
| `ollama.go` | Ollama 本地模型适配器 |

### 核心接口

```go
type LLMProvider interface {
    Chat(ctx context.Context, messages []Message) (*Response, error)
    ChatWithVision(ctx context.Context, messages []Message, images []Image) (*Response, error)
    Name() string
    Available() bool
}

type Gateway struct {
    providers []LLMProvider  // 按优先级排序
}

// Chat 自动选择可用 provider，失败时 fallback 到下一个
func (g *Gateway) Chat(ctx context.Context, messages []Message) (*Response, error)
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
    mu       sync.Mutex
}

// Acquire 尝试获取执行槽位，返回是否成功
func (m *Manager) Acquire(taskID string) (context.Context, error)

// Release 释放槽位
func (m *Manager) Release(taskID string)

// RunningCount 当前运行中的任务数
func (m *Manager) RunningCount() int

// RunningTasks 返回运行中的任务ID列表
func (m *Manager) RunningTasks() []string

// MaxConcurrent 当前用户的并发上限
func (m *Manager) MaxConcurrent() int
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

- [ ] 脚本删除时，已关联的任务如何处理？（级联删除 vs 禁止删除 vs 解除关联）
- [ ] 任务绑定的窗口句柄在系统重启后失效，如何处理？（按窗口标题/进程名重新匹配）
- [ ] 并发执行时，多个任务操作同一个窗口怎么办？（互斥锁 vs 用户自行管理）
- [ ] VIP 等级的验证方式？（本地校验 vs 服务端验证）
- [ ] 任务暂停后恢复时，是从当前步骤继续还是从头开始？
- [ ] Parser 用 LLM 解析还是规则+LLM 混合？
- [ ] Simulator 的标注覆盖图用 Go 端生成还是前端 Canvas 绘制？
- [ ] Recovery 策略的阈值如何配置？
- [ ] LLM 网关是否需要支持流式输出（streaming）用于实时规划？
