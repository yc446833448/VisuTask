import { useState } from 'react'
import { Breadcrumb } from '@/components/layout/Breadcrumb'
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Download, Lock } from 'lucide-react'
import { mockFreeScripts, mockVipScripts } from '@/mock/market'
import { mockUser } from '@/mock/user'
import { toast } from 'sonner'

export default function Market() {
  const [tab, setTab] = useState('free')

  return (
    <div className="p-6 space-y-6">
      <Breadcrumb items={[{ label: '脚本市场' }]} />

      <Tabs value={tab} onValueChange={setTab}>
        <TabsList>
          <TabsTrigger value="free">免费脚本</TabsTrigger>
          <TabsTrigger value="vip">VIP 脚本</TabsTrigger>
          <TabsTrigger value="publish">发布脚本</TabsTrigger>
        </TabsList>

        <TabsContent value="free" className="mt-4">
          <ScriptTable scripts={mockFreeScripts} />
        </TabsContent>

        <TabsContent value="vip" className="mt-4">
          <VipScriptTable scripts={mockVipScripts} userVipLevel={mockUser.vipLevel} />
        </TabsContent>

        <TabsContent value="publish" className="mt-4">
          <PublishForm />
        </TabsContent>
      </Tabs>
    </div>
  )
}

function ScriptTable({ scripts }: { scripts: typeof mockFreeScripts }) {
  return (
    <div className="rounded-lg border">
      <table className="w-full text-sm">
        <thead>
          <tr className="border-b bg-muted/50">
            <th className="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase">名称</th>
            <th className="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase">脚本描述</th>
            <th className="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase w-32">版本号</th>
            <th className="px-4 py-3 text-right text-xs font-medium text-muted-foreground uppercase w-28">操作</th>
          </tr>
        </thead>
        <tbody>
          {scripts.map((script) => (
            <tr key={script.id} className="border-b last:border-0 hover:bg-muted/50 transition-colors">
              <td className="px-4 py-3 font-medium">{script.name}</td>
              <td className="px-4 py-3 text-muted-foreground">{script.description}</td>
              <td className="px-4 py-3">
                <Select defaultValue={script.version}>
                  <SelectTrigger className="h-8 w-24">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {script.versions.map((v) => (
                      <SelectItem key={v} value={v}>{v}</SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </td>
              <td className="px-4 py-3 text-right">
                <Button size="sm" className="h-8" onClick={() => toast.success(`已导入: ${script.name}`)}>
                  <Download className="mr-1 h-4 w-4" /> 下载
                </Button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

function VipScriptTable({ scripts, userVipLevel }: { scripts: typeof mockVipScripts; userVipLevel: number }) {
  return (
    <div className="rounded-lg border">
      <table className="w-full text-sm">
        <thead>
          <tr className="border-b bg-muted/50">
            <th className="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase">名称</th>
            <th className="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase">脚本描述</th>
            <th className="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase w-32">版本号</th>
            <th className="px-4 py-3 text-right text-xs font-medium text-muted-foreground uppercase w-36">操作</th>
          </tr>
        </thead>
        <tbody>
          {scripts.map((script) => {
            const canDownload = userVipLevel >= (script.vipLevel ?? 0)
            return (
              <tr key={script.id} className="border-b last:border-0 hover:bg-muted/50 transition-colors">
                <td className="px-4 py-3">
                  <div className="flex items-center gap-2">
                    <span className="font-medium">{script.name}</span>
                    <Badge variant="secondary" className="text-xs">VIP{script.vipLevel}</Badge>
                  </div>
                </td>
                <td className="px-4 py-3 text-muted-foreground">{script.description}</td>
                <td className="px-4 py-3">
                  <Select defaultValue={script.version}>
                    <SelectTrigger className="h-8 w-24">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      {script.versions.map((v) => (
                        <SelectItem key={v} value={v}>{v}</SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </td>
                <td className="px-4 py-3 text-right">
                  {canDownload ? (
                    <Button size="sm" className="h-8" onClick={() => toast.success(`已导入: ${script.name}`)}>
                      <Download className="mr-1 h-4 w-4" /> VIP{script.vipLevel} 下载
                    </Button>
                  ) : (
                    <Button size="sm" variant="outline" className="h-8" onClick={() => toast('请升级 VIP 等级')}>
                      <Lock className="mr-1 h-4 w-4" /> 升级 VIP{script.vipLevel}
                    </Button>
                  )}
                </td>
              </tr>
            )
          })}
        </tbody>
      </table>
    </div>
  )
}

function PublishForm() {
  return (
    <div className="max-w-md space-y-4">
      <div className="space-y-2">
        <label className="text-sm font-medium">脚本名称</label>
        <input className="flex h-9 w-full rounded-md border bg-background px-3 py-1 text-sm" placeholder="输入脚本名称" />
      </div>
      <div className="space-y-2">
        <label className="text-sm font-medium">脚本描述</label>
        <textarea className="flex w-full rounded-md border bg-background px-3 py-2 text-sm min-h-[80px]" placeholder="描述你的脚本功能" />
      </div>
      <div className="space-y-2">
        <label className="text-sm font-medium">选择脚本</label>
        <Select>
          <SelectTrigger className="h-9">
            <SelectValue placeholder="选择要发布的脚本..." />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="1">数据录入 CRM</SelectItem>
            <SelectItem value="2">安装 Python 3.12</SelectItem>
            <SelectItem value="3">微信消息汇总</SelectItem>
          </SelectContent>
        </Select>
      </div>
      <Button onClick={() => toast.success('发布申请已提交')}>发布到市场</Button>
    </div>
  )
}
