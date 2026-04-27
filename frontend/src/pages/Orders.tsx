import { useState, useMemo, useEffect } from 'react';
import { PageHeader } from '@/components/shared/PageHeader';
import { FilterBar, FilterTabs } from '@/components/shared/FilterBar';
import { StatCard } from '@/components/shared/StatCard';
import { StatusBadge } from '@/components/shared/StatusBadge';
import { Card, CardContent } from '@/components/ui/card';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { getOrderList, type Order as APIOrder } from '@/lib/api';
import { ShoppingCart, Clock, CheckCircle, XCircle } from 'lucide-react';

export default function Orders() {
  const [searchQuery, setSearchQuery] = useState('');
  const [filterStatus, setFilterStatus] = useState<string | number>('all');
  const [orders, setOrders] = useState<APIOrder[]>([]);

  useEffect(() => {
    getOrderList(1, 100).then((res) => {
      if (res.code === 0) setOrders(res.data.list);
    });
  }, []);

  const filtered = useMemo(() => {
    return orders.filter((o) => {
      const matchesSearch = !searchQuery || o.order_no.includes(searchQuery);
      const matchesStatus = filterStatus === 'all' || o.status === Number(filterStatus);
      return matchesSearch && matchesStatus;
    });
  }, [orders, searchQuery, filterStatus]);

  const stats = useMemo(() => {
    return {
      total: orders.length,
      pending: orders.filter((o) => o.status === 0 || o.status === 1).length,
      completed: orders.filter((o) => o.status === 3).length,
      cancelled: orders.filter((o) => o.status === 4).length,
    };
  }, [orders]);

  return (
    <div className="space-y-6">
      <PageHeader title="订单管理" subtitle="管理酒水订单" />

      {/* Stats */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <StatCard title="全部订单" value={stats.total} icon={ShoppingCart} />
        <StatCard title="待处理" value={stats.pending} icon={Clock} />
        <StatCard title="已完成" value={stats.completed} icon={CheckCircle} />
        <StatCard title="已取消" value={stats.cancelled} icon={XCircle} />
      </div>

      <FilterBar
        searchPlaceholder="搜索订单号或消费者..."
        searchValue={searchQuery}
        onSearchChange={setSearchQuery}
      >
        <FilterTabs
          options={[
            { value: 'all', label: '全部' },
            { value: 0, label: '待支付' },
            { value: 1, label: '已支付' },
            { value: 2, label: '制作中' },
            { value: 3, label: '已完成' },
            { value: 4, label: '已取消' },
          ]}
          value={filterStatus}
          onChange={setFilterStatus}
        />
      </FilterBar>

      <Card className="bg-surface-secondary border-gray-100">
        <CardContent className="p-0">
          <div className="overflow-x-auto">
            <Table>
              <TableHeader>
                <TableRow className="border-border hover:bg-transparent">
                  <TableHead className="text-text-muted text-xs font-medium">订单号</TableHead>
                  <TableHead className="text-text-muted text-xs font-medium">消费者</TableHead>
                  <TableHead className="text-text-muted text-xs font-medium">关联股东</TableHead>
                  <TableHead className="text-text-muted text-xs font-medium">订单金额</TableHead>
                  <TableHead className="text-text-muted text-xs font-medium">实付</TableHead>
                  <TableHead className="text-text-muted text-xs font-medium">状态</TableHead>
                  <TableHead className="text-text-muted text-xs font-medium">时间</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {filtered.map((order) => (
                  <TableRow key={order.id} className="border-gray-100 hover:bg-primary/[0.04] transition-colors">
                    <TableCell className="text-sm font-mono text-text-primary">{order.order_no}</TableCell>
                    <TableCell className="text-sm text-text-primary">用户#{order.user_id}</TableCell>
                    <TableCell className="text-sm text-text-secondary">
                      -
                    </TableCell>
                    <TableCell className="text-sm font-mono text-text-secondary">¥{order.total_amount.toFixed(2)}</TableCell>
                    <TableCell className="text-sm font-mono text-gold">¥{order.pay_amount.toFixed(2)}</TableCell>
                    <TableCell>
                      <StatusBadge status={order.status} type="order" />
                    </TableCell>
                    <TableCell className="text-sm text-text-muted">{order.created_at.slice(5, 16)}</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>

          {filtered.length === 0 && (
            <div className="py-16 text-center">
              <p className="text-text-muted text-sm">暂无匹配的订单</p>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
