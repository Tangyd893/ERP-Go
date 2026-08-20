import { useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import {
  Badge,
  Button,
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  FileUpload,
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
import { stubAction } from "@/lib/utils";

interface ImportTask {
  id: string;
  import_type: string;
  store_name: string;
  status: string;
  total_rows: number;
  success_rows: number;
  failed_rows: number;
  created_at: string;
}

export default function OrderImport() {
  const queryClient = useQueryClient();
  const [page, setPage] = useState(1);
  const pageSize = 20;
  const [uploadOpen, setUploadOpen] = useState(false);
  const [storeId, setStoreId] = useState("");

  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ["import-tasks", page],
    queryFn: async () => {
      const res = await apiClient.get<ApiResponse<PageData<ImportTask>>>(
        "/channel/import-tasks",
        { params: { page, page_size: pageSize } }
      );
      return res.data.data;
    },
  });

  const importMutation = useMutation({
    mutationFn: async (file: File) => {
      const formData = new FormData();
      formData.append("file", file);
      if (storeId) formData.append("store_id", storeId);
      const res = await apiClient.post<ApiResponse<ImportTask>>("/channel/import", formData, {
        headers: { "Content-Type": "multipart/form-data" },
      });
      return res.data.data;
    },
    onSuccess: () => {
      toast.success("导入任务已提交");
      setUploadOpen(false);
      queryClient.invalidateQueries({ queryKey: ["import-tasks"] });
    },
    onError: () => toast.error("导入失败"),
  });

  const list = data?.list ?? [];
  const columns: ProColumn<ImportTask>[] = [
    {
      key: "import_type",
      title: "导入方式",
      width: 100,
      render: (row) => (
        <Badge variant={row.import_type === "csv" ? "warning" : "default"}>
          {row.import_type.toUpperCase()}
        </Badge>
      ),
    },
    { key: "store_name", title: "店铺", width: 150 },
    {
      key: "status",
      title: "状态",
      width: 120,
      render: (row) => (
        <Badge variant={row.status === "completed" ? "success" : "warning"}>
          {row.status === "completed" ? "已完成" : "处理中"}
        </Badge>
      ),
    },
    { key: "total_rows", title: "总行数", width: 100, align: "right" },
    {
      key: "success_rows",
      title: "成功",
      width: 100,
      align: "right",
      render: (row) => <span className="text-emerald-600">{row.success_rows}</span>,
    },
    {
      key: "failed_rows",
      title: "失败",
      width: 100,
      align: "right",
      render: (row) => (
        <span className={row.failed_rows > 0 ? "text-destructive" : "text-muted-foreground"}>
          {row.failed_rows}
        </span>
      ),
    },
    { key: "created_at", title: "创建时间", width: 180 },
  ];

  return (
    <PageCard
      title="订单导入任务"
      actions={
        <>
          <Button size="sm" onClick={() => setUploadOpen(true)}>
            上传 CSV 导入
          </Button>
          <Button size="sm" variant="outline" disabled onClick={() => stubAction("API 同步")}>
            API 同步
          </Button>
        </>
      }
    >
      <ProTable
        columns={columns}
        data={list}
        loading={isLoading}
        error={error ? "加载导入任务失败" : undefined}
        total={data?.total ?? list.length}
        page={page}
        pageSize={pageSize}
        rowKey={(row) => row.id}
        onPageChange={setPage}
        onRefresh={() => refetch()}
        actions={(row) => (
          <Button size="sm" variant="outline" disabled={row.status !== "completed"}>
            查看结果
          </Button>
        )}
      />

      <Dialog open={uploadOpen} onOpenChange={setUploadOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>上传 CSV 订单导入</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <div className="space-y-1">
              <Label>店铺</Label>
              <Select value={storeId} onValueChange={setStoreId}>
                <SelectTrigger>
                  <SelectValue placeholder="选择店铺" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="store-1">美国站店铺A</SelectItem>
                  <SelectItem value="store-2">英国站店铺B</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-1">
              <Label>CSV 文件</Label>
              <FileUpload
                accept=".csv"
                label="选择 CSV 文件"
                onFile={async (file) => {
                  await importMutation.mutateAsync(file);
                }}
              />
              <p className="text-xs text-muted-foreground">
                文件需包含: OrderID, SKU, Quantity, BuyerName, Address, Phone 等字段
              </p>
            </div>
            <div className="flex justify-end gap-2">
              <Button variant="outline" onClick={() => setUploadOpen(false)}>
                取消
              </Button>
            </div>
          </div>
        </DialogContent>
      </Dialog>
    </PageCard>
  );
}
