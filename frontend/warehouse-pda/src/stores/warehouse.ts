import { defineStore } from "pinia";
import { ref } from "vue";
import { apiClient } from "@erp/shared";

export interface PickTask {
  id: string;
  outbound_id: string;
  sku_id: string;
  sku_code: string;
  sku_name: string;
  quantity: number;
  picked_quantity: number;
  location_code: string;
  status: string;
}

export interface OutboundOrder {
  id: string;
  order_id: string;
  order_no: string;
  warehouse_id: string;
  status: string;
  items?: Array<{ sku_id: string; sku_name: string; quantity: number }>;
}

export const useWarehouseStore = defineStore("warehouse-pda", () => {
  const outbounds = ref<OutboundOrder[]>([]);
  const pickTasks = ref<PickTask[]>([]);
  const loading = ref(false);

  async function fetchOutbounds() {
    loading.value = true;
    try {
      const res = await apiClient.get("/warehouse/outbounds", {
        params: { page: 1, page_size: 50 },
      });
      const data = res.data?.data ?? res.data;
      outbounds.value = (data.list ?? data.items ?? []).filter(
        (o: OutboundOrder) => o.status !== "shipped"
      );
    } finally {
      loading.value = false;
    }
  }

  async function fetchPickTasks(outboundId: string) {
    loading.value = true;
    try {
      const res = await apiClient.get("/warehouse/pick-tasks", {
        params: { outbound_id: outboundId },
      });
      pickTasks.value = res.data?.data ?? res.data ?? [];
    } finally {
      loading.value = false;
    }
  }

  async function pickScan(taskId: string, quantity: number) {
    await apiClient.post("/warehouse/pick/scan", {
      task_id: taskId,
      quantity,
    });
  }

  async function checkScan(outboundId: string, skuId: string, quantity: number) {
    await apiClient.post("/warehouse/check/scan", {
      outbound_id: outboundId,
      sku_id: skuId,
      quantity,
    });
  }

  async function pack(outboundId: string, weight = 0) {
    await apiClient.post("/warehouse/package", {
      outbound_id: outboundId,
      weight,
    });
  }

  async function weigh(outboundId: string, weight: number) {
    await apiClient.post("/warehouse/weigh", {
      outbound_id: outboundId,
      weight,
    });
  }

  async function confirmShip(outboundId: string, trackingNo = "", carrier = "") {
    await apiClient.post(`/warehouse/outbounds/${outboundId}/ship`, {
      tracking_no: trackingNo,
      carrier,
    });
  }

  return {
    outbounds,
    pickTasks,
    loading,
    fetchOutbounds,
    fetchPickTasks,
    pickScan,
    checkScan,
    pack,
    weigh,
    confirmShip,
  };
});
