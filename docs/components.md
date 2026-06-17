# VisuTask 前端公共组件规划

> 基于 shadcn/ui + Tailwind CSS + Lucide Icons
> 路径：`frontend/src/components/`

---

## 组件目录结构

```
components/
├── layout/                 # 布局组件
│   ├── AppLayout.tsx       # 全局布局 (Header + Content + StatusBar)
│   ├── Header.tsx          # 顶部导航栏
│   ├── StatusBar.tsx       # 底部状态栏
│   └── Breadcrumb.tsx      # 面包屑导航
│
├── screenshot/             # 截图与标注
│   ├── ScreenshotViewer.tsx    # 截图查看器
│   ├── OcrOverlay.tsx          # OCR 标注覆盖层
│   └── OcrAnnotation.tsx      # 单个标注框
│
├── task/                   # 任务相关
│   ├── ConcurrencyBar.tsx      # 并发状态进度条
│   ├── TaskStatusBadge.tsx     # 任务状态徽章
│   ├── StepLog.tsx             # 步骤日志时间线
│   └── TaskDetailExpand.tsx    # 任务表格展开详情
│
├── script/                 # 脚本相关
│   ├── ScriptCard.tsx          # 脚本列表项
│   └── MarkdownPreview.tsx     # Markdown 渲染预览
│
├── form/                   # 表单增强
│   ├── WindowSelector.tsx      # 窗口选择器
│   ├── ModelSelector.tsx       # LLM 模型选择器
│   ├── HotkeyCapture.tsx       # 快捷键录入
│   └── ChatInput.tsx           # AI 对话输入框
│
├── chart/                  # 图表
│   ├── LineChart.tsx           # 折线图 (消费统计)
│   └── BarChart.tsx            # 横向条形图 (用量统计)
│
├── feedback/               # 反馈 (基于 shadcn 封装)
│   ├── ConfirmDialog.tsx       # 确认弹窗
│   └── EmptyState.tsx          # 空状态占位
│
└── ui/                     # shadcn/ui 自动生成 (不手动编辑)
    ├── button.tsx
    ├── card.tsx
    ├── select.tsx
    ├── tabs.tsx
    ├── table.tsx
    ├── progress.tsx
    ├── badge.tsx
    ├── dialog.tsx
    ├── input.tsx
    ├── textarea.tsx
    ├── switch.tsx
    ├── radio-group.tsx
    ├── scroll-area.tsx
    ├── skeleton.tsx
    ├── toast.tsx
    └── alert-dialog.tsx
```

---

## 布局组件 `layout/`

### AppLayout

全局布局容器，包裹所有页面。

```
┌──────────────────────────────────┐
│  Header                          │
├──────────────────────────────────┤
│                                  │
│  <Outlet />  (路由内容)           │
│                                  │
├──────────────────────────────────┤
│  StatusBar                       │
└──────────────────────────────────┘
```

| 属性 | 说明 |
|------|------|
| children | 路由页面内容 |

职责：
- 渲染 Header + StatusBar
- 管理全局状态 (并发数、Agent 状态、OCR 状态)
- 订阅 Wails Events (step:progress, step:error)

---

### Header

```
┌──────────────────────────────────┐
│  [Logo] VisuTask    [⚙] [👤] VIP│
└──────────────────────────────────┘
```

| Props | 类型 | 说明 |
|-------|------|------|
| user | `User \| null` | 用户信息，null 时显示"登录" |
| vipLevel | `number` | VIP 等级，0=普通用户 |

内部状态：无，纯展示 + 路由跳转。

---

### StatusBar

```
┌──────────────────────────────────┐
│ ● Agent 就绪│● OCR│● LLM│2/3    │
└──────────────────────────────────┘
```

