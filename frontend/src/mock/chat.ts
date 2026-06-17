import type { ChatMessage } from '@/types'

export const mockChatMessages: ChatMessage[] = [
  {
    id: 'msg_1',
    role: 'assistant',
    content: '你好！请描述你想要自动化的操作流程。我会帮你规划操作步骤。',
    timestamp: '10:00',
  },
  {
    id: 'msg_2',
    role: 'user',
    content: '把 Excel 里的客户信息逐一录入到网页 CRM 系统中，一共 100 条数据',
    timestamp: '10:01',
  },
  {
    id: 'msg_3',
    role: 'assistant',
    content: '我来为你规划操作步骤。已识别到当前屏幕上有 Excel 和 Chrome 窗口，生成了 8 个操作步骤，请查看右侧预览。如需调整，请告诉我。',
    timestamp: '10:01',
  },
  {
    id: 'msg_4',
    role: 'user',
    content: '第3步改成先输入电话再输入姓名',
    timestamp: '10:02',
  },
  {
    id: 'msg_5',
    role: 'assistant',
    content: '已调整顺序，右侧预览已更新。脚本文件已自动保存。',
    timestamp: '10:02',
  },
]

export const mockScriptMarkdown = `## 数据录入 CRM

**目标应用：** Excel + 网页 CRM
**循环次数：** 100

### 步骤

1. **click** → "新建"按钮
   - 目标OCR: \`新建\`
   - 置信度: 95%

2. **input** → "姓名"输入框
   - 值: \`{{name}}\`
   - 置信度: 92%

3. **input** → "电话"输入框
   - 值: \`{{phone}}\`
   - 置信度: 91%

4. **click** → "保存"按钮
   - 置信度: 98%

5. **verify** → "保存成功"提示
   - 超时: 5000ms
   - 置信度: 88%

### 变量

| 变量名 | 默认值 |
|--------|--------|
| loopCount | 100 |
| startRow | 2 |
| name | 张三 |
| phone | 13800138000 |
`
