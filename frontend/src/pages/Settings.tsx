import { useState, useEffect } from 'react';
import { PageHeader } from '@/components/shared/PageHeader';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Slider } from '@/components/ui/slider';
import { getDrinkList, getSettings, updateSettings, type Drink, type AdminSettings } from '@/lib/api';
import { useToastStore } from '@/stores/toastStore';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { DollarSign, Store, CreditCard } from 'lucide-react';

const defaultSettings: AdminSettings = {
  shareholder_fee: 99.00,
  commission_rate: 0.10,
  free_drink_id: 1,
  review_coupon_amount: 20.00,
  review_coupon_min_order: 0.00,
  review_coupon_valid_days: 30,
  name: 'NoSocial Bar',
  address: '上海市静安区南京西路1266号',
  phone: '021-6288-8888',
  business_hours: '18:00 - 04:00',
  wifi_name: 'NoSocial_Free',
  wifi_password: 'nosocial888',
};

export default function Settings() {
  const { addToast } = useToastStore();
  const [settings, setSettings] = useState<AdminSettings>(defaultSettings);
  const [drinks, setDrinks] = useState<Drink[]>([]);
  const [hasChanges, setHasChanges] = useState(false);
  const [_loading, setLoading] = useState(true);

  useEffect(() => {
    loadSettings();
    loadDrinks();
  }, []);

  async function loadSettings() {
    try {
      const res = await getSettings();
      if (res.code === 0 && res.data) {
        setSettings({ ...defaultSettings, ...res.data });
      }
    } catch (e) {
      const msg = e instanceof Error ? e.message : '加载设置失败';
      addToast({ type: 'error', message: msg });
    } finally {
      setLoading(false);
    }
  }

  async function loadDrinks() {
    const res = await getDrinkList(1, 100);
    if (res.code === 0) setDrinks(res.data.list);
  }

  const handleChange = (field: string, value: unknown) => {
    setSettings((prev) => ({ ...prev, [field]: value }));
    setHasChanges(true);
  };

  const handleSave = async () => {
    try {
      // Convert all values to strings for the API
      const payload: Record<string, string> = {};
      for (const [key, val] of Object.entries(settings)) {
        if (val !== undefined && val !== null) {
          payload[key] = String(val);
        }
      }
      const res = await updateSettings(payload);
      if (res.code === 0) {
        addToast({ type: 'success', message: '设置已保存' });
        setHasChanges(false);
      } else {
        addToast({ type: 'error', message: res.msg || '保存失败' });
      }
    } catch {
      addToast({ type: 'error', message: '保存失败' });
    }
  };

  return (
    <div className="space-y-6">
      <PageHeader title="系统设置" subtitle="配置业务参数和系统选项">
        {hasChanges && (
          <Button
            onClick={handleSave}
            className="bg-gold text-white hover:bg-gold-light hover:shadow-gold-glow"
          >
            保存设置
          </Button>
        )}
      </PageHeader>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Business Settings */}
        <Card className="bg-surface-secondary border-gray-100">
          <CardHeader>
            <div className="flex items-center gap-2">
              <DollarSign className="h-5 w-5 text-gold" />
              <CardTitle className="text-lg text-text-primary">业务设置</CardTitle>
            </div>
          </CardHeader>
          <CardContent className="space-y-5">
            <div className="space-y-2">
              <Label className="text-text-secondary text-xs">股东注册费 (元)</Label>
              <Input
                type="number"
                value={settings.shareholder_fee || 0}
                onChange={(e) => handleChange('shareholder_fee', Number(e.target.value))}
                className="bg-background border-border text-text-primary"
              />
            </div>

            <div className="space-y-3">
              <div className="flex items-center justify-between">
                <Label className="text-text-secondary text-xs">佣金比例</Label>
                <span className="text-sm font-mono text-gold">{((settings.commission_rate || 0) * 100).toFixed(0)}%</span>
              </div>
              <Slider
                value={[(settings.commission_rate || 0) * 100]}
                onValueChange={(v) => handleChange('commission_rate', v[0] / 100)}
                max={30}
                step={1}
              />
              <p className="text-[11px] text-text-muted">
                消费者通过邀请码消费后，上级股东获得的分成比例
              </p>
            </div>

            <div className="space-y-2">
              <Label className="text-text-secondary text-xs">免费赠送酒水</Label>
              <Select
                value={String(settings.free_drink_id || 1)}
                onValueChange={(v) => handleChange('free_drink_id', Number(v))}
              >
                <SelectTrigger className="bg-background border-border text-text-primary">
                  <SelectValue />
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
          </CardContent>
        </Card>

        {/* Store Info */}
        <Card className="bg-surface-secondary border-gray-100">
          <CardHeader>
            <div className="flex items-center gap-2">
              <Store className="h-5 w-5 text-gold" />
              <CardTitle className="text-lg text-text-primary">店铺信息</CardTitle>
            </div>
          </CardHeader>
          <CardContent className="space-y-5">
            <div className="space-y-2">
              <Label className="text-text-secondary text-xs">店铺名称</Label>
              <Input
                value={String(settings.name || '')}
                onChange={(e) => handleChange('name', e.target.value)}
                className="bg-background border-border text-text-primary"
              />
            </div>

            <div className="space-y-2">
              <Label className="text-text-secondary text-xs">店铺地址</Label>
              <Input
                value={String(settings.address || '')}
                onChange={(e) => handleChange('address', e.target.value)}
                className="bg-background border-border text-text-primary"
              />
            </div>

            <div className="space-y-2">
              <Label className="text-text-secondary text-xs">联系电话</Label>
              <Input
                value={String(settings.phone || '')}
                onChange={(e) => handleChange('phone', e.target.value)}
                className="bg-background border-border text-text-primary"
              />
            </div>

            <div className="space-y-2">
              <Label className="text-text-secondary text-xs">营业时间</Label>
              <Input
                value={String(settings.business_hours || '')}
                onChange={(e) => handleChange('business_hours', e.target.value)}
                placeholder="如 18:00 - 04:00"
                className="bg-background border-border text-text-primary"
              />
            </div>
          </CardContent>
        </Card>

        {/* WeChat Pay */}
        <Card className="bg-surface-secondary border-gray-100">
          <CardHeader>
            <div className="flex items-center gap-2">
              <CreditCard className="h-5 w-5 text-gold" />
              <CardTitle className="text-lg text-text-primary">微信支付</CardTitle>
            </div>
          </CardHeader>
          <CardContent className="space-y-5">
            <div className="rounded-md border border-amber-200 bg-amber-50 p-3 text-[12px] leading-5 text-amber-900">
              <p className="font-medium">重要提示</p>
              <p className="mt-1 text-amber-800">
                APIv3 密钥、商户私钥、证书序列号等敏感凭据出于安全合规统一在服务器 <code className="rounded bg-amber-100 px-1">config/config.yaml</code> 中配置，不在后台暴露或存库。
                商户号与 AppID 可在此处查看，修改请同步更新服务器配置并重启服务。
              </p>
            </div>
            <div className="space-y-2">
              <Label className="text-text-secondary text-xs">AppID</Label>
              <Input
                value={String(settings.wx_appid || '')}
                onChange={(e) => handleChange('wx_appid', e.target.value)}
                placeholder="wx开头的小程序 AppID"
                className="bg-background border-border text-text-primary"
              />
            </div>

            <div className="space-y-2">
              <Label className="text-text-secondary text-xs">MCH ID （商户号）</Label>
              <Input
                value={String(settings.wx_mch_id || '')}
                onChange={(e) => handleChange('wx_mch_id', e.target.value)}
                placeholder="微信支付直连商户号"
                className="bg-background border-border text-text-primary"
              />
            </div>

            <div className="space-y-2">
              <Label className="text-text-secondary text-xs">回调地址 Notify URL</Label>
              <Input
                value={String(settings.wx_notify_url || '')}
                onChange={(e) => handleChange('wx_notify_url', e.target.value)}
                placeholder="https://your-domain.com/api/v1/payment/wxpay/notify"
                className="bg-background border-border text-text-primary"
              />
              <p className="text-[11px] text-text-muted">需通过 HTTPS 公网访问，微信支付仅对 443 端口发起回调请求。</p>
            </div>

            <div className="space-y-2">
              <Label className="text-text-secondary text-xs">APIv3 密钥</Label>
              <Input
                type="password"
                value="●●●●●●●●●●●●"
                readOnly
                disabled
                className="bg-background border-border text-text-muted cursor-not-allowed"
              />
              <p className="text-[11px] text-text-muted">由服务器配置文件统一管理，此处不支持在线修改。</p>
            </div>

            <div className="space-y-2">
              <Label className="text-text-secondary text-xs">商户证书与私钥</Label>
              <div className="rounded border border-border bg-background px-3 py-2 text-[12px] text-text-muted">
                证书序列号 / 私钥文件路径由服务器配置提供；平台证书由 SDK 自动轮换。
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {hasChanges && (
        <div className="fixed bottom-6 right-6">
          <Button
            onClick={handleSave}
            className="bg-gold text-white hover:bg-gold-light hover:shadow-gold-glow shadow-lg"
            size="lg"
          >
            保存设置
          </Button>
        </div>
      )}
    </div>
  );
}
