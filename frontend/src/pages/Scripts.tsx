import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Breadcrumb } from '@/components/layout/Breadcrumb'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Button } from '@/components/ui/button'
import { Plus, Play, Pencil, Copy, Trash2, FileText } from 'lucide-react'
import { mockScripts } from '@/mock/scripts'
import { toast } from 'sonner'
import type { Script } from '@/types'

export default function Scripts() {
  const navigate = useNavigate()
  const [tab, setTab] = useState('all')
  const [scripts, setScripts] = useState<Script[]>(mockScripts)

  const handleTabChange = (value: string) => {
    if (value === 'market') {
      navigate('/market')
      return
    }
    setTab(value)
  }

  const filteredScripts = tab === 'mine' ? scripts.filter((s) => s.taskCount && s.taskCount > 0) : scripts

  const handleCreate = () => navigate('/script/new')
  const handleEdit = (id: string) => navigate(`/script/${id}/edit`)
  const handleRun = (id: string) => navigate(`/task/new?scriptId=${id}`)

  const handleCopy = (script: Script) => {
    const copied = { ...script, id: `script_${Date.now()}`, name: `${script.name} (副本)` }
    setScripts([copied, ...scripts])
    toast.success('脚本已复制')
  }

  const handleDelete = (id: string) => {
    setScripts((prev) => prev.filter((s) => s.id !== id))
    toast.success('脚本已删除')
  }

  return (
    <div className="flex h-full flex-col">
      {/* Header */}
      <div className="flex items-center justify-between px-6 pt-4">
        <Breadcrumb items={[{ label: '脚本库' }]} />
        <Button onClick={handleCreate}>
          <Plus className="h-4 w-4" />
          新建脚本
        </Button>
      </div>

      {/* Tabs */}
      <div className="px-6 pt-4">
        <Tabs value={tab} onValueChange={handleTabChange}>
          <TabsList>
            <TabsTrigger value="all">全部</TabsTrigger>
            <TabsTrigger value="mine">我的</TabsTrigger>
            <TabsTrigger value="market">市场</TabsTrigger>
          </TabsList>
        </Tabs>
      </div>

      {/* Script List */}
      <div className="flex-1 overflow-auto px-6 py-4">
        {filteredScripts.length === 0 ? (
          <div className="flex h-40 flex-col items-center justify-center gap-3 text-muted-foreground">
            <FileText className="h-10 w-10" />
            <p className="text-sm">暂无脚本</p>
            <Button variant="outline" size="sm" onClick={handleCreate}>
              <Plus className="h-3.5 w-3.5" />
              创建第一个脚本
            </Button>
          </div>
        ) : (
          <div className="space-y-2">
            {filteredScripts.map((script) => (
              <div
                key={script.id}
                className="flex items-center justify-between rounded-lg border p-3 hover:bg-muted/50 transition-colors cursor-pointer"
                onClick={() => handleEdit(script.id)}
              >
                <div className="flex items-center gap-3 flex-1 min-w-0">
                  <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-md bg-muted">
                    <FileText className="h-4 w-4 text-muted-foreground" />
                  </div>
                  <div className="min-w-0">
                    <p className="text-sm font-medium truncate">{script.name}</p>
                    <p className="text-xs text-muted-foreground truncate">{script.description}</p>
                    <div className="mt-1 flex items-center gap-3 text-xs text-muted-foreground">
                      <span>{script.steps.length} 个步骤</span>
                      <span>{new Date(script.createdAt).toLocaleDateString('zh-CN')}</span>
                      {script.taskCount ? <span>已关联 {script.taskCount} 个任务</span> : null}
                    </div>
                  </div>
                </div>
                {/* Actions */}
                <div className="flex shrink-0 items-center gap-0.5 ml-2" onClick={(e) => e.stopPropagation()}>
                  <Button variant="ghost" size="icon" className="h-8 w-8" onClick={() => handleRun(script.id)} title="执行">
                    <Play className="h-3.5 w-3.5" />
                  </Button>
                  <Button variant="ghost" size="icon" className="h-8 w-8" onClick={() => handleEdit(script.id)} title="编辑">
                    <Pencil className="h-3.5 w-3.5" />
                  </Button>
                  <Button variant="ghost" size="icon" className="h-8 w-8" onClick={() => handleCopy(script)} title="复制">
                    <Copy className="h-3.5 w-3.5" />
                  </Button>
                  <Button variant="ghost" size="icon" className="h-8 w-8" onClick={() => handleDelete(script.id)} title="删除">
                    <Trash2 className="h-3.5 w-3.5 text-destructive" />
                  </Button>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
