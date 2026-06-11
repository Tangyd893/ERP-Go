import { ref, computed } from "vue";
import { apiClient, isDemo } from "@erp/shared";

interface OrderSummary {
  total: number;
  today: number;
  amount: number;
}

interface InventorySummary {
  skuCount: number;
  totalQty: number;
}

interface OutboundSummary {
  total: number;
  today: number;
  shipped: number;
}

export function useDashboard() {
  const orders = ref<OrderSummary>({ total: 0, today: 0, amount: 0 });
  const inventory = ref<InventorySummary>({ skuCount: 0, totalQty: 0 });
  const outbounds = ref<OutboundSummary>({ total: 0, today: 0, shipped: 0 });
  const loading = ref(false);
  const error = ref("");

  async function fetchAll() {
    if (isDemo()) {
      orders.value = { total: 328, today: 45, amount: 12850 };
      inventory.value = { skuCount: 1256, totalQty: 48520 };
      outbounds.value = { total: 186, today: 23, shipped: 152 };
      return;
    }

    loading.value = true;
    error.value = "";
    try {
      const [orderRes, inventoryRes, outboundRes] = await Promise.allSettled([
        apiClient.get("/order/orders", { params: { page: 1, page_size: 1 } }),
        apiClient.get("/inventory/balances", { params: { page: 1, page_size: 1 } }),
        apiClient.get("/warehouse/outbounds", { params: { page: 1, page_size: 1 } }),
      ]);

      if (orderRes.status === "fulfilled") {
        const data = orderRes.value.data?.data ?? orderRes.value.data ?? {};
        orders.value = {
          total: data.total ?? 0,
          today: 0,
          amount: 0,
        };
      }
      if (inventoryRes.status === "fulfilled") {
        const data = inventoryRes.value.data?.data ?? inventoryRes.value.data ?? {};
        inventory.value = {
          skuCount: data.total ?? 0,
          totalQty: 0,
        };
      }
      if (outboundRes.status === "fulfilled") {
        const data = outboundRes.value.data?.data ?? outboundRes.value.data ?? {};
        outbounds.value = {
          total: data.total ?? 0,
          today: 0,
          shipped: 0,
        };
      }
    } catch {
      error.value = "data load failed";
    } finally {
      loading.value = false;
    }
  }

  const orderCount = computed(() => orders.value.total);
  const salesAmount = computed(() => orders.value.amount);
  const outboundCount = computed(() => outbounds.value.total);
  const skuCount = computed(() => inventory.value.skuCount);

  return {
    orders,
    inventory,
    outbounds,
    loading,
    error,
    orderCount,
    salesAmount,
    outboundCount,
    skuCount,
    fetchAll,
  };
}

// 单例 store
let instance: ReturnType<typeof useDashboard> | null = null;
export function useDashboardStore() {
  if (!instance) {
    instance = useDashboard();
  }
  return instance;
}
