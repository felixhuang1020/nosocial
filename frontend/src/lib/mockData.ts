import type {
  User, Shareholder, Drink, DrinkCategory, Order, Review,
  TarotCard, Banner, DashboardStats, RevenueDataPoint, SalesRankItem, CommissionRecord
} from '@/types';

// ===== Dashboard Stats =====
export const dashboardStats: DashboardStats = {
  todayRevenue: 12580,
  todayOrders: 48,
  newShareholders: 5,
  pendingReviews: 8,
  revenueTrend: 8.5,
  ordersTrend: 12,
};

export const revenueData: RevenueDataPoint[] = [
  { date: '04-14', amount: 8200 },
  { date: '04-15', amount: 9600 },
  { date: '04-16', amount: 11200 },
  { date: '04-17', amount: 8900 },
  { date: '04-18', amount: 13400 },
  { date: '04-19', amount: 10800 },
  { date: '04-20', amount: 12580 },
];

export const salesRankData: SalesRankItem[] = [
  { name: '迷雾森林', sales: 156 },
  { name: '午夜飞行', sales: 128 },
  { name: '琥珀时光', sales: 98 },
  { name: '深海之蓝', sales: 87 },
  { name: '烈焰红唇', sales: 76 },
];

// ===== Users =====
export const mockUsers: User[] = [
  { id: 1, openid: 'ox_001', nickname: '酒友小王', avatar: '', phone: '138****1234', birthday: '1995-06-15', isShareholder: true, shareholderLevel: 1, registerFeePaid: true, freeDrinkUsed: true, inviteCode: 'A1B2C3', parentId: null, balance: 286.50, totalEarning: 520.00, status: 1, createdAt: '2025-01-15 10:30:00' },
  { id: 2, openid: 'ox_002', nickname: '微醺女孩', avatar: '', phone: '139****5678', birthday: '1998-03-22', isShareholder: false, shareholderLevel: 0, registerFeePaid: false, freeDrinkUsed: false, inviteCode: '', parentId: null, balance: 0, totalEarning: 0, status: 1, createdAt: '2025-02-01 14:20:00' },
  { id: 3, openid: 'ox_003', nickname: 'NightOwl', avatar: '', phone: '137****9012', birthday: '1992-11-08', isShareholder: true, shareholderLevel: 1, registerFeePaid: true, freeDrinkUsed: false, inviteCode: 'D4E5F6', parentId: 1, balance: 152.00, totalEarning: 310.50, status: 1, createdAt: '2025-01-20 09:15:00' },
  { id: 4, openid: 'ox_004', nickname: 'cocktail_lover', avatar: '', phone: '136****3456', birthday: '2000-01-01', isShareholder: false, shareholderLevel: 0, registerFeePaid: false, freeDrinkUsed: false, inviteCode: '', parentId: null, balance: 0, totalEarning: 0, status: 1, createdAt: '2025-02-10 20:45:00' },
  { id: 5, openid: 'ox_005', nickname: '老张', avatar: '', phone: '135****7890', birthday: '1988-09-30', isShareholder: true, shareholderLevel: 1, registerFeePaid: true, freeDrinkUsed: true, inviteCode: 'G7H8I9', parentId: 1, balance: 450.80, totalEarning: 890.00, status: 1, createdAt: '2025-01-08 16:00:00' },
  { id: 6, openid: 'ox_006', nickname: 'Mojito', avatar: '', phone: '134****2468', birthday: '1997-07-18', isShareholder: false, shareholderLevel: 0, registerFeePaid: false, freeDrinkUsed: false, inviteCode: '', parentId: null, balance: 0, totalEarning: 0, status: 0, createdAt: '2025-03-01 11:30:00' },
  { id: 7, openid: 'ox_007', nickname: '酒鬼阿杰', avatar: '', phone: '133****1357', birthday: '1993-12-25', isShareholder: true, shareholderLevel: 1, registerFeePaid: true, freeDrinkUsed: true, inviteCode: 'J1K2L3', parentId: 3, balance: 89.20, totalEarning: 175.00, status: 1, createdAt: '2025-02-18 22:00:00' },
  { id: 8, openid: 'ox_008', nickname: '小酒窝', avatar: '', phone: '132****9753', birthday: '1999-05-20', isShareholder: false, shareholderLevel: 0, registerFeePaid: false, freeDrinkUsed: false, inviteCode: '', parentId: null, balance: 0, totalEarning: 0, status: 1, createdAt: '2025-03-05 15:45:00' },
];

