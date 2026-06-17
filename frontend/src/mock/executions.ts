import type { Execution } from '@/types'

export const mockExecutions: Execution[] = [
  {
    id: 'exec_001',
    taskId: 'task_a3f8c2',
    taskName: 'task_a3f8c2',
    status: 'success',
    startedAt: '2026-06-17T10:30:00Z',
    finishedAt: '2026-06-17T10:33:22Z',
    duration: '3m 22s',
    stepResults: [
      { stepId: 1, action: 'click', target: '新建', success: true, duration: 0.5 },
      { stepId: 2, action: 'input', target: '姓名', success: true, duration: 1.2 },
      { stepId: 3, action: 'input', target: '电话', success: true, duration: 0.9 },
      { stepId: 4, action: 'click', target: '保存', success: true, duration: 0.8 },
      { stepId: 5, action: 'verify', target: '保存成功', success: true, duration: 1.1 },
    ],
  },
  {
    id: 'exec_002',
    taskId: 'task_d1e5f9',
    taskName: 'task_d1e5f9',
    status: 'failed',
    startedAt: '2026-06-16T14:20:00Z',
    finishedAt: '2026-06-16T14:21:05Z',
    duration: '1m 05s',
    stepResults: [
      { stepId: 1, action: 'click', target: 'Install', success: true, duration: 0.5 },
      { stepId: 2, action: 'click', target: 'Next', success: true, duration: 0.3 },
      { stepId: 3, action: 'click', target: 'Add to PATH', success: true, duration: 0.4 },
      { stepId: 4, action: 'click', target: 'Install Now', success: false, duration: 30.0, error: '操作超时' },
    ],
  },
  {
    id: 'exec_003',
    taskId: 'task_b7c2d1',
    taskName: 'task_b7c2d1',
    status: 'success',
    startedAt: '2026-06-16T09:00:00Z',
    finishedAt: '2026-06-16T09:02:10Z',
    duration: '2m 10s',
    stepResults: [
      { stepId: 1, action: 'click', target: '企业微信', success: true, duration: 0.8 },
      { stepId: 2, action: 'click', target: '消息', success: true, duration: 0.5 },
      { stepId: 3, action: 'click', target: '未读', success: true, duration: 0.6 },
      { stepId: 4, action: 'verify', target: '昨日', success: true, duration: 1.0 },
      { stepId: 5, action: 'click', target: '导出', success: true, duration: 0.4 },
      { stepId: 6, action: 'verify', target: '导出完成', success: true, duration: 0.9 },
    ],
  },
]
