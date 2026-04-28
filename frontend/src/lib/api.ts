// NoSocial Admin API Client
const API_BASE = import.meta.env.VITE_API_BASE || '/api/v1';
const REQUEST_TIMEOUT_MS = 15000;

// Helper: get stored admin token
function getToken(): string | null {
  return localStorage.getItem('admin_token');
}

// 处理 401：清掉本地凭证并引导用户回登录页
// 用 location 而不是 react-router 以防止循环依赖
function handleUnauthorized() {
  try {
    localStorage.removeItem('admin_token');
  } catch {
    // ignore
  }
  if (typeof window !== 'undefined' && !window.location.pathname.endsWith('/login')) {
    window.location.replace('/login');
  }
}

// Helper: unified fetch wrapper
async function request<T>(
  path: string,
  options: RequestInit = {}
): Promise<{ code: number; msg: string; data: T }> {
  const url = `${API_BASE}${path}`;
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options.headers as Record<string, string>),
  };

  const token = getToken();
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  // 超时控制：避免页面挂死在 hung 连接上
  const controller = new AbortController();
  const timer = setTimeout(() => controller.abort(), REQUEST_TIMEOUT_MS);

  let resp: Response;
  try {
    resp = await fetch(url, {
      ...options,
      headers,
      signal: options.signal ?? controller.signal,
    });
  } catch (err) {
    clearTimeout(timer);
    if ((err as Error).name === 'AbortError') {
      throw new Error('请求超时，请稍后重试');
    }
    throw err;
  }
  clearTimeout(timer);

  if (resp.status === 401) {
    handleUnauthorized();
    throw new Error('登录已过期，请重新登录');
  }

  if (!resp.ok) {
    // 尝试从响应体中提取错误消息
    let errorMsg = `HTTP ${resp.status}: ${resp.statusText}`;
    try {
      const body = await resp.json();
      if (body.msg) {
        errorMsg = body.msg;
      }
    } catch {
      // 响应体不是 JSON，使用默认消息
    }
    throw new Error(errorMsg);
  }

  return resp.json();
}

// ==================== Auth ====================
export interface AdminLoginReq {
  username: string;
  password: string;
}

export interface AdminLoginResp {
  token: string;
  expire: number;
  nickname: string;
  role: number;
}

export function adminLogin(req: AdminLoginReq) {
  return request<AdminLoginResp>('/admin/login', {
    method: 'POST',
    body: JSON.stringify(req),
  });
}

// ==================== Dashboard ====================
export interface DashboardStats {
  today_amount: number;
  today_order_count: number;
  new_shareholder_count: number;
  pending_review_count: number;
  revenue_trend?: { date: string; amount: number }[];
}

export function getDashboardStats() {
  return request<DashboardStats>('/admin/dashboard');
}

// ==================== Users ====================
export interface User {
  id: number;
  openid: string;
  unionid?: string;
  nickname?: string;
  avatar?: string;
  phone?: string;
  birthday?: string;
  is_shareholder: number;
  shareholder_level: number;
  shareholder_expire_at?: string;
  register_fee_paid: number;
  free_drink_used: number;
  invite_code?: string;
  parent_id?: number;
  parent_path?: string;
  balance: number;
  total_earning: number;
  status: number;
  team_count?: number;
  created_at: string;
  updated_at: string;
}

export function getUserList(page = 1, size = 20) {
  return request<{ list: User[]; total: number }>(`/admin/users?page=${page}&size=${size}`);
}

export function updateUserStatus(id: number, status: number) {
  return request<null>(`/admin/users/${id}/status`, {
    method: 'PUT',
    body: JSON.stringify({ status }),
  });
}

// ==================== Shareholders ====================
export function getShareholderList(page = 1, size = 20) {
  return request<{ list: User[]; total: number }>(`/admin/shareholders?page=${page}&size=${size}`);
}

export function getShareholderEarnings(id: number, page = 1, size = 20) {
  return request<{ list: unknown[]; total: number }>(
    `/admin/shareholders/${id}/earnings?page=${page}&size=${size}`
  );
}

export function getShareholderTeam(id: number) {
  return request<User[]>(`/admin/shareholders/${id}/team`);
}

// ==================== Drinks ====================
export interface DrinkCategory {
  id: number;
  name: string;
  sort_order: number;
  status: number;
  created_at: string;
}

export interface Drink {
  id: number;
  name: string;
  english_name: string;
  category_id: number;
  price: number;
  cost_price: number;
  alcohol: number;
  volume: string;
  description: string;
  ingredients: string;
  image_url?: string;
  is_recommended: number;
  is_free_drink: number;
  status: number;
  sort_order: number;
  created_at: string;
  updated_at: string;
  category?: DrinkCategory;
}

export function getDrinkList(page = 1, size = 20, categoryId?: number, status?: number) {
  let qs = `page=${page}&size=${size}`;
  if (categoryId !== undefined) qs += `&category_id=${categoryId}`;
  if (status !== undefined) qs += `&status=${status}`;
  return request<{ list: Drink[]; total: number }>(`/admin/drinks?${qs}`);
}

export function createDrink(drink: Partial<Drink>) {
  return request<Drink>('/admin/drinks', {
    method: 'POST',
    body: JSON.stringify(drink),
  });
}

export function updateDrink(id: number, drink: Partial<Drink>) {
  return request<Drink>(`/admin/drinks/${id}`, {
    method: 'PUT',
    body: JSON.stringify(drink),
  });
}

