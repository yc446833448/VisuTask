import { useState } from 'react'
import { Breadcrumb } from '@/components/layout/Breadcrumb'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Switch } from '@/components/ui/switch'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { useParams } from 'react-router-dom'

export default function TaskMonitor() {
  const { id } = useParams()
  const [showOcr, setShowOcr] = useState(true)

  return (
    <div className="flex h-full flex-col">
      <div className="px-6 pt-4">
        <Breadcrumb items={[{ label: '任务管理', path: '/tasks' }, { label: id || '' }, { label: '监控' }]} />
      </div>

      <div className="flex items-center gap-6 px-6 py-3">
        <div className="flex items-center gap-2">
          <Label htmlFor="ocr-toggle" className="text-sm">OCR 标注</Label>
          <Switch id="ocr-toggle" checked={showOcr} onCheckedChange={setShowOcr} />
        </div>
        <div className="flex items-center gap-2">
          <Label className="text-sm text-muted-foreground">采样:</Label>
          <Select defaultValue="agent">
            <SelectTrigger className="h-8 w-36">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="agent">Agent 驱动</SelectItem>
              <SelectItem value="20">定时 20次/分</SelectItem>
              <SelectItem value="10">定时 10次/分</SelectItem>
              <SelectItem value="manual">手动刷新</SelectItem>
            </SelectContent>
          </Select>
        </div>
      </div>

      {/* Screen Capture Area */}
      <div className="mx-6 relative rounded-lg border bg-muted/30 overflow-hidden" style={{ height: '320px' }}>
        <div className="flex h-full items-center justify-center text-muted-foreground text-sm">
          绑定窗口实时画面
        </div>
        {showOcr && (
          <>
            <div className="absolute border-2 border-blue-500 rounded-sm bg-blue-500/15" style={{ left: '45%', top: '55%', width: '80px', height: '30px' }}>
              <span className="absolute -top-5 left-0 text-xs bg-black/70 text-white px-1 rounded">保存 98%</span>
            </div>
            <div className="absolute border-2 border-blue-500 rounded-sm bg-blue-500/15" style={{ left: '60%', top: '55%', width: '80px', height: '30px' }}>
              <span className="absolute -top-5 left-0 text-xs bg-black/70 text-white px-1 rounded">取消 95%</span>
            </div>
          </>
        )}
      </div>

      {/* Bottom: Agent Log + OCR Results */}
      <div className="flex flex-1 gap-0 overflow-hidden mt-4">
        <div className="flex w-1/2 flex-col border-r">
          <div className="px-4 py-2 border-b text-sm font-medium text-muted-foreground">Agent 执行过程</div>
          <ScrollArea className="flex-1 px-4 py-2">
            <div className="space-y-2 text-xs font-mono">
              <div><span className="text-muted-foreground">10:32:05</span> 步骤6: click "新建" <span className="text-green-500">→ ✓ 成功 (0.5s)</span></div>
              <div><span className="text-muted-foreground">10:32:03</span> 步骤5: verify "成功" <span className="text-green-500">→ ✓ 成功 (1.1s)</span></div>
              <div><span className="text-muted-foreground">10:32:01</span> 步骤4: click "保存" <span className="text-green-500">→ ✓ 成功 (0.8s)</span></div>
              <div><span className="text-muted-foreground">10:31:58</span> 步骤3: input "电话" <span className="text-green-500">→ ✓ 成功 (0.9s)</span></div>
              <div><span className="text-muted-foreground">10:31:55</span> 步骤2: input "姓名" <span className="text-green-500">→ ✓ 成功 (1.2s)</span></div>
              <div><span className="text-muted-foreground">10:31:53</span> 步骤1: click "新建" <span className="text-green-500">→ ✓ 成功 (0.5s)</span></div>
              <div className="text-blue-500">● 当前执行中...</div>
            </div>
          </ScrollArea>
        </div>

        <div className="flex w-1/2 flex-col">
          <div className="px-4 py-2 border-b text-sm font-medium text-muted-foreground">OCR 识别结果</div>
          <ScrollArea className="flex-1 px-4 py-2">
            <div className="space-y-3 text-xs">
              <div>
                <div className="font-medium">"保存"</div>
                <div className="text-muted-foreground">x:450 y:320 w:60 h:28 conf: 0.98</div>
              </div>
              <div>
                <div className="font-medium">"取消"</div>
                <div className="text-muted-foreground">x:520 y:320 w:60 h:28 conf: 0.95</div>
              </div>
              <div>
                <div className="font-medium">"用户名"</div>
                <div className="text-muted-foreground">x:200 y:180 w:80 h:24 conf: 0.92</div>
              </div>
              <div>
                <div className="font-medium">"客户信息"</div>
                <div className="text-muted-foreground">x:100 y:50 w:100 h:24 conf: 0.88</div>
              </div>
            </div>
          </ScrollArea>
        </div>
      </div>
    </div>
  )
}
