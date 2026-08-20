import type { AxiosRequestConfig } from "axios";
import {
  DEMO_DASHBOARD,
  DEMO_TOKEN,
  DEMO_USER,
  demoCatalog,
  demoPickScan,
  demoRememberKey,
  demoUpdateOutboundStatus,
  getDemoOutbounds,
  getDemoPickTasks,
  ok,
  pageOf,
} from "./data";

function pathOf(config: AxiosRequestConfig): string {
  const raw = `${config.baseURL ?? ""}${config.url ?? ""}`;
  try {
    if (raw.startsWith("http")) {
      return new URL(raw).pathname.replace(/^\/api\/v1/, "") || "/";
    }
  } catch {
    /* ignore */
  }
  return (config.url ?? "").replace(/^\/api\/v1/, "") || "/";
}

function methodOf(config: AxiosRequestConfig): string {
  return (config.method ?? "get").toLowerCase();
}

function params(config: AxiosRequestConfig): Record<string, unknown> {
  return (config.params ?? {}) as Record<string, unknown>;
}

function body(config: AxiosRequestConfig): Record<string, unknown> {
  if (!config.data) return {};
  if (typeof config.data === "string") {
    try {
      return JSON.parse(config.data) as Record<string, unknown>;
    } catch {
      return {};
    }
  }
  return config.data as Record<string, unknown>;
}

function pageParams(config: AxiosRequestConfig) {
  const p = params(config);
  return {
    page: Number(p.page ?? 1) || 1,
    pageSize: Number(p.page_size ?? 20) || 20,
  };
}

export class DemoHttpError extends Error {
  status: number;
  data: unknown;

  constructor(status: number, message: string, data?: unknown) {
    super(message);
    this.status = status;
    this.data = data ?? ok(null, message);
  }
}

/**
 * Handle a request in demo mode. Returns response body (ApiResponse envelope).
 * Throws DemoHttpError for non-2xx semantics (e.g. 409 duplicate).
 */
