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

interface ChannelStore {
  id: string;
  name: string;
  platform_code: string;
  site: string;
  auth_status: string;
  status: string;
  created_at: string;
}

export default function StoreManagement() {
  const queryClient = useQueryClient();
  const [dialogOpen, setDialogOpen] = useState(false);
  const [form, setForm] = useState({
    name: "",
    platform_code: "amazon_us",
    site: "",
    store_code: "",
  });

  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ["channel-stores"],
    queryFn: async () => {
      const res = await apiClient.get<ApiResponse<PageData<ChannelStore>>>("/channel/stores");
      return res.data.data;
    },
  });

  const createMutation = useMutation({
    mutationFn: async (payload: typeof form) => {
      const res = await apiClient.post<ApiResponse<ChannelStore>>("/channel/stores", payload);
      return res.data.data;
    },
    onSuccess: () => {
      toast.success("店铺添加成功");
      setDialogOpen(false);
      queryClient.invalidateQueries({ queryKey: ["channel-stores"] });
    },
    onError: () => toast.error("添加店铺失败"),
  });

  const list = data?.list ?? [];
  const columns: ProColumn<ChannelStore>[] = [
    { key: "name", title: "店铺名称" },
    { key: "platform_code", title: "平台", width: 150 },
    { key: "site", title: "站点", width: 150 },
    {
      key: "auth_status",
      title: "授权状态",
      width: 120,
      render: (row) => (
        <Badge variant={row.auth_status === "authorized" ? "success" : "warning"}>
          {row.auth_status === "authorized" ? "已授权" : "未授权"}
        </Badge>
      ),
    },
    {
      key: "status",
      title: "店铺状态",
      width: 120,
      render: () => <Badge variant="success">运营中</Badge>,
    },
  ];

  return (
    <PageCard
      title="店铺授权管理"
      actions={
        <Button size="sm" onClick={() => setDialogOpen(true)}>
          添加店铺
        </Button>
      }
    >
      <ProTable
        columns={columns}
        data={list}
        loading={isLoading}
        error={error ? "加载店铺失败" : undefined}
        total={list.length}
        rowKey={(row) => row.id}
        onRefresh={() => refetch()}
        actions={() => (
          <>
            <Button size="sm" variant="outline" disabled onClick={() => stubAction("授权")}>
              授权
            </Button>
            <Button size="sm" variant="outline" disabled onClick={() => stubAction("同步订单")}>
              同步订单
            </Button>
            <Button size="sm" variant="outline" disabled onClick={() => stubAction("编辑")}>
              编辑
            </Button>
          </>
        )}
      />

      <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>添加店铺</DialogTitle>
          </DialogHeader>
          <form
            className="space-y-3"
            onSubmit={(e) => {
              e.preventDefault();
              createMutation.mutate(form);
            }}
          >
            <div className="space-y-1">
              <Label>店铺名称</Label>
              <Input
                required
                placeholder="如: 美国站店铺A"
                value={form.name}
                onChange={(e) => setForm((p) => ({ ...p, name: e.target.value }))}
              />
            </div>
            <div className="space-y-1">
              <Label>平台</Label>
              <Select
                value={form.platform_code}
                onValueChange={(v) => setForm((p) => ({ ...p, platform_code: v }))}
              >
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="amazon_us">Amazon 美国</SelectItem>
                  <SelectItem value="amazon_uk">Amazon 英国</SelectItem>
                  <SelectItem value="amazon_jp">Amazon 日本</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-1">
              <Label>站点</Label>
              <Input
                required
                placeholder="如: Amazon.com"
                value={form.site}
                onChange={(e) => setForm((p) => ({ ...p, site: e.target.value }))}
              />
            </div>
            <div className="space-y-1">
              <Label>店铺编码</Label>
              <Input
                placeholder="店铺唯一标识"
                value={form.store_code}
                onChange={(e) => setForm((p) => ({ ...p, store_code: e.target.value }))}
              />
            </div>
            <div className="flex justify-end gap-2 pt-2">
              <Button type="button" variant="outline" onClick={() => setDialogOpen(false)}>
                取消
              </Button>
              <Button type="submit" disabled={createMutation.isPending}>
                提交
              </Button>
            </div>
          </form>
        </DialogContent>
      </Dialog>
    </PageCard>
  );
}
