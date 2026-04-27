import { create } from 'zustand';
import { adminLogin, type AdminLoginResp } from '@/lib/api';

interface AdminUser {
  username: string;
  nickname: string;
  role: number;
}

interface AuthState {
  isAuthenticated: boolean;
  user: AdminUser | null;
  login: (username: string, password: string) => Promise<boolean>;
  logout: () => void;
}

export const useAuthStore = create<AuthState>((set) => ({
  isAuthenticated: !!localStorage.getItem('admin_token'),
  user: (() => {
    const raw = localStorage.getItem('admin_user');
    return raw ? (JSON.parse(raw) as AdminUser) : null;
  })(),

  login: async (username: string, password: string) => {
    if (!username.trim() || !password.trim()) return false;
    try {
      const res = await adminLogin({ username, password });
      if (res.code !== 0 || !res.data) return false;
      const { token, nickname, role } = res.data;
      localStorage.setItem('admin_token', token);
      const user = { username, nickname, role };
      localStorage.setItem('admin_user', JSON.stringify(user));
      set({ isAuthenticated: true, user });
      return true;
    } catch {
      return false;
    }
  },

  logout: () => {
    localStorage.removeItem('admin_token');
    localStorage.removeItem('admin_user');
    set({ isAuthenticated: false, user: null });
  },
}));
