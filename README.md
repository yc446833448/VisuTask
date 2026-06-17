# VisuTask

> 基于 **OCR 视觉识别** + **Agent 任务分解执行** 的可视化界面流程自动化桌面应用

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)](https://go.dev/)
[![React](https://img.shields.io/badge/React-18-61DAFB?logo=react)](https://react.dev/)
[![Wails](https://img.shields.io/badge/Wails-v3-DF0000?logo=wails)](https://wails.io/)

---

## 📖 项目简介

**VisuTask** 是一款智能化的 GUI 流程自动化桌面应用，采用 **Wails v3（Go 后端 + React 前端）** 框架打包为单一原生可执行文件。它通过 **远端 OCR 服务** 实时感知屏幕内容，结合 **AI Agent 任务规划能力** 将复杂操作分解为可执行的原子步骤，最终实现对任意可视化界面的端到端自动化操作。

与传统 RPA 不同，VisuTask **不需要用户拖拽流程图或编写脚本**。只需用自然语言描述操作，Agent 逐步模拟演示每个步骤（截图 + 标注目标区域），用户确认后保存为**脚本模板**。脚本可绑定到不同窗口创建**执行任务**，多个任务**并发运行**（普通用户最多 3 个，VIP 最高 10 个）。整个过程像"教 AI 做一次操作"——演示一遍，它就能在多个窗口上反复自动执行。

### 核心理念

```
用户描述操作 → Agent 生成脚本 → 逐步模拟确认 → 保存脚本模板
                                                     │
                              创建任务（绑定脚本+窗口句柄+参数）
                                                     │
                              多任务并发执行（最多 3~10 个）
```

- **📝 说**：用自然语言描述要完成的操作，Agent 自动生成脚本
- **👀 看**：Agent 截图标注每一步的操作目标，用户逐条确认
- **🚀 跑**：脚本绑定到不同窗口，多个任务同时并发执行
- **🧠 智**：AI Agent 理解意图、规划步骤、验证结果、异常恢复

---

## 🏗️ 系统架构

VisuTask 基于 **Wails v3** 框架，Go 后端与 React 前端通过 Wails Runtime 进行 IPC 通信，编译后生成单一原生桌面应用。

```
┌──────────────────────────────────────────────────────────────┐
│                  🖥️  VisuTask 桌面应用 (Wails v3)              │
│                                                               │
│  ┌────────────────────────────────────────────────────────┐  │
│  │              🎨 React 前端 (原生 WebView)                │  │
│  │  shadcn/ui + Tailwind CSS + Lucide Icons               │  │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌───────────┐ │  │
│  │  │ 脚本创建  │ │ 脚本库   │ │ 任务管理  │ │ 脚本市场   │ │  │
│  │  │ScriptNew │ │ Scripts  │ │  Tasks   │ │  Market   │ │  │
│  │  └──────────┘ └──────────┘ └──────────┘ └───────────┘ │  │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌───────────┐ │  │
│  │  │ 用量统计  │ │ 新建任务  │ │ 任务监控  │ │   设置    │ │  │
│  │  │  Stats   │ │ TaskNew  │ │ Monitor  │ │ Settings  │ │  │
│  │  └──────────┘ └──────────┘ └──────────┘ └───────────┘ │  │
│  └───────────────────────┬────────────────────────────────┘  │
│                          │  Wails IPC (Bindings + Events)     │
│  ┌───────────────────────▼────────────────────────────────┐  │
│  │                    ⚙️ Go 后端                            │  │
│  │                                                         │  │
│  │  ┌───────────────────────────────────────────────────┐ │  │
│  │  │               🧠 Agent 核心 (Agent Loop)           │ │  │
│  │  │  Session → Loop → Processor → Tool → Memory       │ │  │
│  │  │  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌────────┐  │ │  │
│  │  │  │ Planner │ │Executor │ │Reviewer │ │ Safety │  │ │  │
│  │  │  │ Agent   │ │ Agent   │ │ Agent   │ │  Retry │  │ │  │
│  │  │  └─────────┘ └─────────┘ └─────────┘ └────────┘  │ │  │
│  │  └───────────────────────────────────────────────────┘ │  │
│  │                                                         │  │
│  │  ┌──────────────────┐ ┌───────────────────────────────┐│  │
│  │  │ 👁️ 视觉感知层     │ │ 🖐️ 动作执行层                   ││  │
│  │  │ kbinani/screenshot│ │ Windows user32.dll syscall    ││  │
│  │  │ 远端 OCR 服务     │ │ 鼠标/键盘/窗口管理             ││  │
│  │  └──────────────────┘ └───────────────────────────────┘│  │
│  │                                                         │  │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌───────────┐ │  │
│  │  │ LLM 网关  │ │ 并发管理  │ │ 定时调度  │ │ SQLite    │ │  │
│  │  │Anthropic │ │ Slot+Win │ │  Cron    │ │  GORM     │ │  │
│  │  │+OpenAI   │ │ Conflict │ │          │ │           │ │  │
│  │  └──────────┘ └──────────┘ └──────────┘ └───────────┘ │  │
│  └────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────┘
```

### 核心模块说明

| 模块 | 职责 | 技术实现 |
|------|------|----------|
| **Agent Loop** | while(true) 核心循环：加载上下文 → 调用 LLM → 处理工具 → 评估结果 | 参考 OpenCode 架构 |
| **LLM 网关** | Anthropic 协议为主，兼容 OpenAI/Ollama，自动 fallback | Go HTTP + SSE 流式 |
| **Tool 系统** | 11 个 GUI 专用工具：截图、OCR、点击、输入、验证等 | 声明式定义 + 注册表 |
| **OCR** | 远端 OCR 服务，返回文字 + 坐标 + 置信度 | HTTP API (可配置) |
| **屏幕截图** | 纯 Go 截图，支持全屏/区域/窗口 | kbinani/screenshot |
| **鼠标/键盘** | Windows user32.dll 系统调用，无 CGo 依赖 | syscall + SendInput |
| **窗口管理** | 枚举窗口、聚焦、按标题查找 | EnumWindows + SetForegroundWindow |
| **并发管理** | 任务槽位分配 + 窗口冲突检测 | Go mutex |
| **上下文压缩** | 两阶段：丢弃截图 → LLM 总结旧消息 | 自动触发 |
| **安全机制** | Doom Loop 检测（同一工具连续重复 3 次报警） | 滑动窗口 |
| **持久化** | 脚本、任务、执行记录、用户 | SQLite + GORM |

---

## 📁 项目结构

```
VisuTask/
├── main.go                         # 程序入口 + Wails 应用组装
├── go.mod / go.sum                 # Go 模块
├── config.example.yaml             # 配置模板
│
├── frontend/                       # React 前端
│   ├── src/
│   │   ├── App.tsx                 # 路由配置 (9 页面)
│   │   ├── main.tsx                # 入口
│   │   ├── pages/                  # 页面组件
│   │   │   ├── Home.tsx            # 首页 (5 卡片入口)
│   │   │   ├── ScriptNew.tsx       # 创建脚本 (AI对话+MD预览)
│   │   │   ├── Scripts.tsx         # 脚本库 (Tab+列表)
│   │   │   ├── Market.tsx          # 脚本市场 (Tab+表格)
│   │   │   ├── Stats.tsx           # 用量统计 (图表)
│   │   │   ├── Tasks.tsx           # 任务管理 (可展开表格)
│   │   │   ├── TaskNew.tsx         # 新建任务 (表单)
│   │   │   ├── TaskMonitor.tsx     # 任务监控 (实时画面+日志)
│   │   │   └── Settings.tsx        # 设置 (5 Tab)
│   │   ├── components/
│   │   │   ├── layout/             # Header, StatusBar, Breadcrumb
│   │   │   └── ui/                 # shadcn/ui 组件 (19个)
│   │   ├── mock/                   # Mock 数据
│   │   ├── stores/                 # Zustand 状态
│   │   └── types/                  # TypeScript 类型
│   └── dist/                       # 构建产物 (嵌入到 Go)
│
├── internal/                       # Go 后端
│   ├── agent/                      # Agent 核心
│   │   ├── loop.go                 # Agent Loop 核心循环
│   │   ├── processor.go            # Anthropic SSE 流处理器
│   │   ├── tool.go                 # Tool 定义 + 注册表
│   │   ├── session.go              # Session 管理
│   │   ├── memory.go               # 上下文管理 + 压缩
│   │   ├── safety.go               # Doom Loop 检测
│   │   ├── retry.go                # 指数退避重试
│   │   ├── event.go                # 事件总线 (10 种事件)
│   │   ├── service.go              # 服务入口 (组装子系统)
│   │   ├── builtin.go              # 3 个内置 Agent
│   │   └── tools/tools.go          # 11 个 GUI 工具
│   ├── llm/                        # LLM 网关
│   │   ├── gateway.go              # 统一入口 + Anthropic 消息模型
│   │   ├── anthropic.go            # Anthropic Messages API (主协议)
│   │   └── openai.go               # OpenAI/Ollama 适配器
│   ├── vision/                     # 视觉感知
│   │   ├── engine.go               # 引擎 + 接口定义
│   │   ├── capture.go              # kbinani/screenshot 截图
│   │   ├── ocr.go                  # 远端 OCR HTTP 客户端
│   │   └── stub.go                 # Stub 实现 (降级)
│   ├── action/                     # 动作执行 (Windows syscall)
│   │   ├── engine.go               # 引擎 + 接口定义
│   │   ├── mouse.go                # user32.dll 鼠标
│   │   ├── keyboard.go             # SendInput 键盘 (Unicode)
│   │   ├── window.go               # EnumWindows 窗口
│   │   └── stub.go                 # Stub 实现 (降级)
│   ├── config/                     # 配置加载
│   │   └── config.go               # YAML + 环境变量覆盖
│   ├── concurrency/                # 并发管理
│   │   └── manager.go              # 槽位 + 窗口冲突
│   ├── monitor/                    # 步骤验证
│   │   └── checker.go              # OCR 验证
│   ├── scheduler/                  # 定时调度
│   │   └── cron.go                 # robfig/cron
│   ├── hotkey/                     # 全局快捷键
│   │   └── hook.go                 # 热键注册
│   ├── model/                      # 数据模型 (GORM)
│   │   ├── script.go, task.go, execution.go, user.go, window.go
│   └── store/                      # 持久化
│       ├── sqlite.go               # SQLite 初始化
│       └── repository.go           # CRUD 仓储
│
└── docs/                           # 设计文档
    ├── ui-design.md                # 前端 UI 设计 (9 页面线框图)
    ├── business-modules.md         # 业务模块设计 (Agent 架构)
    ├── design-system.md            # 前端设计规范 (色彩/字号/间距)
    └── components.md               # 公共组件规划
```

---

## 🔧 技术选型

### 前端

| 类别 | 技术 | 说明 |
|------|------|------|
| **框架** | React 18 + TypeScript | 组件化 UI |
| **构建** | Vite | HMR 热更新 |
| **UI 组件** | shadcn/ui + Tailwind CSS | 按需引入，CSS 变量主题 |
| **图标** | Lucide (`lucide-react`) | shadcn/ui 官方配套 |
| **图表** | Recharts | 用量统计折线图/条形图 |
| **Markdown** | react-markdown + remark-gfm | 脚本预览渲染 |
| **状态管理** | Zustand | 轻量响应式状态 |
| **通知** | Sonner | Toast 通知 |

### 后端

| 类别 | 技术 | 说明 |
|------|------|------|
| **语言** | Go 1.26+ | 纯 Go，无 CGo 依赖 |
| **桌面框架** | Wails v3 | Go ↔ React IPC + 嵌入式前端 |
| **LLM** | Anthropic Messages API | 主协议，SSE 流式，兼容 OpenAI |
| **OCR** | 远端 HTTP API | 可配置 endpoint，返回文字+坐标 |
| **截图** | kbinani/screenshot | 纯 Go 跨平台截图 |
| **键鼠** | Windows syscall | user32.dll + SendInput (Unicode) |
| **数据库** | SQLite + GORM | 嵌入式，自动迁移 |
| **定时任务** | robfig/cron | Cron 表达式调度 |
| **重试** | cenkalti/backoff | 指数退避 + retry-after |
| **配置** | YAML + 环境变量 | gopkg.in/yaml.v3 |

---

## ⚙️ 配置

复制 `config.example.yaml` 到 `~/.visutask/config.yaml`：

```yaml
llm:
  primary: anthropic
  anthropic:
    api_key: "sk-ant-xxx"
    model: "claude-sonnet-4-20250514"
  openai:
    api_key: "sk-xxx"
    model: "gpt-4o"
  ollama:
    base_url: "http://localhost:11434"
    model: "llama3"
  fallback_order: [anthropic, openai, ollama]

ocr:
  endpoint: "https://api.deepseek.com/v1/ocr"
  api_key: "sk-xxx"

agent:
  planner_max_steps: 20
  executor_max_steps: 100
  tool_timeout: 30          # 秒
  doom_loop_threshold: 3
  context_window: 128000    # tokens

concurrency:
  default_max: 3
```

**配置优先级：** 默认值 → YAML 文件 → 环境变量（最高）

---

## 🚀 快速开始

### 环境要求

| 工具 | 版本 | 用途 |
|------|------|------|
| **Go** | >= 1.22 | 后端编译 |
| **Node.js** | >= 20 | 前端构建 |
| **pnpm** | >= 9 | 包管理 |
| **Wails CLI** | v3 | `go install github.com/wailsapp/wails/v3/cmd/wails3@latest` |
| **WebView2** | — | Windows 10+ 预装 |

### 安装与运行

```bash
# 克隆仓库
git clone https://github.com/yc446833448/VisuTask.git
cd VisuTask

# 复制配置文件并填入 API Key
cp config.example.yaml ~/.visutask/config.yaml

# 安装前端依赖并构建
cd frontend && pnpm install && pnpm build && cd ..

# 开发模式 (前端 HMR + Go 热重载)
wails3 dev

# 生产构建
wails3 build
# 产物位于 build/bin/
```

### 前端独立开发

```bash
cd frontend
pnpm dev
# 访问 http://localhost:5173 (使用 mock 数据)
```

---

## 🎯 项目路线图

### Phase 1 — 基础能力 ✅
- [x] Wails v3 项目初始化，Go ↔ React 通信
- [x] 屏幕截图 (kbinani/screenshot)
- [x] 远端 OCR 服务集成
- [x] 鼠标/键盘/窗口操作 (Windows syscall)
- [x] SQLite 数据持久化 (GORM)
- [x] 前端 9 页面 + mock 数据

### Phase 2 — Agent 智能化 ✅
- [x] LLM 网关：Anthropic (主) + OpenAI + Ollama
- [x] Agent Loop 核心循环 (参考 OpenCode)
- [x] 11 个 GUI 专用工具
- [x] 上下文管理 + 两阶段压缩
- [x] Doom Loop 检测 + 指数退避重试
- [x] YAML 配置系统

### Phase 3 — 脚本与任务系统
- [ ] 脚本创建向导：Agent 对话 → 生成脚本
- [ ] 截图标注高亮：OCR 结果可视化标注
- [ ] 脚本库管理：CRUD + 复制
- [ ] 任务配置：绑定脚本 + 窗口 + 参数
- [ ] 并发执行：多任务同时运行 + 窗口冲突检测
- [ ] 用户等级：普通 3 / VIP 1-5 (5-10)
- [ ] 执行监控面板：实时画面 + Agent 日志流

### Phase 4 — 生产增强
- [ ] 定时任务：Cron 后台执行 + 系统托盘
- [ ] 全局快捷键：一键触发任务
- [ ] 执行仪表盘：成功率、耗时统计
- [ ] 脚本模板导入/导出
- [ ] 控件检测 (OpenCV / 边缘检测)

### Phase 5 — 生态与分发
- [ ] 脚本市场：社区共享 + 一键导入
- [ ] 交叉编译 + CI/CD
- [ ] Windows / macOS / Linux 三平台发布
- [ ] 自动更新机制

---

## 🤝 贡献指南

1. **Fork** 本仓库
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

### 开发规范

- **前端**：遵循 `docs/design-system.md` 设计规范
- **后端**：`gofmt`，遵循 Effective Go
- **Commit**：遵循 Conventional Commits
- **Agent 工具**：声明式定义，注册到 ToolRegistry

---

## 📄 许可证

MIT License — 详见 [LICENSE](LICENSE)

---

## 🙏 致谢

- [Wails](https://wails.io/) — Go + Web 桌面应用框架
- [OpenCode](https://github.com/anomalyco/opencode) — Agent Loop 架构参考
- [Anthropic](https://www.anthropic.com/) — Claude API
- [shadcn/ui](https://ui.shadcn.com/) — React 组件库
- [Lucide](https://lucide.dev/) — SVG 图标库
- [kbinani/screenshot](https://github.com/kbinani/screenshot) — 纯 Go 截图

---

<p align="center">
  <b>VisuTask</b> — 让 AI 看懂屏幕，替你动手<br>
  单一原生应用 · 双击即用<br>
  Made with ❤️ by <a href="https://github.com/yc446833448">yc446833448</a>
</p>
