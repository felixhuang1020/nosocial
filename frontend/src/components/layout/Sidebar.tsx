import { useLocation, useNavigate } from 'react-router-dom';
import { useAppStore } from '@/stores/appStore';
import { cn } from '@/lib/utils';
import {
  LayoutDashboard,
  Users,
  Award,
  Wine,
  ShoppingCart,
  MessageSquare,
  Sparkles,
  Image,
  Settings,
  ChevronLeft,
  Wine as WineIcon,
} from 'lucide-react';
import type { LucideIcon } from 'lucide-react';

interface NavItem {
  path: string;
  label: string;
  icon: LucideIcon;
}

const navItems: NavItem[] = [
  { path: '/dashboard', label: '数据看板', icon: LayoutDashboard },
  { path: '/users', label: '用户管理', icon: Users },
  { path: '/shareholders', label: '股东管理', icon: Award },
  { path: '/drinks', label: '酒水管理', icon: Wine },
  { path: '/orders', label: '订单管理', icon: ShoppingCart },
  { path: '/reviews', label: '点评审核', icon: MessageSquare },
  { path: '/tarot', label: '塔罗配置', icon: Sparkles },
  { path: '/banners', label: '图片管理', icon: Image },
  { path: '/settings', label: '系统设置', icon: Settings },
];

export function Sidebar() {
  const location = useLocation();
  const navigate = useNavigate();
  const { sidebarCollapsed, toggleSidebar } = useAppStore();

  return (
    <aside
      className={cn(
        'fixed top-0 left-0 z-40 h-full bg-surface-secondary border-r border-gray-100 flex flex-col transition-all duration-300',
        sidebarCollapsed ? 'w-16' : 'w-64'
      )}
    >
      {/* Logo */}
      <div className={cn(
        'h-16 flex items-center border-b border-gray-100 transition-all duration-300',
        sidebarCollapsed ? 'justify-center px-2' : 'px-6'
      )}>
        <div className="flex items-center gap-3">
          <div className="flex items-center justify-center w-9 h-9 rounded-lg bg-gold-dim shrink-0">
            <WineIcon className="h-5 w-5 text-gold" />
          </div>
          {!sidebarCollapsed && (
            <div className="overflow-hidden">
              <span className="text-base font-bold text-gold tracking-[0.05em] whitespace-nowrap">NoSocial</span>
              <span className="block text-[11px] text-text-muted -mt-0.5 whitespace-nowrap">酒吧管理系统</span>
            </div>
          )}
        </div>
      </div>

      {/* Navigation */}
      <nav className="flex-1 py-4 overflow-y-auto">
        <ul className="space-y-1 px-3">
          {navItems.map((item) => {
            const isActive = location.pathname === item.path;
            const Icon = item.icon;
            return (
              <li key={item.path}>
                <button
                  onClick={() => navigate(item.path)}
                  className={cn(
                    'w-full flex items-center gap-3 rounded-lg transition-all duration-200 group relative',
                    sidebarCollapsed ? 'justify-center px-2 py-3' : 'px-4 py-2.5',
                    isActive
                      ? 'bg-gold-dim text-gold'
                      : 'text-text-secondary hover:text-text-primary hover:bg-surface-elevated'
                  )}
                >
                  {/* Active indicator */}
                  {isActive && (
                    <div className="absolute left-0 top-1/2 -translate-y-1/2 w-[3px] h-5 bg-gold rounded-r-full" />
                  )}
                  <Icon className={cn('h-[18px] w-[18px] shrink-0', isActive ? 'text-gold' : '')} />
                  {!sidebarCollapsed && (
                    <span className="text-sm font-medium whitespace-nowrap">{item.label}</span>
                  )}
                  {/* Tooltip for collapsed */}
                  {sidebarCollapsed && (
                    <div className="absolute left-full ml-2 px-2.5 py-1.5 rounded-md bg-surface-elevated text-text-primary text-xs font-medium whitespace-nowrap opacity-0 pointer-events-none group-hover:opacity-100 transition-opacity shadow-modal border border-gray-200">
                      {item.label}
                    </div>
                  )}
                </button>
              </li>
            );
          })}
        </ul>
      </nav>

      {/* Collapse button */}
      <div className="p-3 border-t border-gray-100">
        <button
          onClick={toggleSidebar}
          className={cn(
            'w-full flex items-center rounded-lg py-2 text-text-secondary hover:text-text-primary hover:bg-surface-elevated transition-all duration-200',
            sidebarCollapsed ? 'justify-center px-2' : 'justify-between px-4'
          )}
        >
          {!sidebarCollapsed && <span className="text-xs">收起菜单</span>}
          <ChevronLeft className={cn(
            'h-4 w-4 transition-transform duration-300',
            sidebarCollapsed && 'rotate-180'
          )} />
        </button>
      </div>
    </aside>
  );
}
