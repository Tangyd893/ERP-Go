import { apiClient } from "@erp/shared-react";
import {
  genIdempotencyKey,
  useWarehouseStore,
  type OutboundOrder,
  type PickTask,
  type ScanRecord,
} from "@/stores/warehouse";

function unwrapList<T>(data: unknown): T[] {
  const payload = (data as { data?: unknown })?.data ?? data;
  if (Array.isArray(payload)) return payload as T[];
  const list = (payload as { list?: T[]; items?: T[] })?.list
    ?? (payload as { items?: T[] })?.items
    ?? [];
  return list;
}

export async function fetchOutbounds(): Promise<OutboundOrder[]> {
  const res = await apiClient.get("/warehouse/outbounds", {
    params: { page: 1, page_size: 50 },
  });
  return unwrapList<OutboundOrder>(res.data).filter((o) => o.status !== "shipped");
}

export async function fetchPickTasks(outboundId: string): Promise<PickTask[]> {
  const res = await apiClient.get("/warehouse/pick-tasks", {
    params: { outbound_id: outboundId },
  });
  const payload = res.data?.data ?? res.data;
  return Array.isArray(payload) ? payload : [];
}

async function recordScan(
  record: ScanRecord,
  request: () => Promise<void>,
  duplicateMsg: string,
  failMsg: string
): Promise<ScanRecord> {
  const { addScanRecord } = useWarehouseStore.getState();
  try {
    await request();
    record.success = true;
    record.message = record.message || "操作成功";
  } catch (e: unknown) {
    const status = (e as { response?: { status?: number } })?.response?.status;
    if (status === 409) {
      record.success = true;
      record.message = duplicateMsg;
    } else {
      record.message = failMsg;
      addScanRecord(record);
      throw e;
    }
  } finally {
    if (record.success) {
      addScanRecord(record);
    }
  }
  return record;
}

export async function pickScan(
  taskId: string,
  quantity: number,
  taskLabel = ""
): Promise<ScanRecord> {
  const key = genIdempotencyKey("pick");
  const record: ScanRecord = {
    id: key,
    time: new Date(),
    type: "pick",
    targetId: taskId,
    targetLabel: taskLabel,
    quantity,
    success: false,
    message: "拣货成功",
  };
  return recordScan(
    record,
    () =>
      apiClient.post("/warehouse/pick/scan", {
        task_id: taskId,
        quantity,
        idempotency_key: key,
      }),
    "已拣货（重复扫描）",
    "拣货失败"
  );
}

export async function checkScan(
  outboundId: string,
  skuId: string,
  quantity: number,
  targetLabel = ""
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
    message: "复核成功",
  };
  return recordScan(
    record,
    () =>
      apiClient.post("/warehouse/check/scan", {
        outbound_id: outboundId,
        sku_id: skuId,
        quantity,
        idempotency_key: key,
      }),
    "已复核（重复扫描）",
    "复核失败"
  );
}

export async function packScan(outboundId: string, weight = 0): Promise<ScanRecord> {
  const key = genIdempotencyKey("pack");
  const record: ScanRecord = {
    id: key,
    time: new Date(),
    type: "pack",
    targetId: outboundId,
    targetLabel: outboundId,
    quantity: weight,
    success: false,
    message: "打包完成",
  };
  return recordScan(
    record,
    () =>
      apiClient.post("/warehouse/package", {
        outbound_id: outboundId,
        weight,
        idempotency_key: key,
      }),
    "已打包（重复操作）",
    "打包失败"
  );
}

export async function weighScan(outboundId: string, weight: number): Promise<ScanRecord> {
  const key = genIdempotencyKey("weigh");
  const record: ScanRecord = {
    id: key,
    time: new Date(),
    type: "weigh",
    targetId: outboundId,
    targetLabel: outboundId,
    quantity: weight,
    success: false,
    message: "称重完成",
  };
  return recordScan(
    record,
    () =>
      apiClient.post("/warehouse/weigh", {
        outbound_id: outboundId,
        weight,
        idempotency_key: key,
      }),
    "已称重（重复操作）",
    "称重失败"
  );
}

export async function confirmShip(
  outboundId: string,
  trackingNo = "",
  carrier = "default"
): Promise<void> {
  await apiClient.post(`/warehouse/outbounds/${outboundId}/ship`, {
    tracking_no: trackingNo,
    carrier,
    idempotency_key: genIdempotencyKey("ship"),
  });
}

/** Status filters matching Vue cheat-sheet */
export const STATUS_FILTERS = {
  pick: ["picking", "created"] as const,
  check: ["picked", "picking"] as const,
  pack: ["checked"] as const,
  weigh: ["packed"] as const,
  ship: ["checked", "packed", "weighed", "picked"] as const,
};

export function filterOutbounds(
  outbounds: OutboundOrder[],
  statuses: readonly string[]
): OutboundOrder[] {
  return outbounds.filter((o) => statuses.includes(o.status));
}