| Props | 类型 | 说明 |
|-------|------|------|
| agentStatus | `'idle' \| 'running' \| 'error'` | Agent 状态 |
| ocrConnected | `boolean` | OCR 连接状态 |
| llmModel | `string` | 当前 LLM 模型名 |
| llmAvailable | `boolean` | LLM 是否可用 |
| runningCount | `number` | 运行中任务数 |
| maxConcurrent | `number` | 并发上限 |

---

### Breadcrumb

```
← 首页 / 任务管理 / task_a3f8c2 / 监控
```

| Props | 类型 | 说明 |
|-------|------|------|
| items | `{ label: string, path?: string }[]` | 面包屑层级 |

- 自动根据当前路由生成
- 点击 `← 首页` 返回首页
- 中间层级可点击跳转，最后一级为纯文字

---

## 截图与标注 `screenshot/`

### ScreenshotViewer

截图查看器，支持叠加 OCR 标注覆盖层。用于任务监控页面。

```
┌────────────────────────────────┐
│                                │
│   img (base64)                 │
│   + OcrOverlay (absolute)      │
│                                │
│           h = 320px            │
└────────────────────────────────┘
```

| Props | 类型 | 说明 |
|-------|------|------|
| src | `string` | base64 编码的截图 |
| annotations | `OcrAnnotation[]` | OCR 标注数据 |
| showOverlay | `boolean` | 是否显示标注覆盖层 |
| height | `number` | 容器高度，默认 320 |

---

### OcrOverlay

OCR 标注覆盖层，绝对定位在 ScreenshotViewer 上方。

| Props | 类型 | 说明 |
|-------|------|------|
| annotations | `OcrAnnotation[]` | 标注列表 |
| visible | `boolean` | 是否可见 |

```ts
interface OcrAnnotation {
  text: string
  rect: { x: number; y: number; width: number; height: number }
  confidence: number
  status: 'target' | 'confirmed' | 'failed'  // 决定颜色
}
```

---

### OcrAnnotation

单个标注框。

| Props | 类型 | 说明 |
|-------|------|------|
| text | `string` | OCR 识别文字 |
| rect | `Rect` | 位置和尺寸 |
| confidence | `number` | 置信度 0-1 |
| status | `'target' \| 'confirmed' \| 'failed'` | 状态 |

颜色映射：
- `target` → 蓝 `#3b82f6`
- `confirmed` → 绿 `#22c55e`
- `failed` → 红 `#ef4444`

---

## 任务相关 `task/`

### ConcurrencyBar

并发状态进度条，显示在任务管理页底部。

```
运行中 2/3  ████████████░░░░░░
```

| Props | 类型 | 说明 |
|-------|------|------|
| running | `number` | 当前运行数 |
| max | `number` | 并发上限 |

内部使用 shadcn `Progress` 组件。

---

### TaskStatusBadge

任务状态徽章。

| Props | 类型 | 说明 |
|-------|------|------|
| status | `TaskStatus` | 任务状态 |
| progress | `number` | 运行进度 0-100 (可选) |

```ts
type TaskStatus = 'idle' | 'running' | 'paused' | 'completed' | 'failed'
```

样式映射：
| 状态 | 图标 | 颜色 | 文字 |
|------|------|------|------|
| idle | `Circle` | 灰 | 空闲 |
| running | `Play` | 蓝(脉冲) | 运行中 |
| paused | `Pause` | 黄 | 暂停 |
| completed | `CheckCircle` | 绿 | 成功 |
| failed | `XCircle` | 红 | 失败 |

---

### StepLog

步骤日志时间线，用于任务展开详情和监控页面。

```
✓ 1. click "新建"     0.5s
✓ 2. input "姓名"     1.2s
● 3. click "保存"     执行中...
○ 4. verify "成功"    等待
```

| Props | 类型 | 说明 |
|-------|------|------|
| steps | `StepLogItem[]` | 步骤列表 |
| autoScroll | `boolean` | 是否自动滚动到底部 |
| maxHeight | `number` | 最大高度 |

