import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Breadcrumb } from '@/components/layout/Breadcrumb'
import { Button } from '@/components/ui/button'
import { Progress } from '@/components/ui/progress'
import { Badge } from '@/components/ui/badge'
import { Plus, Play, Pause, Square, Pencil, Trash2, ChevronRight, ChevronDown, Satellite } from 'lucide-react'
import { mockTasks } from '@/mock/tasks'
import { useAppStore } from '@/stores/appStore'
import type { Task } from '@/types'

export default function Tasks() {
  const navigate = useNavigate()
  const [expandedId, setExpandedId] = useState<string | null>(null)
  const { runningCount, maxConcurrent } = useAppStore()

  return (
    <div className="p-6 space-y-6">
      <div className="flex items-center justify-between">
        <Breadcrumb items={[{ label: '任务管理' }]} />
        <Button size="sm" onClick={() => navigate('/task/new')}>
          <Plus className="mr-1 h-4 w-4" /> 新建任务
        </Button>
      </div>

      <div className="rounded-lg border">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b bg-muted/50">
              <th className="w-8 px-2 py-3" />
              <th className="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase">任务名称</th>
              <th className="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase">绑定脚本</th>
              <th className="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase">绑定窗口</th>
              <th className="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase">状态</th>
              <th className="px-4 py-3 text-right text-xs font-medium text-muted-foreground uppercase">操作</th>
            </tr>
          </thead>
          <tbody>
            {mockTasks.map((task) => (
              <TaskRow
                key={task.id}
                task={task}
                expanded={expandedId === task.id}
                onToggle={() => setExpandedId(expandedId === task.id ? null : task.id)}
                onMonitor={() => navigate(`/task/${task.id}/monitor`)}
              />
            ))}
          </tbody>
        </table>
      </div>

      <div className="flex items-center gap-3 text-sm text-muted-foreground">
        <span>运行中 {runningCount}/{maxConcurrent}</span>
        <Progress value={(runningCount / maxConcurrent) * 100} className="h-2 flex-1 max-w-xs" />
      </div>
    </div>
  )
}

function TaskRow({ task, expanded, onToggle, onMonitor }: { task: Task; expanded: boolean; onToggle: () => void; onMonitor: () => void }) {
  const isRunning = task.status === 'running'

  return (
    <>
      <tr className="border-b hover:bg-muted/50 transition-colors cursor-pointer" onClick={onToggle}>
        <td className="px-2 py-3">
          {expanded ? <ChevronDown className="h-4 w-4 text-muted-foreground" /> : <ChevronRight className="h-4 w-4 text-muted-foreground" />}
        </td>
        <td className="px-4 py-3 font-mono text-xs">{task.name}</td>
        <td className="px-4 py-3">{task.scriptName}</td>
        <td className="px-4 py-3 text-muted-foreground">{task.windowTitle || '待绑定'}</td>
        <td className="px-4 py-3">
          {isRunning ? (
            <div className="flex items-center gap-2">
              <Badge variant="outline" className="text-blue-500 border-blue-500">
                <Play className="mr-1 h-3 w-3" /> 运行中
              </Badge>
              <span className="text-xs text-muted-foreground">{task.progress}%</span>
            </div>
          ) : (
            <Badge variant="outline" className="text-zinc-500 border-zinc-500">
              ○ 空闲
            </Badge>
          )}
        </td>
        <td className="px-4 py-3 text-right">
          <div className="flex items-center justify-end gap-1" onClick={(e) => e.stopPropagation()}>
            {isRunning ? (
              <>
                <Button variant="ghost" size="icon" className="h-8 w-8"><Pause className="h-4 w-4" /></Button>
                <Button variant="ghost" size="icon" className="h-8 w-8"><Square className="h-4 w-4" /></Button>
              </>
            ) : (
              <Button variant="ghost" size="icon" className="h-8 w-8"><Play className="h-4 w-4" /></Button>
            )}
            <Button variant="ghost" size="icon" className="h-8 w-8"><Pencil className="h-4 w-4" /></Button>
            <Button variant="ghost" size="icon" className="h-8 w-8"><Trash2 className="h-4 w-4 text-destructive" /></Button>
          </div>
        </td>
      </tr>
      {expanded && (
        <tr className="bg-muted/30">
          <td colSpan={6} className="px-8 py-4">
            <div className="space-y-4 text-sm">
              <div className="grid grid-cols-2 gap-2 text-xs">
                <div><span className="text-muted-foreground">脚本：</span>{task.scriptName} ({task.totalSteps ?? task.scriptName ? '5' : '?'}步)</div>
                <div><span className="text-muted-foreground">窗口：</span>{task.windowTitle || '待绑定'} {task.windowHandle && `(${task.windowHandle})`}</div>
                <div><span className="text-muted-foreground">触发：</span>{task.trigger.type === 'manual' ? '手动' : task.trigger.type === 'cron' ? `定时 ${task.trigger.cron}` : `快捷键 ${task.trigger.hotkey}`}</div>
                <div><span className="text-muted-foreground">参数：</span>{Object.entries(task.parameters).map(([k, v]) => `${k}=${v}`).join(', ') || '无'}</div>
              </div>

              {isRunning && (
                <div className="space-y-1">
                  <div className="flex items-center gap-2 text-xs text-muted-foreground">
                    <span>执行进度</span>
                    <span>{task.currentStep}/{task.totalSteps}</span>
                    <span>耗时 {task.duration}</span>
                  </div>
                  <Progress value={task.progress} className="h-2" />
                </div>
              )}

              <div className="space-y-1 text-xs">
                <div className="text-muted-foreground font-medium mb-1">步骤日志</div>
                <div className="space-y-0.5">
                  <div className="text-green-600">✓ 1. click "新建"    0.5s</div>
                  <div className="text-green-600">✓ 2. input "姓名"   1.2s</div>
                  <div className={isRunning ? 'text-blue-500' : 'text-green-600'}>
                    {isRunning ? '● 3. click "保存"   执行中...' : '✓ 3. click "保存"   0.8s'}
                  </div>
                  <div className="text-zinc-400">○ 4. verify "成功"  等待</div>
                </div>
              </div>

              <div className="space-y-1 text-xs">
                <div className="text-muted-foreground font-medium mb-1">最近执行记录</div>
                <div className="text-green-600">6/17 10:30 ✓ 成功 3m22s</div>
                <div className="text-red-500">6/16 14:20 ✗ 失败 步骤4超时</div>
              </div>

              <Button variant="outline" size="sm" className="h-8" onClick={onMonitor}>
                <Satellite className="mr-1 h-4 w-4" /> 实时监控
              </Button>
            </div>
          </td>
        </tr>
      )}
    </>
  )
}
