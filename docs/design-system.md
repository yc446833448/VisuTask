# VisuTask 前端 UI 设计规范

> 基于 shadcn/ui + Tailwind CSS + Lucide Icons
> 风格：简洁、克制、信息密度适中

---

## 一、色彩系统

### 1.1 基础色

基于 shadcn/ui CSS 变量，支持亮色/暗色自动切换。

| Token | 亮色值 | 暗色值 | Tailwind Class | 用途 |
|-------|--------|--------|----------------|------|
| `--background` | `#ffffff` | `#09090b` | `bg-background` | 页面背景 |
| `--foreground` | `#0a0a0a` | `#fafafa` | `text-foreground` | 主文字 |
| `--card` | `#ffffff` | `#18181b` | `bg-card` | 卡片背景 |
| `--muted` | `#f4f4f5` | `#27272a` | `bg-muted` | 次要背景/区块 |
| `--muted-foreground` | `#71717a` | `#a1a1aa` | `text-muted-foreground` | 次要文字 |
| `--border` | `#e4e4e7` | `#27272a` | `border` | 边框 |
| `--input` | `#e4e4e7` | `#27272a` | `border-input` | 输入框边框 |
| `--ring` | `#3b82f6` | `#60a5fa` | `ring-ring` | 焦点环 |

### 1.2 语义色

| Token | 亮色值 | 暗色值 | 用途 |
|-------|--------|--------|------|
| `--primary` | `#18181b` | `#fafafa` | 主按钮、强调元素 |
| `--primary-foreground` | `#fafafa` | `#18181b` | 主按钮文字 |
| `--secondary` | `#f4f4f5` | `#27272a` | 次按钮 |
| `--secondary-foreground` | `#18181b` | `#fafafa` | 次按钮文字 |
| `--destructive` | `#ef4444` | `#dc2626` | 危险操作（删除、停止） |
| `--accent` | `#f4f4f5` | `#27272a` | hover 高亮 |

### 1.3 状态色

直接使用的固定色值（不随主题切换）：

| 名称 | 色值 | Tailwind | 用途 |
|------|------|----------|------|
| 成功绿 | `#22c55e` | `text-green-500` | 任务成功、已确认、已连接 |
| 运行蓝 | `#3b82f6` | `text-blue-500` | 运行中、当前目标 |
| 警告黄 | `#eab308` | `text-yellow-500` | 暂停、警告 |
| 失败红 | `#ef4444` | `text-red-500` | 失败、错误、断开 |
| 空闲灰 | `#71717a` | `text-zinc-500` | 空闲、未连接 |

### 1.4 标注框色

| 状态 | 边框色 | 背景色 |
|------|--------|--------|
| 当前目标 | `#3b82f6` | `rgba(59, 130, 246, 0.15)` |
| 已确认 | `#22c55e` | `rgba(34, 197, 94, 0.15)` |
| 识别失败 | `#ef4444` | `rgba(239, 68, 68, 0.15)` |

---

## 二、字体排版

### 2.1 字体族

```css
font-family: system-ui, -apple-system, "Segoe UI", "PingFang SC",
             "Microsoft YaHei", sans-serif;
```

使用系统默认字体栈，中英文均清晰。

### 2.2 字号层级

| 层级 | Tailwind | 大小 | 行高 | 用途 |
|------|----------|------|------|------|
| 页面标题 | `text-xl` | 20px | 28px | 页面主标题（面包屑下方） |
| 区块标题 | `text-base` | 16px | 24px | Card 标题、区块标题 |
| 正文 | `text-sm` | 14px | 20px | 正文内容、表格内容、列表项 |
| 辅助文字 | `text-xs` | 12px | 16px | 描述、时间戳、标签、Badge |

**原则：**
- 最大字号不超过 `text-xl` (20px)，保持克制
- 正文统一 `text-sm` (14px)，不使用 16px 作为正文
- 数据密集型区域（表格、日志）统一 `text-sm`

### 2.3 字重

| 权重 | Tailwind | 用途 |
|------|----------|------|
| 400 (regular) | `font-normal` | 正文、描述 |
| 500 (medium) | `font-medium` | 按钮、标签、表格表头 |
| 600 (semibold) | `font-semibold` | 页面标题、区块标题 |

### 2.4 字色规则

| 场景 | Tailwind Class | 效果 |
|------|----------------|------|
| 主标题 | `text-foreground` | 跟随主题的主文字色 |
| 正文 | `text-foreground` | 同上 |
| 描述/副标题 | `text-muted-foreground` | 灰色次要文字 |
| 表格表头 | `text-muted-foreground font-medium` | 灰色中粗 |
| 链接/可点击 | `text-blue-500 hover:text-blue-600` | 蓝色 + hover 加深 |
| 禁用 | `text-muted-foreground/50` | 半透明灰 |

