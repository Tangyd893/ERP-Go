import { defineStore } from "pinia";
import { ref } from "vue";
import apiClient from "@erp/shared";
import type { ApiResponse, PageData } from "@erp/shared";

interface AuditRecord {
  log_id: string;
  username: string;
  action: string;
  resource_type: string;
  resource_id: string;
  detail: string;
  ip: string;
  created_at: string;
}

interface AuditQueryParams {
  page?: number;
  page_size?: number;
  username?: string;
  action?: string;
  start_date?: string;
  end_date?: string;
}

export const useAuditStore = defineStore("audit", () => {
  const logs = ref<AuditRecord[]>([]);
  const total = ref(0);
  const loading = ref(false);

  async function fetchLogs(params: AuditQueryParams = {}) {
    loading.value = true;
    try {
      const res = await apiClient.get<ApiResponse<PageData<AuditRecord>>>(
        "/iam/audit-logs",
        { params }
      );
      logs.value = res.data.data.list;
      total.value = res.data.data.total;
    } finally {
      loading.value = false;
    }
  }

  return {
    logs,
    total,
    loading,
    fetchLogs,
  };
});
