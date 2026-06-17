# VisuTask 前端开发规范

当为 VisuTask 编写前端代码时，必须严格遵守以下规范。

---

## 技术栈

- React 18 + TypeScript + Vite
- shadcn/ui (组件库，不直接安装 antd/mui)
- Tailwind CSS (样式，不写 CSS 文件)
- Lucide React (图标，不引入其他图标库)
- Wails Runtime (IPC 通信)
- Zustand (状态管理)
- recharts (图表)
- react-markdown + remark-gfm (Markdown 渲染)
- sonner (Toast 通知)

---

## 窗口与布局

- 默认窗口 1024×768，最小 800×600
- 无侧边栏，首页卡片导航，子页面面包屑返回
- 全局结构: Header(56px) + Content(自适应) + StatusBar(28px)
- Content 内 padding: `p-6`
- 单断点 900px: `@media (max-width: 899px)` 触发紧凑模式

```tsx
// 页面模板
export default function SomePage() {
  return (
    <div className="p-6">
      <Breadcrumb items={[{ label: "首页", path: "/" }, { label: "当前页" }]} />
      <div className="mt-6 space-y-6">
        {/* 页面内容 */}
      </div>
    </div>
  )
}
```

---

## 色彩规则

- **禁止硬编码颜色值**，全部使用 Tailwind 语义 class
- 主文字: `text-foreground`
- 次要文字: `text-muted-foreground`
- 主背景: `bg-background`
- 卡片背景: `bg-card`
- 区块背景: `bg-muted`
- 边框: `border`

状态色（固定值，不随主题切换）:
- 成功/已连接/已确认: `text-green-500`
- 运行中/当前目标: `text-blue-500`
- 暂停/警告: `text-yellow-500`
- 失败/错误/断开: `text-red-500`
- 空闲/未连接: `text-zinc-500`

---

## 字体规则

```
正文:       text-sm (14px)      ← 默认，大部分场景
页面标题:   text-xl font-semibold
区块标题:   text-base font-semibold
描述文字:   text-sm text-muted-foreground
辅助标签:   text-xs text-muted-foreground
表格表头:   text-xs font-medium text-muted-foreground
```

- 最大字号不超过 `text-xl` (20px)
- 正文统一 `text-sm`，不用 `text-base` 作正文
- 字重只用 3 级: `font-normal`(400) / `font-medium`(500) / `font-semibold`(600)

---

## 间距规则

```
页面 padding:          p-6
区块间距:              space-y-6
卡片 padding:          p-4
卡片 grid 间距:        gap-6
面包屑到内容:          mt-6
标题到描述:            mt-1
标题到内容:            mt-3 或 mt-4
Tab 栏到内容:          mt-4
按钮之间:              gap-2
图标与文字:            gap-2
表单 label 到 input:   space-y-2
表单项之间:            space-y-4
```

---

## 按钮规则

尺寸:
- `h-8 px-3 text-xs` — 表格行内、紧凑区域
- `h-9 px-4 text-sm` — 常规按钮 (默认)
- `h-10 px-6 text-sm` — 重要 CTA
- `h-9 w-9` — 纯图标按钮

变体选择:
- 每个区域最多 1 个 `variant="default"` (主操作)
- 次要操作用 `variant="outline"` 或 `variant="secondary"`
- 弱操作用 `variant="ghost"`
- 危险操作用 `variant="destructive"`

```tsx
// 正确示例
<div className="flex gap-2">
  <Button>加入任务</Button>                          {/* 主操作 */}
  <Button variant="outline">模拟演示</Button>         {/* 次要 */}
  <Button variant="ghost" size="icon"><Trash2 /></Button>  {/* 弱操作 */}
  <Button variant="destructive">停止</Button>         {/* 危险 */}
</div>
```

---

## 图标规则

```
尺寸:  14px(Badge内) / 16px(按钮/表格) / 20px(导航) / 24px(标题旁) / 48px(首页卡片)
颜色:  默认 text-muted-foreground, 状态用状态色
粗细:  strokeWidth 保持默认 2, 不修改
```

从 lucide-react 按需引入，不用全量导入:
```tsx
import { Settings, Play, Pause, Trash2 } from "lucide-react"
```

---

## 圆角规则

- 按钮/输入框/Select: `rounded-md` (6px)
- 卡片/Dialog: `rounded-lg` (8px)
- Badge/Tag: `rounded-sm` (4px)
- 头像/圆形按钮: `rounded-full`

---

## 阴影规则

- 默认不加阴影
- 卡片: `shadow-sm`, hover 时 `shadow`
- 首页卡片: hover 时 `shadow-md` + `-translate-y-0.5`
- Dialog/Toast: `shadow-lg`

---

## 卡片规范

```tsx
// 普通卡片
<Card className="p-4 shadow-sm hover:shadow transition-shadow">
  <h3 className="text-base font-semibold">标题</h3>
  <p className="mt-1 text-sm text-muted-foreground">描述</p>
</Card>

// 首页功能卡片
<Card className="p-6 cursor-pointer shadow-sm hover:shadow-md hover:-translate-y-0.5 transition-all">
  <Sparkles className="h-12 w-12 text-muted-foreground" />
  <h3 className="mt-4 text-lg font-semibold">创建脚本</h3>
  <p className="mt-1 text-sm text-muted-foreground">自然语言 AI 对话</p>
</Card>
```

