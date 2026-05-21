import { defineStore } from "pinia";
import { ref } from "vue";
import apiClient from "@erp/shared";
import type { ApiResponse, PageData } from "@erp/shared";

interface InventoryBalance {
  sku_id: string;
  sku_code: string;
  sku_name: string;
  warehouse_id: string;
  warehouse_name: string;
  qty: number;
  locked_qty: number;
  available_qty: number;
}

interface InventoryJournal {
  journal_id: string;
  sku_id: string;
  sku_name: string;
  change_qty: number;
  before_qty: number;
  after_qty: number;
  biz_type: string;
  biz_no: string;
  created_at: string;
}

export const useInventoryStore = defineStore("inventory", () => {
  const balances = ref<InventoryBalance[]>([]);
  const journals = ref<InventoryJournal[]>([]);
  const loading = ref(false);

  async function fetchBalances(params?: Record<string, unknown>) {
    loading.value = true;
    try {
      const res = await apiClient.get<ApiResponse<PageData<InventoryBalance>>>(
        "/inventory/balances",
        { params }
      );
      balances.value = res.data.data.list;
    } finally {
      loading.value = false;
    }
  }

  async function fetchJournals(params?: Record<string, unknown>) {
    loading.value = true;
    try {
      const res = await apiClient.get<ApiResponse<PageData<InventoryJournal>>>(
        "/inventory/journals",
        { params }
      );
      journals.value = res.data.data.list;
    } finally {
      loading.value = false;
    }
  }

  async function lockInventory(skuId: string, qty: number, orderId: string) {
    const res = await apiClient.post<ApiResponse<InventoryBalance>>(
      "/inventory/lock",
      { sku_id: skuId, qty, order_id: orderId }
    );
    return res.data.data;
  }

  return {
    balances,
    journals,
    loading,
    fetchBalances,
    fetchJournals,
    lockInventory,
  };
});