// ===== Shareholders =====
export const mockShareholders: Shareholder[] = [
  { id: 1, userId: 1, nickname: '酒友小王', avatar: '', inviteCode: 'A1B2C3', teamCount: 12, totalEarning: 520.00, balance: 286.50, status: 1, createdAt: '2025-01-15 10:30:00' },
  { id: 2, userId: 3, nickname: 'NightOwl', avatar: '', inviteCode: 'D4E5F6', teamCount: 8, totalEarning: 310.50, balance: 152.00, status: 1, createdAt: '2025-01-20 09:15:00' },
  { id: 3, userId: 5, nickname: '老张', avatar: '', inviteCode: 'G7H8I9', teamCount: 15, totalEarning: 890.00, balance: 450.80, status: 1, createdAt: '2025-01-08 16:00:00' },
  { id: 4, userId: 7, nickname: '酒鬼阿杰', avatar: '', inviteCode: 'J1K2L3', teamCount: 5, totalEarning: 175.00, balance: 89.20, status: 1, createdAt: '2025-02-18 22:00:00' },
  { id: 5, userId: 9, nickname: 'Lisa Chen', avatar: '', inviteCode: 'M4N5O6', teamCount: 3, totalEarning: 98.50, balance: 45.00, status: 1, createdAt: '2025-03-10 18:30:00' },
];

export const mockCommissionRecords: CommissionRecord[] = [
  { id: 1, shareholderId: 1, consumerName: '微醺女孩', orderId: 101, orderAmount: 268.00, commissionRate: 0.10, commissionAmount: 26.80, status: 1, createdAt: '2025-04-20 21:30:00' },
  { id: 2, shareholderId: 1, consumerName: 'cocktail_lover', orderId: 102, orderAmount: 156.00, commissionRate: 0.10, commissionAmount: 15.60, status: 1, createdAt: '2025-04-19 20:15:00' },
  { id: 3, shareholderId: 3, consumerName: '小酒窝', orderId: 103, orderAmount: 320.00, commissionRate: 0.10, commissionAmount: 32.00, status: 0, createdAt: '2025-04-20 22:00:00' },
  { id: 4, shareholderId: 3, consumerName: 'Mojito', orderId: 104, orderAmount: 88.00, commissionRate: 0.10, commissionAmount: 8.80, status: 1, createdAt: '2025-04-18 19:45:00' },
  { id: 5, shareholderId: 2, consumerName: '老张', orderId: 105, orderAmount: 450.00, commissionRate: 0.10, commissionAmount: 45.00, status: 1, createdAt: '2025-04-17 21:00:00' },
];

// ===== Drink Categories =====
export const mockCategories: DrinkCategory[] = [
  { id: 1, name: '鸡尾酒', icon: 'wine', sortOrder: 1, status: 1 },
  { id: 2, name: '威士忌', icon: 'glass-water', sortOrder: 2, status: 1 },
  { id: 3, name: '精酿啤酒', icon: 'beer', sortOrder: 3, status: 1 },
  { id: 4, name: '伏特加', icon: 'flame', sortOrder: 4, status: 1 },
  { id: 5, name: '特调', icon: 'sparkles', sortOrder: 5, status: 1 },
];

