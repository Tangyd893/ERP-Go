import { create } from "zustand";
import { apiClient } from "@erp/shared-react";
import type { ApiResponse, UserInfo } from "@erp/shared-react";

interface AuthState {
  userInfo: UserInfo | null;
  token: string;
  login: (username: string, password: string, tenantId: string) => Promise<void>;
  logout: () => void;
}

export const useAuthStore = create<AuthState>((set) => ({
  userInfo: null,
  token: localStorage.getItem("access_token") || "",

  async login(username, password, tenantId) {
    const res = await apiClient.post<
      ApiResponse<{ access_token?: string; token?: string; user: UserInfo }>
    >("/iam/login", { username, password, tenant_id: tenantId });
    const data = res.data.data;
    const newToken = data.access_token || data.token || "";
    if (!newToken) {
      throw new Error("登录响应缺少 token");
    }
    localStorage.setItem("access_token", newToken);
    if (tenantId) {
      localStorage.setItem("tenant_id", tenantId);
    }
    set({ token: newToken, userInfo: data.user });
  },

  logout() {
    localStorage.removeItem("access_token");
    localStorage.removeItem("tenant_id");
    set({ token: "", userInfo: null });
  },
}));
