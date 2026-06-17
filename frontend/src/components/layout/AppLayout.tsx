import { Outlet } from 'react-router-dom'
import { Header } from './Header'
import { StatusBar } from './StatusBar'
import { Toaster } from '@/components/ui/sonner'
import { TooltipProvider } from '@/components/ui/tooltip'

export function AppLayout() {
  return (
    <TooltipProvider>
      <div className="flex h-screen flex-col overflow-hidden">
        <Header />
        <main className="flex-1 overflow-auto">
          <Outlet />
        </main>
        <StatusBar />
      </div>
      <Toaster position="top-right" richColors />
    </TooltipProvider>
  )
}