---

## 三、间距系统

使用 Tailwind 默认间距，基础单位 4px。

### 3.1 间距表

| Token | 值 | Tailwind | 用途 |
|-------|-----|----------|------|
| 1 | 4px | `p-1` / `gap-1` | 图标与文字间距、紧凑元素 |
| 2 | 8px | `p-2` / `gap-2` | 按钮内间距、表单项间距 |
| 3 | 12px | `p-3` / `gap-3` | 卡片内 padding、列表项间距 |
| 4 | 16px | `p-4` / `gap-4` | 区块间距、卡片 padding |
| 6 | 24px | `p-6` / `gap-6` | 页面级 padding、大区块间距 |
| 8 | 32px | `p-8` / `gap-8` | 页面主内容区 padding (顶部) |

### 3.2 页面级间距

```
页面内容区 padding:     p-6 (24px)
区块之间间距:           space-y-6 (24px)
卡片内 padding:         p-4 (16px)
表格行高:               h-12 (48px)
```

### 3.3 组件级间距

```
按钮内 padding:         px-4 py-2 (16px / 8px)
图标与文字间距:         gap-2 (8px)
表单 label 与 input:   space-y-2 (8px)
表单项之间:             space-y-4 (16px)
卡片标题与内容:         space-y-3 (12px)
```

---

## 四、圆角

统一使用 shadcn/ui 的圆角变量：

| Token | 值 | Tailwind | 用途 |
|-------|-----|----------|------|
| `--radius` | 6px | `rounded-md` | 按钮、输入框、Select |
| — | 8px | `rounded-lg` | 卡片、Dialog |
| — | 4px | `rounded-sm` | Badge、Tag |
| — | 9999px | `rounded-full` | 头像、圆形按钮 |

**原则：** 一个页面内圆角层级不超过 3 级。

---

## 五、阴影

克制使用阴影，仅用于悬浮和层级区分：

| 层级 | Tailwind | 用途 |
|------|----------|------|
| 无 | `shadow-none` | 默认状态，大部分元素 |
| 微 | `shadow-sm` | 卡片默认状态 |
| 低 | `shadow` | 卡片 hover |
| 中 | `shadow-md` | 下拉菜单、Popover |
| 高 | `shadow-lg` | Dialog、Toast |

**原则：** 默认不加阴影，hover 时才出现。

---

## 六、按钮规范

### 6.1 尺寸

| 尺寸 | 高度 | padding | 字号 | Tailwind | 用途 |
|------|------|---------|------|----------|------|
| sm | 32px | `px-3` | `text-xs` | `h-8 px-3 text-xs` | 表格行内操作、紧凑区域 |
| default | 36px | `px-4` | `text-sm` | `h-9 px-4 text-sm` | 常规按钮 |
| lg | 40px | `px-6` | `text-sm` | `h-10 px-6 text-sm` | 首页卡片、重要 CTA |
| icon | 36px | `p-0` | — | `h-9 w-9` | 纯图标按钮 |

### 6.2 变体

| 变体 | 背景 | 文字 | 用途 | shadcn variant |
|------|------|------|------|----------------|
| default | `bg-primary` | `text-primary-foreground` | 主要操作（加入任务、生成计划） | `default` |
| secondary | `bg-secondary` | `text-secondary-foreground` | 次要操作（取消、返回） | `secondary` |
| outline | `border bg-transparent` | `text-foreground` | 一般操作（模拟演示、刷新） | `outline` |
| ghost | `transparent` | `text-foreground` | 弱操作（编辑、更多） | `ghost` |
| destructive | `bg-destructive` | `white` | 危险操作（删除、停止） | `destructive` |
| link | `transparent` | `text-blue-500 underline` | 文字链接 | `link` |

### 6.3 按钮使用规则

```
每个操作区域最多 1 个 default 按钮（主操作）
其余操作使用 outline / ghost / secondary
危险操作始终使用 destructive
图标按钮使用 ghost variant + Lucide 图标 (size=16)
禁用状态: disabled:opacity-50 disabled:pointer-events-none
```

---

## 七、表单规范

### 7.1 输入框

| 属性 | 值 |
|------|-----|
| 高度 | `h-9` (36px) |
| padding | `px-3 py-1` |
| 边框 | `border rounded-md` |
| 焦点 | `focus-visible:ring-2 focus-visible:ring-ring` |
| 字号 | `text-sm` |
| placeholder | `text-muted-foreground` |

### 7.2 Select

与输入框保持一致的高度和圆角。

### 7.3 表单布局

