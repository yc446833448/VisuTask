import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { AppLayout } from '@/components/layout/AppLayout'
import Home from '@/pages/Home'
import ScriptNew from '@/pages/ScriptNew'
import Scripts from '@/pages/Scripts'
import Market from '@/pages/Market'
import Stats from '@/pages/Stats'
import Tasks from '@/pages/Tasks'
import TaskNew from '@/pages/TaskNew'
import TaskMonitor from '@/pages/TaskMonitor'
import Settings from '@/pages/Settings'

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route element={<AppLayout />}>
          <Route path="/" element={<Home />} />
          <Route path="/script/new" element={<ScriptNew />} />
          <Route path="/scripts" element={<Scripts />} />
          <Route path="/market" element={<Market />} />
          <Route path="/stats" element={<Stats />} />
          <Route path="/tasks" element={<Tasks />} />
          <Route path="/task/new" element={<TaskNew />} />
          <Route path="/task/:id/monitor" element={<TaskMonitor />} />
          <Route path="/settings" element={<Settings />} />
        </Route>
      </Routes>
    </BrowserRouter>
  )
}

export default App
