import { defineStore } from "pinia";
import { ref } from "vue";
import { apiClient } from "@erp/shared";

export const useAuthStore = defineStore("auth-pda", () => {
  const token = ref(localStorage.getItem("access_token") || "");
  const tenantId = ref(localStorage.getItem("tenant_id") || "default");
  const username = ref(localStorage.getItem("username") || "");

  async function login(user: string, password: string, tenant = "default") {
    const res = await apiClient.post("/iam/login", {
      username: user,
      password,
      tenant_id: tenant,
    });
    const data = res.data?.data ?? res.data;
    const newToken = data.access_token || data.token;
    if (!newToken) {
      throw new Error("登录失败");
    }
    token.value = newToken;
    tenantId.value = tenant;
    username.value = user;
    localStorage.setItem("access_token", newToken);
    localStorage.setItem("tenant_id", tenant);
    localStorage.setItem("username", user);
  }

  function logout() {
    token.value = "";
    username.value = "";
    localStorage.removeItem("access_token");
    localStorage.removeItem("tenant_id");
    localStorage.removeItem("username");
  }

  function isLoggedIn() {
    return !!token.value;
  }

  return { token, tenantId, username, login, logout, isLoggedIn };
});
