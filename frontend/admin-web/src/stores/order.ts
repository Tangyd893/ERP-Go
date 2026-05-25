import { defineStore } from "pinia";
import { ref } from "vue";
import { apiClient } from "@erp/shared";
import type { ApiResponse, PageData } from "@erp/shared";

interface OrderRecord {
  id?: string;
  order_id?: string;
  order_no?: string;
  platform_order_no?: string;
  status: string;
  buyer_name?: string;
  customer_name?: string;
  total_amount: number;
  created_at: string;
  items?: OrderItem[];
}

interface OrderItem {
  sku_id: string;
  sku_name: string;
  qty?: number;
  quantity?: number;
  unit_price: number;
}

interface OrderQueryParams {
  page?: number;
  page_size?: number;
  status?: string;
  keyword?: string;
}

export const useOrderStore = defineStore("order", () => {
  const orders = ref<OrderRecord[]>([]);
  const total = ref(0);
  const loading = ref(false);

  async function fetchOrders(params: OrderQueryParams = {}) {
    loading.value = true;
    try {
      const res = await apiClient.get<ApiResponse<PageData<OrderRecord>>>(
        "/order/orders",
        { params }
      );
      orders.value = res.data.data.list;
      total.value = res.data.data.total;
    } finally {
      loading.value = false;
    }
  }

  async function auditOrder(orderId: string, approved: boolean) {
    const res = await apiClient.post<ApiResponse<{ approved: boolean }>>(
      "/order/orders/audit",
      { order_id: orderId, approved }
    );
    return res.data.data;
  }

  async function cancelOrder(orderId: string, reason: string) {
    const res = await apiClient.post<ApiResponse<{ cancelled: boolean }>>(
      "/order/orders/cancel",
      { order_id: orderId, reason }
    );
    return res.data.data;
  }

  return {
    orders,
    total,
    loading,
    fetchOrders,
    auditOrder,
    cancelOrder,
  };
});