---

## 表格规范

```tsx
<Table>
  <TableHeader>
    <TableRow>
      <TableHead className="text-xs font-medium uppercase text-muted-foreground">
        名称
      </TableHead>
    </TableRow>
  </TableHeader>
  <TableBody>
    <TableRow className="hover:bg-muted/50">
      <TableCell className="px-4 py-3 text-sm">内容</TableCell>
    </TableRow>
  </TableBody>
</Table>
```

- 行高 `h-12` (48px)
- 表头 `text-xs font-medium text-muted-foreground`
- 操作列用 `variant="ghost"` 图标按钮 `h-8 w-8`

---

## 表单规范

```tsx
<div className="space-y-4">
  <div className="space-y-2">
    <Label className="text-sm font-medium">选择脚本</Label>
    <Select>
      <SelectTrigger className="h-9">
        <SelectValue placeholder="选择脚本..." />
      </SelectTrigger>
      <SelectContent>
        <SelectItem value="1">数据录入 CRM</SelectItem>
      </SelectContent>
    </Select>
    <p className="text-xs text-muted-foreground">辅助说明</p>
  </div>
</div>
```

- 所有输入控件高度 `h-9`
- Label: `text-sm font-medium`
- 描述: `text-xs text-muted-foreground`
- 错误: `text-xs text-destructive`

---

## 状态徽章

```tsx
// 任务状态
<Badge variant="outline" className="text-green-500 border-green-500">
  <Play className="mr-1 h-3 w-3" /> 运行中
</Badge>

<Badge variant="outline" className="text-zinc-500 border-zinc-500">
  <Circle className="mr-1 h-3 w-3" /> 空闲
</Badge>

// VIP 等级
<Badge variant="secondary">VIP1</Badge>
```

---

## 动效规则

```
过渡时长:    duration-150 (默认) / duration-300 (慢速)
hover 上浮:  hover:-translate-y-0.5 (仅首页卡片)
hover 阴影:  hover:shadow-md (卡片)
hover 背景:  hover:bg-muted/50 (表格行)
按钮按下:    active:scale-[0.98]
运行中脉冲:  animate-pulse (状态指示点)
```

- 页面切换不做动画
- 数据刷新不做渐入渐出
- 展开折叠用 shadcn 内置 animate-in/animate-out

---

## Toast 通知

```tsx
import { toast } from "sonner"

// 成功
toast.success("任务已创建")

// 失败
toast.error("执行失败：步骤4超时")

// 普通
toast("正在执行...", { description: "预计耗时 3 分钟" })
```

---

## 确认弹窗

```tsx
<AlertDialog>
  <AlertDialogTrigger asChild>
    <Button variant="ghost" size="icon"><Trash2 className="h-4 w-4" /></Button>
  </AlertDialogTrigger>
  <AlertDialogContent>
    <AlertDialogHeader>
      <AlertDialogTitle>确认删除</AlertDialogTitle>
      <AlertDialogDescription>删除后无法恢复，确定要删除吗？</AlertDialogDescription>
    </AlertDialogHeader>
    <AlertDialogFooter>
      <AlertDialogCancel>取消</AlertDialogCancel>
      <AlertDialogAction className="bg-destructive text-white">
        删除
      </AlertDialogAction>
    </AlertDialogFooter>
  </AlertDialogContent>
</AlertDialog>
```

---

## 截图标注组件

```tsx
// 标注框样式
<div
  className="absolute border-2 rounded-sm"
  style={{
    left: annotation.rect.x,
    top: annotation.rect.y,
    width: annotation.rect.width,
    height: annotation.rect.height,
    borderColor: statusColor,       // 蓝/绿/红
    backgroundColor: `${statusColor}26`, // 15% opacity
  }}
>
  <span className="absolute -top-5 left-0 text-xs bg-black/70 text-white px-1 rounded">
    {annotation.text} {Math.round(annotation.confidence * 100)}%
  </span>
</div>
```

---

## 禁止事项

- ❌ 不使用 Ant Design / Material UI / Element Plus
- ❌ 不写独立 CSS/SCSS 文件，只用 Tailwind class
- ❌ 不用内联 style (截图标注的动态定位除外)
- ❌ 不硬编码颜色值 (如 `color: #333`)
- ❌ 正文不用 `text-base` (16px)，统一 `text-sm` (14px)
- ❌ 标题不超过 `text-xl` (20px)
- ❌ 不引入其他图标库，只用 Lucide
- ❌ 不做页面切换动画
- ❌ 不在组件内直接调用 Wails API，通过 hook 封装

---

## 页面路由

| 路由 | 页面 | 组件 |
|------|------|------|
| `/` | 首页 | `pages/Home.tsx` |
| `/script/new` | 创建脚本 | `pages/ScriptNew.tsx` |
| `/scripts` | 脚本库 | `pages/Scripts.tsx` |
| `/market` | 脚本市场 | `pages/Market.tsx` |
| `/stats` | 用量统计 | `pages/Stats.tsx` |
| `/tasks` | 任务管理 | `pages/Tasks.tsx` |
| `/task/new` | 新建任务 | `pages/TaskNew.tsx` |
| `/task/:id/monitor` | 任务监控 | `pages/TaskMonitor.tsx` |
| `/settings` | 设置 | `pages/Settings.tsx` |
