import { defineStore } from "pinia";
import { ref } from "vue";
import { apiClient } from "@erp/shared";

export interface OutboundOrder {
  id: string;
  orderNo: string;
  tenantId: string;
  warehouseId: string;
  storeId: string;
  storeName: string;
  recipientName: string;
  recipientPhone: string;
  recipientAddress: string;
  status: string;
  pickerId: string;
  checkerId: string;
  packerId: string;
  items: OutboundItem[];
  createdAt: string;
  updatedAt: string;
}

export interface OutboundItem {
  skuId: string;
  skuCode: string;
  skuName: string;
  quantity: number;
  pickedQty: number;
  checkedQty: number;
  packedQty: number;
}

export const useWarehouseStore = defineStore("warehouse", () => {
  const outbounds = ref<OutboundOrder[]>([]);
  const loading = ref(false);
  const total = ref(0);

  async function fetchOutbounds(params?: Record<string, any>) {
    loading.value = true;
    try {
      const res = await apiClient.get("/warehouse/outbounds", { params });
      const data = res.data?.data ?? res.data;
      outbounds.value = data.list ?? data.items ?? [];
      total.value = data.total ?? 0;
    } finally {
      loading.value = false;
    }
  }

  async function getOutboundDetail(id: string): Promise<OutboundOrder | null> {
    try {
      const res = await apiClient.get(`/warehouse/outbounds/${id}`);
      return res.data?.data ?? res.data;
    } catch {
      return null;
    }
  }

  async function createOutbound(data: Record<string, any>) {
    const res = await apiClient.post("/warehouse/outbounds", data);
    await fetchOutbounds();
    return res.data?.data ?? res.data;
  }

  return {
    outbounds,
    loading,
    total,
    fetchOutbounds,
    getOutboundDetail,
    createOutbound,
  };
});
