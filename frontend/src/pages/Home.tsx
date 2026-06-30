import { Link } from 'react-router-dom'
import { Card } from '@/components/ui/card'
import { Sparkles, Rocket, FileCode2, Store, BarChart3 } from 'lucide-react'

const cards = [
  { icon: Sparkles, title: '创建脚本', desc: '自然语言 AI 对话 新建脚本', path: '/scripts' },
  { icon: Rocket, title: '任务管理', desc: '绑定脚本 并发执行', path: '/tasks' },
  { icon: FileCode2, title: '脚本库', desc: '全部 / 我的 / 市场', path: '/scripts' },
  { icon: Store, title: '脚本市场', desc: '社区共享 一键导入', path: '/market' },
  { icon: BarChart3, title: '用量统计', desc: '消费图表 模型用量', path: '/stats' },
]

export default function Home() {
  return (
    <div className="flex h-full items-center justify-center p-6">
      <div className="grid grid-cols-3 gap-6 w-full max-w-2xl">
        {cards.map((card) => (
          <Link key={card.path} to={card.path}>
            <Card className="p-6 cursor-pointer shadow-sm hover:shadow-md hover:-translate-y-0.5 transition-all flex flex-col items-center text-center">
              <card.icon className="h-12 w-12 text-muted-foreground" />
              <h3 className="mt-4 text-lg font-semibold">{card.title}</h3>
              <p className="mt-1 text-sm text-muted-foreground">{card.desc}</p>
            </Card>
          </Link>
        ))}
      </div>
    </div>
  )
}