export function handleDemoRequest(config: AxiosRequestConfig): unknown {
  const method = methodOf(config);
  const path = pathOf(config).replace(/\/+$/, "") || "/";
  const { page, pageSize } = pageParams(config);
  const b = body(config);

  // Auth
  if (method === "post" && path.endsWith("/iam/login")) {
    const username = String(b.username ?? "admin");
    return ok({
      access_token: DEMO_TOKEN,
      user: { ...DEMO_USER, username, nickname: username === "admin" ? "演示管理员" : username },
    });
  }

  // IAM
  if (method === "get" && path.endsWith("/iam/users")) {
    return ok(pageOf(demoCatalog.users, page, pageSize));
  }
  if (method === "post" && path.endsWith("/iam/users")) {
    const row = {
      user_id: `u-${Date.now()}`,
      username: String(b.username ?? "new"),
      nickname: String(b.nickname ?? ""),
      email: String(b.email ?? ""),
      status: 1,
      roles: [],
      created_at: new Date().toISOString(),
    };
    demoCatalog.users.unshift(row);
    return ok(row);
  }
  if (method === "delete" && /\/iam\/users\//.test(path)) {
    const id = path.split("/").pop();
    const idx = demoCatalog.users.findIndex((u) => u.user_id === id);
    if (idx >= 0) demoCatalog.users.splice(idx, 1);
    return ok(true);
  }
  if (method === "get" && path.endsWith("/iam/roles")) {
    return ok(pageOf(demoCatalog.roles, page, pageSize));
  }
  if (method === "post" && path.endsWith("/iam/roles")) {
    const row = {
      role_id: `r-${Date.now()}`,
      name: String(b.name ?? "新角色"),
      code: String(b.code ?? "new"),
      description: String(b.description ?? ""),
      status: 1,
      created_at: new Date().toISOString(),
    };
    demoCatalog.roles.unshift(row);
    return ok(row);
  }
  if (method === "delete" && /\/iam\/roles\//.test(path)) {
    const id = path.split("/").pop();
    const idx = demoCatalog.roles.findIndex((r) => r.role_id === id);
    if (idx >= 0) demoCatalog.roles.splice(idx, 1);
    return ok(true);
  }
  if (method === "get" && path.endsWith("/iam/audit-logs")) {
    return ok(pageOf(demoCatalog.auditLogs, page, pageSize));
  }

  // Tenant
  if (method === "get" && path.endsWith("/tenant/organizations")) {
    return ok({ list: demoCatalog.orgs });
  }

  // Products
  if (method === "get" && path.endsWith("/products/sku-mappings")) {
    return ok(pageOf(demoCatalog.skuMappings, page, pageSize));
  }
  if (method === "post" && path.endsWith("/products/sku-mappings")) {
    const row = {
      id: `m-${Date.now()}`,
      sku_code: String(b.sku_id ?? b.sku_code ?? ""),
      sku_name: String(b.sku_name ?? ""),
      platform_code: "demo",
      platform_sku: String(b.platform_sku ?? ""),
      asin: String(b.asin ?? ""),
      fnsku: String(b.fnsku ?? ""),
      store: String(b.store_id ?? "Demo Store"),
    };
    demoCatalog.skuMappings.unshift(row);
    return ok(row);
  }
  if (method === "get" && (path.endsWith("/products") || path.endsWith("/products/"))) {
    return ok(pageOf(demoCatalog.products, page, pageSize));
  }

  // Channel
  if (method === "get" && path.endsWith("/channel/stores")) {
    return ok(pageOf(demoCatalog.stores, page, pageSize));
  }
  if (method === "post" && path.endsWith("/channel/stores")) {
    const row = {
      id: `st-${Date.now()}`,
      name: String(b.name ?? "新店铺"),
      platform_code: String(b.platform_code ?? "demo"),
      site: String(b.site ?? "CN"),
      auth_status: "authorized",
      status: "active",
      created_at: new Date().toISOString(),
    };
    demoCatalog.stores.unshift(row);
    return ok(row);
  }
  if (method === "get" && path.endsWith("/channel/import-tasks")) {
    return ok(pageOf(demoCatalog.importTasks, page, pageSize));
  }
  if (method === "post" && path.endsWith("/channel/import")) {
    return ok({ id: `imp-${Date.now()}`, status: "processing" });
  }

  // Orders
  if (method === "get" && path.endsWith("/order/orders")) {
    let list = [...demoCatalog.orders];
    const status = params(config).status;
    const keyword = String(params(config).keyword ?? "");
    if (status) list = list.filter((o) => o.status === status);
    if (keyword) {
      list = list.filter(
        (o) =>
          o.order_no.includes(keyword) ||
          o.platform_order_no.includes(keyword) ||
          o.buyer_name.includes(keyword)
      );
    }
    return ok(pageOf(list, page, pageSize));
  }
  if (method === "post" && path.endsWith("/order/orders/audit")) {
    const order = demoCatalog.orders.find((o) => o.order_id === b.order_id || o.id === b.order_id);
    if (order) order.status = b.approved ? "approved" : "pending";
    return ok({ approved: !!b.approved });
  }
  if (method === "post" && path.endsWith("/order/orders/abnormal")) {
    const order = demoCatalog.orders.find((o) => o.order_id === b.order_id || o.id === b.order_id);
    if (order) order.status = "abnormal";
    return ok(true);
  }
  if (method === "post" && path.endsWith("/order/orders/cancel")) {
    const order = demoCatalog.orders.find((o) => o.order_id === b.order_id || o.id === b.order_id);
    if (order) order.status = "cancelled";
    return ok({ cancelled: true });
  }

  // Inventory
  if (method === "get" && path.endsWith("/inventory/balances")) {
    return ok(pageOf(demoCatalog.inventoryBalances, page, pageSize));
  }
  if (method === "get" && path.endsWith("/inventory/journals")) {
    return ok(pageOf(demoCatalog.inventoryJournals, page, pageSize));
  }

  // Warehouse
  if (method === "get" && /\/warehouse\/outbounds\/[^/]+$/.test(path)) {
    const id = path.split("/").pop()!;
    const ob = getDemoOutbounds().find((o) => o.id === id);
    return ok(ob ?? null);
  }
  if (method === "get" && path.endsWith("/warehouse/outbounds")) {
    const list = getDemoOutbounds();
    return ok({ list, items: list, total: list.length });
  }
  if (method === "get" && path.endsWith("/warehouse/pick-tasks")) {
    const outboundId = String(params(config).outbound_id ?? "");
    return ok(getDemoPickTasks(outboundId));
  }
  if (method === "post" && path.endsWith("/warehouse/pick/scan")) {
    const key = String(b.idempotency_key ?? "");
    if (key && demoRememberKey(key)) {
      throw new DemoHttpError(409, "duplicate");
    }
    demoPickScan(String(b.task_id ?? ""), Number(b.quantity ?? 1));
    return ok(true);
  }
  if (method === "post" && path.endsWith("/warehouse/check/scan")) {
    const key = String(b.idempotency_key ?? "");
    if (key && demoRememberKey(key)) {
      throw new DemoHttpError(409, "duplicate");
    }
    demoUpdateOutboundStatus(String(b.outbound_id ?? ""), "checked");
    return ok(true);
  }
  if (method === "post" && path.endsWith("/warehouse/package")) {
    const key = String(b.idempotency_key ?? "");
    if (key && demoRememberKey(key)) {
      throw new DemoHttpError(409, "duplicate");
    }
    demoUpdateOutboundStatus(String(b.outbound_id ?? ""), "packed");
    return ok(true);
  }
  if (method === "post" && path.endsWith("/warehouse/weigh")) {
    const key = String(b.idempotency_key ?? "");
    if (key && demoRememberKey(key)) {
      throw new DemoHttpError(409, "duplicate");
    }
    demoUpdateOutboundStatus(String(b.outbound_id ?? ""), "weighed");
    return ok(true);
  }
  if (method === "post" && /\/warehouse\/outbounds\/[^/]+\/ship$/.test(path)) {
    const key = String(b.idempotency_key ?? "");
    if (key && demoRememberKey(key)) {
      throw new DemoHttpError(409, "duplicate");
    }
    const id = path.split("/")[path.split("/").length - 2];
    demoUpdateOutboundStatus(id, "shipped");
    return ok(true);
  }

  // Transport
  if (method === "get" && path.endsWith("/transport/carriers")) {
    return ok(demoCatalog.carriers);
  }
  if (method === "get" && path.endsWith("/transport/shipments")) {
    return ok({ list: demoCatalog.shipments, items: demoCatalog.shipments, total: demoCatalog.shipments.length });
  }

  // Notification
  if (method === "get" && path.endsWith("/notification/list")) {
    return ok(pageOf(demoCatalog.notifications, page, pageSize));
  }

  // Purchase / Finance
  if (method === "get" && path.endsWith("/purchase/suppliers")) {
    return ok(pageOf(demoCatalog.suppliers, page, pageSize));
  }
  if (method === "get" && path.endsWith("/purchase/orders")) {
    return ok(pageOf(demoCatalog.purchaseOrders, page, pageSize));
  }
  if (method === "get" && path.endsWith("/finance/settlements")) {
    return ok(pageOf(demoCatalog.settlements, page, pageSize));
  }
  if (method === "get" && path.endsWith("/finance/profit")) {
    return ok(pageOf(demoCatalog.profits, page, pageSize));
  }

  // Report
  if (method === "get" && path.endsWith("/report/dashboard")) {
    return ok(DEMO_DASHBOARD);
  }

  // Default empty page for unknown GET
  if (method === "get") {
    return ok(pageOf([], page, pageSize));
  }
  return ok(true);
}
