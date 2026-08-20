export interface TrendPoint {
  date: string;
  order_count: number;
  sales_amount: number;
}

export interface TimelinessRate {
  within_24h: number;
  within_48h: number;
  overdue: number;
}

export interface DashboardData {
  order_count: number;
  sales_amount: number;
  outbound_count: number;
  sku_count: number;
  trend: TrendPoint[];
  timeliness: TimelinessRate;
}

/** 与 Vue stores/dashboard.ts 一致的硬编码 demo 数据 */
export const DEMO_DASHBOARD: DashboardData = {
  order_count: 328,
  sales_amount: 12850,
  outbound_count: 186,
  sku_count: 1256,
  trend: [
    { date: "周一", order_count: 45, sales_amount: 2100 },
    { date: "周二", order_count: 52, sales_amount: 3200 },
    { date: "周三", order_count: 38, sales_amount: 1850 },
    { date: "周四", order_count: 65, sales_amount: 4200 },
    { date: "周五", order_count: 48, sales_amount: 2800 },
    { date: "周六", order_count: 32, sales_amount: 1500 },
    { date: "周日", order_count: 28, sales_amount: 1200 },
  ],
  timeliness: { within_24h: 152, within_48h: 28, overdue: 6 },
};

export const FALLBACK_TREND = {
  dates: ["周一", "周二", "周三", "周四", "周五", "周六", "周日"],
  orders: [45, 52, 38, 65, 48, 32, 28],
  sales: [2100, 3200, 1850, 4200, 2800, 1500, 1200],
};

export const FALLBACK_TIMELINESS = [
  { value: 152, name: "24h内出库", color: "#67C23A" },
  { value: 28, name: "48h内出库", color: "#E6A23C" },
  { value: 6, name: "超时出库", color: "#F56C6C" },
];

export const CHART_COLORS = {
  orders: "hsl(199, 89%, 37%)",
  sales: "#67C23A",
} as const;
