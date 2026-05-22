import { defineStore } from "pinia";
import { ref } from "vue";
import { apiClient } from "@erp/shared";
import type { ApiResponse, UserInfo } from "@erp/shared";

export const useAuthStore = defineStore("auth", () => {
  const userInfo = ref<UserInfo | null>(null);
  const token = ref<string>(localStorage.getItem("access_token") || "");

  async function login(username: string, password: string, tenantId: string) {
    const res = await apiClient.post<ApiResponse<{ token: string; user: UserInfo }>>(
      "/iam/login",
      { username, password, tenant_id: tenantId }
    );
    const { token: newToken, user } = res.data.data;
    token.value = newToken;
    userInfo.value = user;
    localStorage.setItem("access_token", newToken);
    if (tenantId) {
      localStorage.setItem("tenant_id", tenantId);
    }
  }

  function logout() {
    token.value = "";
    userInfo.value = null;
    localStorage.removeItem("access_token");
    localStorage.removeItem("tenant_id");
  }

  return {
    userInfo,
    token,
    login,
    logout,
  };
});
