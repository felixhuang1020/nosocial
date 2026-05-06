import { useState, useEffect, useCallback } from 'react';
import { PageHeader } from '@/components/shared/PageHeader';
import { FilterTabs } from '@/components/shared/FilterBar';
import { ImageUpload } from '@/components/shared/ImageUpload';
import { Card, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { getBannerList, createBanner, deleteBanner, type Banner } from '@/lib/api';
import { useToastStore } from '@/stores/toastStore';
import { Plus, Trash2, ImagePlus } from 'lucide-react';

const positionOptions = [
  { value: 'all', label: '全部位置' },
  { value: 1, label: '首页轮播' },
  { value: 2, label: '展示页画廊' },
];

export default function Banners() {
  const { addToast } = useToastStore();
  const [banners, setBanners] = useState<Banner[]>([]);
  const [filterPosition, setFilterPosition] = useState<string | number>('all');
  const [uploadModal, setUploadModal] = useState(false);
  const [newBanner, setNewBanner] = useState({
    title: '',
    image_url: '',
    position: 1,
  });

  const loadBanners = useCallback(async () => {
    const res = await getBannerList();
    if (res.code === 0) setBanners(res.data?.list || []);
  }, []);

  useEffect(() => {
    void loadBanners();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const filtered = banners.filter((b) =>
    filterPosition === 'all' ? true : b.position === Number(filterPosition)
  );

  const handleOpenUploadModal = () => {
    const initialPosition = filterPosition !== 'all' ? Number(filterPosition) : 1;
    setNewBanner({
      title: '',
      image_url: '',
      position: initialPosition,
    });
    setUploadModal(true);
  };

  const handleAdd = async () => {
    if (!newBanner.title || !newBanner.image_url) {
      addToast({ type: 'error', message: '请填写标题并上传图片' });
      return;
    }
    try {
      await createBanner(newBanner);
      await loadBanners();
      addToast({ type: 'success', message: '图片已添加' });
      setUploadModal(false);
      setNewBanner({ title: '', image_url: '', position: 1 });
    } catch {
      addToast({ type: 'error', message: '添加失败' });
    }
  };

  const handleDelete = async (id: number) => {
    try {
      await deleteBanner(id);
      await loadBanners();
      addToast({ type: 'success', message: '图片已删除' });
    } catch {
      addToast({ type: 'error', message: '删除失败' });
    }
  };

  return (
    <div className="space-y-6">
      <PageHeader title="图片管理" subtitle="管理首页轮播图和展示页画廊">
        <Button
          onClick={handleOpenUploadModal}
          className="bg-gold text-white hover:bg-gold-light hover:shadow-gold-glow"
        >
          <Plus className="mr-2 h-4 w-4" />
          上传图片
        </Button>
      </PageHeader>

      <FilterTabs
        options={positionOptions}
        value={filterPosition}
        onChange={setFilterPosition}
      />

      {/* Image Grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
        {/* Upload placeholder card */}
        <Card
          className="bg-surface-secondary border-dashed border-2 border-border hover:border-primary hover:bg-primary/[0.04] transition-all duration-200 cursor-pointer min-h-[200px] flex flex-col items-center justify-center"
          onClick={handleOpenUploadModal}
        >
          <ImagePlus className="h-10 w-10 text-text-muted mb-2" />
          <span className="text-sm text-text-muted">点击上传新图片</span>
        </Card>

        {filtered.map((banner) => (
          <Card
            key={banner.id}
            className="bg-surface-secondary border-gray-100 overflow-hidden group"
          >
            <div className="relative aspect-[4/3]">
              <img
                src={banner.image_url}
                alt={banner.title}
                className="w-full h-full object-cover"
              />
              {/* Hover overlay */}
              <div className="absolute inset-0 bg-black/30 flex items-center justify-center gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                <Button
                  variant="ghost"
                  size="icon"
                  onClick={() => handleDelete(banner.id)}
                  className="h-9 w-9 bg-red-500/20 text-red-400 hover:bg-red-500/30 hover:text-red-400"
                >
                  <Trash2 className="h-5 w-5" />
                </Button>
              </div>
            </div>
            <CardContent className="p-3">
              <p className="text-sm font-medium text-text-primary truncate">{banner.title}</p>
              <p className="text-xs text-text-muted mt-0.5">
                {banner.position === 1 ? '首页轮播' : '展示页画廊'} · 排序 {banner.sort_order}
              </p>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Upload Modal */}
      <Dialog open={uploadModal} onOpenChange={setUploadModal}>
        <DialogContent className="max-w-[480px] bg-surface-secondary border-gray-100">
          <DialogHeader>
            <DialogTitle className="text-text-primary">上传图片</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2">
              <Label className="text-text-secondary text-xs">标题</Label>
              <Input
                value={newBanner.title}
                onChange={(e) => setNewBanner({ ...newBanner, title: e.target.value })}
                placeholder="输入图片标题"
                className="bg-background border-border text-text-primary placeholder:text-text-muted"
              />
            </div>

            <div className="space-y-2">
              <Label className="text-text-secondary text-xs">位置</Label>
              <div className="flex gap-2">
                <button
                  onClick={() => setNewBanner({ ...newBanner, position: 1 })}
                  className={`flex-1 py-2 rounded-lg text-sm font-medium transition-all ${
                    newBanner.position === 1
                      ? 'bg-gold text-white'
                      : 'bg-background text-text-secondary hover:text-text-primary'
                  }`}
                >
                  首页轮播
                </button>
                <button
                  onClick={() => setNewBanner({ ...newBanner, position: 2 })}
                  className={`flex-1 py-2 rounded-lg text-sm font-medium transition-all ${
                    newBanner.position === 2
                      ? 'bg-gold text-white'
                      : 'bg-background text-text-secondary hover:text-text-primary'
                  }`}
                >
                  展示页画廊
                </button>
              </div>
            </div>

            <div className="space-y-2">
              <Label className="text-text-secondary text-xs">图片</Label>
              <ImageUpload
                value={newBanner.image_url}
                onChange={(url) => setNewBanner({ ...newBanner, image_url: url })}
                className="w-full h-40"
                dir="banners"
              />
            </div>

            <div className="flex justify-end gap-3 pt-2">
              <Button
                variant="outline"
                onClick={() => { setUploadModal(false); setNewBanner({ title: '', image_url: '', position: 1 }); }}
                className="border-border bg-transparent text-text-primary hover:bg-surface-elevated"
              >
                取消
              </Button>
              <Button
                onClick={handleAdd}
                className="bg-gold text-white hover:bg-gold-light hover:shadow-gold-glow"
              >
                确认上传
              </Button>
            </div>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
}
