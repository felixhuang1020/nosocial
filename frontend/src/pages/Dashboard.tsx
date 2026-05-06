import { useNavigate } from 'react-router-dom';
import { useState, useEffect } from 'react';
import { PageHeader } from '@/components/shared/PageHeader';
import { StatCard } from '@/components/shared/StatCard';
import { RevenueChart } from '@/components/charts/RevenueChart';
import { StatusBadge } from '@/components/shared/StatusBadge';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { getDashboardStats, getOrderList, getReviewList, type Order, type Review } from '@/lib/api';
import type { RevenueDataPoint } from '@/types';
import {
  DollarSign,
  ShoppingCart,
  Users,
  MessageSquare,
  ArrowRight,
  CheckCircle,
} from 'lucide-react';
import { useToastStore } from '@/stores/toastStore';

export default function Dashboard() {
  const navigate = useNavigate();
  const addToast = useToastStore((s) => s.addToast);
  const [stats, setStats] = useState({
    today_amount: 0,
    today_order_count: 0,
    new_shareholder_count: 0,
    pending_review_count: 0,
  });
  const [revenueTrend, setRevenueTrend] = useState<RevenueDataPoint[]>([]);
  const [orders, setOrders] = useState<Order[]>([]);
  const [reviews, setReviews] = useState<Review[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function load() {
      try {
        const [statsRes, orderRes, reviewRes] = await Promise.all([
          getDashboardStats(),
          getOrderList(1, 5),
          getReviewList(1, 5, 0),
        ]);
        if (statsRes.code === 0) {
          setStats(statsRes.data);
          if (statsRes.data.revenue_trend) {
            setRevenueTrend(statsRes.data.revenue_trend);
          }
        }
        if (orderRes.code === 0) setOrders(orderRes.data.list);
        if (reviewRes.code === 0) setReviews(reviewRes.data.list);
      } catch (e) {
        const msg = e instanceof Error ? e.message : '加载数据失败';
        addToast({ type: 'error', message: msg });
      } finally {
        setLoading(false);
      }
    }
    load();
  }, [addToast]);

  if (loading) {
    return (
      <div className="flex items-center justify-center py-24">
        <p className="text-text-muted text-sm">加载中...</p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <PageHeader title="数据看板" subtitle="实时监控酒吧经营状况" />

      {/* Stats Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <StatCard
          title="今日营业额"
          value={stats.today_amount}
          prefix="¥"
          icon={DollarSign}
        />
        <StatCard
          title="今日订单"
          value={stats.today_order_count}
          icon={ShoppingCart}
        />
        <StatCard
          title="新增股东"
          value={stats.new_shareholder_count}
          icon={Users}
        />
        <StatCard
          title="待审核点评"
          value={stats.pending_review_count}
          icon={MessageSquare}
        />
      </div>

      {/* Charts Row */}
      <div className="grid grid-cols-1 gap-6">
        {/* Revenue Trend */}
        <Card className="bg-surface-secondary border-gray-100">
          <CardHeader className="pb-2">
            <div className="flex items-center justify-between">
              <CardTitle className="text-lg text-text-primary">营业额趋势（近7天）</CardTitle>
            </div>
          </CardHeader>
          <CardContent>
            {revenueTrend.length > 0 ? (
              <RevenueChart data={revenueTrend} />
            ) : (
              <div className="py-16 text-center">
                <p className="text-text-muted text-sm">暂无数据</p>
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Bottom Lists */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-6">
        {/* Recent Orders */}
        <Card className="md:col-span-1 lg:col-span-3 bg-surface-secondary border-gray-100">
          <CardHeader className="pb-3">
            <div className="flex items-center justify-between">
              <CardTitle className="text-lg text-text-primary">最近订单</CardTitle>
              <Button
                variant="ghost"
                size="sm"
                onClick={() => navigate('/orders')}
                className="text-text-secondary hover:text-gold hover:bg-primary/10"
              >
                查看全部
                <ArrowRight className="ml-1 h-4 w-4" />
              </Button>
            </div>
          </CardHeader>
          <CardContent className="pt-0">
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr className="border-b border-border">
                    <th className="text-left text-xs text-text-muted font-medium py-2.5 px-3">订单号</th>
                    <th className="text-left text-xs text-text-muted font-medium py-2.5 px-3">消费者</th>
                    <th className="text-left text-xs text-text-muted font-medium py-2.5 px-3">金额</th>
                    <th className="text-left text-xs text-text-muted font-medium py-2.5 px-3">状态</th>
                    <th className="text-left text-xs text-text-muted font-medium py-2.5 px-3">时间</th>
                  </tr>
                </thead>
                <tbody>
                  {orders.map((order: Order) => (
                    <tr key={order.id} className="border-b border-gray-100 hover:bg-primary/[0.04] transition-colors">
                      <td className="py-3 px-3 text-sm font-mono text-text-secondary">{order.order_no}</td>
                      <td className="py-3 px-3 text-sm text-text-primary">用户#{order.user_id}</td>
                      <td className="py-3 px-3 text-sm font-mono text-gold">¥{order.pay_amount.toFixed(2)}</td>
                      <td className="py-3 px-3">
                        <StatusBadge status={order.status} type="order" />
                      </td>
                      <td className="py-3 px-3 text-sm text-text-muted">{order.created_at.slice(5, 16)}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </CardContent>
        </Card>

        {/* Pending Reviews */}
        <Card className="md:col-span-1 lg:col-span-2 bg-surface-secondary border-gray-100">
          <CardHeader className="pb-3">
            <div className="flex items-center justify-between">
              <CardTitle className="text-lg text-text-primary">待审核点评</CardTitle>
              <Button
                variant="ghost"
                size="sm"
                onClick={() => navigate('/reviews')}
                className="text-text-secondary hover:text-gold hover:bg-primary/10"
              >
                去审核
                <CheckCircle className="ml-1 h-4 w-4" />
              </Button>
            </div>
          </CardHeader>
          <CardContent className="pt-0">
            <div className="space-y-3">
              {reviews.map((review: Review) => (
                <div
                  key={review.id}
                  onClick={() => navigate('/reviews')}
                  className="flex items-center gap-3 p-3 rounded-lg bg-background hover:bg-surface-elevated transition-colors cursor-pointer"
                >
                  <div className="w-9 h-9 rounded-full bg-gold-dim flex items-center justify-center shrink-0">
                    <span className="text-sm font-medium text-gold">
                      {review.user_id}
                    </span>
                  </div>
                  <div className="flex-1 min-w-0">
                    <p className="text-sm text-text-primary truncate">用户#{review.user_id}</p>
                    <p className="text-xs text-text-muted">{review.created_at.slice(5, 16)}</p>
                  </div>
                  <StatusBadge status={review.status} type="review" />
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
