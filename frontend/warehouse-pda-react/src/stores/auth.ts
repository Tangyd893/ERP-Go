import { create } from "zustand";
import { apiClient } from "@erp/shared-react";

interface AuthState {
  token: string;
  tenantId: string;
  username: string;
  login: (username: string, password: string, tenantId?: string) => Promise<void>;
  logout: () => void;
  isLoggedIn: () => boolean;
}

export const useAuthStore = create<AuthState>((set, get) => ({
  token: localStorage.getItem("access_token") || "",
  tenantId: localStorage.getItem("tenant_id") || "default",
  username: localStorage.getItem("username") || "",

  async login(username, password, tenantId = "default") {
    const res = await apiClient.post("/iam/login", {
      username,
      password,
      tenant_id: tenantId,
    });
    const data = res.data?.data ?? res.data;
    const newToken = data.access_token || data.token;
    if (!newToken) {
      throw new Error("登录失败");
    }
    localStorage.setItem("access_token", newToken);
    localStorage.setItem("tenant_id", tenantId);
    localStorage.setItem("username", username);
    set({ token: newToken, tenantId, username });
  },

  logout() {
    localStorage.removeItem("access_token");
    localStorage.removeItem("tenant_id");
    localStorage.removeItem("username");
    set({ token: "", tenantId: "default", username: "" });
  },

  isLoggedIn() {
    return !!get().token;
  },
}));
