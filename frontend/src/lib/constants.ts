// 订单状态
export const ORDER_STATUS = {
  UNPAID: 0,
  PAID: 1,
  PREPARING: 2,
  COMPLETED: 3,
  CANCELLED: 4,
} as const;

export const ORDER_STATUS_LABELS: Record<number, string> = {
  [ORDER_STATUS.UNPAID]: '待付款',
  [ORDER_STATUS.PAID]: '已付款',
  [ORDER_STATUS.PREPARING]: '制作中',
  [ORDER_STATUS.COMPLETED]: '已完成',
  [ORDER_STATUS.CANCELLED]: '已取消',
};

// 评价状态
export const REVIEW_STATUS = {
  PENDING: 0,
  APPROVED: 1,
  REJECTED: 2,
} as const;

export const REVIEW_STATUS_LABELS: Record<number, string> = {
  [REVIEW_STATUS.PENDING]: '待审核',
  [REVIEW_STATUS.APPROVED]: '已通过',
  [REVIEW_STATUS.REJECTED]: '已拒绝',
};

// 用户状态
export const USER_STATUS = {
  DISABLED: 0,
  ACTIVE: 1,
} as const;

// 提现状态
export const WITHDRAWAL_STATUS = {
  PENDING: 0,
  APPROVED: 1,
  REJECTED: 2,
} as const;