// ===== Drinks =====
export const mockDrinks: Drink[] = [
  { id: 1, name: '迷雾森林', englishName: 'Misty Forest', categoryId: 1, categoryName: '鸡尾酒', price: 88.00, costPrice: 25.00, alcohol: 25.5, volume: '350ml', description: '以烟熏威士忌为基底的创新鸡尾酒，融入松针与蜂蜜的风味', ingredients: '烟熏威士忌、松针糖浆、蜂蜜、柠檬汁', imageUrl: '', isRecommended: true, isFreeDrink: true, status: 1, sortOrder: 1 },
  { id: 2, name: '午夜飞行', englishName: 'Midnight Flight', categoryId: 1, categoryName: '鸡尾酒', price: 98.00, costPrice: 30.00, alcohol: 28.0, volume: '350ml', description: '深蓝紫色调的神秘鸡尾酒，口感层次丰富', ingredients: '金酒、紫罗兰利口酒、柠檬汁、蝶豆花茶', imageUrl: '', isRecommended: true, isFreeDrink: false, status: 1, sortOrder: 2 },
  { id: 3, name: '琥珀时光', englishName: 'Amber Time', categoryId: 2, categoryName: '威士忌', price: 128.00, costPrice: 45.00, alcohol: 40.0, volume: '45ml', description: '陈酿12年单一麦芽威士忌，琥珀色泽，口感醇厚', ingredients: '12年单一麦芽威士忌', imageUrl: '', isRecommended: true, isFreeDrink: false, status: 1, sortOrder: 3 },
  { id: 4, name: '深海之蓝', englishName: 'Deep Blue', categoryId: 1, categoryName: '鸡尾酒', price: 78.00, costPrice: 22.00, alcohol: 18.0, volume: '350ml', description: '清爽海洋风味的蓝色鸡尾酒', ingredients: '伏特加、蓝柑橘、椰奶、菠萝汁', imageUrl: '', isRecommended: false, isFreeDrink: false, status: 1, sortOrder: 4 },
  { id: 5, name: '烈焰红唇', englishName: 'Red Lips', categoryId: 1, categoryName: '鸡尾酒', price: 85.00, costPrice: 24.00, alcohol: 22.0, volume: '350ml', description: '热情似火的红色鸡尾酒，微辣口感', ingredients: '龙舌兰、辣椒酱、石榴糖浆、青柠汁', imageUrl: '', isRecommended: false, isFreeDrink: false, status: 1, sortOrder: 5 },
  { id: 6, name: '山崎1923', englishName: 'Yamazaki 1923', categoryId: 2, categoryName: '威士忌', price: 168.00, costPrice: 85.00, alcohol: 43.0, volume: '45ml', description: '日本经典威士忌，果香与橡木桶香的完美平衡', ingredients: '山崎1923威士忌', imageUrl: '', isRecommended: true, isFreeDrink: false, status: 1, sortOrder: 6 },
  { id: 7, name: 'IPA精酿', englishName: 'Craft IPA', categoryId: 3, categoryName: '精酿啤酒', price: 58.00, costPrice: 20.00, alcohol: 6.5, volume: '500ml', description: '酒花香气浓郁的印度淡色艾尔', ingredients: '水、大麦芽、啤酒花、酵母', imageUrl: '', isRecommended: false, isFreeDrink: false, status: 1, sortOrder: 7 },
  { id: 8, name: '绝对伏特加', englishName: 'Absolut Vodka', categoryId: 4, categoryName: '伏特加', price: 68.00, costPrice: 28.00, alcohol: 40.0, volume: '45ml', description: '瑞典经典伏特加，纯净顺滑', ingredients: '伏特加', imageUrl: '', isRecommended: false, isFreeDrink: false, status: 0, sortOrder: 8 },
];

// ===== Orders =====
export const mockOrders: Order[] = [
  { id: 1, orderNo: 'NO20250420001', userId: 2, userName: '微醺女孩', shareholderId: 1, shareholderName: '酒友小王', totalAmount: 268.00, discountAmount: 0, payAmount: 268.00, couponId: null, status: 3, payTime: '2025-04-20 21:30:00', createdAt: '2025-04-20 21:25:00' },
  { id: 2, orderNo: 'NO20250420002', userId: 4, userName: 'cocktail_lover', shareholderId: 1, shareholderName: '酒友小王', totalAmount: 156.00, discountAmount: 20.00, payAmount: 136.00, couponId: 1, status: 3, payTime: '2025-04-20 20:15:00', createdAt: '2025-04-20 20:10:00' },
  { id: 3, orderNo: 'NO20250420003', userId: 8, userName: '小酒窝', shareholderId: 3, shareholderName: '老张', totalAmount: 320.00, discountAmount: 0, payAmount: 320.00, couponId: null, status: 2, payTime: '2025-04-20 22:00:00', createdAt: '2025-04-20 21:55:00' },
  { id: 4, orderNo: 'NO20250419001', userId: 6, userName: 'Mojito', shareholderId: 3, shareholderName: '老张', totalAmount: 88.00, discountAmount: 0, payAmount: 88.00, couponId: null, status: 3, payTime: '2025-04-19 19:45:00', createdAt: '2025-04-19 19:40:00' },
  { id: 5, orderNo: 'NO20250419002', userId: 5, userName: '老张', shareholderId: 2, shareholderName: 'NightOwl', totalAmount: 450.00, discountAmount: 50.00, payAmount: 400.00, couponId: 2, status: 3, payTime: '2025-04-19 21:00:00', createdAt: '2025-04-19 20:55:00' },
  { id: 6, orderNo: 'NO20250420004', userId: 2, userName: '微醺女孩', shareholderId: null, shareholderName: null, totalAmount: 176.00, discountAmount: 0, payAmount: 176.00, couponId: null, status: 1, payTime: '2025-04-20 23:00:00', createdAt: '2025-04-20 22:55:00' },
  { id: 7, orderNo: 'NO20250420005', userId: 7, userName: '酒鬼阿杰', shareholderId: 1, shareholderName: '酒友小王', totalAmount: 98.00, discountAmount: 0, payAmount: 98.00, couponId: null, status: 0, payTime: null, createdAt: '2025-04-20 23:30:00' },
  { id: 8, orderNo: 'NO20250418001', userId: 3, userName: 'NightOwl', shareholderId: null, shareholderName: null, totalAmount: 216.00, discountAmount: 0, payAmount: 216.00, couponId: null, status: 4, payTime: null, createdAt: '2025-04-18 20:00:00' },
];

