import { cn } from '@/lib/utils';

interface StatusBadgeProps {
  status: number;
  type?: 'user' | 'drink' | 'order' | 'review' | 'shareholder';
  text?: string;
}

const statusConfig: Record<string, Record<number, { label: string; className: string }>> = {
  user: {
    0: { label: '禁用', className: 'bg-gray-500/12 text-gray-400' },
    1: { label: '正常', className: 'bg-green-500/12 text-green-500' },
  },
  drink: {
    0: { label: '下架', className: 'bg-gray-500/12 text-gray-400' },
    1: { label: '上架', className: 'bg-green-500/12 text-green-500' },
  },
  order: {
    0: { label: '待支付', className: 'bg-blue-500/12 text-blue-500' },
    1: { label: '已支付', className: 'bg-blue-500/12 text-blue-500' },
    2: { label: '制作中', className: 'bg-amber-500/12 text-amber-500' },
    3: { label: '已完成', className: 'bg-green-500/12 text-green-500' },
    4: { label: '已取消', className: 'bg-red-500/12 text-red-500' },
  },
  review: {
    0: { label: '待审核', className: 'bg-amber-500/12 text-amber-500' },
    1: { label: '已通过', className: 'bg-green-500/12 text-green-500' },
    2: { label: '已拒绝', className: 'bg-red-500/12 text-red-500' },
  },
  shareholder: {
    0: { label: '禁用', className: 'bg-gray-500/12 text-gray-400' },
    1: { label: '正常', className: 'bg-green-500/12 text-green-500' },
  },
};

export function StatusBadge({ status, type = 'user', text }: StatusBadgeProps) {
  const config = statusConfig[type]?.[status] || { label: '未知', className: 'bg-gray-500/12 text-gray-400' };

  return (
    <span
      className={cn(
        'inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium',
        config.className
      )}
    >
      {text || config.label}
    </span>
  );
}
