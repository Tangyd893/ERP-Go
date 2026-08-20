import { useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import {
  Button,
  ProTable,
  type ProColumn,
  apiClient,
  type ApiResponse,
  type PageData,
} from "@erp/shared-react";
import { PageCard } from "@/components/PageCard";
import { orderId } from "@/lib/utils";

interface OrderRecord {
  id?: string;
  order_id?: string;
  platform_order_no?: string;
  order_no?: string;
  store?: string;
  buyer_name?: string;
  total_amount: number;
  currency?: string;
  status: string;
  created_at: string;
}

export default function OrderAudit() {
  const queryClient = useQueryClient();
  const [page, setPage] = useState(1);

  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ["orders-audit", page],
    queryFn: async () => {
      const res = await apiClient.get<ApiResponse<PageData<OrderRecord>>>("/order/orders", {
        params: { status: "pending", page, page_size: 20 },
      });
      return res.data.data;
    },
  });

  const auditMutation = useMutation({
    mutationFn: async ({ id, approved }: { id: string; approved: boolean }) => {
      await apiClient.post("/order/orders/audit", { order_id: id, approved });
    },
    onSuccess: (_data, vars) => {
      toast.success(vars.approved ? "审核通过" : "已驳回");
      queryClient.invalidateQueries({ queryKey: ["orders-audit"] });
    },
    onError: () => toast.error("审核操作失败"),
  });

  const list = (data?.list ?? []).filter((o) => o.status === "pending");

  const columns: ProColumn<OrderRecord>[] = [
    { key: "platform_order_no", title: "平台订单号", width: 200 },
    { key: "store", title: "店铺", width: 150 },
    { key: "buyer_name", title: "买家", width: 120 },
    {
      key: "total_amount",
      title: "金额",
      width: 120,
      align: "right",
      render: (row) => `${row.currency || "¥"} ${Number(row.total_amount || 0).toFixed(2)}`,
    },
    { key: "created_at", title: "创建时间", width: 170 },
  ];

  return (
    <PageCard title="订单审核">
      <ProTable
        columns={columns}
        data={list}
        loading={isLoading}
        error={error ? "加载待审核订单失败" : undefined}
        total={list.length}
        page={page}
        pageSize={20}
        rowKey={(row) => orderId(row)}
        onPageChange={setPage}
        onRefresh={() => refetch()}
        actions={(row) => {
          const id = orderId(row);
          const label = row.platform_order_no || row.order_no || id;
          return (
            <>
              <Button
                size="sm"
                onClick={() => auditMutation.mutate({ id, approved: true })}
              >
                审核通过
              </Button>
              <Button
                size="sm"
                variant="destructive"
                onClick={() => auditMutation.mutate({ id, approved: false })}
              >
                驳回
              </Button>
              <Button
                size="sm"
                variant="outline"
                onClick={() => toast.info(`查看订单 ${label} 详情`)}
              >
                查看详情
              </Button>
            </>
          );
        }}
      />
    </PageCard>
  );
}
