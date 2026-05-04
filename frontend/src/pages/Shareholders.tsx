import { useState, useEffect } from 'react';
import { PageHeader } from '@/components/shared/PageHeader';
import { FilterBar } from '@/components/shared/FilterBar';
import { StatCard } from '@/components/shared/StatCard';
import { Card, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { getShareholderList, getShareholderEarnings, getShareholderTeam, type User, type EarningsRecord } from '@/lib/api';
import { Pagination } from '@/components/shared/Pagination';
import { Users, DollarSign, Award, Eye, GitBranch } from 'lucide-react';

export default function Shareholders() {
  const [searchQuery, setSearchQuery] = useState('');
  const [debouncedSearch, setDebouncedSearch] = useState('');
  const [shareholders, setShareholders] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const pageSize = 20;
  const [dialog, setDialog] = useState<{
    type: null | 'earnings' | 'team';
    shareholder: User | null;
    data: EarningsRecord[] | User[];
  }>({
    type: null,
    shareholder: null,
    data: [],
  });

  // debounce 搜索
  useEffect(() => {
    const timer = setTimeout(() => {
      setDebouncedSearch(searchQuery);
      setPage(1);
    }, 500);
    return () => clearTimeout(timer);
  }, [searchQuery]);

  // 数据获取
  useEffect(() => {
    setLoading(true);
    getShareholderList(page, pageSize, debouncedSearch || undefined).then((res) => {
      if (res.code === 0) {
        setShareholders(res.data.list || []);
        setTotal(res.data.total || 0);
      }
    }).finally(() => setLoading(false));
  }, [page, debouncedSearch]);

  const openEarnings = async (sh: User) => {
    const res = await getShareholderEarnings(sh.id, 1, 100);
    setDialog({
      type: 'earnings',
      shareholder: sh,
      data: res.code === 0 ? (res.data?.list || []) : [],
    });
  };

  const openTeam = async (sh: User) => {
    const res = await getShareholderTeam(sh.id);
    setDialog({
      type: 'team',
      shareholder: sh,
      data: res.code === 0 ? (res.data || []) : [],
    });
  };

  const closeDialog = () => setDialog({ type: null, shareholder: null, data: [] });

  return (
    <div className="space-y-6">
      <PageHeader title="股东管理" subtitle="管理共享股东及佣金" />

      {/* Stats */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <StatCard title="总股东数" value={total} icon={Users} />
        <StatCard title="本月新增" value={0} icon={Award} />
        <StatCard title="待结算佣金" value={0} prefix="¥" icon={DollarSign} />
        <StatCard title="总发放佣金" value={0} prefix="¥" icon={DollarSign} />
      </div>

      <FilterBar
        searchPlaceholder="搜索昵称或邀请码..."
        searchValue={searchQuery}
        onSearchChange={setSearchQuery}
      />

      <Card className="bg-surface-secondary border-gray-100">
        <CardContent className="p-0">
          {loading ? (
            <div className="py-16 text-center text-gray-400">加载中...</div>
          ) : shareholders.length === 0 ? (
            <div className="py-16 text-center">
              <p className="text-text-muted text-sm">暂无股东数据</p>
            </div>
          ) : (
          <div className="overflow-x-auto">
            <Table>
              <TableHeader>
                <TableRow className="border-border hover:bg-transparent">
                  <TableHead className="text-text-muted text-xs font-medium">股东</TableHead>
                  <TableHead className="text-text-muted text-xs font-medium">邀请码</TableHead>
                  <TableHead className="text-text-muted text-xs font-medium">团队人数</TableHead>
                  <TableHead className="text-text-muted text-xs font-medium">累计收益</TableHead>
                  <TableHead className="text-text-muted text-xs font-medium">余额</TableHead>
                  <TableHead className="text-text-muted text-xs font-medium text-right">操作</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
              {shareholders.map((sh) => (
                  <TableRow key={sh.id} className="border-gray-100 hover:bg-primary/[0.04] transition-colors">
                    <TableCell>
                      <div className="flex items-center gap-3">
                        <div className="w-9 h-9 rounded-full bg-gold-dim flex items-center justify-center shrink-0">
                          <span className="text-sm font-medium text-gold">{(sh.nickname || '?').charAt(0)}</span>
                        </div>
                        <span className="text-sm font-medium text-text-primary">{sh.nickname || '未知'}</span>
                      </div>
                    </TableCell>
                    <TableCell className="text-sm font-mono text-gold">{sh.invite_code || '-'}</TableCell>
                    <TableCell className="text-sm text-text-primary">{sh.team_count ?? '-'} 人</TableCell>
                    <TableCell className="text-sm font-mono text-gold">¥{sh.total_earning?.toFixed(2) || '0.00'}</TableCell>
                    <TableCell className="text-sm font-mono text-text-primary">¥{sh.balance?.toFixed(2) || '0.00'}</TableCell>
                    <TableCell className="text-right">
                      <div className="flex items-center justify-end gap-1">
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => openTeam(sh)}
                          className="h-8 w-8 text-text-secondary hover:text-gold hover:bg-primary/10"
                        >
                          <GitBranch className="h-4 w-4" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => openEarnings(sh)}
                          className="h-8 w-8 text-text-secondary hover:text-text-primary hover:bg-surface-elevated"
                        >
                          <Eye className="h-4 w-4" />
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>
          )}
        </CardContent>
      </Card>

      <Pagination page={page} total={total} pageSize={pageSize} onPageChange={setPage} />

      {/* Earnings Dialog */}
      <Dialog open={dialog.type === 'earnings'} onOpenChange={() => closeDialog()}>
        <DialogContent className="max-w-[640px] bg-surface-secondary border-gray-100">
          <DialogHeader>
            <DialogTitle className="text-text-primary">
              {dialog.shareholder?.nickname || '股东'} 的收益明细
            </DialogTitle>
          </DialogHeader>
          <div className="overflow-x-auto">
            <Table>
              <TableHeader>
                <TableRow className="border-border hover:bg-transparent">
                  <TableHead className="text-text-muted text-xs">类型</TableHead>
                  <TableHead className="text-text-muted text-xs">金额</TableHead>
                  <TableHead className="text-text-muted text-xs">状态</TableHead>
                  <TableHead className="text-text-muted text-xs">时间</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {(dialog.type === 'earnings' ? (dialog.data as EarningsRecord[]) : []).map((record: EarningsRecord) => (
                  <TableRow key={record.id} className="border-gray-100">
                    <TableCell className="text-sm text-text-primary">{record.type || '-'}</TableCell>
                    <TableCell className="text-sm font-mono text-gold">¥{Number(record.amount || 0).toFixed(2)}</TableCell>
                    <TableCell>
                      <span className={`text-xs ${record.status === 1 ? 'text-green-500' : 'text-amber-500'}`}>
                        {record.status === 1 ? '已结算' : '待结算'}
                      </span>
                    </TableCell>
                    <TableCell className="text-sm text-text-muted">{record.created_at || '-'}</TableCell>
                  </TableRow>
                ))}
                {dialog.type === 'earnings' && dialog.data.length === 0 && (
                  <TableRow>
                    <TableCell colSpan={4} className="text-center py-8 text-text-muted text-sm">暂无收益记录</TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
          </div>
        </DialogContent>
      </Dialog>

      {/* Team Dialog */}
      <Dialog open={dialog.type === 'team'} onOpenChange={() => closeDialog()}>
        <DialogContent className="max-w-[480px] bg-surface-secondary border-gray-100">
          <DialogHeader>
            <DialogTitle className="text-text-primary">
              {dialog.shareholder?.nickname} 的团队 ({dialog.type === 'team' ? dialog.data.length : 0}人)
            </DialogTitle>
          </DialogHeader>
          <div className="py-2">
            {dialog.type === 'team' && dialog.data.length > 0 ? (
              <div className="space-y-2">
                {(dialog.data as User[]).map((member) => (
                  <div key={member.id} className="flex items-center gap-3 p-2 rounded-lg bg-background">
                    <div className="w-8 h-8 rounded-full bg-gold-dim flex items-center justify-center shrink-0">
                      <span className="text-xs font-medium text-gold">{(member.nickname || '?').charAt(0)}</span>
                    </div>
                    <div className="flex-1">
                      <p className="text-sm text-text-primary">{member.nickname || '未知'}</p>
                      <p className="text-xs text-text-muted">
                        {member.is_shareholder === 1 ? '股东' : '普通用户'}
                        {member.created_at && ` · ${member.created_at.slice(0, 10)}`}
                      </p>
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <div className="py-8 text-center">
                <GitBranch className="h-10 w-10 text-text-muted mx-auto mb-2" />
                <p className="text-text-muted text-sm">暂无团队成员</p>
              </div>
            )}
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
}
