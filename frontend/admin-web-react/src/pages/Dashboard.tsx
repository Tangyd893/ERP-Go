import { useQuery } from "@tanstack/react-query";
import { Card, CardContent, CardHeader, CardTitle, Skeleton } from "@erp/shared-react";
import { apiClient, type ApiResponse, type PageData } from "@erp/shared-react";

interface OrderRow {
  status: string;
}

interface InventoryBalance {
  quantity?: number;
  alert_threshold?: number;
  available_qty?: number;
}

export default function Dashboard() {
  const ordersQuery = useQuery({
    queryKey: ["dashboard", "orders"],
    queryFn: async () => {
      const res = await apiClient.get<ApiResponse<PageData<OrderRow>>>("/order/orders");
      return res.data.data;
    },
  });

  const inventoryQuery = useQuery({
    queryKey: ["dashboard", "inventory"],
    queryFn: async () => {
      const res = await apiClient.get<ApiResponse<PageData<InventoryBalance>>>(
        "/inventory/balances"
      );
      return res.data.data;
    },
  });

  const warehouseQuery = useQuery({
    queryKey: ["dashboard", "warehouse"],
    queryFn: async () => {
      const res = await apiClient.get("/warehouse/outbounds");
      const data = res.data?.data ?? res.data;
      return { total: data?.total ?? 0 };
    },
  });

  const orders = ordersQuery.data?.list ?? [];
  const todayOrders = ordersQuery.data?.total ?? 0;
  const pendingAudit = orders.filter((o) => o.status === "pending").length;
  const pendingOutbound = warehouseQuery.data?.total ?? 0;
  const balances = inventoryQuery.data?.list ?? [];
  const alertCount = balances.filter(
    (b) => (b.available_qty ?? b.quantity ?? 0) < (b.alert_threshold ?? 10)
  ).length;

  const stats = [
    { label: "今日订单", value: todayOrders, color: "text-primary" },
    { label: "待审核", value: pendingAudit, color: "text-amber-600" },
    { label: "待出库", value: pendingOutbound, color: "text-emerald-600" },
    { label: "库存告警", value: alertCount, color: "text-destructive" },
  ];

  const loading = ordersQuery.isLoading || inventoryQuery.isLoading || warehouseQuery.isLoading;

  return (
    <div className="space-y-4">
      <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
        {stats.map((stat) => (
          <Card key={stat.label}>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium text-muted-foreground">
                {stat.label}
              </CardTitle>
            </CardHeader>
            <CardContent>
              {loading ? (
                <Skeleton className="h-9 w-16" />
              ) : (
                <div className={`text-3xl font-bold tabular-nums ${stat.color}`}>{stat.value}</div>
              )}
            </CardContent>
          </Card>
        ))}
      </div>

      <Card>
        <CardHeader>
          <CardTitle className="text-base">欢迎使用 ERP-Go 跨境电商 ERP 系统</CardTitle>
        </CardHeader>
        <CardContent className="space-y-2 text-sm text-muted-foreground">
          <p>系统当前处于 MVP 阶段，功能正在持续构建中。</p>
          <p>首期功能：商品管理、订单管理、库存管理、仓储出库、物流发货。</p>
        </CardContent>
      </Card>
    </div>
  );
}
