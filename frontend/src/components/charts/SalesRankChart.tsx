import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Cell,
} from 'recharts';
import type { SalesRankItem } from '@/types';

interface SalesRankChartProps {
  data: SalesRankItem[];
}

function CustomTooltip({ active, payload, label }: { active?: boolean; payload?: Array<{ value: number }>; label?: string }) {
  if (active && payload && payload.length) {
    return (
      <div className="rounded-lg bg-surface-elevated border border-gray-200 px-4 py-3 shadow-modal">
        <p className="text-xs text-text-muted mb-1">{label}</p>
        <p className="text-sm font-semibold text-text-primary font-mono">
          {payload[0].value} 杯
        </p>
      </div>
    );
  }
  return null;
}

export function SalesRankChart({ data }: SalesRankChartProps) {
  const maxSales = Math.max(...data.map((d) => d.sales));

  return (
    <ResponsiveContainer width="100%" height={280}>
      <BarChart data={data} layout="vertical" margin={{ top: 5, right: 30, left: 10, bottom: 5 }}>
        <CartesianGrid strokeDasharray="3 3" stroke="#E5E7EB" horizontal={false} />
        <XAxis
          type="number"
          tick={{ fill: '#9CA3AF', fontSize: 11 }}
          axisLine={false}
          tickLine={false}
          domain={[0, maxSales * 1.2]}
        />
        <YAxis
          type="category"
          dataKey="name"
          tick={{ fill: '#1A1D23', fontSize: 12 }}
          axisLine={false}
          tickLine={false}
          width={80}
        />
        <Tooltip content={<CustomTooltip />} />
        <Bar dataKey="sales" radius={[0, 6, 6, 0]} barSize={24}>
          {data.map((_, index) => (
            <Cell
              key={`cell-${index}`}
              fill={index === 0 ? '#7C9A92' : index === 1 ? '#7C9A92cc' : index === 2 ? '#7C9A9299' : '#7C9A9266'}
            />
          ))}
        </Bar>
      </BarChart>
    </ResponsiveContainer>
  );
}
