import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Breadcrumb } from '@/components/layout/Breadcrumb'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Button } from '@/components/ui/button'
import { Play, Pencil, Copy, Trash2 } from 'lucide-react'
import { mockScripts } from '@/mock/scripts'
import { toast } from 'sonner'

export default function Scripts() {
  const navigate = useNavigate()
  const [tab, setTab] = useState('all')

  const handleTabChange = (value: string) => {
    if (value === 'market') {
      navigate('/market')
      return
    }
    if (value === 'create') {
      navigate('/script/new')
      return
    }
    setTab(value)
  }

  const filteredScripts = tab === 'mine' ? mockScripts : mockScripts

  return (
    <div className="p-6 space-y-6">
      <Breadcrumb items={[{ label: '脚本库' }]} />

      <Tabs value={tab} onValueChange={handleTabChange}>
        <TabsList>
          <TabsTrigger value="all">全部</TabsTrigger>
          <TabsTrigger value="mine">我的</TabsTrigger>
          <TabsTrigger value="market">市场</TabsTrigger>
          <TabsTrigger value="create">✨ 创建新脚本</TabsTrigger>
        </TabsList>
      </Tabs>

      <div className="space-y-3">
        {filteredScripts.map((script) => (
          <div
            key={script.id}
            className="flex items-center justify-between rounded-lg border p-4 hover:bg-muted/50 transition-colors"
          >
            <div className="flex-1">
              <div className="flex items-center gap-2">
                <FileIcon />
                <span className="text-sm font-medium">{script.name}</span>
              </div>
              <div className="mt-1 flex items-center gap-4 text-xs text-muted-foreground">
                <span>{script.steps.length} 步</span>
                <span>创建: {new Date(script.createdAt).toLocaleDateString('zh-CN', { month: 'numeric', day: 'numeric' })}</span>
                <span>已创建 {script.taskCount} 个任务</span>
              </div>
            </div>
            <div className="flex items-center gap-1">
              <Button variant="ghost" size="icon" className="h-8 w-8" onClick={() => navigate('/task/new')}>
                <Play className="h-4 w-4" />
              </Button>
              <Button variant="ghost" size="icon" className="h-8 w-8">
                <Pencil className="h-4 w-4" />
              </Button>
              <Button variant="ghost" size="icon" className="h-8 w-8">
                <Copy className="h-4 w-4" />
              </Button>
              <Button variant="ghost" size="icon" className="h-8 w-8" onClick={() => toast.success('已删除')}>
                <Trash2 className="h-4 w-4 text-destructive" />
              </Button>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}

function FileIcon() {
  return (
    <div className="flex h-8 w-8 items-center justify-center rounded-md bg-muted">
      <span className="text-sm">📦</span>
    </div>
  )
}