// ===== Reviews =====
export const mockReviews: Review[] = [
  { id: 1, userId: 2, userName: '微醺女孩', userAvatar: '', dianpingUrl: 'https://dianping.com', screenshotUrl: 'https://images.unsplash.com/photo-1566417713940-fe7c737a9ef2?w=400&h=600&fit=crop', reviewContent: '环境超级棒！调酒师很专业，推荐迷雾森林，口感层次很丰富。下次还会来！', status: 0, rejectReason: null, couponId: null, auditTime: null, createdAt: '2025-04-20 18:30:00' },
  { id: 2, userId: 4, userName: 'cocktail_lover', userAvatar: '', dianpingUrl: 'https://dianping.com', screenshotUrl: 'https://images.unsplash.com/photo-1551024709-8f23befc6f87?w=400&h=600&fit=crop', reviewContent: '酒很好喝，氛围也很棒，适合约会', status: 0, rejectReason: null, couponId: null, auditTime: null, createdAt: '2025-04-20 17:15:00' },
  { id: 3, userId: 8, userName: '小酒窝', userAvatar: '', dianpingUrl: 'https://dianping.com', screenshotUrl: 'https://images.unsplash.com/photo-1470337458703-46ad1756a187?w=400&h=600&fit=crop', reviewContent: '第一次来就被惊艳到了，酒单设计很有品味', status: 0, rejectReason: null, couponId: null, auditTime: null, createdAt: '2025-04-19 22:00:00' },
  { id: 4, userId: 6, userName: 'Mojito', userAvatar: '', dianpingUrl: 'https://dianping.com', screenshotUrl: 'https://images.unsplash.com/photo-1514362545857-3bc16c4c7d1b?w=400&h=600&fit=crop', reviewContent: '不错的体验', status: 1, rejectReason: null, couponId: 1, auditTime: '2025-04-19 10:00:00', createdAt: '2025-04-18 21:30:00' },
  { id: 5, userId: 3, userName: 'NightOwl', userAvatar: '', dianpingUrl: 'https://dianping.com', screenshotUrl: 'https://images.unsplash.com/photo-1572116469696-31de0f17cc34?w=400&h=600&fit=crop', reviewContent: '已经是第三次来了，每次都有惊喜', status: 1, rejectReason: null, couponId: 2, auditTime: '2025-04-18 09:30:00', createdAt: '2025-04-17 23:00:00' },
  { id: 6, userId: 7, userName: '酒鬼阿杰', userAvatar: '', dianpingUrl: 'https://dianping.com', screenshotUrl: 'https://images.unsplash.com/photo-1560508179-b2c9a3f8e92b?w=400&h=600&fit=crop', reviewContent: '截图不太清楚', status: 2, rejectReason: '截图无法识别，请重新上传', couponId: null, auditTime: '2025-04-17 11:00:00', createdAt: '2025-04-16 20:00:00' },
];

// ===== Tarot Cards =====
const arcanaMajor = [
  '愚者', '魔术师', '女祭司', '皇后', '皇帝', '教皇', '恋人', '战车',
  '力量', '隐者', '命运之轮', '正义', '倒吊人', '死神', '节制', '恶魔',
  '塔', '星星', '月亮', '太阳', '审判', '世界'
];
const arcanaMajorEn = [
  'The Fool', 'The Magician', 'The High Priestess', 'The Empress', 'The Emperor',
  'The Hierophant', 'The Lovers', 'The Chariot', 'Strength', 'The Hermit',
  'Wheel of Fortune', 'Justice', 'The Hanged Man', 'Death', 'Temperance',
  'The Devil', 'The Tower', 'The Star', 'The Moon', 'The Sun',
  'Judgement', 'The World'
];

