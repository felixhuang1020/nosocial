import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuthStore } from '@/stores/authStore';
import { useToastStore } from '@/stores/toastStore';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Wine, User, Lock, Eye, EyeOff } from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';

export default function Login() {
  const navigate = useNavigate();
  const { login } = useAuthStore();
  const { addToast } = useToastStore();
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [loading, setLoading] = useState(false);
  const [shake, setShake] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!username.trim() || !password.trim()) {
      setShake(true);
      setTimeout(() => setShake(false), 400);
      addToast({ type: 'error', message: '请输入用户名和密码' });
      return;
    }

    setLoading(true);
    const success = await login(username, password);
    setLoading(false);

    if (success) {
      addToast({ type: 'success', message: '登录成功' });
      navigate('/dashboard');
    } else {
      setShake(true);
      setTimeout(() => setShake(false), 400);
      addToast({ type: 'error', message: '用户名或密码错误' });
    }
  };

  return (
    <div className="min-h-screen flex">
      {/* Left - Brand */}
      <div className="hidden lg:flex lg:w-3/5 relative overflow-hidden gradient-hero items-center justify-center">
        <div className="absolute inset-0 gradient-gold-glow" />
        <div className="relative z-10 text-center px-12">
          <div className="flex items-center justify-center gap-3 mb-6">
            <div className="flex items-center justify-center w-14 h-14 rounded-xl bg-gold-dim">
              <Wine className="h-8 w-8 text-gold" />
            </div>
          </div>
          <h1 className="text-4xl font-bold text-gold tracking-[0.05em] mb-3">NoSocial</h1>
          <p className="text-lg text-text-secondary mb-8">酒吧管理系统</p>
          <div className="flex items-center gap-8 justify-center text-text-muted text-sm">
            <span>共享股东</span>
            <span className="w-1 h-1 rounded-full bg-text-muted" />
            <span>塔罗推荐</span>
            <span className="w-1 h-1 rounded-full bg-text-muted" />
            <span>生日营销</span>
            <span className="w-1 h-1 rounded-full bg-text-muted" />
            <span>点评返券</span>
          </div>
        </div>
        {/* Decorative lines */}
        <div className="absolute bottom-0 left-0 right-0 h-px bg-gradient-to-r from-transparent via-gold/20 to-transparent" />
      </div>

      {/* Right - Form */}
      <div className="flex-1 flex items-center justify-center px-6 py-12 bg-background relative">
        <AnimatePresence>
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5, ease: 'easeOut' }}
            className="w-full max-w-[400px]"
          >
            {/* Mobile logo */}
            <div className="lg:hidden flex items-center justify-center gap-3 mb-8">
              <div className="flex items-center justify-center w-10 h-10 rounded-lg bg-gold-dim">
                <Wine className="h-6 w-6 text-gold" />
              </div>
              <span className="text-xl font-bold text-gold tracking-[0.05em]">NoSocial</span>
            </div>

            <div className="mb-8">
              <h2 className="text-2xl font-semibold text-text-primary mb-1">欢迎回来</h2>
              <p className="text-sm text-text-secondary">请登录您的管理账户</p>
            </div>

            <motion.form
              onSubmit={handleSubmit}
              animate={shake ? { x: [0, -8, 8, -4, 0] } : {}}
              transition={{ duration: 0.4 }}
              className="space-y-5"
            >
              <div className="space-y-2">
                <Label htmlFor="username" className="text-sm text-text-secondary">
                  用户名
                </Label>
                <div className="relative">
                  <User className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-text-muted" />
                  <Input
                    id="username"
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                    placeholder="请输入用户名"
                    className="pl-10 bg-surface-secondary border-border text-text-primary placeholder:text-text-muted focus:border-gold focus:ring-gold/15"
                  />
                </div>
              </div>

              <div className="space-y-2">
                <Label htmlFor="password" className="text-sm text-text-secondary">
                  密码
                </Label>
                <div className="relative">
                  <Lock className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-text-muted" />
                  <Input
                    id="password"
                    type={showPassword ? 'text' : 'password'}
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    placeholder="请输入密码"
                    className="pl-10 pr-10 bg-surface-secondary border-border text-text-primary placeholder:text-text-muted focus:border-gold focus:ring-gold/15"
                  />
                  <button
                    type="button"
                    onClick={() => setShowPassword(!showPassword)}
                    className="absolute right-3 top-1/2 -translate-y-1/2 text-text-muted hover:text-text-secondary transition-colors"
                  >
                    {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                  </button>
                </div>
              </div>

              <Button
                type="submit"
                disabled={loading}
                className="w-full h-11 bg-gold text-white font-semibold hover:bg-gold-light hover:shadow-gold-glow transition-all duration-200"
              >
                {loading ? (
                  <div className="h-5 w-5 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                ) : (
                  '登录'
                )}
              </Button>
            </motion.form>

            <p className="mt-6 text-center text-xs text-text-muted">
              默认账号: admin / admin123
            </p>
          </motion.div>
        </AnimatePresence>
      </div>
    </div>
  );
}
