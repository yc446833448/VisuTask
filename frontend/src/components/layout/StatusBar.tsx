import { useAppStore } from '@/stores/appStore'

export function StatusBar() {
  const { agentStatus, ocrConnected, llmModel, llmAvailable, runningCount, maxConcurrent } = useAppStore()

  return (
    <div className="flex h-7 items-center gap-4 border-t px-4 text-xs text-muted-foreground">
      <span className="flex items-center gap-1">
        <span className={`inline-block h-2 w-2 rounded-full ${
          agentStatus === 'idle' ? 'bg-green-500' :
          agentStatus === 'running' ? 'bg-blue-500 animate-pulse' :
          'bg-red-500'
        }`} />
        Agent {agentStatus === 'idle' ? '就绪' : agentStatus === 'running' ? '执行中' : '错误'}
      </span>

      <span className="flex items-center gap-1">
        <span className={`inline-block h-2 w-2 rounded-full ${ocrConnected ? 'bg-green-500' : 'bg-zinc-400'}`} />
        OCR {ocrConnected ? '已连接' : '断开'}
      </span>

      <span className="flex items-center gap-1">
        <span className={`inline-block h-2 w-2 rounded-full ${llmAvailable ? 'bg-green-500' : 'bg-red-500'}`} />
        LLM: {llmModel}
      </span>

      <span className="ml-auto">
        {runningCount}/{maxConcurrent}
      </span>
    </div>
  )
}
