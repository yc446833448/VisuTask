# VisuTask

> 基于 **OCR 视觉识别** + **Agent 任务分解执行** 的可视化界面流程自动化桌面应用

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)](https://go.dev/)
[![React](https://img.shields.io/badge/React-18-61DAFB?logo=react)](https://react.dev/)
[![Wails](https://img.shields.io/badge/Wails-v3-DF0000?logo=wails)](https://wails.io/)

---

## 📖 项目简介

**VisuTask** 是一款智能化的 GUI 流程自动化桌面应用，采用 **Wails（Go 后端 + React 前端）** 框架打包为单一原生可执行文件。它通过 **OCR（光学字符识别）** 实时感知屏幕内容，结合 **AI Agent 任务规划能力** 将复杂操作分解为可执行的原子步骤，最终实现对任意可视化界面的端到端自动化操作。

与传统 RPA 不同，VisuTask **不需要用户拖拽流程图或编写脚本**。只需用自然语言描述任务，Agent 就会逐步模拟演示每个操作（截图 + 标注目标区域），用户逐条确认后自动生成可复用的任务模板。整个过程像"教 AI 做一次操作"——演示一遍，它就能反复自动执行。

### 核心理念

```
用户描述任务 → Agent 规划步骤 → 逐步模拟确认 → 自动生成文档 → 正式执行验证
```

- **📝 说**：用自然语言描述要完成的任务，无需学习任何编排工具
- **👀 看**：Agent 截图标注每一步的操作目标，让用户直观确认
- **🧠 想**：AI Agent 理解意图，将高层任务分解为可执行的原子步骤
- **🖐️ 动**：按确认后的步骤模拟键鼠操作，实时验证每步结果

### 为什么选择桌面应用？

| 优势 | 说明 |
|------|------|
| **原生性能** | Go 直接调用系统 API，无浏览器沙箱限制，键鼠操作延迟更低 |
| **单文件分发** | Wails 打包为单一 exe/dmg/AppImage，用户无需安装任何运行时 |
| **系统级能力** | 直接访问屏幕截图、全局快捷键、窗口句柄、剪贴板等 OS 能力 |
| **离线可用** | 搭配本地 OCR 引擎和本地 LLM，无需网络即可完成自动化 |
| **轻量小巧** | 使用系统原生 WebView，不捆绑 Chromium，安装包 < 20MB |

---

## 🏗️ 系统架构

VisuTask 基于 **Wails** 框架，Go 后端与 React 前端通过 Wails Runtime 进行 IPC 通信，编译后生成单一原生桌面应用。

```
┌──────────────────────────────────────────────────────────────┐
│                  🖥️  VisuTask 桌面应用 (Wails)                 │
│                                                               │
│  ┌────────────────────────────────────────────────────────┐  │
│  │              🎨 React 前端 (原生 WebView)                │  │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌───────────┐ │  │
│  │  │ 任务创建  │ │ 实时预览  │ │ 执行日志  │ │ 任务列表   │ │  │
│  │  │Designer  │ │Live View │ │   Logs   │ │ TaskList  │ │  │
│  │  └──────────┘ └──────────┘ └──────────┘ └───────────┘ │  │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────────────────┐   │  │
│  │  │ 定时调度  │ │ 录制回放  │ │ 执行监控面板          │   │  │
│  │  │Scheduler │ │ Recorder │ │  实时进度+截图流       │   │  │
│  │  └──────────┘ └──────────┘ └──────────────────────┘   │  │
│  └───────────────────────┬────────────────────────────────┘  │
│                          │  Wails IPC (Bindings)              │
│  ┌───────────────────────▼────────────────────────────────┐  │
│  │                    ⚙️ Go 后端                            │  │
│  │                                                         │  │
│  │  ┌───────────────────────────────────────────────────┐ │  │
│  │  │               🧠 Agent 调度核心                      │ │  │
│  │  │  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌──────────┐│ │  │
│  │  │  │ 任务解析 │ │ 计划生成 │ │ 模拟演示 │ │ 步骤执行  ││ │  │
│  │  │  │ Parser  │ │ Planner │ │Simulate │ │ Executor ││ │  │
│  │  │  └─────────┘ └─────────┘ └─────────┘ └──────────┘│ │  │
│  │  │  ┌─────────┐                                      │ │  │
│  │  │  │ 异常恢复 │                                      │ │  │
│  │  │  │ Recovery │                                      │ │  │
│  │  │  └─────────┘                                      │ │  │
│  │  └───────────────────────────────────────────────────┘ │  │
│  │                                                         │  │
│  │  ┌───────────────────────────────────────────────────┐ │  │
│  │  │               👁️ 视觉感知层                          │ │  │
│  │  │  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌──────────┐│ │  │
│  │  │  │ 屏幕截图 │ │ OCR 引擎 │ │ 控件检测 │ │ 布局分析  ││ │  │
│  │  │  │Capture  │ │   OCR   │ │Detection│ │ Layout   ││ │  │
│  │  │  └─────────┘ └─────────┘ └─────────┘ └──────────┘│ │  │
│  │  └───────────────────────────────────────────────────┘ │  │
│  │                                                         │  │
│  │  ┌───────────────────────────────────────────────────┐ │  │
│  │  │               🖐️ 动作执行层                          │ │  │
│  │  │  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌──────────┐│ │  │
│  │  │  │ 鼠标操作 │ │ 键盘输入 │ │ 窗口管理 │ │ 剪贴板   ││ │  │
│  │  │  │  Mouse  │ │Keyboard │ │ Window  │ │  Clip    ││ │  │
│  │  │  └─────────┘ └─────────┘ └─────────┘ └──────────┘│ │  │
│  │  └───────────────────────────────────────────────────┘ │  │
│  │                                                         │  │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌───────────┐ │  │
│  │  │ 本地存储  │ │ 任务调度  │ │ LLM 网关  │ │ 全局快捷键 │ │  │
│  │  │ SQLite   │ │  Cron    │ │ LLM GW   │ │  Hotkey   │ │  │
│  │  └──────────┘ └──────────┘ └──────────┘ └───────────┘ │  │
│  └────────────────────────────────────────────────────────┘  │
│                                                               │
│  ┌────────────────────────────────────────────────────────┐  │
│  │                 🖥️ 操作系统层                             │  │
│  │   Windows GDI / macOS CGWindow / Linux X11              │  │
│  └────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────┘
```

### 核心模块说明

| 模块 | 职责 | 技术实现 |
|------|------|----------|
| **Wails Runtime** | Go ↔ React 双向绑定通信 | Wails v3 原生 IPC |
| **任务解析器 (Parser)** | 将自然语言描述的任务转为结构化指令 | Go + LLM API |
| **计划生成器 (Planner)** | 根据屏幕状态和任务目标规划操作步骤序列 | Agent 推理链 / 分层规划 |
| **模拟演示器 (Simulator)** | 逐步模拟操作（截图标注目标，不真正执行），供用户确认 | 截图 + OCR 标注 + 高亮覆盖 |
| **OCR 引擎** | 识别屏幕文字及其位置坐标（返回给前端标注） | Tesseract CGo / PaddleOCR gRPC |
| **控件检测** | 定位按钮、输入框、下拉菜单等可交互元素 | GoCV (OpenCV) / 边缘检测 / 轮廓分析 |
| **屏幕截图** | 高性能屏幕捕获，实时帧推送到前端预览 | Windows GDI / macOS CGWindow / X11 SHM |
| **动作执行器** | 鼠标点击、键盘输入、拖拽等底层操作 | robotgo / Win32 user32.dll / XTest / CGEvent |
| **执行监控 (Monitor)** | 每步操作后截图对比 + OCR 校验，判断成功/失败 | GoCV 直方图对比 + OCR diff |
| **LLM 网关** | 统一管理多模型（云端 + 本地），自动 fallback | Go HTTP Client + 速率限制 |
| **任务调度** | 定时触发任务、后台静默执行 | robfig/cron |
| **全局快捷键** | 用户可注册全局热键，一键触发任务 | robotgo / 系统级 Hook |
| **本地存储** | 任务模板、执行历史、配置的持久化 | SQLite (GORM) |

### 数据流

```
[用户] 在 VisuTask 桌面窗口中用自然语言描述任务
   │
   ▼  (Wails IPC: Go 函数调用)
[Parser] 接收自然语言 → 调用 LLM 解析意图
   │
   ▼
[Planner] 截取当前屏幕 → OCR 识别 → 结合意图生成操作计划
   │
   ▼  (Wails Events: Go → React 实时状态推送)
[React 前端] 展示执行计划的步骤列表，等待用户确认
   │
   ▼  (用户点击"模拟执行")
[Agent] 逐步模拟演示操作序列
   │
   ├── 每步：截图 → OCR 标注 → 高亮目标区域 → 暂停等待确认
   ├── 用户可确认、修改、跳过或插入新步骤
   ├── Wails Events 实时推送当前步骤截图和 OCR 结果到前端
   │
   ▼  (用户确认全部步骤)
[保存] 结构化任务模板存入 SQLite，可复用、可编辑
   │
   ▼  (用户点击"正式执行")
[Executor] 按已确认的步骤序列自动执行
   │
   ├── 每步执行前截图 → OCR 定位目标 → 执行键鼠操作 → 截图验证结果
   ├── Wails Events 实时推送执行进度、当前步骤截图到前端
   ├── 失败时自动重试 / 回退 / 提示用户介入
   │
   ▼
[完成] 执行记录存入 SQLite → 前端展示结果报告和截图回放
```

---

## 🚀 典型使用场景

### 场景一：跨系统数据录入

```
用户: "把 Excel 里的 100 条客户信息逐一录入到网页 CRM 系统中"

VisuTask:
  1. Agent 解析意图，规划操作步骤
  2. 模拟演示第一步：截图 Excel，OCR 标注第一行数据 → 用户确认
  3. 模拟演示第二步：截图 CRM，标注"新增"按钮和各输入框 → 用户确认
  4. 模拟演示第三步：标注"保存"按钮，验证成功提示 → 用户确认
  5. 用户确认所有步骤 → 保存为任务模板
  6. 正式执行：循环 100 次，每次录入一条数据
```

### 场景二：软件安装与配置

```
用户: "帮我安装 Python 3.12 并配置好环境变量"

VisuTask:
  1. Agent 规划安装步骤
  2. 逐步模拟：打开安装程序 → 勾选 Add to PATH → Install → Finish
  3. 用户确认每步操作目标和顺序
  4. 保存为可复用模板
  5. 正式执行 → 截图验证安装结果
```

### 场景三：定时办公自动化

```
用户: "每天 9 点打开企业微信，找到昨日未读消息并汇总"

VisuTask:
  1. Agent 规划：打开企业微信 → 定位未读消息 → 逐条读取 → 汇总
  2. 逐步模拟并标注每个操作目标 → 用户确认
  3. 保存模板 + 配置定时触发
  4. 后台静默执行 → 生成汇总报告 → 弹窗通知
```

### 场景四：桌面应用回归测试

```
用户: "回归测试这个桌面应用的所有表单提交功能"

VisuTask:
  1. Agent 截图识别界面中所有表单和按钮
  2. 自动生成测试用例（正常值、边界值、异常值）
  3. 逐步模拟每个表单的填写和提交 → 用户确认
  4. 正式执行 → OCR 校验提交结果 → 记录 Pass/Fail
```

---

## 🔧 技术选型

### 桌面壳框架

| 技术 | 说明 |
|------|------|
| **Wails v3** | Go + Web 前端 → 单一原生二进制，使用系统 WebView（WebView2/WKWebView/WebKitGTK），安装包 < 15MB |
| **WebView2** (Windows) | Windows 10+ 预装，无需额外依赖 |
| **WKWebView** (macOS) | macOS 系统内置 |
| **WebKitGTK** (Linux) | 多数 Linux 发行版已预装 |

### 前端 (React)

| 类别 | 技术 | 说明 |
|------|------|------|
| **框架** | React 18 + TypeScript | 组件化 UI |
| **构建** | Vite | Wails 原生集成，HMR 热更新 |
| **UI 组件** | shadcn/ui + Tailwind CSS | 现代化可定制组件，按需引入 |
| **图标** | Lucide (`lucide-react`) | shadcn/ui 官方配套，Tree-shakeable 按需加载 |
| **任务创建** | Agent 模拟 + 对话确认 | 逐步演示、确认、自动生成任务文档 |
| **状态管理** | Zustand | 轻量响应式状态 |
| **实时通信** | Wails Runtime (Events) | Go → React 事件推送，无网络开销 |
| **图表** | ECharts / AntV | 执行统计仪表盘 |

### 后端 (Go)

| 类别 | 技术 | 说明 |
|------|------|------|
| **语言** | Go 1.22+ | 高性能、交叉编译 |
| **IPC 框架** | Wails v3 Bindings | 暴露 Go 函数给前端直接调用 |
| **OCR** | Tesseract (CGo) / PaddleOCR gRPC | 离线文字识别 + 坐标 |
| **视觉处理** | GoCV (OpenCV 绑定) | 模板匹配、轮廓检测、图像对比 |
| **键鼠模拟** | robotgo + syscall | 跨平台鼠标/键盘/窗口操作 |
| **屏幕捕获** | screenshot + GDI/X11/CG | 截图帧获取，并行推前端 |
| **数据库** | SQLite (GORM) | 嵌入式，零配置，单文件存储 |
| **任务调度** | robfig/cron | 定时任务，毫秒级精度 |
| **LLM 集成** | 自建 LLM Gateway | 统一路由 OpenAI / Claude / Ollama 本地模型 |
| **配置** | Viper | YAML 配置文件 |
| **日志** | zap | 结构化高性能日志 |
| **构建** | Wails build + NSIS (Win) / DMG (Mac) / AppImage (Linux) | |

### 打包与分发

| 平台 | 格式 | 体积 |
|------|------|------|
| **Windows** | `.exe` (NSIS 安装包) | ~15MB |
| **macOS** | `.dmg` / `.app` | ~18MB |
| **Linux** | `.AppImage` / `.deb` | ~16MB |

---

## 📁 项目结构

```
VisuTask/
├── app.go                          # Wails 应用入口，注册 Go ↔ 前端绑定
├── main.go                         # 程序主入口
├── wails.json                      # Wails 项目配置
│
├── frontend/                       # React 前端 (Wails 自动识别)
│   └── src/
│       ├── App.tsx                 # 根组件
│       ├── main.tsx                # 前端入口
│       ├── pages/
│       │   ├── Dashboard/          # 仪表盘首页（任务概览、统计）
│       │   ├── Designer/           # 任务创建向导（Agent 模拟 + 确认）
│       │   ├── Recorder/           # 操作录制器
│       │   ├── TaskList/           # 任务列表管理
│       │   ├── ExecutionLog/       # 执行日志与截图回放
│       │   ├── ScreenLive/         # 实时屏幕预览
│       │   └── Settings/           # 系统设置（LLM / OCR / 快捷键）
│       ├── components/
│       │   ├── TaskWizard/         # 任务创建向导（步骤确认面板）
│       │   ├── ScreenViewer/       # 屏幕查看器（标注 OCR 结果）
│       │   ├── StepRecorder/       # 步骤录制悬浮窗
│       │   ├── ResultDiff/         # 执行前后截图对比
│       │   └── StatusBar/          # 底部状态栏
│       ├── hooks/
│       │   ├── useAgent.ts         # 调用 Go Agent 绑定方法
│       │   ├── useScreenCapture.ts # 订阅屏幕截图事件流
│       │   └── useHotkey.ts        # 全局快捷键注册
│       ├── stores/                 # Zustand 状态管理
│       ├── bindings/               # Wails 自动生成的 Go 绑定类型
│       └── assets/                 # 静态资源
│
├── internal/                       # Go 后端核心逻辑
│   ├── agent/                      # Agent 调度核心
│   │   ├── parser.go               # 任务解析器 (NL → 结构化计划)
│   │   ├── planner.go              # 计划生成器
│   │   ├── simulator.go            # 模拟演示器（截图标注，不执行）
│   │   ├── executor.go             # 步骤执行器
│   │   ├── memory.go               # 短期/长期记忆
│   │   └── recovery.go             # 异常恢复策略
│   ├── vision/                     # 视觉感知模块
│   │   ├── capture.go              # 屏幕截图
│   │   ├── ocr.go                  # OCR 识别引擎
│   │   ├── detection.go            # UI 控件检测
│   │   ├── matching.go             # 图像模板匹配
│   │   └── layout.go               # 布局结构分析
│   ├── action/                     # 动作执行模块
│   │   ├── mouse.go                # 鼠标操作
│   │   ├── keyboard.go             # 键盘操作
│   │   ├── window.go               # 窗口查找与管理
│   │   └── clipboard.go            # 剪贴板读写
│   ├── monitor/                    # 监控与验证
│   │   ├── checker.go              # 步骤结果校验
│   │   └── logger.go               # 执行日志与截图序列
│   ├── llm/                        # LLM 网关
│   │   ├── gateway.go              # 多模型路由 + 重试
│   │   ├── openai.go               # OpenAI 适配器
│   │   ├── claude.go               # Claude 适配器
│   │   └── ollama.go               # Ollama 本地模型适配器
│   ├── scheduler/                  # 定时任务
│   │   └── cron.go                 # Cron 调度管理
│   ├── hotkey/                     # 全局快捷键
│   │   └── hook.go                 # 系统级热键注册
│   ├── model/                      # 数据模型
│   │   ├── task.go                 # 任务定义
│   │   ├── step.go                 # 操作步骤
│   │   └── execution.go            # 执行记录
│   └── store/                      # 持久化层
│       ├── sqlite.go               # SQLite GORM 初始化
│       └── repository.go           # 仓储接口
│
├── pkg/                            # 可复用公共包
│   ├── screenshot/                 # 跨平台高性能截图
│   ├── keymouse/                   # 跨平台键鼠封装
│   └── ocrclient/                  # OCR 客户端抽象
│
├── scripts/                        # 构建/签名脚本
│   ├── build.bat                   # Windows 构建
│   ├── build.sh                    # macOS/Linux 构建
│   └── installer.nsi               # NSIS 安装包脚本
├── config/                         # 默认配置
│   └── config.example.yaml
├── docs/                           # 文档
├── go.mod
├── go.sum
├── .gitignore
├── README.md
└── LICENSE
```

---

## 🎯 项目路线图

### Phase 1 — 基础能力（MVP）
- [ ] Wails 项目初始化，Go ↔ React 通信贯通
- [ ] 屏幕截图与区域选择（前端标注 OCR 结果）
- [ ] 集成 Tesseract OCR，实现屏幕文字识别与坐标定位
- [ ] 基础键鼠操作（点击、输入、拖拽、快捷键、滚轮）
- [ ] 任务脚本 JSON 定义与顺序执行
- [ ] 执行日志与截图序列存储
- [ ] 前端任务列表 + 手动创建任务界面

### Phase 2 — Agent 智能化
- [ ] LLM 网关：统一接入 OpenAI / Claude / Ollama 本地模型
- [ ] 自然语言任务解析（NL Prompt → 结构化操作序列）
- [ ] Agent 分层规划：屏幕感知 → 意图理解 → 步骤分解
- [ ] 步骤级截图验证 + 自动重试/回退机制
- [ ] 前端实时执行进度展示（步骤动画 + 截图流）

### Phase 3 — 智能任务创建
- [ ] 任务创建向导：Agent 逐步模拟演示 + 用户确认交互
- [ ] 截图标注高亮：OCR 识别结果 + 目标区域可视化标注
- [ ] 步骤编辑器：用户可修改、插入、删除、跳过步骤
- [ ] 任务模板保存：结构化 JSON 持久化，支持编辑和复用
- [ ] 操作录制：记录用户键鼠操作，自动生成任务步骤
- [ ] 录制回放：按录制轨迹回放并 OCR 校验

### Phase 4 — 生产增强
- [ ] 定时任务：Cron 后台静默执行 + 系统托盘驻留
- [ ] 全局快捷键：一键触发指定任务
- [ ] 执行仪表盘：成功率、耗时统计、热力分析 (ECharts)
- [ ] 任务模板导入/导出（JSON/YAML）
- [ ] 系统托盘 + 后台运行模式

### Phase 5 — 生态与分发
- [ ] 插件体系（gRPC 接口，自定义 OCR/执行器/LLM）
- [ ] 任务市场（社区共享模板 + 一键安装）
- [ ] 交叉编译 + 自动构建流水线 (GitHub Actions)
- [ ] Windows (.exe) / macOS (.dmg) / Linux (.AppImage) 三平台发布
- [ ] 自动更新机制

---

## 🔨 快速开始

> ⚠️ 项目处于早期规划阶段，以下为预期开发流程。

### 开发环境要求

| 工具 | 版本 | 用途 |
|------|------|------|
| **Go** | >= 1.22 | 后端开发 |
| **Node.js** | >= 20 LTS | 前端构建 |
| **pnpm** | >= 9 | 包管理器 |
| **Wails CLI** | v3 | 项目脚手架、热重载、构建打包 |
| **Tesseract** | >= 5.0 | OCR 引擎开发库 |
| **GCC / MinGW** | (Windows) | CGo 编译 Tesseract |
| **WebView2** | (Windows) | Windows 10+ 预装 |

### 安装 Wails

```bash
# 安装 Wails CLI
go install github.com/wailsapp/wails/v3/cmd/wails3@latest

# 验证安装
wails3 doctor
```

### 安装 OCR 引擎

```bash
# Windows: 下载安装
# https://github.com/UB-Mannheim/tesseract/wiki

# macOS
brew install tesseract

# Linux (Debian/Ubuntu)
sudo apt install tesseract-ocr libtesseract-dev

# 下载中文语言包
# https://github.com/tesseract-ocr/tessdata
```

### 克隆并运行

```bash
# 克隆仓库
git clone https://github.com/yc446833448/VisuTask.git
cd VisuTask

# 安装前端依赖
cd frontend && pnpm install && cd ..

# 开发模式（热重载，前端 HMR + Go 自动重编译）
wails3 dev

# 构建生产版本
wails3 build

# 构建产物位于 build/bin/ 目录
```

### 开发模式体验

```bash
wails3 dev
```

- Go 代码修改 → 自动重编译重启
- React 代码修改 → Vite HMR 毫秒级热更新
- Wails 自动打开桌面窗口，内嵌 WebView 加载前端页面
- 前端通过 `wailsjs/go` 模块直接调用 Go 函数（类型安全 + IDE 自动补全）

---

### 基础使用示例

#### Go 后端 — 暴露给前端的绑定方法

```go
// internal/agent/service.go
package agent

import (
    "context"
    "github.com/wailsapp/wails/v3/pkg/application"
)

type AgentService struct {
    parser   *Parser
    planner  *Planner
    executor *Executor
    ocr      *vision.OCREngine
}

// CreateTaskPlan 根据自然语言描述生成任务计划，供前端展示确认
func (a *AgentService) CreateTaskPlan(description string) (*TaskPlan, error) {
    // 1. LLM 解析自然语言意图
    intent, err := a.parser.Parse(context.Background(), description)
    if err != nil {
        return nil, err
    }

    // 2. 截图 + OCR 感知当前屏幕状态
    screenshot := captureScreen()
    ocrResults := a.ocr.Recognize(screenshot)

    // 3. 结合意图和屏幕状态生成操作计划
    plan, err := a.planner.GeneratePlan(context.Background(), intent, screenshot, ocrResults)
    if err != nil {
        return nil, err
    }

    return plan, nil
}

// SimulateStep 模拟执行单个步骤（不真正操作，仅截图标注目标）
func (a *AgentService) SimulateStep(stepIndex int, plan *TaskPlan) (*StepSimulation, error) {
    step := plan.Steps[stepIndex]

    // 截图 → OCR → 定位目标区域 → 生成标注图
    screenshot := captureScreen()
    ocrResults := a.ocr.Recognize(screenshot)
    target := a.planner.LocateTarget(screenshot, ocrResults, step.Target)

    return &StepSimulation{
        StepIndex:  stepIndex,
        Action:     step.Action,
        Target:     step.Target,
        Screenshot: screenshot,    // base64 截图
        Annotation: target.Overlay, // 标注框高亮目标区域
        Confidence: target.Score,   // 匹配置信度
    }, nil
}

// ConfirmAndSave 用户确认所有步骤后保存为任务模板
func (a *AgentService) ConfirmAndSave(plan *TaskPlan, name string) (*Task, error) {
    task := &Task{
        Name:    name,
        Steps:   plan.Steps,
        Created: time.Now(),
    }
    return a.store.SaveTask(task)
}

// RunTask 正式执行已确认的任务
func (a *AgentService) RunTask(taskID string) (*ExecutionResult, error) {
    task, err := a.store.GetTask(taskID)
    if err != nil {
        return nil, err
    }

    for i, step := range task.Steps {
        // 推送进度到前端
        application.Get().EmitEvent("step:progress", StepProgress{
            Index:      i + 1,
            Total:      len(task.Steps),
            Action:     step.Action,
            Target:     step.Target,
            Screenshot: captureScreen(),
        })

        if err := a.executor.ExecuteStep(ctx, step); err != nil {
            application.Get().EmitEvent("step:error", err)
        }
    }

    return result, nil
}
```

#### React 前端 — 任务创建向导

```tsx
// pages/Designer/index.tsx
import {
  CreateTaskPlan,
  SimulateStep,
  ConfirmAndSave,
  RunTask,
} from "@wailsjs/go/agent/AgentService";
import { EventsOn } from "@wailsjs/runtime/runtime";

function TaskWizard() {
  const [description, setDescription] = useState("");
  const [plan, setPlan] = useState<TaskPlan | null>(null);
  const [currentStep, setCurrentStep] = useState(0);
  const [simulation, setSimulation] = useState<StepSimulation | null>(null);

  // 第一步：用户描述任务，Agent 生成计划
  const handleCreate = async () => {
    const plan = await CreateTaskPlan(description);
    setPlan(plan);
    setCurrentStep(0);
  };

  // 第二步：逐步模拟演示，用户确认
  const handleSimulate = async () => {
    if (!plan) return;
    const sim = await SimulateStep(currentStep, plan);
    setSimulation(sim);
  };

  // 用户确认当前步骤，进入下一步
  const handleConfirmStep = () => {
    if (plan && currentStep < plan.steps.length - 1) {
      setCurrentStep(currentStep + 1);
      handleSimulate();
    }
  };

  // 全部确认，保存任务模板
  const handleSave = async () => {
    if (!plan) return;
    const task = await ConfirmAndSave(plan, description);
    console.log("任务已保存:", task);
  };

  // 正式执行
  const handleRun = async () => {
    const result = await RunTask(plan?.taskID);
    console.log("执行完成:", result);
  };

  // 监听执行进度
  EventsOn("step:progress", (p: StepProgress) => {
    console.log(`执行中: ${p.index}/${p.total} - ${p.action}`);
  });

  return (
    <div>
      {/* 输入区 */}
      <Textarea
        value={description}
        onChange={(e) => setDescription(e.target.value)}
        placeholder="描述你的任务，例如：打开浏览器访问 GitHub 搜索 RPA"
      />
      <Button onClick={handleCreate}>生成计划</Button>

      {/* 步骤确认区 */}
      {plan && (
        <div>
          <h3>
            步骤 {currentStep + 1}/{plan.steps.length}：
            {plan.steps[currentStep].action} - {plan.steps[currentStep].target}
          </h3>
          <Button onClick={handleSimulate}>模拟演示</Button>

          {simulation && (
            <div>
              {/* 截图 + OCR 标注高亮 */}
              <img src={`data:image/png;base64,${simulation.screenshot}`} />
              <img
                src={`data:image/png;base64,${simulation.annotation}`}
                style={{ position: "absolute" }}
              />
              <p>置信度: {(simulation.confidence * 100).toFixed(1)}%</p>
            </div>
          )}

          <Button onClick={handleConfirmStep}>确认此步</Button>
          <Button onClick={handleSave}>保存任务</Button>
          <Button onClick={handleRun}>正式执行</Button>
        </div>
      )}
    </div>
  );
}
```

#### 结构化任务模板（Agent 自动生成）

```json
{
  "id": "task_001",
  "name": "数据录入",
  "description": "将客户信息从 Excel 录入到 CRM 系统",
  "trigger": {
    "type": "manual",
    "hotkey": "Ctrl+Shift+D"
  },
  "variables": {
    "name": "张三",
    "phone": "13800138000"
  },
  "steps": [
    {
      "id": 1,
      "action": "click",
      "target": "新建按钮",
      "targetOCR": "新建",
      "confidence": 0.95,
      "confirmed": true
    },
    {
      "id": 2,
      "action": "input",
      "target": "姓名输入框",
      "targetOCR": "姓名",
      "value": "{{name}}",
      "confidence": 0.92,
      "confirmed": true
    },
    {
      "id": 3,
      "action": "input",
      "target": "电话输入框",
      "targetOCR": "电话",
      "value": "{{phone}}",
      "confidence": 0.91,
      "confirmed": true
    },
    {
      "id": 4,
      "action": "click",
      "target": "保存",
      "targetOCR": "保存",
      "confidence": 0.98,
      "confirmed": true
    },
    {
      "id": 5,
      "action": "verify",
      "target": "保存成功",
      "targetOCR": "保存成功",
      "timeout": 5000,
      "confidence": 0.88,
      "confirmed": true
    }
  ],
  "createdAt": "2026-06-17T10:30:00Z",
  "updatedAt": "2026-06-17T10:30:00Z"
}
```

---

## 🤝 贡献指南

本项目处于早期开发阶段，欢迎各种形式的贡献：

1. **Fork** 本仓库
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

### 开发规范

- **前端**：ESLint + Prettier，遵循 React 最佳实践
- **后端**：`gofmt` + `golangci-lint`，遵循 Effective Go
- **Commit**：遵循 Conventional Commits
- **Wails 绑定**：Go 导出函数遵循 Wails 命名约定，自动生成前端类型

---

## 📄 许可证

本项目基于 MIT 许可证开源，详见 [LICENSE](LICENSE) 文件。

---

## 🙏 致谢

- [Wails](https://wails.io/) — Go + Web 前端的桌面应用框架
- [Tesseract OCR](https://github.com/tesseract-ocr/tesseract) — Google 开源 OCR 引擎
- [GoCV](https://github.com/hybridgroup/gocv) — Go 语言 OpenCV 绑定
- [robotgo](https://github.com/go-vgo/robotgo) — Go 跨平台桌面自动化
- [shadcn/ui](https://ui.shadcn.com/) — 现代化可定制 React 组件库
- [Lucide](https://lucide.dev/) — 现代 SVG 图标库，shadcn/ui 官方配套
- [AutoGPT](https://github.com/Significant-Gravitas/AutoGPT) — 自主 AI Agent
- [OmniParser](https://microsoft.github.io/OmniParser/) — 微软屏幕理解模型

---

<p align="center">
  <b>VisuTask</b> — 让 AI 看懂屏幕，替你动手<br>
  单一原生应用 · 双击即用 · 离线可用<br>
  Made with ❤️ by <a href="https://github.com/yc446833448">yc446833448</a>
</p>
