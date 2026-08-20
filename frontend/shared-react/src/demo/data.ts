import type { ApiResponse, PageData, UserInfo } from "../types";

export const DEMO_TOKEN = "demo-access-token";

export const DEMO_USER: UserInfo = {
  user_id: "u-demo-admin",
  username: "admin",
  nickname: "演示管理员",
  email: "admin@demo.local",
  tenant_id: "default",
  roles: ["admin"],
  permissions: ["*"],
};

function pageOf<T>(list: T[], page = 1, pageSize = 20): PageData<T> {
  const start = (page - 1) * pageSize;
  const slice = list.slice(start, start + pageSize);
  return {
    list: slice,
    total: list.length,
    page,
    page_size: pageSize,
    total_pages: Math.max(1, Math.ceil(list.length / pageSize)),
  };
}

export function ok<T>(data: T, message = "ok"): ApiResponse<T> {
  return { code: 0, message, data };
}

const now = () => new Date().toISOString();

/** Mutable warehouse outbounds for PDA demo flows */
export type DemoOutbound = {
  id: string;
  order_no: string;
  orderNo?: string;
  status: string;
  warehouse_id?: string;
  store_name?: string;
  storeName?: string;
  recipient_name?: string;
  items?: Array<{
    skuId: string;
    skuCode: string;
    skuName: string;
    quantity: number;
    pickedQty?: number;
    checkedQty?: number;
    packedQty?: number;
  }>;
  created_at?: string;
  createdAt?: string;
};

export type DemoPickTask = {
  id: string;
  outbound_id: string;
  sku_id: string;
  sku_code: string;
  sku_name: string;
  location_code: string;
  quantity: number;
  picked_quantity: number;
  status: string;
};

const demoState = {
  outbounds: [] as DemoOutbound[],
  pickTasks: {} as Record<string, DemoPickTask[]>,
  scanKeys: new Set<string>(),
};

function resetWarehouseDemo() {
  demoState.outbounds = [
    {
      id: "ob-1001",
      order_no: "SO20260820001",
      orderNo: "SO20260820001",
      status: "created",
      store_name: "Amazon US",
      storeName: "Amazon US",
      recipient_name: "Alice",
      items: [
        {
          skuId: "sku-blue",
          skuCode: "SKU-BLUE-M",
          skuName: "蓝色 T 恤 M",
          quantity: 10,
          pickedQty: 0,
        },
        {
          skuId: "sku-red",
          skuCode: "SKU-RED-L",
          skuName: "红色卫衣 L",
          quantity: 5,
          pickedQty: 0,
        },
      ],
      created_at: now(),
      createdAt: now(),
    },
    {
      id: "ob-1002",
      order_no: "SO20260820002",
      orderNo: "SO20260820002",
      status: "picking",
      store_name: "Shopify CN",
      storeName: "Shopify CN",
      recipient_name: "Bob",
      items: [
        {
          skuId: "sku-xl",
          skuCode: "SKU-XL-01",
          skuName: "大号收纳箱",
          quantity: 20,
          pickedQty: 8,
        },
      ],
      created_at: now(),
      createdAt: now(),
    },
    {
      id: "ob-1003",
      order_no: "SO20260820003",
      orderNo: "SO20260820003",
      status: "picked",
      store_name: "Amazon US",
      storeName: "Amazon US",
      recipient_name: "Carol",
      items: [
        {
          skuId: "sku-blue",
          skuCode: "SKU-BLUE-M",
          skuName: "蓝色 T 恤 M",
          quantity: 15,
          pickedQty: 15,
          checkedQty: 0,
        },
      ],
      created_at: now(),
      createdAt: now(),
    },
    {
      id: "ob-1004",
      order_no: "SO20260820004",
      orderNo: "SO20260820004",
      status: "checked",
      store_name: "TikTok Shop",
      storeName: "TikTok Shop",
      recipient_name: "Dave",
      items: [
        {
          skuId: "sku-sm",
          skuCode: "SKU-SM-02",
          skuName: "小号配件包",
          quantity: 8,
          pickedQty: 8,
          checkedQty: 8,
        },
      ],
      created_at: now(),
      createdAt: now(),
    },
    {
      id: "ob-1005",
      order_no: "SO20260820005",
      orderNo: "SO20260820005",
      status: "packed",
      store_name: "Amazon US",
      storeName: "Amazon US",
      recipient_name: "Eve",
      items: [
        {
          skuId: "sku-xl",
          skuCode: "SKU-XL-01",
          skuName: "大号收纳箱",
          quantity: 6,
          pickedQty: 6,
          checkedQty: 6,
          packedQty: 6,
        },
      ],
      created_at: now(),
      createdAt: now(),
    },
  ];

  demoState.pickTasks = {
    "ob-1001": [
      {
        id: "pt-1",
        outbound_id: "ob-1001",
        sku_id: "sku-blue",
        sku_code: "SKU-BLUE-M",
        sku_name: "蓝色 T 恤 M",
        location_code: "A-01-02",
        quantity: 10,
        picked_quantity: 0,
        status: "pending",
      },
      {
        id: "pt-2",
        outbound_id: "ob-1001",
        sku_id: "sku-red",
        sku_code: "SKU-RED-L",
        sku_name: "红色卫衣 L",
        location_code: "B-02-01",
        quantity: 5,
        picked_quantity: 0,
        status: "pending",
      },
    ],
    "ob-1002": [
      {
        id: "pt-3",
        outbound_id: "ob-1002",
        sku_id: "sku-xl",
        sku_code: "SKU-XL-01",
        sku_name: "大号收纳箱",
        location_code: "C-03-01",
        quantity: 20,
        picked_quantity: 8,
        status: "pending",
      },
    ],
  };
  demoState.scanKeys.clear();
}

