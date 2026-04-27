import { Outlet } from 'react-router-dom';
import { useAppStore } from '@/stores/appStore';
import { cn } from '@/lib/utils';
import { Sidebar } from './Sidebar';
import { Topbar } from './Topbar';
import { ToastContainer } from '@/components/shared/Toast';

export function Layout() {
  const { sidebarCollapsed } = useAppStore();

  return (
    <div className="min-h-screen bg-background">
      <Sidebar />
      <Topbar />

      {/* Main content */}
      <main
        className={cn(
          'pt-16 min-h-screen transition-all duration-300',
          sidebarCollapsed ? 'pl-16' : 'pl-64'
        )}
      >
        <div className="p-6 max-w-[1440px] mx-auto">
          <Outlet />
        </div>
      </main>

      <ToastContainer />
    </div>
  );
}
