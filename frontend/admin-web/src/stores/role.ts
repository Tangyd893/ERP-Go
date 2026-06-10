import { defineStore } from "pinia";
import { ref } from "vue";
import { apiClient } from "@erp/shared";
import type { ApiResponse, PageData } from "@erp/shared";

interface RoleRecord {
  role_id: string;
  name: string;
  code: string;
  description: string;
  status: number;
  created_at: string;
}

export const useRoleStore = defineStore("role", () => {
  const roles = ref<RoleRecord[]>([]);
  const total = ref(0);
  const loading = ref(false);

  async function fetchRoles(page = 1, pageSize = 20) {
    loading.value = true;
    try {
      const res = await apiClient.get<ApiResponse<PageData<RoleRecord>>>(
        "/iam/roles",
        { params: { page, page_size: pageSize } }
      );
      roles.value = res.data.data.list;
      total.value = res.data.data.total;
    } finally {
      loading.value = false;
    }
  }

  async function createRole(data: Partial<RoleRecord>) {
    const res = await apiClient.post<ApiResponse<RoleRecord>>("/iam/roles", data);
    return res.data.data;
  }

  async function updateRole(id: string, data: Partial<RoleRecord>) {
    const res = await apiClient.put<ApiResponse<RoleRecord>>(
      `/iam/roles/${id}`,
      data
    );
    return res.data.data;
  }

  return {
    roles,
    total,
    loading,
    fetchRoles,
    createRole,
    updateRole,
    deleteRole,
  };
});

async function deleteRole(id: string) {
  await apiClient.delete(`/iam/roles/${id}`);
}
