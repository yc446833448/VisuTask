export type TaskStatus = 'idle' | 'running' | 'paused' | 'completed' | 'failed'

export interface Script {
  id: string
  name: string
  description: string
  steps: Step[]
  variables: Record<string, string>
  createdAt: string
  updatedAt: string
  taskCount?: number
}

export interface Step {
  id: number
  action: string
  target: string
  targetOCR: string
  value?: string
  timeout?: number
  confidence: number
  confirmed: boolean
}

export interface Task {
  id: string
  name: string
  scriptId: string
  scriptName: string
  windowHandle: string
  windowTitle: string
  parameters: Record<string, string>
  trigger: Trigger
  status: TaskStatus
  progress?: number
  currentStep?: number
  totalSteps?: number
  duration?: string
  createdAt: string
}

export interface Trigger {
  type: 'manual' | 'cron' | 'hotkey'
  cron?: string
  hotkey?: string
}

export interface Execution {
  id: string
  taskId: string
  taskName: string
  status: 'success' | 'failed' | 'cancelled'
  startedAt: string
  finishedAt: string
  duration: string
  stepResults: StepResult[]
}

export interface StepResult {
  stepId: number
  action: string
  target: string
  success: boolean
  duration: number
  error?: string
}

export interface StepLogItem {
  index: number
  action: string
  target: string
  status: 'done' | 'running' | 'waiting' | 'failed'
  duration?: number
  error?: string
}

export interface OcrAnnotation {
  text: string
  rect: { x: number; y: number; width: number; height: number }
  confidence: number
  status: 'target' | 'confirmed' | 'failed'
}

export interface MarketScript {
  id: string
  name: string
  description: string
  version: string
  versions: string[]
  downloads: number
  author: string
  isVip: boolean
  vipLevel?: number
}

export interface ChatMessage {
  id: string
  role: 'user' | 'assistant'
  content: string
  timestamp: string
}

export interface User {
  id: string
  vipLevel: number
  maxConcurrent: number
  balance: number
  avatar?: string
}

export interface WindowInfo {
  handle: string
  title: string
  process: string
}

export interface StatsData {
  balance: number
  monthSpending: number
  lastMonthChange: number
  dailySpending: { date: string; value: number }[]
  llmUsage: { label: string; value: number }[]
  ocrUsage: { label: string; value: number }[]
}

export interface WalletRecord {
  id: string
  description: string
  amount: number
  date: string
}
