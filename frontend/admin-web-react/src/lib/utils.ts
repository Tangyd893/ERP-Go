import { toast } from "sonner";

export function stubAction(label?: string) {
  toast.info(label ? `${label}：即将支持` : "即将支持");
}

export function orderId(row: { id?: string; order_id?: string }) {
  return row.order_id || row.id || "";
}

export function unwrapList<T>(data: unknown): T[] {
  if (Array.isArray(data)) return data;
  if (data && typeof data === "object") {
    const obj = data as Record<string, unknown>;
    if (Array.isArray(obj.list)) return obj.list as T[];
    if (Array.isArray(obj.items)) return obj.items as T[];
  }
  return [];
}

export function unwrapPage<T>(data: unknown): { list: T[]; total: number } {
  if (data && typeof data === "object") {
    const obj = data as Record<string, unknown>;
    return {
      list: unwrapList<T>(data),
      total: Number(obj.total ?? unwrapList<T>(data).length),
    };
  }
  return { list: [], total: 0 };
}
