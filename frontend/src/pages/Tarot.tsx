import { useState, useMemo, useEffect } from 'react';
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
import { getDrinkList, getTarotMappings, updateTarotMapping, getPublicTarotCards, type Drink, type TarotMapping, type TarotCard as ApiTarotCard } from '@/lib/api';
import { useToastStore } from '@/stores/toastStore';
import { Settings, Sparkles } from 'lucide-react';

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
    loadCards();
    loadDrinks();
    loadMappings();
  }, []);

  async function loadCards() {
    const res = await getPublicTarotCards();
    if (res.code === 0) setCards(res.data);
  }

  async function loadDrinks() {
    const res = await getDrinkList(1, 100);
    if (res.code === 0) setDrinks(res.data.list);
  }

  async function loadMappings() {
    const res = await getTarotMappings();
    if (res.code === 0) setMappings(res.data);
  }

  const filteredCards = useMemo(() => {
    if (activeSuit === 'all') return cards;
    if (activeSuit === 'major') return cards.filter((c) => c.arcana_type === 1);
    return cards.filter((c) => c.suit === activeSuit);
  }, [cards, activeSuit]);

  const openConfig = (card: ApiTarotCard) => {
    const existing = mappings.find((m) => m.card_no === card.card_no && m.is_reversed === 0);
    setConfigModal({ open: true, card });
    setMappingData({
      drink_id: existing?.drink_id || 0,
      match_score: [existing?.match_score || 80],
      reason_template: existing?.reason_template || `${card.name}代表{{keyword}}，这杯{{drink_name}}正适合{{scene}}。`,
      is_reversed: false,
    });
  };

  const handleSaveMapping = async () => {
    if (!configModal.card) return;
    try {
      await updateTarotMapping([{
        card_no: configModal.card.card_no,
        card_name: configModal.card.name,
        drink_id: mappingData.drink_id,
        drink_name: drinks.find((d) => d.id === mappingData.drink_id)?.name || '',
        is_reversed: mappingData.is_reversed ? 1 : 0,
        match_score: mappingData.match_score[0],
        reason_template: mappingData.reason_template,
      }]);
      await loadMappings();
      addToast({ type: 'success', message: '映射配置已保存' });
      setConfigModal({ open: false, card: null });
    } catch {
      addToast({ type: 'error', message: '保存失败' });
    }
  };

  return (
    <div className="space-y-6">
      <PageHeader title="塔罗配置" subtitle="管理78张塔罗牌与酒水推荐映射" />

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
      <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 xl:grid-cols-8 gap-4">
        {filteredCards.map((card) => (
          <Card
            key={card.id}
            className="bg-surface-secondary border-gray-100 overflow-hidden group cursor-pointer hover:border-primary/30 transition-all duration-300"
            onClick={() => openConfig(card)}
          >
            <div className="aspect-[2/3] bg-gradient-to-b from-surface-elevated to-background flex flex-col items-center justify-center p-4 relative">
              {/* Card number */}
              <span className="absolute top-2 left-2 text-[10px] font-mono text-text-muted">
                {card.card_no}
              </span>

              {/* Element badge */}
              {card.element && (
                <span className="absolute top-2 right-2 text-[9px] px-1.5 py-0.5 rounded-full bg-gold-dim text-gold">
                  {card.element}
                </span>
              )}

              {/* Card icon */}
              <div className="w-12 h-12 rounded-full bg-gold-dim flex items-center justify-center mb-3">
                <Sparkles className="h-6 w-6 text-gold" />
              </div>

              {/* Card name */}
              <h3 className="text-sm font-semibold text-text-primary text-center">{card.name}</h3>
              <p className="text-[10px] text-text-muted mt-1 text-center">{card.name_en}</p>

              {/* Suit badge */}
              {card.suit && (
                <span className="mt-2 text-[10px] text-text-muted">{card.suit}</span>
              )}

              {/* Hover overlay */}
              <div className="absolute inset-0 bg-gold/10 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
                <Button
                  size="sm"
                  variant="secondary"
                  className="bg-gold text-white hover:bg-gold-light text-xs"
                  onClick={(e) => { e.stopPropagation(); openConfig(card); }}
                >
                  <Settings className="mr-1 h-3 w-3" />
                  配置
                </Button>
              </div>
            </div>
          </Card>
        ))}
      </div>

      {/* Config Modal */}
      <Dialog open={configModal.open} onOpenChange={(open) => { if (!open) setConfigModal({ open: false, card: configModal.card }); }}>
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
                onClick={() => setMappingData({ ...mappingData, is_reversed: false })}
                className={`flex-1 py-2 rounded-lg text-sm font-medium transition-all ${
                  !mappingData.is_reversed
                    ? 'bg-gold text-white'
                    : 'bg-background text-text-secondary hover:text-text-primary'
                }`}
              >
                正位推荐
              </button>
              <button
                onClick={() => setMappingData({ ...mappingData, is_reversed: true })}
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
                  {drinks.filter((d) => d.status === 1).map((drink) => (
                    <SelectItem key={drink.id} value={String(drink.id)} className="text-text-primary focus:bg-gold-dim focus:text-gold">
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
                onChange={(e) => setMappingData({ ...mappingData, reason_template: e.target.value })}
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
