import { defineStore } from "pinia";
import { ref } from "vue";
import { apiClient } from "@erp/shared";
import type { ApiResponse, PageData } from "@erp/shared";

interface UserRecord {
  user_id: string;
  username: string;
  nickname: string;
  email: string;
  status: number;
  roles: string[];
  created_at: string;
}

export const useUserStore = defineStore("user", () => {
  const users = ref<UserRecord[]>([]);
  const total = ref(0);
  const loading = ref(false);

  async function fetchUsers(page: number, pageSize: number) {
    loading.value = true;
    try {
      const res = await apiClient.get<ApiResponse<PageData<UserRecord>>>(
        "/iam/users",
        { params: { page, page_size: pageSize } }
      );
      users.value = res.data.data.list;
      total.value = res.data.data.total;
    } finally {
      loading.value = false;
    }
  }

  async function createUser(data: Partial<UserRecord>) {
    const res = await apiClient.post<ApiResponse<UserRecord>>("/iam/users", data);
    return res.data.data;
  }

  async function updateUser(userId: string, data: Partial<UserRecord>) {
    const res = await apiClient.put<ApiResponse<UserRecord>>(
      `/iam/users/${userId}`,
      data
    );
    return res.data.data;
  }

  async function deleteUser(userId: string) {
    await apiClient.delete(`/iam/users/${userId}`);
  }

  return {
    users,
    total,
    loading,
    fetchUsers,
    createUser,
    updateUser,
    deleteUser,
  };
});