```
Label:      text-sm font-medium text-foreground
间距:       Label 与 Input 之间 space-y-2 (8px)
表单项之间:  space-y-4 (16px)
描述文字:    text-xs text-muted-foreground (Input 下方)
错误提示:    text-xs text-destructive (Input 下方)
```

---

## 八、表格规范

| 属性 | 值 |
|------|-----|
| 表头 | `text-xs font-medium text-muted-foreground uppercase` |
| 行高 | `h-12` (48px) |
| 单元格 padding | `px-4 py-3` |
| 边框 | `border-b` (行间分隔线) |
| hover | `hover:bg-muted/50` |
| 展开行 | `bg-muted/30` |
| 操作列 | 图标按钮 `ghost` variant, `h-8 w-8` |

---

## 九、卡片规范

| 属性 | 值 |
|------|-----|
| 背景 | `bg-card` |
| 边框 | `border rounded-lg` |
| 阴影 | `shadow-sm` → hover `shadow` |
| padding | `p-4` (16px) |
| 标题 | `text-base font-semibold` |
| 描述 | `text-sm text-muted-foreground` |

首页功能卡片特殊规格：
```
padding:     p-6 (24px)
图标:        48px, text-muted-foreground
标题:        text-lg font-semibold, mt-4
描述:        text-sm text-muted-foreground, mt-1
hover:       shadow-md + -translate-y-0.5 transition
cursor:      pointer
```

---

## 十、图标规范

### 10.1 尺寸

| 尺寸 | 值 | 用途 |
|------|-----|------|
| xs | 14px | Badge 内、辅助文字旁 |
| sm | 16px | 按钮内、表格操作列、表单项 |
| md | 20px | 导航、面包屑 |
| lg | 24px | 区块标题旁 |
| xl | 48px | 首页功能卡片 |

### 10.2 颜色

| 场景 | 颜色 |
|------|------|
| 默认 | `text-muted-foreground` |
| 可点击 | `text-foreground hover:text-foreground` |
| 状态指示 | 使用状态色（绿/蓝/黄/红/灰） |
| 首页卡片 | `text-muted-foreground` |

### 10.3 stroke 宽度

统一使用 Lucide 默认 `strokeWidth: 2`，不修改。

---

## 十一、动效

### 11.1 过渡

| 属性 | 值 | Tailwind |
|------|-----|----------|
| 默认 | 150ms ease | `transition-all duration-150` |
| 慢速 | 300ms ease | `transition-all duration-300` |

### 11.2 可动效属性

```
hover 上浮:      hover:-translate-y-0.5 (首页卡片)
hover 阴影:      hover:shadow-md (卡片)
hover 背景:      hover:bg-muted/50 (表格行)
按钮 press:      active:scale-[0.98]
展开/折叠:       data-[state=open]:animate-in data-[state=closed]:animate-out
脉冲 (运行中):   animate-pulse (状态指示点)
```

### 11.3 不使用动效的场景

- 页面切换（直接跳转，不做过渡动画）
- 数据刷新（直接替换，不做渐入渐出）

---

## 十二、响应式断点

针对 1024→800 缩放的处理：

| 断点 | 条件 | 处理 |
|------|------|------|
| 标准 | ≥ 900px | 完整布局 |
| 紧凑 | < 900px | 缩减 padding, 隐藏次要列, 调整比例 |

**不使用 Tailwind 默认断点**（sm/md/lg/xl），仅使用自定义的 `900px` 断点。

```css
@media (max-width: 899px) {
  /* 紧凑模式 */
}
```

---

## 十三、组件间距速查

```
Header ↔ Content:            0 (紧贴)
Content ↔ StatusBar:         0 (紧贴)
Content 内部 padding:        p-6
面包屑 ↔ 内容:               mb-6 (24px)
Tab 栏 ↔ 内容:               mt-4 (16px)
标题 ↔ 描述:                 mt-1 (4px)
标题 ↔ 内容:                 mt-3 (12px)
区块 ↔ 区块:                 mt-6 (24px)
卡片内 padding:              p-4
卡片 ↔ 卡片 (grid):          gap-6 (24px)
表格行内 padding:            px-4 py-3
按钮 ↔ 按钮:                 gap-2 (8px)
图标 ↔ 文字:                 gap-2 (8px)
```

---

## 十四、暗色模式注意事项

- 所有颜色使用 CSS 变量（`bg-background`, `text-foreground`），不硬编码
- 截图标注框背景使用 `rgba` 半透明色，亮暗模式通用
- 阴影在暗色模式下不可见时，改用 `border` 区分层级
- 状态色（绿/蓝/黄/红）在暗色模式下使用稍亮的变体