resetWarehouseDemo();

export function getDemoOutbounds() {
  return demoState.outbounds;
}

export function getDemoPickTasks(outboundId: string) {
  return demoState.pickTasks[outboundId] ?? [];
}

export function demoRememberKey(key: string): boolean {
  if (demoState.scanKeys.has(key)) return true;
  demoState.scanKeys.add(key);
  return false;
}

export function demoUpdateOutboundStatus(id: string, status: string) {
  const ob = demoState.outbounds.find((o) => o.id === id);
  if (ob) ob.status = status;
}

export function demoPickScan(taskId: string, quantity: number) {
  for (const tasks of Object.values(demoState.pickTasks)) {
    const task = tasks.find((t) => t.id === taskId);
    if (!task) continue;
    task.picked_quantity = Math.min(task.quantity, task.picked_quantity + quantity);
    if (task.picked_quantity >= task.quantity) task.status = "done";
    const allDone = tasks.every((t) => t.picked_quantity >= t.quantity);
    if (allDone) demoUpdateOutboundStatus(task.outbound_id, "picked");
    else demoUpdateOutboundStatus(task.outbound_id, "picking");
    return;
  }
}

export const DEMO_DASHBOARD = {
  order_count: 328,
  sales_amount: 12850,
  outbound_count: 186,
  sku_count: 1256,
  trend: [
    { date: "周一", order_count: 45, sales_amount: 2100 },
    { date: "周二", order_count: 52, sales_amount: 3200 },
    { date: "周三", order_count: 38, sales_amount: 1850 },
    { date: "周四", order_count: 65, sales_amount: 4200 },
    { date: "周五", order_count: 48, sales_amount: 2800 },
    { date: "周六", order_count: 32, sales_amount: 1500 },
    { date: "周日", order_count: 28, sales_amount: 1200 },
  ],
  timeliness: { within_24h: 152, within_48h: 28, overdue: 6 },
};

