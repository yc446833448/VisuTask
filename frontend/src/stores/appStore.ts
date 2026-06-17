import { create } from 'zustand'

interface AppState {
  theme: 'light' | 'dark' | 'system'
  agentStatus: 'idle' | 'running' | 'error'
  ocrConnected: boolean
  llmModel: string
  llmAvailable: boolean
  runningCount: number
  maxConcurrent: number
  setTheme: (theme: 'light' | 'dark' | 'system') => void
  toggleTheme: () => void
}

export const useAppStore = create<AppState>((set) => ({
  theme: 'light',
  agentStatus: 'idle',
  ocrConnected: true,
  llmModel: 'GPT-4o',
  llmAvailable: true,
  runningCount: 1,
  maxConcurrent: 5,
  setTheme: (theme) => {
    set({ theme })
    if (theme === 'dark') {
      document.documentElement.classList.add('dark')
    } else if (theme === 'light') {
      document.documentElement.classList.remove('dark')
    } else {
      if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
        document.documentElement.classList.add('dark')
      } else {
        document.documentElement.classList.remove('dark')
      }
    }
  },
  toggleTheme: () => {
    set((state) => {
      const next = state.theme === 'light' ? 'dark' : 'light'
      if (next === 'dark') {
        document.documentElement.classList.add('dark')
      } else {
        document.documentElement.classList.remove('dark')
      }
      return { theme: next }
    })
  },
}))
