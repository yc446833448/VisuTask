import type { Script } from '@/types'

export const mockScripts: Script[] = [
  {
    id: 'script_001',
    name: '数据录入 CRM',
    description: '将 Excel 客户信息逐条录入到网页 CRM 系统',
    steps: [
      { id: 1, action: 'click', target: '新建按钮', targetOCR: '新建', confidence: 0.95, confirmed: true },
      { id: 2, action: 'input', target: '姓名输入框', targetOCR: '姓名', value: '{{name}}', confidence: 0.92, confirmed: true },
      { id: 3, action: 'input', target: '电话输入框', targetOCR: '电话', value: '{{phone}}', confidence: 0.91, confirmed: true },
      { id: 4, action: 'click', target: '保存按钮', targetOCR: '保存', confidence: 0.98, confirmed: true },
      { id: 5, action: 'verify', target: '保存成功提示', targetOCR: '保存成功', timeout: 5000, confidence: 0.88, confirmed: true },
    ],
    variables: { loopCount: '100', startRow: '2' },
    createdAt: '2026-06-17T10:30:00Z',
    updatedAt: '2026-06-17T10:30:00Z',
    taskCount: 3,
  },
  {
    id: 'script_002',
    name: '安装 Python 3.12',
    description: '自动安装 Python 并配置环境变量',
    steps: [
      { id: 1, action: 'click', target: '安装程序', targetOCR: 'Install', confidence: 0.96, confirmed: true },
      { id: 2, action: 'click', target: 'Next', targetOCR: 'Next', confidence: 0.99, confirmed: true },
      { id: 3, action: 'click', target: 'Add to PATH', targetOCR: 'Add Python to PATH', confidence: 0.94, confirmed: true },
      { id: 4, action: 'click', target: 'Install Now', targetOCR: 'Install Now', confidence: 0.97, confirmed: true },
      { id: 5, action: 'click', target: 'Finish', targetOCR: 'Finish', confidence: 0.99, confirmed: true },
    ],
    variables: {},
    createdAt: '2026-06-15T14:20:00Z',
    updatedAt: '2026-06-15T14:20:00Z',
    taskCount: 1,
  },
  {
    id: 'script_003',
    name: '微信消息汇总',
    description: '打开企业微信，读取昨日未读消息并汇总',
    steps: [
      { id: 1, action: 'click', target: '企业微信图标', targetOCR: '企业微信', confidence: 0.90, confirmed: true },
      { id: 2, action: 'click', target: '消息列表', targetOCR: '消息', confidence: 0.88, confirmed: true },
      { id: 3, action: 'click', target: '未读消息', targetOCR: '未读', confidence: 0.85, confirmed: true },
      { id: 4, action: 'verify', target: '消息内容', targetOCR: '昨日', confidence: 0.82, confirmed: true },
      { id: 5, action: 'click', target: '导出', targetOCR: '导出', confidence: 0.91, confirmed: true },
      { id: 6, action: 'verify', target: '导出成功', targetOCR: '导出完成', confidence: 0.87, confirmed: true },
    ],
    variables: {},
    createdAt: '2026-06-10T09:00:00Z',
    updatedAt: '2026-06-12T16:00:00Z',
    taskCount: 2,
  },
]
