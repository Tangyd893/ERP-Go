import axios, {
  type AxiosAdapter,
  type AxiosInstance,
  type InternalAxiosRequestConfig,
  type AxiosResponse,
} from "axios";
import { isDemo } from "./demo";
import { DemoHttpError, handleDemoRequest } from "./demo/handler";

const demoAdapter: AxiosAdapter = async (config) => {
  try {
    const data = handleDemoRequest(config);
    const response: AxiosResponse = {
      data,
      status: 200,
      statusText: "OK",
      headers: { "content-type": "application/json" },
      config: config as InternalAxiosRequestConfig,
    };
    return response;
  } catch (err) {
    if (err instanceof DemoHttpError) {
      return Promise.reject(
        Object.assign(new Error(err.message), {
          isAxiosError: true,
          config,
          response: {
            data: err.data,
            status: err.status,
            statusText: err.message,
            headers: {},
            config,
          },
        })
      );
    }
    return Promise.reject(err);
  }
};

const apiClient: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || "/api/v1",
  timeout: 30000,
  headers: {
    "Content-Type": "application/json",
  },
});

apiClient.interceptors.request.use((config) => {
  // Demo adapter is development-only (isDemo() is false in PROD builds).
  if (isDemo()) {
    config.adapter = demoAdapter;
  }

  const token = localStorage.getItem("access_token");
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  const tenantId = localStorage.getItem("tenant_id");
  if (tenantId) {
    config.headers["X-Tenant-ID"] = tenantId;
  }
  return config;
});

apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (isDemo()) {
      return Promise.reject(error);
    }
    if (error.response?.status === 401) {
      localStorage.removeItem("access_token");
      const redirectOn401 = import.meta.env.VITE_AUTH_REDIRECT !== "false";
      const onLoginPage = globalThis.location.pathname === "/login";
      if (redirectOn401 && !onLoginPage) {
        globalThis.location.href = "/login";
      }
    }
    return Promise.reject(error);
  }
);

export default apiClient;
