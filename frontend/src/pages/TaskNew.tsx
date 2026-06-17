import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Breadcrumb } from '@/components/layout/Breadcrumb'
import { Button } from '@/components/ui/button'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group'
import { Label } from '@/components/ui/label'
import { RefreshCw } from 'lucide-react'
import { mockScripts } from '@/mock/scripts'
import { mockWindows } from '@/mock/tasks'
import { toast } from 'sonner'

export default function TaskNew() {
  const navigate = useNavigate()
  const [triggerType, setTriggerType] = useState('manual')

  const handleCreate = () => {
    toast.success('任务已创建')
    navigate('/tasks')
  }

  return (
    <div className="p-6 space-y-6 max-w-lg">
      <Breadcrumb items={[{ label: '任务管理', path: '/tasks' }, { label: '新建任务' }]} />

      <div className="space-y-4">
        <div className="space-y-2">
          <Label className="text-sm font-medium">选择脚本</Label>
          <Select>
            <SelectTrigger className="h-9">
              <SelectValue placeholder="选择脚本..." />
            </SelectTrigger>
            <SelectContent>
              {mockScripts.map((s) => (
                <SelectItem key={s.id} value={s.id}>{s.name}</SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        <div className="space-y-2">
          <Label className="text-sm font-medium">绑定窗口</Label>
          <Select>
            <SelectTrigger className="h-9">
              <SelectValue placeholder="选择目标窗口..." />
            </SelectTrigger>
            <SelectContent>
              {mockWindows.map((w) => (
                <SelectItem key={w.handle} value={w.handle}>
                  {w.process} - {w.title}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          <Button variant="ghost" size="sm" className="h-7 text-xs text-muted-foreground">
            <RefreshCw className="mr-1 h-3 w-3" /> 刷新窗口列表
          </Button>
        </div>

        <div className="space-y-2">
          <Label className="text-sm font-medium">触发方式</Label>
          <RadioGroup value={triggerType} onValueChange={setTriggerType} className="space-y-2">
            <div className="flex items-center gap-2">
              <RadioGroupItem value="manual" id="manual" />
              <Label htmlFor="manual" className="text-sm cursor-pointer">手动执行</Label>
            </div>
            <div className="flex items-center gap-2">
              <RadioGroupItem value="cron" id="cron" />
              <Label htmlFor="cron" className="text-sm cursor-pointer">定时执行</Label>
              {triggerType === 'cron' && (
                <div className="flex items-center gap-2 ml-2">
                  <Select defaultValue="daily">
                    <SelectTrigger className="h-8 w-20"><SelectValue /></SelectTrigger>
                    <SelectContent>
                      <SelectItem value="daily">每天</SelectItem>
                      <SelectItem value="weekly">每周</SelectItem>
                    </SelectContent>
                  </Select>
                  <input type="time" className="h-8 rounded-md border bg-background px-2 text-sm" defaultValue="09:00" />
                </div>
              )}
            </div>
            <div className="flex items-center gap-2">
              <RadioGroupItem value="hotkey" id="hotkey" />
              <Label htmlFor="hotkey" className="text-sm cursor-pointer">快捷键</Label>
              {triggerType === 'hotkey' && (
                <input className="h-8 rounded-md border bg-background px-2 text-sm ml-2 w-40" placeholder="Ctrl+Shift+__" />
              )}
            </div>
          </RadioGroup>
        </div>
      </div>

      <div className="flex justify-end">
        <Button onClick={handleCreate}>加入任务</Button>
      </div>
    </div>
  )
}
