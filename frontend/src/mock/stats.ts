import type { StatsData } from '@/types'

export const mockStats: StatsData = {
  balance: 128.50,
  monthSpending: 45.80,
  lastMonthChange: -12,
  dailySpending: [
    { date: '6/11', value: 3.2 },
    { date: '6/12', value: 5.8 },
    { date: '6/13', value: 8.4 },
    { date: '6/14', value: 4.2 },
    { date: '6/15', value: 12.0 },
    { date: '6/16', value: 7.6 },
    { date: '6/17', value: 4.6 },
  ],
  llmUsage: [
    { label: 'GPT-4o', value: 3240 },
    { label: 'Claude 3.5', value: 1856 },
    { label: 'Ollama', value: 892 },
  ],
  ocrUsage: [
    { label: '本月总调用', value: 12480 },
    { label: '中文识别', value: 8320 },
    { label: '英文识别', value: 4160 },
  ],
}
