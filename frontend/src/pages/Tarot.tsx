import { useState, useMemo, useEffect, useRef, useCallback } from 'react';
import { PageHeader } from '@/components/shared/PageHeader';
import { Card } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Slider } from '@/components/ui/slider';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import {
  getDrinkList,
  getTarotMappings,
  updateTarotMapping,
  getAdminTarotCards,
  updateTarotCard,
  type Drink,
  type TarotMapping,
  type TarotCard as ApiTarotCard,
} from '@/lib/api';
import { uploadToOSS, validateImageFile } from '@/lib/upload';
import { useToastStore } from '@/stores/toastStore';
import { Settings, Sparkles, ImagePlus, Loader2, X, Upload, Check } from 'lucide-react';

const suitFilters = [
  { value: 'all', label: '全部' },
  { value: 'major', label: '大阿卡纳' },
  { value: '权杖', label: '权杖' },
  { value: '圣杯', label: '圣杯' },
  { value: '宝剑', label: '宝剑' },
  { value: '星币', label: '星币' },
];

export default function Tarot() {
  const { addToast } = useToastStore();
  const [activeSuit, setActiveSuit] = useState('all');
  const [drinks, setDrinks] = useState<Drink[]>([]);
  const [mappings, setMappings] = useState<TarotMapping[]>([]);
  const [cards, setCards] = useState<ApiTarotCard[]>([]);
  const [loading, setLoading] = useState(true);

  // Image upload state
  const [imageModal, setImageModal] = useState<{ open: boolean; card: ApiTarotCard | null }>({
    open: false,
    card: null,
  });
  const [uploading, setUploading] = useState(false);
  const imageInputRef = useRef<HTMLInputElement>(null);

  // Mapping config state
  const [configModal, setConfigModal] = useState<{ open: boolean; card: ApiTarotCard | null }>({
    open: false,
    card: null,
  });
  const [mappingData, setMappingData] = useState({
    drink_id: 0,
    match_score: [80],
    reason_template: '',
    is_reversed: false,
  });

  useEffect(() => {
    loadAll();
  }, []);

  async function loadAll() {
    setLoading(true);
    try {
      const [cardsRes, drinksRes, mappingsRes] = await Promise.all([
        getAdminTarotCards(),
        getDrinkList(1, 100),
        getTarotMappings(),
      ]);
      if (cardsRes.code === 0) setCards(cardsRes.data);
      if (drinksRes.code === 0) setDrinks(drinksRes.data.list);
      if (mappingsRes.code === 0) setMappings(mappingsRes.data);
    } finally {
      setLoading(false);
    }
  }

  const filteredCards = useMemo(() => {
    if (activeSuit === 'all') return cards;
    if (activeSuit === 'major') return cards.filter((c) => c.arcana_type === 1);
    return cards.filter((c) => c.suit === activeSuit);
  }, [cards, activeSuit]);

  // --- Image Upload ---
  const openImageModal = (card: ApiTarotCard, e: React.MouseEvent) => {
    e.stopPropagation();
    setImageModal({ open: true, card });
  };

  const handleImageFile = useCallback(
    async (file: File) => {
      if (!imageModal.card) return;
      const err = validateImageFile(file);
      if (err) {
        addToast({ type: 'error', message: err });
        return;
      }
      setUploading(true);
      try {
        const url = await uploadToOSS(file, 'tarot');
        await updateTarotCard(imageModal.card.id, url);
        // Update local state
        setCards((prev) =>
          prev.map((c) => (c.id === imageModal.card!.id ? { ...c, image_url: url } : c))
        );
        addToast({ type: 'success', message: `「${imageModal.card.name}」图片已更新` });
        setImageModal({ open: false, card: null });
      } catch (e) {
        addToast({ type: 'error', message: e instanceof Error ? e.message : '上传失败' });
      } finally {
        setUploading(false);
      }
    },
    [imageModal.card, addToast]
  );

  const handleDeleteImage = async () => {
    if (!imageModal.card) return;
    try {
      await updateTarotCard(imageModal.card.id, '');
      setCards((prev) =>
        prev.map((c) => (c.id === imageModal.card!.id ? { ...c, image_url: '' } : c))
      );
      addToast({ type: 'success', message: `「${imageModal.card.name}」图片已清除` });
      setImageModal({ open: false, card: null });
    } catch {
      addToast({ type: 'error', message: '操作失败' });
    }
  };

  // --- Mapping Config ---
  const openConfig = (card: ApiTarotCard) => {
    const existing = mappings.find((m) => m.card_no === card.card_no && m.is_reversed === 0);
    setConfigModal({ open: true, card });
    setMappingData({
      drink_id: existing?.drink_id || 0,
      match_score: [existing?.match_score || 80],
      reason_template:
        existing?.reason_template ||
        `${card.name}代表{{keyword}}，这杯{{drink_name}}正适合{{scene}}。`,
      is_reversed: false,
    });
  };

  const handleSaveMapping = async () => {
    if (!configModal.card) return;
    try {
      await updateTarotMapping([
        {
          card_no: configModal.card.card_no,
          card_name: configModal.card.name,
          drink_id: mappingData.drink_id,
          drink_name: drinks.find((d) => d.id === mappingData.drink_id)?.name || '',
          is_reversed: mappingData.is_reversed ? 1 : 0,
          match_score: mappingData.match_score[0],
          reason_template: mappingData.reason_template,
        },
      ]);
      const res = await getTarotMappings();
      if (res.code === 0) setMappings(res.data);
      addToast({ type: 'success', message: '映射配置已保存' });
      setConfigModal({ open: false, card: null });
    } catch {
      addToast({ type: 'error', message: '保存失败' });
    }
  };

  // Check if a card has a mapping configured
  const hasMapping = (cardNo: number) => mappings.some((m) => m.card_no === cardNo);

  // Stats
  const totalCards = cards.length;
  const withImage = cards.filter((c) => c.image_url).length;
  const withMapping = new Set(mappings.map((m) => m.card_no)).size;

  return (
    <div className="space-y-6">
      <PageHeader title="塔罗配置" subtitle="管理78张塔罗牌图片与酒水推荐映射" />

      {/* Stats */}
      <div className="grid grid-cols-3 gap-4">
        <div className="bg-surface-secondary border border-gray-100 rounded-xl px-4 py-3">
          <p className="text-[11px] text-text-muted">总牌数</p>
          <p className="text-xl font-semibold text-text-primary">{totalCards}</p>
        </div>
        <div className="bg-surface-secondary border border-gray-100 rounded-xl px-4 py-3">
          <p className="text-[11px] text-text-muted">已配图</p>
          <p className="text-xl font-semibold text-gold">
            {withImage}
            <span className="text-xs text-text-muted font-normal ml-1">/ {totalCards}</span>
          </p>
        </div>
        <div className="bg-surface-secondary border border-gray-100 rounded-xl px-4 py-3">
          <p className="text-[11px] text-text-muted">已映射酒水</p>
          <p className="text-xl font-semibold text-emerald-600">
            {withMapping}
            <span className="text-xs text-text-muted font-normal ml-1">/ {totalCards}</span>
          </p>
        </div>
      </div>

      {/* Suit Filters */}
      <div className="flex flex-wrap gap-2">
        {suitFilters.map((filter) => (
          <button
            key={filter.value}
            onClick={() => setActiveSuit(filter.value)}
            className={`px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200 ${
              activeSuit === filter.value
                ? 'bg-gold text-white'
                : 'bg-surface-secondary text-text-secondary hover:text-text-primary hover:bg-surface-elevated border border-gray-100'
            }`}
          >
            {filter.label}
          </button>
        ))}
      </div>

      {/* Cards Grid */}
      {loading ? (
        <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 xl:grid-cols-8 gap-4">
          {Array.from({ length: 16 }).map((_, i) => (
            <div key={i} className="aspect-[2/3] rounded-xl bg-surface-elevated animate-pulse" />
          ))}
        </div>
      ) : (
        <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 xl:grid-cols-8 gap-4">
          {filteredCards.map((card) => (
            <Card
              key={card.id}
              className="bg-surface-secondary border-gray-100 overflow-hidden group cursor-pointer hover:border-primary/30 transition-all duration-300"
              onClick={() => openConfig(card)}
            >
              <div className="aspect-[2/3] relative flex flex-col">
                {/* Card image or placeholder */}
                {card.image_url ? (
                  <img
                    src={card.image_url}
                    alt={card.name}
                    className="absolute inset-0 w-full h-full object-cover"
                  />
                ) : (
                  <div className="absolute inset-0 bg-gradient-to-b from-surface-elevated to-background flex flex-col items-center justify-center p-3">
                    <div className="w-10 h-10 rounded-full bg-gold-dim flex items-center justify-center mb-2">
                      <Sparkles className="h-5 w-5 text-gold" />
                    </div>
                    <span className="text-[10px] text-text-muted text-center">未配图</span>
                  </div>
                )}

                {/* Card number badge */}
                <span className="absolute top-1.5 left-1.5 text-[9px] font-mono bg-black/40 text-white px-1.5 py-0.5 rounded">
                  {card.card_no}
                </span>

                {/* Status badges */}
                <div className="absolute top-1.5 right-1.5 flex flex-col gap-1">
                  {card.element && (
                    <span className="text-[8px] px-1.5 py-0.5 rounded-full bg-gold-dim text-gold text-center">
                      {card.element}
                    </span>
                  )}
                  {hasMapping(card.card_no) && (
                    <span className="text-[8px] px-1.5 py-0.5 rounded-full bg-emerald-100 text-emerald-700 text-center">
                      <Check className="h-2.5 w-2.5 inline" />
                    </span>
                  )}
                </div>

                {/* Hover overlay with actions */}
                <div className="absolute inset-0 bg-black/50 flex flex-col items-center justify-center gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                  <Button
                    size="sm"
                    className="bg-gold text-white hover:bg-gold-light text-[11px] h-7 px-3"
                    onClick={(e) => openImageModal(card, e)}
                  >
                    <Upload className="mr-1 h-3 w-3" />
                    {card.image_url ? '换图' : '上传图片'}
                  </Button>
                  <Button
                    size="sm"
                    variant="secondary"
                    className="bg-white/90 text-gray-700 hover:bg-white text-[11px] h-7 px-3"
                    onClick={(e) => {
                      e.stopPropagation();
                      openConfig(card);
                    }}
                  >
                    <Settings className="mr-1 h-3 w-3" />
                    酒水映射
                  </Button>
                </div>
              </div>

              {/* Card info footer */}
              <div className="px-2.5 py-2 border-t border-gray-100/50">
                <h3 className="text-xs font-semibold text-text-primary truncate">{card.name}</h3>
                <p className="text-[10px] text-text-muted truncate">{card.name_en}</p>
                {card.suit && (
                  <span className="text-[9px] text-text-muted">{card.suit}</span>
                )}
              </div>
            </Card>
          ))}
        </div>
      )}

      {/* Image Upload Modal */}
      <Dialog
        open={imageModal.open}
        onOpenChange={(open) => {
          if (!open) setImageModal({ open: false, card: null });
        }}
      >
        <DialogContent className="max-w-[440px] bg-surface-secondary border-gray-100">
          <DialogHeader>
            <DialogTitle className="text-text-primary flex items-center gap-2">
              <ImagePlus className="h-5 w-5 text-gold" />
              「{imageModal.card?.name}」牌面图片
            </DialogTitle>
          </DialogHeader>

          <div className="space-y-4">
            {/* Current image preview */}
            {imageModal.card?.image_url ? (
              <div className="relative mx-auto w-48 rounded-xl overflow-hidden border border-gray-200 shadow-sm">
                <img
                  src={imageModal.card.image_url}
                  alt={imageModal.card.name}
                  className="w-full aspect-[2/3] object-cover"
                />
                <button
                  onClick={handleDeleteImage}
                  className="absolute top-2 right-2 p-1.5 rounded-full bg-red-500/80 text-white hover:bg-red-600 transition-colors"
                >
                  <X className="h-3.5 w-3.5" />
                </button>
              </div>
            ) : (
              <div className="mx-auto w-48 aspect-[2/3] rounded-xl border-2 border-dashed border-border bg-background flex flex-col items-center justify-center gap-2">
                <Sparkles className="h-10 w-10 text-text-muted/30" />
                <span className="text-xs text-text-muted">暂无牌面图片</span>
              </div>
            )}

            {/* Upload area */}
            <div
              onClick={() => imageInputRef.current?.click()}
              className={`rounded-xl border-2 border-dashed border-border bg-background cursor-pointer px-4 py-5 flex flex-col items-center gap-2 transition-all hover:border-gold hover:bg-gold/[0.03] ${
                uploading ? 'opacity-60 pointer-events-none' : ''
              }`}
            >
              {uploading ? (
                <Loader2 className="h-6 w-6 text-gold animate-spin" />
              ) : (
                <Upload className="h-6 w-6 text-text-muted" />
              )}
              <span className="text-sm text-text-secondary">
                {uploading ? '上传中...' : '点击选择新图片'}
              </span>
              <span className="text-[10px] text-text-muted">支持 JPG / PNG / WebP，≤ 10MB</span>
              <input
                ref={imageInputRef}
                type="file"
                accept="image/*"
                onChange={(e) => {
                  const f = e.target.files?.[0];
                  if (f) handleImageFile(f);
                  if (imageInputRef.current) imageInputRef.current.value = '';
                }}
                className="hidden"
              />
            </div>

            {/* Card info */}
            <div className="text-xs text-text-muted space-y-1 bg-background rounded-lg p-3">
              <p>
                <span className="text-text-secondary font-medium">英文名：</span>
                {imageModal.card?.name_en || '-'}
              </p>
              <p>
                <span className="text-text-secondary font-medium">花色：</span>
                {imageModal.card?.suit || '大阿卡纳'}
              </p>
              {imageModal.card?.keywords && (
                <p>
                  <span className="text-text-secondary font-medium">关键词：</span>
                  {imageModal.card.keywords}
                </p>
              )}
            </div>
          </div>
        </DialogContent>
      </Dialog>

      {/* Mapping Config Modal */}
      <Dialog
        open={configModal.open}
        onOpenChange={(open) => {
          if (!open) setConfigModal({ open: false, card: configModal.card });
        }}
      >
        <DialogContent className="max-w-[520px] bg-surface-secondary border-gray-100">
          <DialogHeader>
            <DialogTitle className="text-text-primary flex items-center gap-2">
              <Sparkles className="h-5 w-5 text-gold" />
              配置「{configModal.card?.name}」的酒水映射
            </DialogTitle>
          </DialogHeader>

          <div className="space-y-5">
            {/* Position toggle */}
            <div className="flex items-center gap-4">
              <button
                onClick={() => {
                  setMappingData({ ...mappingData, is_reversed: false });
                  if (configModal.card) {
                    const existing = mappings.find(
                      (m) => m.card_no === configModal.card!.card_no && m.is_reversed === 0
                    );
                    if (existing) {
                      setMappingData({
                        drink_id: existing.drink_id,
                        match_score: [existing.match_score],
                        reason_template: existing.reason_template,
                        is_reversed: false,
                      });
                    }
                  }
                }}
                className={`flex-1 py-2 rounded-lg text-sm font-medium transition-all ${
                  !mappingData.is_reversed
                    ? 'bg-gold text-white'
                    : 'bg-background text-text-secondary hover:text-text-primary'
                }`}
              >
                正位推荐
              </button>
              <button
                onClick={() => {
                  setMappingData({ ...mappingData, is_reversed: true });
                  if (configModal.card) {
                    const existing = mappings.find(
                      (m) => m.card_no === configModal.card!.card_no && m.is_reversed === 1
                    );
                    if (existing) {
                      setMappingData({
                        drink_id: existing.drink_id,
                        match_score: [existing.match_score],
                        reason_template: existing.reason_template,
                        is_reversed: true,
                      });
                    }
                  }
                }}
                className={`flex-1 py-2 rounded-lg text-sm font-medium transition-all ${
                  mappingData.is_reversed
                    ? 'bg-amber-600 text-white'
                    : 'bg-background text-text-secondary hover:text-text-primary'
                }`}
              >
                逆位推荐
              </button>
            </div>

            {/* Drink selector */}
            <div className="space-y-2">
              <Label className="text-text-secondary text-xs">推荐酒水</Label>
              <Select
                value={String(mappingData.drink_id || '')}
                onValueChange={(v) => setMappingData({ ...mappingData, drink_id: Number(v) })}
              >
                <SelectTrigger className="bg-background border-border text-text-primary">
                  <SelectValue placeholder="选择酒水" />
                </SelectTrigger>
                <SelectContent className="bg-surface-elevated border-gray-200">
                  {drinks
                    .filter((d) => d.status === 1)
                    .map((drink) => (
                      <SelectItem
                        key={drink.id}
                        value={String(drink.id)}
                        className="text-text-primary focus:bg-gold-dim focus:text-gold"
                      >
                        {drink.name} - ¥{drink.price.toFixed(2)}
                      </SelectItem>
                    ))}
                </SelectContent>
              </Select>
            </div>

            {/* Match score */}
            <div className="space-y-3">
              <div className="flex items-center justify-between">
                <Label className="text-text-secondary text-xs">匹配分数</Label>
                <span className="text-sm font-mono text-gold">{mappingData.match_score[0]}</span>
              </div>
              <Slider
                value={mappingData.match_score}
                onValueChange={(v) => setMappingData({ ...mappingData, match_score: v })}
                max={100}
                step={1}
                className="w-full"
              />
            </div>

            {/* Reason template */}
            <div className="space-y-2">
              <Label className="text-text-secondary text-xs">推荐理由模板</Label>
              <Textarea
                value={mappingData.reason_template}
                onChange={(e) =>
                  setMappingData({ ...mappingData, reason_template: e.target.value })
                }
                rows={3}
                className="bg-background border-border text-text-primary"
                placeholder="使用 {{card_name}} {{drink_name}} {{keyword}} 等变量"
              />
              <p className="text-[11px] text-text-muted">
                可用变量: {'{{card_name}}'} {'{{drink_name}}'} {'{{keyword}}'} {'{{scene}}'}
              </p>
            </div>

            <div className="flex justify-end gap-3 pt-2">
              <Button
                variant="outline"
                onClick={() => setConfigModal({ open: false, card: null })}
                className="border-border bg-transparent text-text-primary hover:bg-surface-elevated"
              >
                取消
              </Button>
              <Button
                onClick={handleSaveMapping}
                className="bg-gold text-white hover:bg-gold-light hover:shadow-gold-glow"
              >
                保存配置
              </Button>
            </div>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
}
