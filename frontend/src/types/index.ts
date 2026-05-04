// 主要类型定义统一在 src/lib/api.ts 中管理
// 以下为仅前端使用的辅助类型（不直接对应 API 返回）

export interface RevenueDataPoint {
  date: string;
  amount: number;
}

export interface SalesRankItem {
  name: string;
  sales: number;
}
