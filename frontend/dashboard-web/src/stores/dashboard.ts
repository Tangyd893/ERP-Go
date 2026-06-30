import { ref, computed, onUnmounted } from "vue";
import { apiClient, isDemo } from "@erp/shared";

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

interface DashboardResponse {
  order_count: number;
  sales_amount: number;
  outbound_count: number;
  sku_count: number;
  trend: TrendPoint[];
  timeliness: TimelinessRate;
}

export function useDashboard() {
  const orderCount = ref(0);
  const salesAmount = ref(0);
  const outboundCount = ref(0);
  const skuCount = ref(0);
  const trend = ref<TrendPoint[]>([]);
  const timeliness = ref<TimelinessRate>({ within_24h: 0, within_48h: 0, overdue: 0 });
  const loading = ref(false);
  const error = ref("");
  const authRequired = ref(false);

  let pollTimer: ReturnType<typeof setInterval> | null = null;

  /** 硬编码 demo 数据 */
  function applyDemo() {
    orderCount.value = 328;
    salesAmount.value = 12850;
    outboundCount.value = 186;
    skuCount.value = 1256;
    trend.value = [
      { date: "周一", order_count: 45, sales_amount: 2100 },
      { date: "周二", order_count: 52, sales_amount: 3200 },
      { date: "周三", order_count: 38, sales_amount: 1850 },
      { date: "周四", order_count: 65, sales_amount: 4200 },
      { date: "周五", order_count: 48, sales_amount: 2800 },
      { date: "周六", order_count: 32, sales_amount: 1500 },
      { date: "周日", order_count: 28, sales_amount: 1200 },
    ];
    timeliness.value = { within_24h: 152, within_48h: 28, overdue: 6 };
    authRequired.value = false;
    error.value = "";
  }

  async function fetchAll() {
    if (isDemo()) {
      applyDemo();
      return;
    }

    loading.value = true;
    error.value = "";
    authRequired.value = false;

    try {
      const res = await apiClient.get("/report/dashboard");
      const data: DashboardResponse = res.data?.data ?? res.data;

      orderCount.value = data.order_count ?? 0;
      salesAmount.value = data.sales_amount ?? 0;
      outboundCount.value = data.outbound_count ?? 0;
      skuCount.value = data.sku_count ?? 0;
      trend.value = data.trend ?? [];
      timeliness.value = data.timeliness ?? { within_24h: 0, within_48h: 0, overdue: 0 };
    } catch (err: unknown) {
      const e = err as { response?: { status?: number } };
      if (e?.response?.status === 401) {
        authRequired.value = true;
        error.value = "请先登录管理后台获取访问令牌";
      } else {
        error.value = "数据加载失败";
      }
    } finally {
      loading.value = false;
    }
  }

  /** 启动轮询（T-638） */
  function startPolling(intervalMs = 30000) {
    stopPolling();
    pollTimer = setInterval(fetchAll, intervalMs);
  }

  /** 停止轮询 */
  function stopPolling() {
    if (pollTimer !== null) {
      clearInterval(pollTimer);
      pollTimer = null;
    }
  }

  // 组件卸载时自动清理
  onUnmounted(stopPolling);

  return {
    orderCount,
    salesAmount,
    outboundCount,
    skuCount,
    trend,
    timeliness,
    loading,
    error,
    authRequired,
    fetchAll,
    applyDemo,
    startPolling,
    stopPolling,
  };
}

// 单例 store
let instance: ReturnType<typeof useDashboard> | null = null;
export function useDashboardStore() {
  instance ??= useDashboard();
  return instance;
}
