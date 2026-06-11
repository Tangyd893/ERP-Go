import { defineStore } from "pinia";
import { ref } from "vue";
import { apiClient } from "@erp/shared";

// ── Types ──────────────────────────────────────────────

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
  type: "pick" | "check";
  targetId: string;
  targetLabel: string;
  quantity: number;
  success: boolean;
  message: string;
}

// ── Helpers ────────────────────────────────────────────

let _scanSeq = 0;
function genIdempotencyKey(prefix: string): string {
  _scanSeq++;
  return `${prefix}-${Date.now()}-${_scanSeq}-${Math.random().toString(36).slice(2, 8)}`;
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

export async function confirmShip(
  outboundId: string,
  trackingNo = "",
  carrier = "",
) {
  await apiClient.post(`/warehouse/outbounds/${outboundId}/ship`, {
    tracking_no: trackingNo,
    carrier,
    idempotency_key: genIdempotencyKey("ship"),
  });
}

// ── Store ──────────────────────────────────────────────

export const useWarehouseStore = defineStore("warehouse-pda", () => {
  const outbounds = ref<OutboundOrder[]>([]);
  const pickTasks = ref<PickTask[]>([]);
  const loading = ref(false);
  const scanHistory = ref<ScanRecord[]>([]);
  const online = ref(navigator.onLine);

  // Track online/offline
  if (typeof globalThis.window !== "undefined") {
    globalThis.window.addEventListener("online", () => { online.value = true; });
    globalThis.window.addEventListener("offline", () => { online.value = false; });
  }

  function addScanRecord(r: ScanRecord) {
    scanHistory.value.unshift(r);
    // Keep last 50 records
    if (scanHistory.value.length > 50) {
      scanHistory.value = scanHistory.value.slice(0, 50);
    }
  }

  async function fetchOutbounds() {
    loading.value = true;
    try {
      const res = await apiClient.get("/warehouse/outbounds", {
        params: { page: 1, page_size: 50 },
      });
      const data = res.data?.data ?? res.data;
      outbounds.value = (data.list ?? data.items ?? []).filter(
        (o: OutboundOrder) => o.status !== "shipped"
      );
    } finally {
      loading.value = false;
    }
  }

  async function fetchPickTasks(outboundId: string) {
    loading.value = true;
    try {
      const res = await apiClient.get("/warehouse/pick-tasks", {
        params: { outbound_id: outboundId },
      });
      pickTasks.value = res.data?.data ?? res.data ?? [];
    } finally {
      loading.value = false;
    }
  }

  // ── Scan operations with idempotency ─────────────────

  async function pickScan(taskId: string, quantity: number, taskLabel = ""): Promise<ScanRecord> {
    const key = genIdempotencyKey("pick");
    const record: ScanRecord = {
      id: key,
      time: new Date(),
      type: "pick",
      targetId: taskId,
      targetLabel: taskLabel,
      quantity,
      success: false,
      message: "",
    };
    try {
      await apiClient.post("/warehouse/pick/scan", {
        task_id: taskId,
        quantity,
        idempotency_key: key,
      });
      record.success = true;
      record.message = "拣货成功";
    } catch (e: unknown) {
      if (isDuplicateError(e)) {
        record.success = true;
        record.message = "已拣货（重复扫描）";
      } else {
        record.message = getErrorMessage(e, "拣货失败");
      }
      throw e;
    } finally {
      addScanRecord(record);
    }
    return record;
  }

  async function checkScan(
    outboundId: string,
    skuId: string,
    quantity: number,
    targetLabel = "",
  ): Promise<ScanRecord> {
    const key = genIdempotencyKey("check");
    const record: ScanRecord = {
      id: key,
      time: new Date(),
      type: "check",
      targetId: outboundId,
      targetLabel: targetLabel || skuId || outboundId,
      quantity,
      success: false,
      message: "",
    };
    try {
      await apiClient.post("/warehouse/check/scan", {
        outbound_id: outboundId,
        sku_id: skuId,
        quantity,
        idempotency_key: key,
      });
      record.success = true;
      record.message = "复核成功";
    } catch (e: unknown) {
      if (isDuplicateError(e)) {
        record.success = true;
        record.message = "已复核（重复扫描）";
      } else {
        record.message = getErrorMessage(e, "复核失败");
      }
      throw e;
    } finally {
      addScanRecord(record);
    }
    return record;
  }

  async function pack(outboundId: string, weight = 0): Promise<ScanRecord> {
    const key = genIdempotencyKey("pack");
    const record: ScanRecord = {
      id: key,
      time: new Date(),
      type: "check",
      targetId: outboundId,
      targetLabel: outboundId,
      quantity: weight,
      success: false,
      message: "",
    };
    try {
      await apiClient.post("/warehouse/package", {
        outbound_id: outboundId,
        weight,
        idempotency_key: key,
      });
      record.success = true;
      record.message = "打包完成";
    } catch (e: unknown) {
      if (isDuplicateError(e)) {
        record.success = true;
        record.message = "已打包（重复操作）";
      } else {
        record.message = getErrorMessage(e, "打包失败");
      }
      throw e;
    } finally {
      addScanRecord(record);
    }
    return record;
  }

  async function weigh(outboundId: string, weight: number): Promise<ScanRecord> {
    const key = genIdempotencyKey("weigh");
    const record: ScanRecord = {
      id: key,
      time: new Date(),
      type: "check",
      targetId: outboundId,
      targetLabel: outboundId,
      quantity: weight,
      success: false,
      message: "",
    };
    try {
      await apiClient.post("/warehouse/weigh", {
        outbound_id: outboundId,
        weight,
        idempotency_key: key,
      });
      record.success = true;
      record.message = "称重完成";
    } catch (e: unknown) {
      if (isDuplicateError(e)) {
        record.success = true;
        record.message = "已称重（重复操作）";
      } else {
        record.message = getErrorMessage(e, "称重失败");
      }
      throw e;
    } finally {
      addScanRecord(record);
    }
    return record;
  }

  function clearScanHistory() {
    scanHistory.value = [];
  }

  return {
    outbounds,
    pickTasks,
    loading,
    scanHistory,
    online,
    fetchOutbounds,
    fetchPickTasks,
    pickScan,
    checkScan,
    pack,
    weigh,
    clearScanHistory,
  };
});
