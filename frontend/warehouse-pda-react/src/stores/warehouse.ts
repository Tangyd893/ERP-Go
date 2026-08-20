import { create } from "zustand";

export interface PickTask {
  id: string;
  outbound_id: string;
  sku_id: string;
  sku_code: string;
  sku_name: string;
  quantity: number;
  picked_quantity: number;
  location_code: string;
  status: string;
}

export interface OutboundOrder {
  id: string;
  order_id: string;
  order_no: string;
  warehouse_id: string;
  status: string;
  items?: Array<{ sku_id: string; sku_name: string; quantity: number }>;
}

export interface ScanRecord {
  id: string;
  time: Date;
  type: "pick" | "check" | "pack" | "weigh";
  targetId: string;
  targetLabel: string;
  quantity: number;
  success: boolean;
  message: string;
}

interface WarehouseUiState {
  scanHistory: ScanRecord[];
  online: boolean;
  addScanRecord: (record: ScanRecord) => void;
  clearScanHistory: () => void;
}

export const useWarehouseStore = create<WarehouseUiState>((set) => {
  const online = typeof navigator !== "undefined" ? navigator.onLine : true;

  if (typeof window !== "undefined") {
    window.addEventListener("online", () => set({ online: true }));
    window.addEventListener("offline", () => set({ online: false }));
  }

  return {
    scanHistory: [],
    online,
    addScanRecord(record) {
      set((state) => {
        const next = [record, ...state.scanHistory].slice(0, 50);
        return { scanHistory: next };
      });
    },
    clearScanHistory() {
      set({ scanHistory: [] });
    },
  };
});

export function genIdempotencyKey(prefix: string): string {
  return `${prefix}-${crypto.randomUUID()}`;
}

export function getErrorMessage(e: unknown, defaultMsg = "操作失败，请重试"): string {
  const status = (e as { response?: { status?: number } })?.response?.status;
  if (status === 409) return "该步骤已完成或重复操作，请刷新页面";
  if (status === 404) return "未找到对应数据，请确认信息正确";
  if (status === 401) return "登录已过期，请重新登录";
  if (status && status >= 500) return "服务器繁忙，请稍后重试";
  if (!status) return "网络不可用，请检查连接后重试";
  return defaultMsg;
}

export function isDuplicateError(e: unknown): boolean {
  return (e as { response?: { status?: number } })?.response?.status === 409;
}

export function isNetworkError(e: unknown): boolean {
  return !(e as { response?: { status?: number } })?.response?.status;
}
