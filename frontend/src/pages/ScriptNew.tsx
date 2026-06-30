import { useState, useEffect } from 'react'
import { useParams } from 'react-router-dom'
import { Breadcrumb } from '@/components/layout/Breadcrumb'
import { Button } from '@/components/ui/button'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group'
import { Label } from '@/components/ui/label'
import { Save, Send } from 'lucide-react'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import { mockChatMessages, mockScriptMarkdown } from '@/mock/chat'
import { mockScripts } from '@/mock/scripts'
import { mockWindows } from '@/mock/tasks'
import { toast } from 'sonner'
import type { ChatMessage, Script } from '@/types'

export default function ScriptNew() {
  const { id } = useParams<{ id?: string }>()
  const isEdit = Boolean(id)

  const [scriptName, setScriptName] = useState('')
  const [messages, setMessages] = useState<ChatMessage[]>([])
  const [input, setInput] = useState('')
  const [scriptContent, setScriptContent] = useState('')

  // Load script data if editing
  useEffect(() => {
    if (id) {
      const script = mockScripts.find((s) => s.id === id)
      if (script) {
        setScriptName(script.name)
        setScriptContent(buildMarkdown(script))
        setMessages([
          {
            id: 'init_1',
            role: 'assistant',
            content: `正在编辑脚本「${script.name}」。你可以在左侧对话框调整步骤，右侧预览会实时更新。`,
            timestamp: new Date().toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' }),
          },
        ])
      }
    } else {
      // Create mode: show the demo chat
      setMessages(mockChatMessages)
      setScriptContent(mockScriptMarkdown)
    }
  }, [id])

  const handleSend = () => {
    if (!input.trim()) return
    const userMsg: ChatMessage = {
      id: `msg_${Date.now()}`,
      role: 'user',
      content: input,
      timestamp: new Date().toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' }),
    }
    setMessages((prev) => [...prev, userMsg])
    setInput('')

    // Simulate AI response
    setTimeout(() => {
      const aiMsg: ChatMessage = {
        id: `msg_${Date.now() + 1}`,
        role: 'assistant',
        content: '收到，我正在更新脚本方案。右侧预览已同步更新，脚本文件已自动保存。',
        timestamp: new Date().toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' }),
      }
      setMessages((prev) => [...prev, aiMsg])
    }, 800)
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleSend()
    }
  }

  const handleSave = () => {
    toast.success(isEdit ? '脚本已更新' : '脚本已创建')
  }

  return (
    <div className="flex h-full flex-col">
      {/* Header */}
      <div className="flex items-center justify-between px-6 pt-4">
        <Breadcrumb
          items={[
            { label: '脚本库', path: '/scripts' },
            { label: isEdit ? `编辑: ${scriptName}` : '创建脚本' },
          ]}
        />
        <Button size="sm" onClick={handleSave}>
          <Save className="h-4 w-4" />
          {isEdit ? '更新脚本' : '保存脚本'}
        </Button>
      </div>

      <div className="flex flex-1 gap-0 overflow-hidden mt-4">
        {/* Left: Chat Area */}
        <div className="flex w-[60%] flex-col border-r">
          <ScrollArea className="flex-1 px-6">
            <div className="space-y-4 pb-4">
              {messages.map((msg) => (
                <div
                  key={msg.id}
                  className={`flex gap-3 ${msg.role === 'user' ? 'flex-row-reverse' : ''}`}
                >
                  <div
                    className={`flex h-8 w-8 shrink-0 items-center justify-center rounded-full text-xs font-medium ${
                      msg.role === 'assistant'
                        ? 'bg-primary text-primary-foreground'
                        : 'bg-muted'
                    }`}
                  >
                    {msg.role === 'assistant' ? '🤖' : '👤'}
                  </div>
                  <div
                    className={`max-w-[75%] rounded-lg px-3 py-2 text-sm ${
                      msg.role === 'assistant'
                        ? 'bg-muted'
                        : 'bg-primary text-primary-foreground'
                    }`}
                  >
                    {msg.content}
                  </div>
                </div>
              ))}
            </div>
          </ScrollArea>

          {/* Input Area */}
          <div className="border-t px-6 py-3 space-y-3">
            <div className="flex gap-2">
              <textarea
                value={input}
                onChange={(e) => setInput(e.target.value)}
                onKeyDown={handleKeyDown}
                placeholder="描述需要修改的内容..."
                className="flex-1 resize-none rounded-md border bg-background px-3 py-2 text-sm min-h-[40px] max-h-[100px] focus:outline-none focus:ring-2 focus:ring-ring"
                rows={1}
              />
              <button
                onClick={handleSend}
                className="flex h-9 w-9 items-center justify-center rounded-md bg-primary text-primary-foreground hover:bg-primary/90 transition-colors"
              >
                <Send className="h-4 w-4" />
              </button>
            </div>

            <div className="flex flex-wrap items-center gap-4 text-sm">
              <div className="flex items-center gap-2">
                <Select>
                  <SelectTrigger className="h-8 w-44">
                    <SelectValue placeholder="绑定窗口..." />
                  </SelectTrigger>
                  <SelectContent>
                    {mockWindows.map((w) => (
                      <SelectItem key={w.handle} value={w.handle}>
                        {w.title}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              <div className="flex items-center gap-2">
                <RadioGroup defaultValue="new" className="flex items-center gap-3">
                  <div className="flex items-center gap-1">
                    <RadioGroupItem value="new" id="new" className="h-3 w-3" />
                    <Label htmlFor="new" className="text-sm cursor-pointer">新建</Label>
                  </div>
                  <div className="flex items-center gap-1">
                    <RadioGroupItem value="existing" id="existing" className="h-3 w-3" />
                    <Label htmlFor="existing" className="text-sm cursor-pointer">已有</Label>
                  </div>
                </RadioGroup>
              </div>

              <div className="flex items-center gap-2">
                <Select defaultValue="gpt4o">
                  <SelectTrigger className="h-8 w-32">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="gpt4o">GPT-4o</SelectItem>
                    <SelectItem value="claude">Claude 3.5</SelectItem>
                    <SelectItem value="ollama">Ollama</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>
          </div>
        </div>

        {/* Right: Markdown Preview */}
        <div className="flex w-[40%] flex-col">
          <div className="flex items-center gap-2 border-b px-6 py-2 shrink-0">
            <span className="text-sm font-medium">脚本预览</span>
            {isEdit && <span className="text-xs text-muted-foreground">{scriptName}</span>}
          </div>
          <ScrollArea className="flex-1 px-6 py-4">
            <div className="prose prose-sm dark:prose-invert max-w-none">
              <ReactMarkdown remarkPlugins={[remarkGfm]}>
                {scriptContent}
              </ReactMarkdown>
            </div>
          </ScrollArea>
        </div>
      </div>
    </div>
  )
}

// buildMarkdown generates a markdown preview from a Script object
function buildMarkdown(script: Script): string {
  const steps = script.steps
    .map(
      (s, i) =>
        `${i + 1}. **${s.action}** → "${s.target}"\n   - 目标OCR: \`${s.targetOCR}\`\n   - 置信度: ${Math.round(s.confidence * 100)}%${s.value ? `\n   - 值: \`${s.value}\`` : ''}`
    )
    .join('\n\n')

  return `## ${script.name}\n\n**描述：** ${script.description}\n\n### 步骤\n\n${steps}\n${
    Object.keys(script.variables).length > 0
      ? `\n### 变量\n\n| 变量名 | 默认值 |\n|--------|--------|\n${Object.entries(script.variables)
          .map(([k, v]) => `| ${k} | ${v} |`)
          .join('\n')}`
      : ''
  }`
}
