import { useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import {
  Button,
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
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
import { stubAction } from "@/lib/utils";

interface SKUMappingRecord {
  id: string;
  sku_code: string;
  sku_name: string;
  platform_code: string;
  platform_sku: string;
  asin: string;
  fnsku: string;
  store: string;
}

export default function SKUMapping() {
  const queryClient = useQueryClient();
  const [dialogOpen, setDialogOpen] = useState(false);
  const [form, setForm] = useState({
    sku_id: "",
    store_id: "",
    platform_sku: "",
    asin: "",
    fnsku: "",
  });

  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ["sku-mappings"],
    queryFn: async () => {
      const res = await apiClient.get<ApiResponse<PageData<SKUMappingRecord>>>(
        "/products/sku-mappings"
      );
      return res.data.data;
    },
  });

  const createMutation = useMutation({
    mutationFn: async (payload: typeof form) => {
      const res = await apiClient.post<ApiResponse<SKUMappingRecord>>(
        "/products/sku-mappings",
        payload
      );
      return res.data.data;
    },
    onSuccess: () => {
      toast.success("映射创建成功");
      setDialogOpen(false);
      queryClient.invalidateQueries({ queryKey: ["sku-mappings"] });
    },
    onError: () => toast.error("创建映射失败"),
  });

  const list = data?.list ?? [];
  const columns: ProColumn<SKUMappingRecord>[] = [
    { key: "sku_code", title: "内部SKU", width: 140 },
    { key: "sku_name", title: "SKU名称" },
    { key: "store", title: "店铺", width: 150 },
    { key: "platform_code", title: "平台", width: 120 },
    { key: "platform_sku", title: "平台SKU", width: 140 },
    { key: "asin", title: "ASIN", width: 120 },
    { key: "fnsku", title: "FNSKU", width: 120 },
  ];

  return (
    <PageCard
      title="SKU 映射管理"
      actions={
        <Button size="sm" onClick={() => setDialogOpen(true)}>
          新建映射
        </Button>
      }
    >
      <ProTable
        columns={columns}
        data={list}
        loading={isLoading}
        error={error ? "加载映射失败" : undefined}
        total={list.length}
        rowKey={(row) => row.id}
        onRefresh={() => refetch()}
        actions={() => (
          <>
            <Button size="sm" variant="outline" disabled onClick={() => stubAction("编辑")}>
              编辑
            </Button>
            <Button size="sm" variant="destructive" disabled onClick={() => stubAction("删除")}>
              删除
            </Button>
          </>
        )}
      />

      <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>新建 SKU 映射</DialogTitle>
          </DialogHeader>
          <form
            className="space-y-3"
            onSubmit={(e) => {
              e.preventDefault();
              createMutation.mutate(form);
            }}
          >
            <div className="space-y-1">
              <Label>内部SKU</Label>
              <Select value={form.sku_id} onValueChange={(v) => setForm((p) => ({ ...p, sku_id: v }))}>
                <SelectTrigger>
                  <SelectValue placeholder="选择SKU" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="1">TSHIRT-001</SelectItem>
                  <SelectItem value="3">MUG-001</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-1">
              <Label>店铺</Label>
              <Select
                value={form.store_id}
                onValueChange={(v) => setForm((p) => ({ ...p, store_id: v }))}
              >
                <SelectTrigger>
                  <SelectValue placeholder="选择店铺" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="store-1">美国站店铺A</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-1">
              <Label>平台SKU</Label>
              <Input
                value={form.platform_sku}
                onChange={(e) => setForm((p) => ({ ...p, platform_sku: e.target.value }))}
              />
            </div>
            <div className="space-y-1">
              <Label>ASIN</Label>
              <Input value={form.asin} onChange={(e) => setForm((p) => ({ ...p, asin: e.target.value }))} />
            </div>
            <div className="space-y-1">
              <Label>FNSKU</Label>
              <Input
                value={form.fnsku}
                onChange={(e) => setForm((p) => ({ ...p, fnsku: e.target.value }))}
              />
            </div>
            <div className="flex justify-end gap-2 pt-2">
              <Button type="button" variant="outline" onClick={() => setDialogOpen(false)}>
                取消
              </Button>
              <Button type="submit" disabled={createMutation.isPending}>
                确定
              </Button>
            </div>
          </form>
        </DialogContent>
      </Dialog>
    </PageCard>
  );
}