export const mockTarotCards: TarotCard[] = [
  ...arcanaMajor.map((name, i) => ({
    id: i + 1,
    cardNo: i,
    name,
    nameEn: arcanaMajorEn[i],
    arcanaType: 1,
    suit: null,
    imageUrl: '',
    uprightMeaning: `${name}正位含义...`,
    reversedMeaning: `${name}逆位含义...`,
    keywords: '关键词',
    element: ['火', '水', '风', '土'][i % 4],
  })),
  // 小阿卡纳简化为各花色14张
  ...Array.from({ length: 14 }, (_, i) => ({
    id: 23 + i,
    cardNo: 22 + i,
    name: `权杖${i < 10 ? i + 1 : ['侍从', '骑士', '皇后', '国王'][i - 10]}`,
    nameEn: `Wands ${i < 10 ? i + 1 : ['Page', 'Knight', 'Queen', 'King'][i - 10]}`,
    arcanaType: 2,
    suit: '权杖',
    imageUrl: '',
    uprightMeaning: '权杖正位含义...',
    reversedMeaning: '权杖逆位含义...',
    keywords: '创造力、行动',
    element: '火',
  })),
  ...Array.from({ length: 14 }, (_, i) => ({
    id: 37 + i,
    cardNo: 36 + i,
    name: `圣杯${i < 10 ? i + 1 : ['侍从', '骑士', '皇后', '国王'][i - 10]}`,
    nameEn: `Cups ${i < 10 ? i + 1 : ['Page', 'Knight', 'Queen', 'King'][i - 10]}`,
    arcanaType: 2,
    suit: '圣杯',
    imageUrl: '',
    uprightMeaning: '圣杯正位含义...',
    reversedMeaning: '圣杯逆位含义...',
    keywords: '情感、直觉',
    element: '水',
  })),
  ...Array.from({ length: 14 }, (_, i) => ({
    id: 51 + i,
    cardNo: 50 + i,
    name: `宝剑${i < 10 ? i + 1 : ['侍从', '骑士', '皇后', '国王'][i - 10]}`,
    nameEn: `Swords ${i < 10 ? i + 1 : ['Page', 'Knight', 'Queen', 'King'][i - 10]}`,
    arcanaType: 2,
    suit: '宝剑',
    imageUrl: '',
    uprightMeaning: '宝剑正位含义...',
    reversedMeaning: '宝剑逆位含义...',
    keywords: '智慧、冲突',
    element: '风',
  })),
  ...Array.from({ length: 14 }, (_, i) => ({
    id: 65 + i,
    cardNo: 64 + i,
    name: `星币${i < 10 ? i + 1 : ['侍从', '骑士', '皇后', '国王'][i - 10]}`,
    nameEn: `Pentacles ${i < 10 ? i + 1 : ['Page', 'Knight', 'Queen', 'King'][i - 10]}`,
    arcanaType: 2,
    suit: '星币',
    imageUrl: '',
    uprightMeaning: '星币正位含义...',
    reversedMeaning: '星币逆位含义...',
    keywords: '物质、财富',
    element: '土',
  })),
];

// ===== Banners =====
export const mockBanners: Banner[] = [
  { id: 1, title: '春季特调上新', imageUrl: 'https://images.unsplash.com/photo-1551024709-8f23befc6f87?w=800&h=400&fit=crop', linkType: 0, linkValue: null, position: 1, sortOrder: 1, status: 1 },
  { id: 2, title: '共享股东招募中', imageUrl: 'https://images.unsplash.com/photo-1572116469696-31de0f17cc34?w=800&h=400&fit=crop', linkType: 0, linkValue: null, position: 1, sortOrder: 2, status: 1 },
  { id: 3, title: '酒吧环境展示', imageUrl: 'https://images.unsplash.com/photo-1566417713940-fe7c737a9ef2?w=600&h=800&fit=crop', linkType: 0, linkValue: null, position: 2, sortOrder: 1, status: 1 },
  { id: 4, title: '调酒师特写', imageUrl: 'https://images.unsplash.com/photo-1514362545857-3bc16c4c7d1b?w=600&h=800&fit=crop', linkType: 0, linkValue: null, position: 2, sortOrder: 2, status: 1 },
  { id: 5, title: 'VIP包厢', imageUrl: 'https://images.unsplash.com/photo-1470337458703-46ad1756a187?w=600&h=800&fit=crop', linkType: 0, linkValue: null, position: 2, sortOrder: 3, status: 1 },
];

// ===== System Settings =====
export const mockSettings = {
  shareholderFee: 99.00,
  commissionRate: 0.10,
  freeDrinkId: 1,
  storeName: 'NoSocial Bar',
  storeAddress: '上海市静安区南京西路1266号',
  storePhone: '021-6288-8888',
  businessHours: '18:00 - 04:00',
  wxAppId: 'wx1234567890',
  wxMchId: '1234567890',
};