```ts
interface StepLogItem {
  index: number
  action: string
  target: string
  status: 'done' | 'running' | 'waiting' | 'failed'
  duration?: number  // 秒
  error?: string
}
```

---

### TaskDetailExpand

任务表格展开行详情，封装展开区域的所有内容。

| Props | 类型 | 说明 |
|-------|------|------|
| task | `Task` | 任务数据 |
| execution | `Execution \| null` | 当前执行记录 |
| onMonitor | `() => void` | 点击"实时监控"回调 |

内部组合：任务信息 + Progress + StepLog + 最近执行记录 + 监控按钮。

---

## 脚本相关 `script/`

### ScriptCard

脚本列表项，用于脚本库页面。

```
📦 数据录入 CRM
   8 步  │  创建: 6/17  │  已创建 3 个任务
   [▶ 创建任务] [✏ 编辑] [📋 复制] [🗑 删除]
```

| Props | 类型 | 说明 |
|-------|------|------|
| script | `Script` | 脚本数据 |
| onRun | `() => void` | 创建任务 |
| onEdit | `() => void` | 编辑 |
| onCopy | `() => void` | 复制 |
| onDelete | `() => void` | 删除 |

---

### MarkdownPreview

Markdown 渲染预览，用于创建脚本页面右侧。

| Props | 类型 | 说明 |
|-------|------|------|
| content | `string` | Markdown 文本 |
| className | `string` | 自定义样式 |

依赖：`react-markdown` + `remark-gfm`

---

## 表单增强 `form/`

### WindowSelector

窗口选择器，列出系统当前可见窗口供绑定。

```
┌────────────────────────────┐
│ 选择目标窗口...         [▾] │
│                            │
│ · Excel.exe - 客户数据.xlsx │
│ · Chrome - CRM 系统        │
│ · 企业微信                  │
└────────────────────────────┘
[🔄 刷新窗口列表]
```

| Props | 类型 | 说明 |
|-------|------|------|
| value | `string \| null` | 选中的窗口句柄 |
| onChange | `(handle: string, title: string) => void` | 选择回调 |

内部：
- 调用 `TaskService.ListWindows()` 获取窗口列表
- 刷新按钮重新获取列表
- Select 下拉显示 `进程名 - 窗口标题`

---

### ModelSelector

LLM 模型选择器，用于创建脚本页面。

```
🤖 模型: [GPT-4o ▾]
```

| Props | 类型 | 说明 |
|-------|------|------|
| value | `string` | 当前选中模型 |
| onChange | `(model: string) => void` | 切换回调 |

模型列表从环境变量 / 构建配置读取。

---

### HotkeyCapture

快捷键录入组件，用于设置页。

```
启动任务：  [Ctrl+Shift+E]  [修改]
```

| Props | 类型 | 说明 |
|-------|------|------|
| value | `string` | 当前快捷键 |
| onChange | `(hotkey: string) => void` | 修改回调 |

交互：
- 点击"修改"后进入录入模式
- 监听键盘事件，捕获组合键
- Enter 确认，Escape 取消

---

### ChatInput

AI 对话输入框，用于创建脚本页面。

```
┌─────────────────────────────────┐
│ 输入消息...                      │
└─────────────────────────────────┘
```

| Props | 类型 | 说明 |
|-------|------|------|
| onSend | `(message: string) => void` | 发送回调 |
| disabled | `boolean` | 是否禁用 |
| placeholder | `string` | 占位文字 |

- Enter 发送，Shift+Enter 换行
- 发送后自动清空

---

## 图表 `chart/`

### LineChart

折线图，用于消费金额统计。

| Props | 类型 | 说明 |
|-------|------|------|
| data | `{ date: string; value: number }[]` | 数据点 |
| height | `number` | 图表高度 |
| timeRange | `'7d' \| '30d' \| '90d'` | 时间范围 |
| onRangeChange | `(range) => void` | 范围切换回调 |

依赖：`recharts`

