import { useNavigate } from 'react-router-dom';
import { useAuthStore } from '@/stores/authStore';
import { Breadcrumb } from './Breadcrumb';
import { Bell, LogOut, Menu, User } from 'lucide-react';
import { useAppStore } from '@/stores/appStore';
import { cn } from '@/lib/utils';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';

export function Topbar() {
  const navigate = useNavigate();
  const { user, logout } = useAuthStore();
  const { sidebarCollapsed, toggleSidebar } = useAppStore();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <header className={cn(
      'fixed top-0 right-0 z-30 h-16 bg-surface-secondary/80 backdrop-blur-lg border-b border-gray-100 flex items-center justify-between px-6 transition-all duration-300',
      sidebarCollapsed ? 'left-16' : 'left-64'
    )}>
      <div className="flex items-center gap-4">
        <button
          onClick={toggleSidebar}
          className="p-2 rounded-lg text-text-secondary hover:text-text-primary hover:bg-surface-elevated transition-colors"
        >
          <Menu className="h-5 w-5" />
        </button>
        <Breadcrumb />
      </div>

      <div className="flex items-center gap-3">
        {/* Notifications */}
        <button className="relative p-2 rounded-lg text-text-secondary hover:text-text-primary hover:bg-surface-elevated transition-colors">
          <Bell className="h-5 w-5" />
          <span className="absolute top-1.5 right-1.5 w-2 h-2 rounded-full bg-red-500" />
        </button>

        {/* User dropdown */}
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <button className="flex items-center gap-2.5 pl-2 pr-1 py-1 rounded-lg hover:bg-surface-elevated transition-colors">
              <Avatar className="h-8 w-8">
                <AvatarFallback className="bg-gold-dim text-gold text-sm font-medium">
                  <User className="h-4 w-4" />
                </AvatarFallback>
              </Avatar>
              <div className="hidden sm:block text-left">
                <p className="text-sm font-medium text-text-primary leading-tight">{user?.nickname || '管理员'}</p>
                <p className="text-[11px] text-text-muted leading-tight">店长</p>
              </div>
            </button>
          </DropdownMenuTrigger>
          <DropdownMenuContent
            align="end"
            className="w-48 bg-surface-elevated border-gray-200"
          >
            <DropdownMenuItem
              onClick={handleLogout}
              className="text-text-secondary hover:text-red-400 hover:bg-red-500/10 cursor-pointer focus:bg-red-500/10 focus:text-red-400"
            >
              <LogOut className="h-4 w-4 mr-2" />
              退出登录
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </header>
  );
}