export const demoCatalog = {
  users: [
    {
      user_id: "u-1",
      username: "admin",
      nickname: "演示管理员",
      email: "admin@demo.local",
      status: 1,
      roles: ["admin"],
      created_at: "2026-08-01T08:00:00Z",
    },
    {
      user_id: "u-2",
      username: "ops01",
      nickname: "运营小王",
      email: "ops01@demo.local",
      status: 1,
      roles: ["ops"],
      created_at: "2026-08-05T09:00:00Z",
    },
    {
      user_id: "u-3",
      username: "wh01",
      nickname: "仓管小李",
      email: "wh01@demo.local",
      status: 1,
      roles: ["warehouse"],
      created_at: "2026-08-10T10:00:00Z",
    },
  ],
  roles: [
    {
      role_id: "r-1",
      name: "管理员",
      code: "admin",
      description: "全部权限",
      status: 1,
      created_at: "2026-08-01T08:00:00Z",
    },
    {
      role_id: "r-2",
      name: "运营",
      code: "ops",
      description: "订单与商品",
      status: 1,
      created_at: "2026-08-01T08:00:00Z",
    },
    {
      role_id: "r-3",
      name: "仓管",
      code: "warehouse",
      description: "仓储作业",
      status: 1,
      created_at: "2026-08-01T08:00:00Z",
    },
  ],
  orgs: [
    {
      id: "org-1",
      name: "演示集团",
      code: "DEMO",
      children: [
        { id: "org-2", name: "运营中心", code: "OPS" },
        { id: "org-3", name: "华东仓", code: "WH-EAST" },
      ],
    },
  ],
  auditLogs: [
    {
      log_id: "a-1",
      username: "admin",
      action: "login",
      resource_type: "session",
      resource_id: "-",
      detail: "演示登录",
      ip: "127.0.0.1",
      created_at: now(),
    },
    {
      log_id: "a-2",
      username: "ops01",
      action: "order.audit",
      resource_type: "order",
      resource_id: "ord-1001",
      detail: "审核通过",
      ip: "10.0.0.8",
      created_at: now(),
    },
  ],
  products: [
    {
      id: "p-1",
      code: "SPU-TEE",
      name: "基础 T 恤",
      spu_name: "基础 T 恤",
      barcode: "6900001",
      weight: 0.2,
      sale_price: 59,
      currency: "CNY",
      status: "active",
    },
    {
      id: "p-2",
      code: "SPU-BOX",
      name: "收纳箱",
      spu_name: "收纳箱",
      barcode: "6900002",
      weight: 1.5,
      sale_price: 128,
      currency: "CNY",
      status: "active",
    },
    {
      id: "p-3",
      code: "SPU-HOOD",
      name: "卫衣",
      spu_name: "卫衣",
      barcode: "6900003",
      weight: 0.6,
      sale_price: 199,
      currency: "CNY",
      status: "active",
    },
  ],
  skuMappings: [
    {
      id: "m-1",
      sku_code: "SKU-BLUE-M",
      sku_name: "蓝色 T 恤 M",
      platform_code: "amazon",
      platform_sku: "AMZ-BLUE-M",
      asin: "B0DEMO001",
      fnsku: "X00DEMO1",
      store: "Amazon US",
    },
    {
      id: "m-2",
      sku_code: "SKU-XL-01",
      sku_name: "大号收纳箱",
      platform_code: "shopify",
      platform_sku: "SHP-XL-01",
      asin: "",
      fnsku: "",
      store: "Shopify CN",
    },
  ],
  stores: [
    {
      id: "st-1",
      name: "Amazon US",
      platform_code: "amazon",
      site: "US",
      auth_status: "authorized",
      status: "active",
      created_at: "2026-07-01T00:00:00Z",
    },
    {
      id: "st-2",
      name: "Shopify CN",
      platform_code: "shopify",
      site: "CN",
      auth_status: "authorized",
      status: "active",
      created_at: "2026-07-15T00:00:00Z",
    },
  ],
  importTasks: [
    {
      id: "imp-1",
      import_type: "order",
      store_name: "Amazon US",
      status: "success",
      total_rows: 120,
      success_rows: 118,
      failed_rows: 2,
      created_at: now(),
    },
  ],
  orders: [
    {
      id: "ord-1001",
      order_id: "ord-1001",
      order_no: "PO-1001",
      platform_order_no: "AMZ-77881",
      store: "Amazon US",
      buyer_name: "Alice",
      total_amount: 258,
      currency: "USD",
      status: "pending",
      created_at: now(),
      items: [{ sku_id: "sku-blue", sku_name: "蓝色 T 恤 M", qty: 2, unit_price: 29 }],
    },
    {
      id: "ord-1002",
      order_id: "ord-1002",
      order_no: "PO-1002",
      platform_order_no: "SHP-9921",
      store: "Shopify CN",
      buyer_name: "Bob",
      total_amount: 399,
      currency: "CNY",
      status: "approved",
      created_at: now(),
      items: [{ sku_id: "sku-xl", sku_name: "大号收纳箱", qty: 1, unit_price: 128 }],
    },
    {
      id: "ord-1003",
      order_id: "ord-1003",
      order_no: "PO-1003",
      platform_order_no: "AMZ-77890",
      store: "Amazon US",
      buyer_name: "Carol",
      total_amount: 89,
      currency: "USD",
      status: "picking",
      created_at: now(),
    },
    {
      id: "ord-1004",
      order_id: "ord-1004",
      order_no: "PO-1004",
      platform_order_no: "TT-331",
      store: "TikTok Shop",
      buyer_name: "Dave",
      total_amount: 66,
      currency: "CNY",
      status: "shipped",
      created_at: now(),
    },
    {
      id: "ord-1005",
      order_id: "ord-1005",
      order_no: "PO-1005",
      platform_order_no: "AMZ-77901",
      store: "Amazon US",
      buyer_name: "Eve",
      total_amount: 120,
      currency: "USD",
      status: "abnormal",
      created_at: now(),
    },
  ],
  inventoryBalances: [
    {
      sku_id: "sku-blue",
      sku_code: "SKU-BLUE-M",
      sku_name: "蓝色 T 恤 M",
      warehouse_id: "wh-1",
      warehouse_name: "华东仓",
      qty: 120,
      locked_qty: 15,
      available_qty: 105,
      quantity: 120,
      alert_threshold: 20,
    },
    {
      sku_id: "sku-red",
      sku_code: "SKU-RED-L",
      sku_name: "红色卫衣 L",
      warehouse_id: "wh-1",
      warehouse_name: "华东仓",
      qty: 18,
      locked_qty: 5,
      available_qty: 13,
      quantity: 18,
      alert_threshold: 20,
    },
    {
      sku_id: "sku-xl",
      sku_code: "SKU-XL-01",
      sku_name: "大号收纳箱",
      warehouse_id: "wh-1",
      warehouse_name: "华东仓",
      qty: 60,
      locked_qty: 10,
      available_qty: 50,
      quantity: 60,
      alert_threshold: 10,
    },
  ],
  inventoryJournals: [
    {
      journal_id: "j-1",
      sku_id: "sku-blue",
      sku_name: "蓝色 T 恤 M",
      change_qty: -2,
      before_qty: 122,
      after_qty: 120,
      biz_type: "outbound",
      biz_no: "SO20260820001",
      created_at: now(),
    },
  ],
  carriers: [
    {
      id: "c-1",
      name: "演示快递",
      code: "DEMO_EXP",
      contactPerson: "运力经理",
      contactPhone: "13800000000",
      status: "active",
    },
  ],
  shipments: [
    {
      id: "sh-1",
      orderId: "ord-1004",
      orderNo: "PO-1004",
      outboundId: "ob-shipped",
      trackingNo: "YT000DEMO1",
      carrierId: "c-1",
      carrierName: "演示快递",
      status: "shipped",
      labelUrl: "",
      shippingCost: 12,
      createdAt: now(),
      tracking_no: "YT000DEMO1",
      carrier: "演示快递",
      created_at: now(),
    },
  ],
  notifications: [
    {
      id: "n-1",
      title: "演示环境已启用",
      content: "当前为 VITE_DEMO_MODE，数据为本地样例，不请求真实后端。",
      type: "info",
      read: false,
      created_at: now(),
    },
    {
      id: "n-2",
      title: "低库存提醒",
      content: "SKU-RED-L 可用库存低于阈值。",
      type: "warning",
      read: false,
      created_at: now(),
    },
  ],
  suppliers: [
    {
      id: "sup-1",
      name: "演示供应商 A",
      code: "SUP-A",
      contact: "张三",
      phone: "13900000001",
      status: "active",
    },
  ],
  purchaseOrders: [
    {
      id: "po-1",
      order_no: "PR-20260801",
      supplier_name: "演示供应商 A",
      status: "approved",
      total_amount: 5600,
      currency: "CNY",
      created_at: now(),
    },
  ],
  settlements: [
    {
      id: "set-1",
      store: "Amazon US",
      period: "2026-07",
      total_sales: 12850,
      commission: 1285,
      net_amount: 11565,
      status: "matched",
      currency: "USD",
    },
  ],
  profits: [
    {
      order_no: "PO-1004",
      sale_amount: 66,
      purchase_cost: 28,
      shipping_cost: 8,
      commission: 6,
      total_cost: 42,
      gross_profit: 24,
      profit_margin: 0.36,
    },
  ],
};

export { pageOf, resetWarehouseDemo };
