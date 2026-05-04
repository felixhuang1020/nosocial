import { create } from 'zustand';
import { adminLogin, adminLogout, getDashboardStats } from '@/lib/api';

interface AdminUser {
  username: string;
  nickname: string;
  role: number;
}

interface AuthState {
  isAuthenticated: boolean;
  user: AdminUser | null;
  login: (username: string, password: string) => Promise<boolean>;
  logout: () => Promise<void>;
  checkAuth: () => Promise<boolean>;
}

export const useAuthStore = create<AuthState>((set) => ({
  isAuthenticated: false,
  user: null,

  login: async (username: string, password: string) => {
    if (!username.trim() || !password.trim()) return false;
    try {
      const res = await adminLogin({ username, password });
      if (res.code !== 0 || !res.data) return false;
      const { nickname, role } = res.data;
      const user: AdminUser = { username, nickname, role };
      set({ isAuthenticated: true, user });
      return true;
    } catch {
      return false;
    }
  },

  logout: async () => {
    try {
      await adminLogout();
    } catch {
      // ignore — Cookie may already be expired
    }
    set({ isAuthenticated: false, user: null });
  },

  checkAuth: async () => {
    try {
      const res = await getDashboardStats();
      if (res.code === 0) {
        set({ isAuthenticated: true });
        return true;
      }
    } catch {
      // 401 or network error — not authenticated
    }
    set({ isAuthenticated: false, user: null });
    return false;
  },
}));
