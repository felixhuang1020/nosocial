/**
 * 常量定义
 */

// 股东等级
const SHAREHOLDER_LEVELS = {
  0: { name: '普通用户', color: '#8a8a8a' },
  1: { name: '初级股东', color: '#D4AF37' },
  2: { name: '高级股东', color: '#FF6B6B' }
};

// 订单状态
const ORDER_STATUS = {
  0: { name: '待支付', color: '#ffc107' },
  1: { name: '已支付', color: '#1890ff' },
  2: { name: '制作中', color: '#722ed1' },
  3: { name: '已完成', color: '#52c41a' },
  4: { name: '已取消', color: '#ff4d4f' }
};

// 佣金状态
const COMMISSION_STATUS = {
  0: { name: '待结算', color: '#ffc107' },
  1: { name: '已结算', color: '#52c41a' },
  2: { name: '已提现', color: '#1890ff' }
};

// 优惠券类型
const COUPON_TYPES = {
  1: { name: '点评返券', color: '#D4AF37' },
  2: { name: '生日券', color: '#ff6b6b' },
  3: { name: '活动券', color: '#1890ff' }
};

// 优惠券状态
const COUPON_STATUS = {
  0: { name: '未使用', color: '#52c41a' },
  1: { name: '已使用', color: '#8a8a8a' },
  2: { name: '已过期', color: '#ff4d4f' }
};

// 点评审核状态
const REVIEW_STATUS = {
  0: { name: '待审核', color: '#ffc107' },
  1: { name: '已通过', color: '#52c41a' },
  2: { name: '已拒绝', color: '#ff4d4f' }
};

// 支付方式
const PAY_TYPES = {
  WX_PAY: 'wx_pay',
  BALANCE: 'balance'
};

// 存储键名
const STORAGE_KEYS = {
  TOKEN: 'token',
  USER_INFO: 'user_info',
  INVITE_CODE: 'invite_code',
  SETTINGS: 'settings'
};

// 页面路径
const PAGES = {
  HOME: '/pages/index/index',
  SHOWCASE: '/pages/showcase/showcase',
  MENU: '/pages/menu/menu',
  PROFILE: '/pages/profile/profile',
  SHAREHOLDER_APPLY: '/pages/shareholder/apply',
  SHAREHOLDER_EARNINGS: '/pages/shareholder/earnings',
  SHAREHOLDER_TEAM: '/pages/shareholder/team',
  TAROT_DIVINE: '/pages/tarot/divine',
  TAROT_RESULT: '/pages/tarot/result',
  TAROT_HISTORY: '/pages/tarot/history',
  BIRTHDAY_GIFT: '/pages/birthday/gift',
  REVIEW_SUBMIT: '/pages/review/submit',
  REVIEW_GUIDE: '/pages/review/guide',
  DRINK_DETAIL: '/pages/drink/detail',
  COUPON_LIST: '/pages/coupon/list',
  WEBVIEW: '/pages/webview/webview'
};

module.exports = {
  SHAREHOLDER_LEVELS,
  ORDER_STATUS,
  COMMISSION_STATUS,
  COUPON_TYPES,
  COUPON_STATUS,
  REVIEW_STATUS,
  PAY_TYPES,
  STORAGE_KEYS,
  PAGES
};
