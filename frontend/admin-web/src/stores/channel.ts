import { defineStore } from "pinia";
import { ref } from "vue";
import { apiClient } from "@erp/shared";
import type { ApiResponse, PageData } from "@erp/shared";

interface ChannelStore {
  id: string;
  name: string;
  platform_code: string;
  site: string;
  auth_status: string;
  status: string;
  created_at: string;
}

interface ImportTask {
  id: string;
  import_type: string;
  store_name: string;
  status: string;
  total_rows: number;
  success_rows: number;
  failed_rows: number;
  created_at: string;
}

export const useChannelStore = defineStore("channel", () => {
  const stores = ref<ChannelStore[]>([]);
  const importTasks = ref<ImportTask[]>([]);
  const loading = ref(false);

  async function fetchStores() {
    loading.value = true;
    try {
      const res = await apiClient.get<ApiResponse<PageData<ChannelStore>>>(
        "/channel/stores"
      );
      stores.value = res.data.data.list;
    } finally {
      loading.value = false;
    }
  }

  async function createStore(data: Partial<ChannelStore>) {
    const res = await apiClient.post<ApiResponse<ChannelStore>>(
      "/channel/stores",
      data
    );
    return res.data.data;
  }

  async function importOrders(file: File) {
    const formData = new FormData();
    formData.append("file", file);
    const res = await apiClient.post<ApiResponse<ImportTask>>(
      "/channel/import",
      formData
    );
    return res.data.data;
  }

  async function fetchImportTasks(page: number, pageSize: number) {
    loading.value = true;
    try {
      const res = await apiClient.get<ApiResponse<PageData<ImportTask>>>(
        "/channel/import-tasks",
        { params: { page, page_size: pageSize } }
      );
      importTasks.value = res.data.data.list;
    } finally {
      loading.value = false;
    }
  }

  return {
    stores,
    importTasks,
    loading,
    fetchStores,
    createStore,
    importOrders,
    fetchImportTasks,
  };
});
