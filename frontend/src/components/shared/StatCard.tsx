import type { LucideIcon } from 'lucide-react';
import { TrendingUp, TrendingDown } from 'lucide-react';
import { useEffect, useState } from 'react';

interface StatCardProps {
  title: string;
  value: number;
  prefix?: string;
  suffix?: string;
  trend?: number;
  icon: LucideIcon;
}

function useCountUp(end: number, duration: number = 600) {
  const [count, setCount] = useState(0);
  useEffect(() => {
    let startTime: number | null = null;
    const animate = (currentTime: number) => {
      if (!startTime) startTime = currentTime;
      const progress = Math.min((currentTime - startTime) / duration, 1);
      setCount(Math.floor(progress * end));
      if (progress < 1) requestAnimationFrame(animate);
    };
    requestAnimationFrame(animate);
  }, [end, duration]);
  return count;
}

export function StatCard({ title, value, prefix = '', suffix = '', trend, icon: Icon }: StatCardProps) {
  const animatedValue = useCountUp(value);

  return (
    <div className="relative rounded-xl bg-surface-secondary border border-gray-100 p-6 transition-all duration-300 hover:border-primary/10 hover:shadow-card group overflow-hidden">
      {/* Gold top line */}
      <div className="absolute top-0 left-4 right-4 h-[2px] gradient-gold-line opacity-60" />

      <div className="flex items-start justify-between">
        <div className="flex-1">
          <p className="text-sm text-text-secondary mb-2">{title}</p>
          <p className="text-[28px] font-bold text-text-primary font-mono tracking-tight">
            {prefix}{animatedValue.toLocaleString()}{suffix}
          </p>
          {trend !== undefined && (
            <div className="flex items-center gap-1 mt-2">
              {trend >= 0 ? (
                <>
                  <TrendingUp className="h-3.5 w-3.5 text-green-500" />
                  <span className="text-xs font-medium text-green-500">+{trend}%</span>
                </>
              ) : (
                <>
                  <TrendingDown className="h-3.5 w-3.5 text-red-500" />
                  <span className="text-xs font-medium text-red-500">{trend}%</span>
                </>
              )}
              <span className="text-xs text-text-muted ml-1">较昨日</span>
            </div>
          )}
        </div>
        <div className="flex items-center justify-center w-10 h-10 rounded-[10px] bg-gold-dim shrink-0">
          <Icon className="h-5 w-5 text-gold" />
        </div>
      </div>
    </div>
  );
}