---

### BarChart

横向条形图，用于 LLM / OCR 用量统计。

| Props | 类型 | 说明 |
|-------|------|------|
| data | `{ label: string; value: number }[]` | 数据项 |
| height | `number` | 图表高度 |

---

## 反馈 `feedback/`

### ConfirmDialog

确认弹窗，基于 shadcn AlertDialog 封装。

| Props | 类型 | 说明 |
|-------|------|------|
| open | `boolean` | 是否打开 |
| onOpenChange | `(open: boolean) => void` | 状态回调 |
| title | `string` | 标题 |
| description | `string` | 描述 |
| confirmText | `string` | 确认按钮文字 |
| cancelText | `string` | 取消按钮文字 |
| variant | `'default' \| 'destructive'` | 样式变体 |
| onConfirm | `() => void` | 确认回调 |

---

### EmptyState

空状态占位，用于列表无数据时。

```
        📭
    暂无数据
  去创建一个脚本吧
```

| Props | 类型 | 说明 |
|-------|------|------|
| icon | `LucideIcon` | 图标 |
| title | `string` | 标题 |
| description | `string` | 描述 |
| action | `{ label: string; onClick: () => void }` | 操作按钮 (可选) |

---

## shadcn/ui 组件清单

需要安装的 shadcn/ui 组件：

```bash
npx shadcn@latest add button card select tabs table progress
npx shadcn@latest add badge dialog input textarea switch
npx shadcn@latest add radio-group scroll-area skeleton toast
npx shadcn@latest add alert-dropdown separator tooltip
```

| 组件 | 用途 |
|------|------|
| `button` | 所有按钮 |
| `card` | 首页卡片、设置区块 |
| `select` | 脚本/窗口/模型/版本选择 |
| `tabs` | 脚本库/市场/设置 Tab 切换 |
| `table` | 任务管理/脚本市场表格 |
| `progress` | 并发进度条、执行进度 |
| `badge` | 任务状态、VIP 等级 |
| `dialog` | 新建任务 (已改为子页面，备用) |
| `input` | 表单输入 |
| `textarea` | AI 对话输入 |
| `switch` | OCR 标注开关、设置开关 |
| `radio-group` | 触发方式、脚本选择、主题 |
| `scroll-area` | 对话区域、步骤日志滚动 |
| `skeleton` | 加载占位 |
| `toast` (sonner) | 操作反馈通知 |
| `alert-dialog` | 删除/停止确认 |
| `dropdown-menu` | 更多操作菜单 |
| `separator` | 分隔线 |
| `tooltip` | 图标按钮提示 |

---

## 其他第三方依赖

| 包名 | 用途 | 使用组件 |
|------|------|---------|
| `react-markdown` | Markdown 渲染 | MarkdownPreview |
| `remark-gfm` | GFM 语法支持 | MarkdownPreview |
| `recharts` | 图表库 | LineChart, BarChart |
| `sonner` | Toast 通知 | 全局通知 |
| `lucide-react` | 图标库 | 全局 |

---

## Hook 规划

```
hooks/
├── useConcurrency.ts       # 并发状态 (运行数/上限)
├── useAgent.ts             # Agent 状态 + Wails 绑定调用
├── useWailsEvent.ts        # 订阅 Wails Events
├── useWindows.ts           # 获取系统窗口列表
└── useTheme.ts             # 主题切换
```

| Hook | 返回值 | 说明 |
|------|--------|------|
| `useConcurrency` | `{ running, max, status }` | 轮询或事件驱动的并发状态 |
| `useAgent` | `{ createPlan, simulate, save, run, stop }` | Agent 操作封装 |
| `useWailsEvent` | `void` | 订阅/取消订阅 Wails 事件 |
| `useWindows` | `{ windows, refresh }` | 系统窗口列表 + 刷新 |
| `useTheme` | `{ theme, setTheme }` | 主题状态管理 |
