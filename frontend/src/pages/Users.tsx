import { useState, useMemo, useEffect } from 'react';
import { PageHeader } from '@/components/shared/PageHeader';
import { FilterBar, FilterTabs } from '@/components/shared/FilterBar';
import { StatusBadge } from '@/components/shared/StatusBadge';
import { ConfirmDialog } from '@/components/shared/ConfirmDialog';
import { Card, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { getUserList, updateUserStatus, type User as APIUser } from '@/lib/api';
import { Download, Ban, CheckCircle, Eye } from 'lucide-react';
import { useToastStore } from '@/stores/toastStore';

export default function Users() {
  const { addToast } = useToastStore();
  const [searchQuery, setSearchQuery] = useState('');
  const [filterRole, setFilterRole] = useState<string | number>('all');
  const [users, setUsers] = useState<APIUser[]>([]);
  const [loading, setLoading] = useState(true);
  const [confirmDialog, setConfirmDialog] = useState<{ open: boolean; user: APIUser | null; action: 'disable' | 'enable' }>({
    open: false,
    user: null,
    action: 'disable',
  });

  useEffect(() => {
    getUserList(1, 100).then((res) => {
      if (res.code === 0) setUsers(res.data.list);
      setLoading(false);
    });
  }, []);

  const filteredUsers = useMemo(() => {
    return users.filter((user) => {
      const matchesSearch =
        !searchQuery ||
        (user.nickname && user.nickname.toLowerCase().includes(searchQuery.toLowerCase())) ||
        (user.phone && user.phone.includes(searchQuery));
      const matchesRole =
        filterRole === 'all'
          ? true
          : filterRole === 'shareholder'
            ? user.is_shareholder === 1
            : user.is_shareholder === 0;
      return matchesSearch && matchesRole;
    });
  }, [users, searchQuery, filterRole]);

  const handleToggleStatus = (user: APIUser) => {
    const action = user.status === 1 ? 'disable' : 'enable';
    setConfirmDialog({ open: true, user, action });
  };

  const confirmToggle = async () => {
    if (!confirmDialog.user) return;
    const newStatus = confirmDialog.user.status === 1 ? 0 : 1;
    try {
      await updateUserStatus(confirmDialog.user.id, newStatus);
      setUsers((prev) =>
        prev.map((u) =>
          u.id === confirmDialog.user!.id ? { ...u, status: newStatus } : u
        )
      );
      addToast({
        type: 'success',
        message: confirmDialog.action === 'disable' ? '用户已禁用' : '用户已启用',
      });
    } catch {
      addToast({ type: 'error', message: '操作失败' });
    }
    setConfirmDialog({ open: false, user: null, action: 'disable' });
  };

  return (
    <div className="space-y-6">
      <PageHeader title="用户管理" subtitle="管理小程序注册用户">
        <Button variant="outline" className="border-border bg-transparent text-text-secondary hover:text-text-primary hover:bg-surface-elevated">
          <Download className="mr-2 h-4 w-4" />
          导出数据
        </Button>
      </PageHeader>

      <FilterBar
        searchPlaceholder="搜索昵称或手机号..."
        searchValue={searchQuery}
        onSearchChange={setSearchQuery}
      >
        <FilterTabs
          options={[
            { value: 'all', label: '全部' },
            { value: 'shareholder', label: '股东' },
            { value: 'normal', label: '普通用户' },
          ]}
          value={filterRole}
          onChange={setFilterRole}
        />
      </FilterBar>

      <Card className="bg-surface-secondary border-gray-100">
        <CardContent className="p-0">
          <div className="overflow-x-auto">
            <Table>
              <TableHeader>
                <TableRow className="border-border hover:bg-transparent">
                  <TableHead className="text-text-muted text-xs font-medium">用户信息</TableHead>
                  <TableHead className="text-text-muted text-xs font-medium">手机号</TableHead>
                  <TableHead className="text-text-muted text-xs font-medium">生日</TableHead>
                  <TableHead className="text-text-muted text-xs font-medium">股东</TableHead>
                  <TableHead className="text-text-muted text-xs font-medium">注册时间</TableHead>
                  <TableHead className="text-text-muted text-xs font-medium">状态</TableHead>
                  <TableHead className="text-text-muted text-xs font-medium text-right">操作</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {filteredUsers.map((user) => (
                  <TableRow
                    key={user.id}
                    className="border-gray-100 hover:bg-primary/[0.04] transition-colors"
                  >
                    <TableCell>
                      <div className="flex items-center gap-3">
                        <div className="w-9 h-9 rounded-full bg-gold-dim flex items-center justify-center shrink-0">
                          <span className="text-sm font-medium text-gold">
                            {user.nickname?.charAt(0) || '?'}
                          </span>
                        </div>
                        <div>
                          <p className="text-sm font-medium text-text-primary">{user.nickname || '未命名'}</p>
                          <p className="text-xs text-text-muted">ID: {user.id}</p>
                        </div>
                      </div>
                    </TableCell>
                    <TableCell className="text-sm font-mono text-text-secondary">{user.phone}</TableCell>
                    <TableCell className="text-sm text-text-secondary">{user.birthday?.slice(0,10) || '-'}</TableCell>
                    <TableCell>
                      {user.is_shareholder === 1 ? (
                        <span className="inline-flex items-center gap-1 text-xs text-green-500">
                          <CheckCircle className="h-3.5 w-3.5" />
                          是
                        </span>
                      ) : (
                        <span className="text-xs text-text-muted">-</span>
                      )}
                    </TableCell>
                    <TableCell className="text-sm text-text-muted">{user.created_at?.slice(0, 16)}</TableCell>
                    <TableCell>
                      <StatusBadge status={user.status} type="user" />
                    </TableCell>
                    <TableCell className="text-right">
                      <div className="flex items-center justify-end gap-1">
                        <Button
                          variant="ghost"
                          size="icon"
                          className="h-8 w-8 text-text-secondary hover:text-text-primary hover:bg-surface-elevated"
                        >
                          <Eye className="h-4 w-4" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => handleToggleStatus(user)}
                          className={`h-8 w-8 ${
                            user.status === 1
                              ? 'text-text-secondary hover:text-red-400 hover:bg-red-500/10'
                              : 'text-text-secondary hover:text-green-500 hover:bg-green-500/10'
                          }`}
                        >
                          {user.status === 1 ? <Ban className="h-4 w-4" /> : <CheckCircle className="h-4 w-4" />}
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>

          {filteredUsers.length === 0 && (
            <div className="py-16 text-center">
              <p className="text-text-muted text-sm">暂无匹配的用户</p>
            </div>
          )}
        </CardContent>
      </Card>

      <ConfirmDialog
        open={confirmDialog.open}
        onOpenChange={(open) => setConfirmDialog({ ...confirmDialog, open })}
        title={confirmDialog.action === 'disable' ? '禁用用户' : '启用用户'}
        description={`确定要${confirmDialog.action === 'disable' ? '禁用' : '启用'}用户「${confirmDialog.user?.nickname || ''}」吗？`}
        type={confirmDialog.action === 'disable' ? 'warning' : 'success'}
        confirmText={confirmDialog.action === 'disable' ? '禁用' : '启用'}
        onConfirm={confirmToggle}
      />
    </div>
  );
}
