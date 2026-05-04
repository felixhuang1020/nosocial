import { Routes, Route, Navigate, useNavigate } from 'react-router-dom';
import { useAuthStore } from '@/stores/authStore';
import { authEvents } from '@/lib/api';
import { Layout } from '@/components/layout/Layout';
import { lazy, Suspense, useEffect, useState } from 'react';
import { Skeleton } from '@/components/ui/skeleton';

// Lazy load pages
const Login = lazy(() => import('./pages/Login'));
const Dashboard = lazy(() => import('./pages/Dashboard'));
const Users = lazy(() => import('./pages/Users'));
const Shareholders = lazy(() => import('./pages/Shareholders'));
const Drinks = lazy(() => import('./pages/Drinks'));
const Orders = lazy(() => import('./pages/Orders'));
const Reviews = lazy(() => import('./pages/Reviews'));
const Tarot = lazy(() => import('./pages/Tarot'));
const Banners = lazy(() => import('./pages/Banners'));
const Settings = lazy(() => import('./pages/Settings'));

function PageLoader() {
  return (
    <div className="space-y-4 p-6">
      <Skeleton className="h-8 w-48 bg-surface-elevated" />
      <div className="grid grid-cols-4 gap-4">
        {[1, 2, 3, 4].map((i) => (
          <Skeleton key={i} className="h-28 bg-surface-elevated rounded-xl" />
        ))}
      </div>
      <Skeleton className="h-80 bg-surface-elevated rounded-xl" />
    </div>
  );
}

function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const { isAuthenticated } = useAuthStore();
  return isAuthenticated ? <>{children}</> : <Navigate to="/login" replace />;
}

function PublicRoute({ children }: { children: React.ReactNode }) {
  const { isAuthenticated } = useAuthStore();
  return !isAuthenticated ? <>{children}</> : <Navigate to="/dashboard" replace />;
}

// 监听 API 层 401 事件，通过 React Router 跳转到登录页
function AuthEventListener() {
  const navigate = useNavigate();
  useEffect(() => {
    const handler = () => {
      navigate('/login', { replace: true });
    };
    authEvents.addEventListener('unauthorized', handler);
    return () => authEvents.removeEventListener('unauthorized', handler);
  }, [navigate]);
  return null;
}

export default function App() {
  const checkAuth = useAuthStore((s) => s.checkAuth);
  const [authChecked, setAuthChecked] = useState(false);

  useEffect(() => {
    checkAuth().finally(() => setAuthChecked(true));
  }, []);

  if (!authChecked) {
    return <div className="flex items-center justify-center h-screen">加载中...</div>;
  }

  return (
    <Suspense fallback={<PageLoader />}>
      <AuthEventListener />
      <Routes>
        <Route
          path="/login"
          element={
            <PublicRoute>
              <Login />
            </PublicRoute>
          }
        />
        <Route
          path="/"
          element={
            <ProtectedRoute>
              <Layout />
            </ProtectedRoute>
          }
        >
          <Route index element={<Navigate to="/dashboard" replace />} />
          <Route path="dashboard" element={<Dashboard />} />
          <Route path="users" element={<Users />} />
          <Route path="shareholders" element={<Shareholders />} />
          <Route path="drinks" element={<Drinks />} />
          <Route path="orders" element={<Orders />} />
          <Route path="reviews" element={<Reviews />} />
          <Route path="tarot" element={<Tarot />} />
          <Route path="banners" element={<Banners />} />
          <Route path="settings" element={<Settings />} />
        </Route>
        <Route path="*" element={<Navigate to="/dashboard" replace />} />
      </Routes>
    </Suspense>
  );
}
