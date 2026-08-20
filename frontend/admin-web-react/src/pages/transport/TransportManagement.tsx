import { useQuery } from "@tanstack/react-query";
import {
  Badge,
  Button,
  ProTable,
  type ProColumn,
  apiClient,
} from "@erp/shared-react";
import { PageCard } from "@/components/PageCard";
import { stubAction, unwrapList, unwrapPage } from "@/lib/utils";

interface Carrier {
  id: string;
  name: string;
  code: string;
  service?: string;
  status: string;
}

interface Shipment {
  id: string;
  tracking_no?: string;
  trackingNo?: string;
  carrier?: string;
  carrierName?: string;
  order_no?: string;
  orderNo?: string;
  weight?: number;
  status: string;
  created_at?: string;
  createdAt?: string;
}

const statusLabels: Record<string, { variant: "secondary" | "warning" | "success" | "destructive" | "default"; label: string }> = {
  pending: { variant: "secondary", label: "待获取面单" },
  labeled: { variant: "warning", label: "已获取面单" },
  shipped: { variant: "success", label: "已发货" },
  in_transit: { variant: "default", label: "运输中" },
  delivered: { variant: "success", label: "已签收" },
  failed: { variant: "destructive", label: "发货失败" },
};

export default function TransportManagement() {
  const carriersQuery = useQuery({
    queryKey: ["transport-carriers"],
    queryFn: async () => {
      const res = await apiClient.get("/transport/carriers");
      return unwrapList<Carrier>(res.data?.data ?? res.data);
    },
  });

  const shipmentsQuery = useQuery({
    queryKey: ["transport-shipments"],
    queryFn: async () => {
      const res = await apiClient.get("/transport/shipments");
      return unwrapPage<Shipment>(res.data?.data ?? res.data);
    },
  });

  const carrierColumns: ProColumn<Carrier>[] = [
    { key: "name", title: "物流商", width: 120 },
    { key: "code", title: "编码", width: 100 },
    { key: "service", title: "物流产品" },
    {
      key: "status",
      title: "状态",
      width: 100,
      render: (row) => (
        <Badge variant={row.status === "active" ? "success" : "destructive"}>
          {row.status === "active" ? "启用" : "禁用"}
        </Badge>
      ),
    },
  ];

  const shipmentColumns: ProColumn<Shipment>[] = [
    {
      key: "tracking_no",
      title: "运单号",
      width: 200,
      render: (row) => row.tracking_no ?? row.trackingNo ?? "-",
    },
    {
      key: "carrier",
      title: "物流商",
      width: 100,
      render: (row) => row.carrier ?? row.carrierName ?? "-",
    },
    {
      key: "order_no",
      title: "订单号",
      width: 200,
      render: (row) => row.order_no ?? row.orderNo ?? "-",
    },
    { key: "weight", title: "重量(kg)", width: 100 },
    {
      key: "status",
      title: "状态",
      width: 120,
      render: (row) => {
        const meta = statusLabels[row.status];
        return <Badge variant={meta?.variant ?? "secondary"}>{meta?.label ?? row.status}</Badge>;
      },
    },
    {
      key: "created_at",
      title: "创建时间",
      width: 180,
      render: (row) => row.created_at ?? row.createdAt ?? "-",
    },
  ];

  const carriers = carriersQuery.data ?? [];
  const shipments = shipmentsQuery.data?.list ?? [];

  return (
    <div className="space-y-4">
      <PageCard
        title="物流配置"
        actions={
          <Button size="sm" disabled onClick={() => stubAction("添加物流商")}>
            添加物流商
          </Button>
        }
      >
        <ProTable
          columns={carrierColumns}
          data={carriers}
          loading={carriersQuery.isLoading}
          error={carriersQuery.error ? "加载物流商失败" : undefined}
          total={carriers.length}
          pageSize={carriers.length || 10}
          rowKey={(row) => row.id}
          onRefresh={() => carriersQuery.refetch()}
        />
      </PageCard>

      <PageCard title="发运记录">
        <ProTable
          columns={shipmentColumns}
          data={shipments}
          loading={shipmentsQuery.isLoading}
          error={shipmentsQuery.error ? "加载发运记录失败" : undefined}
          total={shipments.length}
          pageSize={shipments.length || 10}
          rowKey={(row) => row.id}
          onRefresh={() => shipmentsQuery.refetch()}
        />
      </PageCard>
    </div>
  );
}
