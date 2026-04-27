export interface User {
  id: number;
  openid: string;
  nickname: string;
  avatar: string;
  phone: string;
  birthday: string;
  isShareholder: boolean;
  shareholderLevel: number;
  registerFeePaid: boolean;
  freeDrinkUsed: boolean;
  inviteCode: string;
  parentId: number | null;
  balance: number;
  totalEarning: number;
  status: number;
  createdAt: string;
}

export interface Shareholder {
  id: number;
  userId: number;
  nickname: string;
  avatar: string;
  inviteCode: string;
  teamCount: number;
  totalEarning: number;
  balance: number;
  status: number;
  createdAt: string;
}

export interface CommissionRecord {
  id: number;
  shareholderId: number;
  consumerName: string;
  orderId: number;
  orderAmount: number;
  commissionRate: number;
  commissionAmount: number;
  status: number;
  createdAt: string;
}

export interface Drink {
  id: number;
  name: string;
  englishName: string;
  categoryId: number;
  categoryName: string;
  price: number;
  costPrice: number;
  alcohol: number;
  volume: string;
  description: string;
  ingredients: string;
  imageUrl: string;
  isRecommended: boolean;
  isFreeDrink: boolean;
  status: number;
  sortOrder: number;
}

export interface DrinkCategory {
  id: number;
  name: string;
  icon: string;
  sortOrder: number;
  status: number;
}

export interface Order {
  id: number;
  orderNo: string;
  userId: number;
  userName: string;
  shareholderId: number | null;
  shareholderName: string | null;
  totalAmount: number;
  discountAmount: number;
  payAmount: number;
  couponId: number | null;
  status: number;
  payTime: string | null;
  createdAt: string;
}

export interface Review {
  id: number;
  userId: number;
  userName: string;
  userAvatar: string;
  dianpingUrl: string;
  screenshotUrl: string;
  reviewContent: string;
  status: number;
  rejectReason: string | null;
  couponId: number | null;
  auditTime: string | null;
  createdAt: string;
}

export interface TarotCard {
  id: number;
  cardNo: number;
  name: string;
  nameEn: string;
  arcanaType: number;
  suit: string | null;
  imageUrl: string;
  uprightMeaning: string;
  reversedMeaning: string;
  keywords: string;
  element: string | null;
}

export interface TarotDrinkMapping {
  id: number;
  cardNo: number;
  cardName: string;
  drinkId: number;
  drinkName: string;
  isReversed: boolean;
  matchScore: number;
  reasonTemplate: string;
}

export interface Banner {
  id: number;
  title: string;
  imageUrl: string;
  linkType: number;
  linkValue: string | null;
  position: number;
  sortOrder: number;
  status: number;
}

export interface Coupon {
  id: number;
  userId: number;
  couponNo: string;
  type: number;
  name: string;
  amount: number;
  minOrderAmount: number;
  status: number;
  validStart: string;
  validEnd: string;
}

export interface AdminUser {
  id: number;
  username: string;
  nickname: string;
  role: number;
  status: number;
  lastLoginAt: string | null;
}

export interface SystemSettings {
  shareholderFee: number;
  commissionRate: number;
  freeDrinkId: number;
  storeName: string;
  storeAddress: string;
  storePhone: string;
  businessHours: string;
  wxAppId: string;
  wxMchId: string;
}

export interface DashboardStats {
  todayRevenue: number;
  todayOrders: number;
  newShareholders: number;
  pendingReviews: number;
  revenueTrend: number;
  ordersTrend: number;
}

export interface RevenueDataPoint {
  date: string;
  amount: number;
}

export interface SalesRankItem {
  name: string;
  sales: number;
}
