import { useQuery } from "@tanstack/react-query";
import {
  Badge,
  Button,
  ProTable,
  type ProColumn,
  apiClient,
} from "@erp/shared-react";
import { PageCard } from "@/components/PageCard";
import { stubAction, unwrapList } from "@/lib/utils";

interface Supplier {
  id: string;
  name: string;
  code: string;
  contact: string;
  phone: string;
  status: string;
}

interface PurchaseOrder {
  id: string;
  order_no: string;
  supplier_name: string;
  status: string;
  total_amount: number;
  currency: string;
  created_at: string;
}

const statusMap: Record<string, { variant: "secondary" | "warning" | "success" | "default"; label: string }> = {
  draft: { variant: "secondary", label: "草稿" },
  pending: { variant: "warning", label: "待审核" },
  approved: { variant: "default", label: "已审核" },
  ordered: { variant: "success", label: "已下单" },
  completed: { variant: "success", label: "已完成" },
};

export default function PurchaseManagement() {
  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ["purchase"],
    queryFn: async () => {
      const [supRes, ordRes] = await Promise.allSettled([
        apiClient.get("/purchase/suppliers", { params: { page: 1, page_size: 100 } }),
        apiClient.get("/purchase/orders", { params: { page: 1, page_size: 100 } }),
      ]);
      return {
        suppliers:
          supRes.status === "fulfilled"
            ? unwrapList<Supplier>(supRes.value.data?.data ?? supRes.value.data)
            : [],
        orders:
          ordRes.status === "fulfilled"
            ? unwrapList<PurchaseOrder>(ordRes.value.data?.data ?? ordRes.value.data)
            : [],
      };
    },
  });

  const suppliers = data?.suppliers ?? [];
  const orders = data?.orders ?? [];

  const supplierColumns: ProColumn<Supplier>[] = [
    { key: "name", title: "名称", width: 180 },
    { key: "code", title: "编码", width: 120 },
    { key: "contact", title: "联系人" },
    { key: "phone", title: "电话", width: 150 },
  ];

  const orderColumns: ProColumn<PurchaseOrder>[] = [
    { key: "order_no", title: "采购单号", width: 200 },
    { key: "supplier_name", title: "供应商", width: 150 },
    {
      key: "status",
      title: "状态",
      width: 120,
      render: (row) => {
        const meta = statusMap[row.status];
        return <Badge variant={meta?.variant ?? "secondary"}>{meta?.label ?? row.status}</Badge>;
      },
    },
    {
      key: "total_amount",
      title: "金额",
      width: 120,
      render: (row) => `${row.currency} ${row.total_amount}`,
    },
    { key: "created_at", title: "创建时间" },
  ];

  return (
    <div className="space-y-4">
      <PageCard
        title="供应商"
        actions={
          <Button size="sm" disabled onClick={() => stubAction("添加供应商")}>
            添加供应商
          </Button>
        }
      >
        <ProTable
          columns={supplierColumns}
          data={suppliers}
          loading={isLoading}
          error={error ? "加载采购数据失败" : undefined}
          total={suppliers.length}
          pageSize={suppliers.length || 10}
          rowKey={(row) => row.id}
          onRefresh={() => refetch()}
          actions={() => (
            <>
              <Button size="sm" variant="outline" disabled onClick={() => stubAction("编辑")}>
                编辑
              </Button>
              <Button size="sm" variant="destructive" disabled onClick={() => stubAction("禁用")}>
                禁用
              </Button>
            </>
          )}
        />
      </PageCard>

      <PageCard
        title="采购单"
        actions={
          <Button size="sm" disabled onClick={() => stubAction("新建采购单")}>
            新建采购单
          </Button>
        }
      >
        <ProTable
          columns={orderColumns}
          data={orders}
          loading={isLoading}
          total={orders.length}
          pageSize={orders.length || 10}
          rowKey={(row) => row.id}
          onRefresh={() => refetch()}
          actions={() => (
            <>
              <Button size="sm" disabled onClick={() => stubAction("审核")}>
                审核
              </Button>
              <Button size="sm" variant="outline" disabled onClick={() => stubAction("收货")}>
                收货
              </Button>
            </>
          )}
        />
      </PageCard>
    </div>
  );
}
