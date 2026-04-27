import {
  AreaChart,
  Area,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from 'recharts';
import type { RevenueDataPoint } from '@/types';

interface RevenueChartProps {
  data: RevenueDataPoint[];
}

function CustomTooltip({ active, payload, label }: { active?: boolean; payload?: Array<{ value: number }>; label?: string }) {
  if (active && payload && payload.length) {
    return (
      <div className="rounded-lg bg-surface-elevated border border-gray-200 px-4 py-3 shadow-modal">
        <p className="text-xs text-text-muted mb-1">{label}</p>
        <p className="text-sm font-semibold text-text-primary font-mono">
          ¥{payload[0].value.toLocaleString()}
        </p>
      </div>
    );
  }
  return null;
}

export function RevenueChart({ data }: RevenueChartProps) {
  return (
    <ResponsiveContainer width="100%" height={300}>
      <AreaChart data={data} margin={{ top: 10, right: 10, left: -20, bottom: 0 }}>
        <defs>
          <linearGradient id="revenueGradient" x1="0" y1="0" x2="0" y2="1">
            <stop offset="0%" stopColor="#7C9A92" stopOpacity={0.3} />
            <stop offset="100%" stopColor="#7C9A92" stopOpacity={0} />
          </linearGradient>
        </defs>
        <CartesianGrid strokeDasharray="3 3" stroke="#E5E7EB" vertical={false} />
        <XAxis
          dataKey="date"
          tick={{ fill: '#9CA3AF', fontSize: 12 }}
          axisLine={{ stroke: '#E5E7EB' }}
          tickLine={false}
        />
        <YAxis
          tick={{ fill: '#9CA3AF', fontSize: 12 }}
          axisLine={false}
          tickLine={false}
          tickFormatter={(v: number) => `¥${(v / 1000).toFixed(0)}k`}
        />
        <Tooltip content={<CustomTooltip />} />
        <Area
          type="monotone"
          dataKey="amount"
          stroke="#7C9A92"
          strokeWidth={2}
          fill="url(#revenueGradient)"
          dot={{ fill: '#7C9A92', strokeWidth: 0, r: 4 }}
          activeDot={{ r: 6, fill: '#7C9A92', stroke: '#F5F6F8', strokeWidth: 2 }}
        />
      </AreaChart>
    </ResponsiveContainer>
  );
}
