import { useState, useMemo, useEffect, useCallback } from 'react';
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
import { ORDER_STATUS } from '@/lib/constants';
import { Pagination } from '@/components/shared/Pagination';
import { ShoppingCart, Clock, CheckCircle, XCircle } from 'lucide-react';

export default function Orders() {
  const [searchQuery, setSearchQuery] = useState('');
  const [filterStatus, setFilterStatus] = useState<string | number>('all');
  const [orders, setOrders] = useState<APIOrder[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const pageSize = 20;

  // status 切换时重置页码
  const handleStatusChange = useCallback((val: string | number) => {
    setFilterStatus(val);
    setPage(1);
  }, []);

  useEffect(() => {
    let cancelled = false;
    const fetchData = async () => {
      setLoading(true);
      try {
        const status = filterStatus === 'all' ? undefined : Number(filterStatus);
        const res = await getOrderList(page, pageSize, status);
        if (!cancelled && res.code === 0) {
          setOrders(res.data.list || []);
          setTotal(res.data.total || 0);
        }
      } catch {
        // handle error
      } finally {
        if (!cancelled) setLoading(false);
      }
    };
    fetchData();
    return () => { cancelled = true; };
  }, [page, filterStatus]);

  // 订单号前端搜索（后端不支持 order_no 搜索）
  const filtered = useMemo(() => {
    if (!searchQuery) return orders;
    return orders.filter((o) => o.order_no.includes(searchQuery));
  }, [orders, searchQuery]);

  const stats = useMemo(() => {
    return {
      total: orders.length,
      pending: orders.filter((o) => o.status === ORDER_STATUS.UNPAID || o.status === ORDER_STATUS.PAID).length,
      completed: orders.filter((o) => o.status === ORDER_STATUS.COMPLETED).length,
      cancelled: orders.filter((o) => o.status === ORDER_STATUS.CANCELLED).length,
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
          onChange={handleStatusChange}
        />
      </FilterBar>

      <Card className="bg-surface-secondary border-gray-100">
        <CardContent className="p-0">
          {loading ? (
            <div className="py-16 text-center text-gray-400">加载中...</div>
          ) : filtered.length === 0 ? (
            <div className="py-16 text-center">
              <p className="text-text-muted text-sm">暂无订单数据</p>
            </div>
          ) : (
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
          )}
        </CardContent>
      </Card>

      <Pagination page={page} total={total} pageSize={pageSize} onPageChange={setPage} />
    </div>
  );
}
