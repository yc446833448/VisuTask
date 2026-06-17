import { Breadcrumb } from '@/components/layout/Breadcrumb'
import { Card } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { mockStats } from '@/mock/stats'
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, ResponsiveContainer, LineChart as ReLineChart, Line, Tooltip } from 'recharts'

export default function Stats() {
  return (
    <div className="p-6 space-y-6">
      <Breadcrumb items={[{ label: '用量统计' }]} />

      {/* Overview Cards */}
      <div className="grid grid-cols-2 gap-6 max-w-lg">
        <Card className="p-4">
          <div className="text-xs text-muted-foreground">💰 钱包余额</div>
          <div className="mt-2 text-2xl font-semibold">¥ {mockStats.balance.toFixed(2)}</div>
          <Button variant="outline" size="sm" className="mt-3 h-8">充值</Button>
        </Card>
        <Card className="p-4">
          <div className="text-xs text-muted-foreground">📊 本月消费</div>
          <div className="mt-2 text-2xl font-semibold">¥ {mockStats.monthSpending.toFixed(2)}</div>
          <div className="mt-1 text-xs text-muted-foreground">
            较上月 <span className={mockStats.lastMonthChange < 0 ? 'text-green-500' : 'text-red-500'}>
              {mockStats.lastMonthChange > 0 ? '+' : ''}{mockStats.lastMonthChange}%
            </span>
          </div>
        </Card>
      </div>

      {/* Spending Chart */}
      <div className="space-y-3">
        <div className="flex items-center justify-between">
          <h3 className="text-sm font-semibold">消费金额统计</h3>
          <Select defaultValue="7d">
            <SelectTrigger className="h-8 w-24"><SelectValue /></SelectTrigger>
            <SelectContent>
              <SelectItem value="7d">7天</SelectItem>
              <SelectItem value="30d">30天</SelectItem>
              <SelectItem value="90d">90天</SelectItem>
            </SelectContent>
          </Select>
        </div>
        <Card className="p-4">
          <ResponsiveContainer width="100%" height={200}>
            <ReLineChart data={mockStats.dailySpending}>
              <CartesianGrid strokeDasharray="3 3" className="stroke-border" />
              <XAxis dataKey="date" tick={{ fontSize: 12 }} className="text-muted-foreground" />
              <YAxis tick={{ fontSize: 12 }} className="text-muted-foreground" />
              <Tooltip />
              <Line type="monotone" dataKey="value" stroke="#3b82f6" strokeWidth={2} dot={{ r: 3 }} />
            </ReLineChart>
          </ResponsiveContainer>
        </Card>
      </div>

      {/* LLM Usage */}
      <div className="space-y-3">
        <h3 className="text-sm font-semibold">LLM 模型用量</h3>
        <Card className="p-4">
          <ResponsiveContainer width="100%" height={150}>
            <BarChart data={mockStats.llmUsage} layout="vertical">
              <CartesianGrid strokeDasharray="3 3" className="stroke-border" />
              <XAxis type="number" tick={{ fontSize: 12 }} />
              <YAxis type="category" dataKey="label" tick={{ fontSize: 12 }} width={80} />
              <Tooltip />
              <Bar dataKey="value" fill="#3b82f6" radius={[0, 4, 4, 0]} />
            </BarChart>
          </ResponsiveContainer>
        </Card>
      </div>

      {/* OCR Usage */}
      <div className="space-y-3">
        <h3 className="text-sm font-semibold">OCR 识别用量</h3>
        <Card className="p-4">
          <ResponsiveContainer width="100%" height={150}>
            <BarChart data={mockStats.ocrUsage} layout="vertical">
              <CartesianGrid strokeDasharray="3 3" className="stroke-border" />
              <XAxis type="number" tick={{ fontSize: 12 }} />
              <YAxis type="category" dataKey="label" tick={{ fontSize: 12 }} width={80} />
              <Tooltip />
              <Bar dataKey="value" fill="#22c55e" radius={[0, 4, 4, 0]} />
            </BarChart>
          </ResponsiveContainer>
        </Card>
      </div>
    </div>
  )
}
