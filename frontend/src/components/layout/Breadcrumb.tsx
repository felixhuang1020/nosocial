import { useNavigate, useLocation } from 'react-router-dom';
import { ChevronRight, Home } from 'lucide-react';

const routeNames: Record<string, string> = {
  '/dashboard': '数据看板',
  '/users': '用户管理',
  '/shareholders': '股东管理',
  '/drinks': '酒水管理',
  '/orders': '订单管理',
  '/reviews': '点评审核',
  '/tarot': '塔罗配置',
  '/banners': '图片管理',
  '/settings': '系统设置',
};

export function Breadcrumb() {
  const location = useLocation();
  const navigate = useNavigate();
  const currentName = routeNames[location.pathname] || '页面';

  return (
    <nav className="flex items-center gap-2 text-sm">
      <button
        onClick={() => navigate('/dashboard')}
        className="flex items-center gap-1 text-text-muted hover:text-gold transition-colors"
      >
        <Home className="h-3.5 w-3.5" />
        <span>首页</span>
      </button>
      <ChevronRight className="h-3.5 w-3.5 text-text-muted" />
      <span className="text-text-primary font-medium">{currentName}</span>
    </nav>
  );
}
