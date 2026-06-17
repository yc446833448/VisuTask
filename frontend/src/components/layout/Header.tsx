import { Link } from 'react-router-dom'
import { Settings } from 'lucide-react'
import { Badge } from '@/components/ui/badge'
import { useAppStore } from '@/stores/appStore'
import { mockUser } from '@/mock/user'

export function Header() {
  const user = mockUser

  return (
    <header className="flex h-14 items-center justify-between border-b px-6">
      <Link to="/" className="flex items-center gap-2">
        <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-primary text-primary-foreground text-sm font-bold">
          V
        </div>
        <span className="text-lg font-semibold">VisuTask</span>
      </Link>

      <div className="flex items-center gap-3">
        <Link to="/settings">
          <Settings className="h-5 w-5 text-muted-foreground hover:text-foreground transition-colors" />
        </Link>

        <div className="h-8 w-8 rounded-full bg-muted flex items-center justify-center text-sm font-medium">
          {user.avatar ? (
            <img src={user.avatar} alt="avatar" className="h-8 w-8 rounded-full" />
          ) : (
            <span>U</span>
          )}
        </div>

        {user.vipLevel > 0 && (
          <Badge variant="secondary" className="text-xs">
            VIP{user.vipLevel}
          </Badge>
        )}
      </div>
    </header>
  )
}
