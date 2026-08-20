import { useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useNavigate } from "react-router-dom";
import { toast } from "sonner";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
  Badge,
  Button,
  Input,
  Label,
  ProTable,
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
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
  store?: string;
  buyer_name?: string;
  total_amount: number;
  currency?: string;
  items?: unknown[] | number;
  outbound_no?: string;
  status: string;
  shipped_at?: string;
  created_at: string;
}

const statusLabels: Record<string, { variant: "secondary" | "warning" | "success" | "destructive" | "default"; label: string }> = {
  pending: { variant: "secondary", label: "待审核" },
  approved: { variant: "warning", label: "已审核" },
  locked: { variant: "default", label: "已锁定" },
  picking: { variant: "default", label: "拣货中" },
  packed: { variant: "default", label: "已打包" },
  shipped: { variant: "success", label: "已发货" },
  abnormal: { variant: "destructive", label: "异常" },
  cancelled: { variant: "secondary", label: "已取消" },
};

export default function OrderList() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [page, setPage] = useState(1);
  const [keyword, setKeyword] = useState("");
  const [status, setStatus] = useState("");
  const [abnormalReason, setAbnormalReason] = useState("");
  const [abnormalTarget, setAbnormalTarget] = useState<string | null>(null);

  const filters = { keyword, status };

  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ["orders", page, filters],
    queryFn: async () => {
      const res = await apiClient.get<ApiResponse<PageData<OrderRecord>>>("/order/orders", {
        params: {
          page,
          page_size: 20,
          keyword: keyword || undefined,
          status: status || undefined,
        },
      });
      return res.data.data;
    },
  });

  const auditMutation = useMutation({
    mutationFn: async (id: string) => {
      await apiClient.post("/order/orders/audit", { order_id: id, approved: true });
    },
    onSuccess: () => {
      toast.success("审核通过");
      queryClient.invalidateQueries({ queryKey: ["orders"] });
    },
    onError: () => toast.error("审核失败"),
  });

  const abnormalMutation = useMutation({
    mutationFn: async ({ id, reason }: { id: string; reason: string }) => {
      await apiClient.post("/order/orders/abnormal", { order_id: id, reason });
    },
    onSuccess: () => {
      toast.success("已标记为异常");
      setAbnormalTarget(null);
      setAbnormalReason("");
      queryClient.invalidateQueries({ queryKey: ["orders"] });
    },
    onError: () => toast.error("标记异常失败"),
  });

  const cancelMutation = useMutation({
    mutationFn: async (id: string) => {
      await apiClient.post("/order/orders/cancel", { order_id: id, reason: "手动取消" });
    },
    onSuccess: () => {
      toast.success("订单已取消");
      queryClient.invalidateQueries({ queryKey: ["orders"] });
    },
    onError: () => toast.error("取消订单失败"),
  });

  const columns: ProColumn<OrderRecord>[] = [
    { key: "platform_order_no", title: "平台订单号", width: 200 },
    { key: "store", title: "店铺", width: 150 },
    { key: "buyer_name", title: "买家", width: 130 },
    {
      key: "total_amount",
      title: "金额",
      width: 100,
      align: "right",
      render: (row) => `${row.currency || "¥"} ${Number(row.total_amount || 0).toFixed(2)}`,
    },
    {
      key: "items",
      title: "件数",
      width: 70,
      align: "center",
      render: (row) =>
        Array.isArray(row.items) ? row.items.length : Number(row.items || 0),
    },
    { key: "outbound_no", title: "出库单号", width: 150 },
    {
      key: "status",
      title: "状态",
      width: 100,
      render: (row) => {
        const meta = statusLabels[row.status];
        return <Badge variant={meta?.variant ?? "secondary"}>{meta?.label ?? row.status}</Badge>;
      },
    },
    { key: "shipped_at", title: "发货时间", width: 170 },
    { key: "created_at", title: "创建时间", width: 170 },
  ];

  return (
    <PageCard
      title="订单列表"
      actions={
        <Button size="sm" onClick={() => navigate("/channel/import")}>
          导入订单
        </Button>
      }
    >
      <div className="mb-4 flex flex-wrap items-end gap-3">
        <div className="space-y-1">
          <Label>搜索</Label>
          <Input
            placeholder="订单号 / 买家"
            value={keyword}
            onChange={(e) => setKeyword(e.target.value)}
            className="w-52"
          />
        </div>
        <div className="space-y-1">
          <Label>状态</Label>
          <Select value={status || "all"} onValueChange={(v) => setStatus(v === "all" ? "" : v)}>
            <SelectTrigger className="w-36">
              <SelectValue placeholder="全部" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">全部</SelectItem>
              {Object.entries(statusLabels).map(([key, meta]) => (
                <SelectItem key={key} value={key}>
                  {meta.label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
        <Button
          size="sm"
          onClick={() => {
            setPage(1);
            refetch();
          }}
        >
          查询
        </Button>
        <Button
          size="sm"
          variant="outline"
          onClick={() => {
            setKeyword("");
            setStatus("");
            setPage(1);
          }}
        >
          重置
        </Button>
      </div>

      <ProTable
        columns={columns}
        data={data?.list ?? []}
        loading={isLoading}
        error={error ? "加载订单失败" : undefined}
        total={data?.total ?? 0}
        page={page}
        pageSize={20}
        rowKey={(row) => orderId(row)}
        onPageChange={setPage}
        onRefresh={() => refetch()}
        actions={(row) => {
          const id = orderId(row);
          return (
            <>
              <Button
                size="sm"
                disabled={row.status !== "pending"}
                onClick={() => auditMutation.mutate(id)}
              >
                审核
              </Button>
              <Button
                size="sm"
                variant="outline"
                disabled={["shipped", "cancelled", "abnormal"].includes(row.status)}
                onClick={() => setAbnormalTarget(id)}
              >
                异常
              </Button>
              <AlertDialog>
                <AlertDialogTrigger asChild>
                  <Button
                    size="sm"
                    variant="destructive"
                    disabled={["cancelled", "shipped"].includes(row.status)}
                  >
                    取消
                  </Button>
                </AlertDialogTrigger>
                <AlertDialogContent>
                  <AlertDialogHeader>
                    <AlertDialogTitle>确认取消该订单？</AlertDialogTitle>
                    <AlertDialogDescription>此操作不可撤销。</AlertDialogDescription>
                  </AlertDialogHeader>
                  <AlertDialogFooter>
                    <AlertDialogCancel>返回</AlertDialogCancel>
                    <AlertDialogAction onClick={() => cancelMutation.mutate(id)}>
                      确认取消
                    </AlertDialogAction>
                  </AlertDialogFooter>
                </AlertDialogContent>
              </AlertDialog>
            </>
          );
        }}
      />

      <AlertDialog open={!!abnormalTarget} onOpenChange={(open) => !open && setAbnormalTarget(null)}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>标记异常</AlertDialogTitle>
            <AlertDialogDescription>请输入异常原因</AlertDialogDescription>
          </AlertDialogHeader>
          <textarea
            className="min-h-24 w-full rounded-md border bg-background px-3 py-2 text-sm"
            value={abnormalReason}
            onChange={(e) => setAbnormalReason(e.target.value)}
          />
          <AlertDialogFooter>
            <AlertDialogCancel>取消</AlertDialogCancel>
            <AlertDialogAction
              onClick={() => {
                if (abnormalTarget && abnormalReason.trim()) {
                  abnormalMutation.mutate({ id: abnormalTarget, reason: abnormalReason.trim() });
                }
              }}
            >
              确定
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </PageCard>
  );
}
