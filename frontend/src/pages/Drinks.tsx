import { useState, useMemo, useEffect } from 'react';
import { PageHeader } from '@/components/shared/PageHeader';
import { FilterBar, FilterTabs } from '@/components/shared/FilterBar';
import { StatusBadge } from '@/components/shared/StatusBadge';
import { ImageUpload } from '@/components/shared/ImageUpload';
import { Card, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Switch } from '@/components/ui/switch';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import {
  getDrinkList,
  getCategories,
  createDrink,
  updateDrink,
  deleteDrink,
  createCategory,
  updateCategory,
  deleteCategory,
  type Drink,
  type DrinkCategory,
} from '@/lib/api';
import { Plus, Pencil, Trash2, Wine, LayoutGrid, List } from 'lucide-react';
import { useToastStore } from '@/stores/toastStore';

const emptyDrink: Partial<Drink> = {
  name: '',
  english_name: '',
  category_id: 1,
  price: 0,
  cost_price: 0,
  alcohol: 0,
  volume: '',
  description: '',
  ingredients: '',
  image_url: '',
  is_recommended: 0,
  is_free_drink: 0,
  status: 1,
};

export default function Drinks() {
  const { addToast } = useToastStore();
  const [searchQuery, setSearchQuery] = useState('');
  const [filterCategory, setFilterCategory] = useState<string | number>('all');
  const [filterStatus, setFilterStatus] = useState<string | number>('all');
  const [viewMode, setViewMode] = useState<'table' | 'grid'>('table');
  const [drinks, setDrinks] = useState<Drink[]>([]);
  const [categories, setCategories] = useState<DrinkCategory[]>([]);
  const [modalOpen, setModalOpen] = useState(false);
  const [editingDrink, setEditingDrink] = useState<Partial<Drink> | null>(null);
  const [categoryModalOpen, setCategoryModalOpen] = useState(false);
  const [editingCategory, setEditingCategory] = useState<Partial<DrinkCategory> | null>(null);

  useEffect(() => {
    loadDrinks();
    loadCategories();
  }, []);

  async function loadDrinks() {
    const res = await getDrinkList(1, 100);
    if (res.code === 0) setDrinks(res.data.list);
  }

  async function loadCategories() {
    const res = await getCategories();
    if (res.code === 0) setCategories(res.data);
  }

  const filtered = useMemo(() => {
    return drinks.filter((d) => {
      const matchesSearch = !searchQuery || d.name.toLowerCase().includes(searchQuery.toLowerCase());
      const matchesCategory = filterCategory === 'all' || d.category_id === Number(filterCategory);
      const matchesStatus = filterStatus === 'all' || d.status === Number(filterStatus);
      return matchesSearch && matchesCategory && matchesStatus;
    });
  }, [drinks, searchQuery, filterCategory, filterStatus]);

  const openAdd = () => {
    setEditingDrink({ ...emptyDrink });
    setModalOpen(true);
  };

  const openEdit = (drink: Drink) => {
    setEditingDrink({ ...drink });
    setModalOpen(true);
  };

  const handleSave = async () => {
    if (!editingDrink?.name) {
      addToast({ type: 'error', message: '请输入酒水名称' });
      return;
    }
    try {
      if (editingDrink.id) {
        await updateDrink(editingDrink.id, editingDrink);
        addToast({ type: 'success', message: '酒水已更新' });
      } else {
        await createDrink(editingDrink);
        addToast({ type: 'success', message: '酒水已添加' });
      }
      await loadDrinks();
      setModalOpen(false);
      setEditingDrink(null);
    } catch {
      addToast({ type: 'error', message: '操作失败' });
    }
  };

  const handleDelete = async (id: number) => {
    try {
      await deleteDrink(id);
      await loadDrinks();
      addToast({ type: 'success', message: '酒水已删除' });
    } catch {
      addToast({ type: 'error', message: '删除失败' });
    }
  };

  const emptyCategory = { name: '', sort_order: 0, status: 1 };

  const openAddCategory = () => {
    setEditingCategory({ ...emptyCategory });
    setCategoryModalOpen(true);
  };

  const openEditCategory = (cat: DrinkCategory) => {
    setEditingCategory({ ...cat });
    setCategoryModalOpen(true);
  };

  const handleSaveCategory = async () => {
    if (!editingCategory?.name) {
      addToast({ type: 'error', message: '请输入分类名称' });
      return;
    }
    try {
      if (editingCategory.id) {
        await updateCategory(editingCategory.id, editingCategory);
        addToast({ type: 'success', message: '分类已更新' });
      } else {
        await createCategory(editingCategory);
        addToast({ type: 'success', message: '分类已添加' });
      }
      await loadCategories();
      setCategoryModalOpen(false);
      setEditingCategory(null);
    } catch {
      addToast({ type: 'error', message: '操作失败' });
    }
  };

  const handleDeleteCategory = async (id: number) => {
    if (!confirm('确定要删除该分类吗？')) return;
    try {
      await deleteCategory(id);
      await loadCategories();
      addToast({ type: 'success', message: '分类已删除' });
    } catch {
      addToast({ type: 'error', message: '删除失败' });
    }
  };

  return (
    <div className="space-y-6">
      <PageHeader title="酒水管理" subtitle="管理酒水菜单和分类">
        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="icon"
            onClick={() => setViewMode(viewMode === 'table' ? 'grid' : 'table')}
            className="border-border bg-transparent text-text-secondary hover:text-text-primary hover:bg-surface-elevated"
          >
            {viewMode === 'table' ? <LayoutGrid className="h-4 w-4" /> : <List className="h-4 w-4" />}
          </Button>
          <Button
            onClick={openAddCategory}
            className="bg-gold text-white hover:bg-gold-light hover:shadow-gold-glow"
          >
            <Plus className="mr-1 h-4 w-4" />
            分类管理
          </Button>
          <Button
            onClick={openAdd}
            className="bg-gold text-white hover:bg-gold-light hover:shadow-gold-glow"
          >
            <Plus className="mr-2 h-4 w-4" />
            新增酒水
          </Button>
        </div>
      </PageHeader>

      <FilterBar
        searchPlaceholder="搜索酒水名称..."
        searchValue={searchQuery}
        onSearchChange={setSearchQuery}
      >
        <FilterTabs
          options={[{ value: 'all', label: '全部分类' }, ...categories.map((c) => ({ value: c.id, label: c.name }))]}
          value={filterCategory}
          onChange={setFilterCategory}
        />
        <FilterTabs
          options={[
            { value: 'all', label: '全部状态' },
            { value: 1, label: '上架' },
            { value: 0, label: '下架' },
          ]}
          value={filterStatus}
          onChange={setFilterStatus}
        />
      </FilterBar>

      {viewMode === 'table' ? (
        <Card className="bg-surface-secondary border-gray-100">
          <CardContent className="p-0">
            <div className="overflow-x-auto">
              <Table>
                <TableHeader>
                  <TableRow className="border-border hover:bg-transparent">
                    <TableHead className="text-text-muted text-xs font-medium">酒水</TableHead>
                    <TableHead className="text-text-muted text-xs font-medium">分类</TableHead>
                    <TableHead className="text-text-muted text-xs font-medium">价格</TableHead>
                    <TableHead className="text-text-muted text-xs font-medium">酒精度</TableHead>
                    <TableHead className="text-text-muted text-xs font-medium">推荐</TableHead>
                    <TableHead className="text-text-muted text-xs font-medium">状态</TableHead>
                    <TableHead className="text-text-muted text-xs font-medium text-right">操作</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {filtered.map((drink) => (
                    <TableRow key={drink.id} className="border-gray-100 hover:bg-primary/[0.04] transition-colors">
                      <TableCell>
                        <div className="flex items-center gap-3">
                          <div className="w-10 h-10 rounded-lg bg-surface-elevated flex items-center justify-center">
                            {drink.image_url ? (
                              <img src={drink.image_url} alt="" className="w-full h-full object-cover rounded-lg" />
                            ) : (
                              <Wine className="h-5 w-5 text-text-muted" />
                            )}
                          </div>
                          <div>
                            <p className="text-sm font-medium text-text-primary">{drink.name}</p>
                            <p className="text-xs text-text-muted">{drink.english_name}</p>
                          </div>
                        </div>
                      </TableCell>
                      <TableCell className="text-sm text-text-secondary">{drink.category?.name || '-'}</TableCell>
                      <TableCell className="text-sm font-mono text-gold">¥{drink.price.toFixed(2)}</TableCell>
                      <TableCell className="text-sm text-text-secondary">{drink.alcohol}%</TableCell>
                      <TableCell>
                        {drink.is_recommended ? (
                          <span className="text-xs text-gold">是</span>
                        ) : (
                          <span className="text-xs text-text-muted">否</span>
                        )}
                      </TableCell>
                      <TableCell>
                        <StatusBadge status={drink.status} type="drink" />
                      </TableCell>
                      <TableCell className="text-right">
                        <div className="flex items-center justify-end gap-1">
                          <Button
                            variant="ghost"
                            size="icon"
                            onClick={() => openEdit(drink)}
                            className="h-8 w-8 text-text-secondary hover:text-gold hover:bg-primary/10"
                          >
                            <Pencil className="h-4 w-4" />
                          </Button>
                          <Button
                            variant="ghost"
                            size="icon"
                            onClick={() => handleDelete(drink.id)}
                            className="h-8 w-8 text-text-secondary hover:text-red-400 hover:bg-red-500/10"
                          >
                            <Trash2 className="h-4 w-4" />
                          </Button>
                        </div>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </div>
          </CardContent>
        </Card>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          {filtered.map((drink) => (
            <Card key={drink.id} className="bg-surface-secondary border-gray-100 overflow-hidden group">
              <div className="aspect-[16/10] bg-surface-elevated flex items-center justify-center">
                {drink.image_url ? (
                  <img src={drink.image_url} alt="" className="w-full h-full object-cover" />
                ) : (
                  <Wine className="h-10 w-10 text-text-muted" />
                )}
              </div>
              <CardContent className="p-4">
                <h3 className="text-sm font-semibold text-text-primary mb-1">{drink.name}</h3>
                <p className="text-xs text-text-muted mb-3">{drink.english_name}</p>
                <div className="flex items-center justify-between">
                  <span className="text-sm font-mono text-gold">¥{drink.price.toFixed(2)}</span>
                  <span className="text-xs text-text-muted">{drink.alcohol}%</span>
                </div>
                <div className="flex items-center gap-2 mt-3">
                  <StatusBadge status={drink.status} type="drink" />
                  <div className="flex-1" />
                  <Button variant="ghost" size="icon" onClick={() => openEdit(drink)} className="h-7 w-7 text-text-secondary hover:text-gold hover:bg-primary/10">
                    <Pencil className="h-3.5 w-3.5" />
                  </Button>
                  <Button variant="ghost" size="icon" onClick={() => handleDelete(drink.id)} className="h-7 w-7 text-text-secondary hover:text-red-400 hover:bg-red-500/10">
                    <Trash2 className="h-3.5 w-3.5" />
                  </Button>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}

      {/* Add/Edit Modal */}
      <Dialog open={modalOpen} onOpenChange={(open) => { if (!open) setEditingDrink(null); setModalOpen(open); }}>
        <DialogContent className="max-w-[560px] max-h-[85vh] overflow-y-auto bg-surface-secondary border-gray-100">
          <DialogHeader>
            <DialogTitle className="text-text-primary">
              {editingDrink?.id ? '编辑酒水' : '新增酒水'}
            </DialogTitle>
          </DialogHeader>
          {editingDrink && (
            <div className="space-y-4">
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label className="text-text-secondary text-xs">名称</Label>
                  <Input
                    value={editingDrink.name || ''}
                    onChange={(e) => setEditingDrink({ ...editingDrink, name: e.target.value })}
                    className="bg-background border-border text-text-primary"
                  />
                </div>
                <div className="space-y-2">
                  <Label className="text-text-secondary text-xs">英文名</Label>
                  <Input
                    value={editingDrink.english_name || ''}
                    onChange={(e) => setEditingDrink({ ...editingDrink, english_name: e.target.value })}
                    className="bg-background border-border text-text-primary"
                  />
                </div>
              </div>

              <div className="space-y-2">
                <Label className="text-text-secondary text-xs">分类</Label>
                <Select
                  value={String(editingDrink.category_id || 1)}
                  onValueChange={(v) => setEditingDrink({ ...editingDrink, category_id: Number(v) })}
                >
                  <SelectTrigger className="bg-background border-border text-text-primary">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent className="bg-surface-elevated border-gray-200">
                    {categories.map((c) => (
                      <SelectItem key={c.id} value={String(c.id)} className="text-text-primary focus:bg-gold-dim focus:text-gold">
                        {c.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              <div className="grid grid-cols-3 gap-4">
                <div className="space-y-2">
                  <Label className="text-text-secondary text-xs">售价</Label>
                  <Input
                    type="number"
                    value={editingDrink.price || ''}
                    onChange={(e) => setEditingDrink({ ...editingDrink, price: Number(e.target.value) })}
                    className="bg-background border-border text-text-primary"
                  />
                </div>
                <div className="space-y-2">
                  <Label className="text-text-secondary text-xs">成本价</Label>
                  <Input
                    type="number"
                    value={editingDrink.cost_price || ''}
                    onChange={(e) => setEditingDrink({ ...editingDrink, cost_price: Number(e.target.value) })}
                    className="bg-background border-border text-text-primary"
                  />
                </div>
                <div className="space-y-2">
                  <Label className="text-text-secondary text-xs">酒精度(%)</Label>
                  <Input
                    type="number"
                    step="0.1"
                    value={editingDrink.alcohol || ''}
                    onChange={(e) => setEditingDrink({ ...editingDrink, alcohol: Number(e.target.value) })}
                    className="bg-background border-border text-text-primary"
                  />
                </div>
              </div>

              <div className="space-y-2">
                <Label className="text-text-secondary text-xs">容量</Label>
                <Input
                  value={editingDrink.volume || ''}
                  onChange={(e) => setEditingDrink({ ...editingDrink, volume: e.target.value })}
                  placeholder="如 350ml"
                  className="bg-background border-border text-text-primary"
                />
              </div>

              <div className="space-y-2">
                <Label className="text-text-secondary text-xs">配料</Label>
                <Input
                  value={editingDrink.ingredients || ''}
                  onChange={(e) => setEditingDrink({ ...editingDrink, ingredients: e.target.value })}
                  placeholder="逗号分隔"
                  className="bg-background border-border text-text-primary"
                />
              </div>

              <div className="space-y-2">
                <Label className="text-text-secondary text-xs">描述</Label>
                <Textarea
                  value={editingDrink.description || ''}
                  onChange={(e) => setEditingDrink({ ...editingDrink, description: e.target.value })}
                  rows={3}
                  className="bg-background border-border text-text-primary"
                />
              </div>

              <div className="space-y-2">
                <Label className="text-text-secondary text-xs">图片</Label>
                <ImageUpload
                  value={editingDrink.image_url || ''}
                  onChange={(url) => setEditingDrink({ ...editingDrink, image_url: url })}
                  dir="drinks"
                />
              </div>

              <div className="flex items-center gap-6 pt-2">
                <div className="flex items-center gap-2">
                  <Switch
                    checked={!!editingDrink.is_recommended}
                    onCheckedChange={(v) => setEditingDrink({ ...editingDrink, is_recommended: v ? 1 : 0 })}
                  />
                  <Label className="text-text-secondary text-xs cursor-pointer">推荐</Label>
                </div>
                <div className="flex items-center gap-2">
                  <Switch
                    checked={!!editingDrink.is_free_drink}
                    onCheckedChange={(v) => setEditingDrink({ ...editingDrink, is_free_drink: v ? 1 : 0 })}
                  />
                  <Label className="text-text-secondary text-xs cursor-pointer">免费酒</Label>
                </div>
                <div className="flex items-center gap-2">
                  <Switch
                    checked={editingDrink.status === 1}
                    onCheckedChange={(v) => setEditingDrink({ ...editingDrink, status: v ? 1 : 0 })}
                  />
                  <Label className="text-text-secondary text-xs cursor-pointer">上架</Label>
                </div>
              </div>

              <div className="flex justify-end gap-3 pt-4">
                <Button
                  variant="outline"
                  onClick={() => { setModalOpen(false); setEditingDrink(null); }}
                  className="border-border bg-transparent text-text-primary hover:bg-surface-elevated"
                >
                  取消
                </Button>
                <Button
                  onClick={handleSave}
                  className="bg-gold text-white hover:bg-gold-light hover:shadow-gold-glow"
                >
                  保存
                </Button>
              </div>
            </div>
          )}
        </DialogContent>
      </Dialog>

      {/* 分类管理弹窗 */}
      <Dialog open={categoryModalOpen} onOpenChange={(open) => { if (!open) setEditingCategory(null); setCategoryModalOpen(open); }}>
        <DialogContent className="max-w-md bg-surface-secondary border-gray-100">
          <DialogHeader>
            <DialogTitle className="text-text-primary flex items-center justify-between">
              <span>{editingCategory?.id ? '编辑分类' : '新增分类'}</span>
              <div className="flex items-center gap-1">
                {!editingCategory?.id && (
                  <Button size="sm" variant="ghost" onClick={openAddCategory} className="text-gold text-xs h-6 px-2">+ 新增</Button>
                )}
                <Button size="sm" variant="ghost" onClick={() => { setCategoryModalOpen(false); setEditingCategory(null); }} className="text-text-muted h-6 px-2">✕</Button>
              </div>
            </DialogTitle>
          </DialogHeader>

          {/* 分类列表 */}
          {!editingCategory?.id && (
            <div className="space-y-2 max-h-60 overflow-y-auto">
              {categories.length === 0 && (
                <p className="text-text-muted text-sm text-center py-4">暂无分类</p>
              )}
              {categories.map((cat) => (
                <div key={cat.id} className="flex items-center justify-between p-3 rounded-lg bg-surface-elevated border border-border">
                  <div>
                    <p className="text-sm text-text-primary">{cat.name}</p>
                    <p className="text-xs text-text-muted">排序 {cat.sort_order} · {cat.status === 1 ? '启用' : '禁用'}</p>
                  </div>
                  <div className="flex items-center gap-1">
                    <Button size="icon" variant="ghost" className="h-7 w-7" onClick={() => openEditCategory(cat)}>
                      <Pencil className="h-3.5 w-3.5 text-text-secondary" />
                    </Button>
                    <Button size="icon" variant="ghost" className="h-7 w-7" onClick={() => handleDeleteCategory(cat.id)}>
                      <Trash2 className="h-3.5 w-3.5 text-red-400" />
                    </Button>
                  </div>
                </div>
              ))}
            </div>
          )}

          {/* 新增/编辑表单 */}
          {editingCategory && (
            <div className="space-y-4">
              <div className="space-y-2">
                <Label className="text-text-secondary text-xs">分类名称</Label>
                <Input
                  value={editingCategory.name || ''}
                  onChange={(e) => setEditingCategory({ ...editingCategory, name: e.target.value })}
                  placeholder="如：鸡尾酒、啤酒"
                  className="bg-background border-border text-text-primary"
                />
              </div>
              <div className="space-y-2">
                <Label className="text-text-secondary text-xs">排序值</Label>
                <Input
                  type="number"
                  value={editingCategory.sort_order ?? 0}
                  onChange={(e) => setEditingCategory({ ...editingCategory, sort_order: Number(e.target.value) })}
                  className="bg-background border-border text-text-primary"
                />
              </div>
              <div className="flex items-center gap-2">
                <Switch
                  checked={editingCategory.status === 1}
                  onCheckedChange={(v) => setEditingCategory({ ...editingCategory, status: v ? 1 : 0 })}
                />
                <Label className="text-text-secondary text-xs cursor-pointer">启用</Label>
              </div>
              <div className="flex justify-end pt-2">
                <Button
                  onClick={handleSaveCategory}
                  className="bg-gold text-white hover:bg-gold-light hover:shadow-gold-glow"
                >
                  保存
                </Button>
              </div>
            </div>
          )}
        </DialogContent>
      </Dialog>
    </div>
  );
}
