import { useState } from 'react'
import { Breadcrumb } from '@/components/layout/Breadcrumb'
import { Card } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs'
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { mockUser, mockWalletRecords } from '@/mock/user'
import { useAppStore } from '@/stores/appStore'
import { toast } from 'sonner'

export default function Settings() {
  const [tab, setTab] = useState('account')
  const { theme, setTheme } = useAppStore()

  return (
    <div className="p-6 space-y-6">
      <Breadcrumb items={[{ label: '设置' }]} />

      <Tabs value={tab} onValueChange={setTab}>
        <TabsList>
          <TabsTrigger value="account">账户</TabsTrigger>
          <TabsTrigger value="wallet">钱包</TabsTrigger>
          <TabsTrigger value="hotkey">快捷键</TabsTrigger>
          <TabsTrigger value="appearance">外观</TabsTrigger>
          <TabsTrigger value="about">关于</TabsTrigger>
        </TabsList>

        <TabsContent value="account" className="mt-4">
          <Card className="p-4 max-w-md space-y-3">
            <div className="text-sm">
              <span className="text-muted-foreground">用户等级：</span>
              {mockUser.vipLevel > 0 ? <Badge variant="secondary">VIP{mockUser.vipLevel}</Badge> : <span>普通用户</span>}
            </div>
            <div className="text-sm">
              <span className="text-muted-foreground">并发上限：</span>{mockUser.maxConcurrent} 个任务
            </div>
            <div className="rounded-md bg-muted p-3 text-xs space-y-1">
              <div>VIP 1 → 5   VIP 2 → 6   VIP 3 → 7</div>
              <div>VIP 4 → 8   VIP 5 → 10</div>
            </div>
            <Button variant="outline" size="sm" onClick={() => toast('请联系客服升级')}>升级 VIP</Button>
          </Card>
        </TabsContent>

        <TabsContent value="wallet" className="mt-4">
          <Card className="p-4 max-w-md space-y-3">
            <div className="text-sm">
              <span className="text-muted-foreground">余额：</span>
              <span className="text-lg font-semibold">¥ {mockUser.balance.toFixed(2)}</span>
            </div>
            <div className="space-y-2">
              <div className="text-xs font-medium text-muted-foreground">消费记录</div>
              {mockWalletRecords.map((r) => (
                <div key={r.id} className="flex items-center justify-between text-xs">
                  <span>{r.description}</span>
                  <div className="flex items-center gap-3">
                    <span className={r.amount > 0 ? 'text-green-500' : 'text-foreground'}>
                      {r.amount > 0 ? '+' : ''}{r.amount.toFixed(2)}
                    </span>
                    <span className="text-muted-foreground">{r.date}</span>
                  </div>
                </div>
              ))}
            </div>
            <div className="flex gap-2">
              <Button size="sm" onClick={() => toast('充值功能开发中')}>充值</Button>
              <Button variant="outline" size="sm">查看全部记录</Button>
            </div>
          </Card>
        </TabsContent>

        <TabsContent value="hotkey" className="mt-4">
          <Card className="p-4 max-w-md space-y-3">
            {[
              { label: '启动任务', key: 'Ctrl+Shift+E' },
              { label: '暂停全部', key: 'Ctrl+Shift+P' },
              { label: '停止全部', key: 'Ctrl+Shift+S' },
            ].map((item) => (
              <div key={item.label} className="flex items-center justify-between">
                <span className="text-sm">{item.label}</span>
                <div className="flex items-center gap-2">
                  <kbd className="rounded-md border bg-muted px-2 py-1 text-xs font-mono">{item.key}</kbd>
                  <Button variant="ghost" size="sm" className="h-7 text-xs">修改</Button>
                </div>
              </div>
            ))}
          </Card>
        </TabsContent>

        <TabsContent value="appearance" className="mt-4">
          <Card className="p-4 max-w-md space-y-3">
            <div className="text-sm font-medium">主题</div>
            <RadioGroup value={theme} onValueChange={(v) => setTheme(v as 'light' | 'dark' | 'system')}>
              <div className="flex items-center gap-2">
                <RadioGroupItem value="light" id="light" />
                <Label htmlFor="light" className="text-sm cursor-pointer">浅色</Label>
              </div>
              <div className="flex items-center gap-2">
                <RadioGroupItem value="dark" id="dark" />
                <Label htmlFor="dark" className="text-sm cursor-pointer">深色</Label>
              </div>
              <div className="flex items-center gap-2">
                <RadioGroupItem value="system" id="system" />
                <Label htmlFor="system" className="text-sm cursor-pointer">跟随系统</Label>
              </div>
            </RadioGroup>
          </Card>
        </TabsContent>

        <TabsContent value="about" className="mt-4">
          <Card className="p-4 max-w-md space-y-3">
            <div className="text-sm font-semibold">VisuTask v0.1.0</div>
            <div className="text-xs text-muted-foreground">已是最新版本</div>
            <Button variant="outline" size="sm" onClick={() => toast('正在检查更新...')}>检查更新</Button>
            <div className="flex gap-3 text-xs text-muted-foreground pt-2">
              <span className="hover:text-foreground cursor-pointer transition-colors">用户协议</span>
              <span className="hover:text-foreground cursor-pointer transition-colors">隐私政策</span>
              <span className="hover:text-foreground cursor-pointer transition-colors">开源许可</span>
            </div>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}
