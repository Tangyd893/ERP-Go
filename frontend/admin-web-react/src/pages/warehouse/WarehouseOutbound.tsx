import { useMemo, useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { toast } from "sonner";
import {
  Badge,
  Button,
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  ProTable,
  type ProColumn,
  apiClient,
} from "@erp/shared-react";
import { PageCard } from "@/components/PageCard";
import { unwrapPage } from "@/lib/utils";

interface OutboundOrder {
  id: string;
  orderNo?: string;
  order_no?: string;
  status: string;
  items?: unknown[];
  createdAt?: string;
  created_at?: string;
}

const statusLabels: Record<string, { variant: "secondary" | "warning" | "success" | "destructive" | "default"; label: string }> = {
  created: { variant: "secondary", label: "待波次" },
  picking: { variant: "default", label: "拣货中" },
  picked: { variant: "warning", label: "已拣货" },
  checking: { variant: "default", label: "复核中" },
  checked: { variant: "warning", label: "已复核" },
  packed: { variant: "default", label: "已打包" },
  weighed: { variant: "default", label: "已称重" },
  shipped: { variant: "success", label: "已出库" },
  abnormal: { variant: "destructive", label: "异常" },
};

export default function WarehouseOutbound() {
  const [page, setPage] = useState(1);

  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ["warehouse-outbounds", page],
    queryFn: async () => {
      const res = await apiClient.get("/warehouse/outbounds", { params: { page } });
      return unwrapPage<OutboundOrder>(res.data?.data ?? res.data);
    },
  });

  const outbounds = data?.list ?? [];

  const stats = useMemo(() => {
    return {
      picking: outbounds.filter((o) => o.status === "picking" || o.status === "created").length,
      checking: outbounds.filter((o) =>
        ["picked", "checking", "checked"].includes(o.status)
      ).length,
      packing: outbounds.filter((o) => ["packed", "weighed"].includes(o.status)).length,
      total: data?.total ?? outbounds.length,
    };
  }, [outbounds, data?.total]);

  const columns: ProColumn<OutboundOrder>[] = [
    {
      key: "order_no",
      title: "关联订单",
      width: 200,
      render: (row) => row.order_no ?? row.orderNo ?? "-",
    },
    { key: "id", title: "出库单号", width: 180 },
    {
      key: "status",
      title: "出库状态",
      width: 110,
      render: (row) => {
        const meta = statusLabels[row.status];
        return <Badge variant={meta?.variant ?? "secondary"}>{meta?.label ?? row.status}</Badge>;
      },
    },
    {
      key: "items",
      title: "商品数",
      width: 80,
      align: "center",
      render: (row) => (Array.isArray(row.items) ? row.items.length : "-"),
    },
    {
      key: "created_at",
      title: "创建时间",
      width: 180,
      render: (row) => row.created_at ?? row.createdAt ?? "-",
    },
  ];

  function handleProcess(row: OutboundOrder) {
    const status = row.status;
    if (status === "created") toast.info("出库单已创建，等待波次分配");
    else if (status === "picking" || status === "picked") toast.info("拣货进行中，查看 PDA 端进度");
    else if (status === "checked" || status === "packed") toast.info("请继续打包/称重操作");
    else toast.info(`当前状态: ${statusLabels[status]?.label || status}`);
  }

  return (
    <div className="space-y-4">
      <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
        {[
          { label: "全部出库单", value: stats.total, color: "text-foreground" },
          { label: "拣货任务", value: stats.picking, color: "text-primary" },
          { label: "复核任务", value: stats.checking, color: "text-amber-600" },
          { label: "待打包", value: stats.packing, color: "text-emerald-600" },
        ].map((item) => (
          <Card key={item.label}>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium text-muted-foreground">
                {item.label}
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className={`text-center text-3xl font-bold tabular-nums ${item.color}`}>
                {item.value}
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      <PageCard
        title="出库单列表"
        actions={
          <Button size="sm" variant="outline" onClick={() => toast.info("创建出库单功能：请从已审核订单触发出库")}>
            创建出库单
          </Button>
        }
      >
        <ProTable
          columns={columns}
          data={outbounds}
          loading={isLoading}
          error={error ? "加载出库单失败" : undefined}
          total={data?.total ?? 0}
          page={page}
          pageSize={20}
          rowKey={(row) => row.id}
          onPageChange={setPage}
          onRefresh={() => refetch()}
          actions={(row) => (
            <>
              <Button
                size="sm"
                disabled={row.status === "shipped"}
                onClick={() => handleProcess(row)}
              >
                处理
              </Button>
              <Button
                size="sm"
                variant="destructive"
                disabled={row.status === "shipped" || row.status === "abnormal"}
                onClick={() => toast.warning(`出库单 ${row.id} 标记异常 — 功能待接入`)}
              >
                异常
              </Button>
            </>
          )}
        />
      </PageCard>
    </div>
  );
}