export function deleteDrink(id: number) {
  return request<null>(`/admin/drinks/${id}`, {
    method: 'DELETE',
  });
}

export function getCategories() {
  return request<DrinkCategory[]>('/admin/categories');
}

export function createCategory(category: Partial<DrinkCategory>) {
  return request<DrinkCategory>('/admin/categories', {
    method: 'POST',
    body: JSON.stringify(category),
  });
}

export function updateCategory(id: number, category: Partial<DrinkCategory>) {
  return request<DrinkCategory>(`/admin/categories/${id}`, {
    method: 'PUT',
    body: JSON.stringify(category),
  });
}

export function deleteCategory(id: number) {
  return request<null>(`/admin/categories/${id}`, {
    method: 'DELETE',
  });
}

// ==================== Orders ====================
export interface Order {
  id: number;
  order_no: string;
  user_id: number;
  total_amount: number;
  discount_amount: number;
  pay_amount: number;
  status: number;
  pay_time?: string;
  created_at: string;
}

export function getOrderList(page = 1, size = 20, status?: number) {
  let qs = `page=${page}&size=${size}`;
  if (status !== undefined) qs += `&status=${status}`;
  return request<{ list: Order[]; total: number }>(`/admin/orders?${qs}`);
}

export function updateOrderStatus(id: number, status: number) {
  return request<null>(`/admin/orders/${id}/status`, {
    method: 'PUT',
    body: JSON.stringify({ status }),
  });
}

// ==================== Reviews ====================
export interface Review {
  id: number;
  user_id: number;
  screenshot_url: string;
  image_urls?: string[];
  review_content: string;
  status: number;
  reject_reason?: string;
  coupon_id?: number;
  audit_time?: string;
  audit_admin_id?: number;
  created_at: string;
}

export function getReviewList(page = 1, size = 20, status?: number) {
  let qs = `page=${page}&size=${size}`;
  if (status !== undefined) qs += `&status=${status}`;
  return request<{ list: Review[]; total: number }>(`/admin/reviews?${qs}`);
}

export function auditReview(id: number, pass: boolean, rejectReason?: string) {
  return request<null>(`/admin/reviews/${id}/audit`, {
    method: 'POST',
    body: JSON.stringify({ pass, reject_reason: rejectReason }),
  });
}

// ==================== Banners ====================
export interface Banner {
  id: number;
  title: string;
  image_url: string;
  link_type: number;
  link_value?: string;
  position: number;
  sort_order: number;
  status: number;
  created_at: string;
}

export function getBannerList(page = 1, size = 20) {
  return request<{ list: Banner[]; total: number }>(`/admin/banners?page=${page}&size=${size}`);
}

export function createBanner(banner: Partial<Banner>) {
  return request<Banner>('/admin/banners', {
    method: 'POST',
    body: JSON.stringify(banner),
  });
}

export function updateBanner(id: number, banner: Partial<Banner>) {
  return request<Banner>(`/admin/banners/${id}`, {
    method: 'PUT',
    body: JSON.stringify(banner),
  });
}

export function deleteBanner(id: number) {
  return request<null>(`/admin/banners/${id}`, {
    method: 'DELETE',
  });
}

// ==================== Tarot ====================
export interface TarotMapping {
  id: number;
  card_no: number;
  card_name: string;
  drink_id: number;
  drink_name: string;
  is_reversed: number;
  match_score: number;
  reason_template: string;
}

export function getTarotMappings() {
  return request<TarotMapping[]>('/admin/tarot/mappings');
}

export function updateTarotMapping(mappings: Partial<TarotMapping>[]) {
  return request<null>('/admin/tarot/mappings', {
    method: 'PUT',
    body: JSON.stringify({ mappings }),
  });
}

// ==================== Tarot Cards (Public API) ====================
export interface TarotCard {
  id: number;
  card_no: number;
  name: string;
  name_en?: string;
  arcana_type: number;
  suit?: string;
  image_url: string;
  upright_meaning?: string;
  reversed_meaning?: string;
  keywords?: string;
  element?: string;
}

export function getPublicTarotCards() {
  return request<TarotCard[]>('/public/tarot/cards');
}

// ==================== Settings ====================
export interface AdminSettings {
  shareholder_fee?: number;
  commission_rate?: number;
  free_drink_id?: number;
  review_coupon_amount?: number;
  review_coupon_min_order?: number;
  review_coupon_valid_days?: number;
  name?: string;
  address?: string;
  phone?: string;
  business_hours?: string;
  wifi_name?: string;
  wifi_password?: string;
  wx_appid?: string;
  wx_mch_id?: string;
  wx_notify_url?: string;
  [key: string]: unknown;
}

export function getSettings() {
  return request<AdminSettings>('/admin/settings');
}

export function updateSettings(settings: Record<string, string>) {
  return request<null>('/admin/settings', {
    method: 'PUT',
    body: JSON.stringify(settings),
  });
}

// ==================== Upload ====================
export interface OSSSignature {
  host: string;
  key: string;
  policy: string;
  OSSAccessKeyId: string;
  signature: string;
  expires: string;
  endpoint: string;
  bucket: string;
}

export function getOSSSignature(dir = 'uploads', ext = '') {
  const params = new URLSearchParams({ dir });
  if (ext) params.set('ext', ext);
  return request<OSSSignature>(`/admin/upload/signature?${params.toString()}`);
}

export function fixObjectInline(key: string) {
  return request<{ key: string; status: string }>(`/admin/upload/inline?key=${encodeURIComponent(key)}`, { method: 'PUT' });
}
