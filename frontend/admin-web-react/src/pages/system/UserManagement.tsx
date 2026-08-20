import { useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
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
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  Input,
  Label,
  ProTable,
  type ProColumn,
  apiClient,
  type ApiResponse,
  type PageData,
} from "@erp/shared-react";
import { PageCard } from "@/components/PageCard";
import { stubAction } from "@/lib/utils";

interface UserRecord {
  user_id: string;
  username: string;
  nickname: string;
  email: string;
  status: number;
  created_at: string;
}

export default function UserManagement() {
  const queryClient = useQueryClient();
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [dialogOpen, setDialogOpen] = useState(false);
  const [form, setForm] = useState({
    username: "",
    password: "",
    nickname: "",
    email: "",
    phone: "",
  });

  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ["users", page, pageSize],
    queryFn: async () => {
      const res = await apiClient.get<ApiResponse<PageData<UserRecord>>>("/iam/users", {
        params: { page, page_size: pageSize },
      });
      return res.data.data;
    },
  });

  const createMutation = useMutation({
    mutationFn: async (payload: typeof form) => {
      const res = await apiClient.post<ApiResponse<UserRecord>>("/iam/users", payload);
      return res.data.data;
    },
    onSuccess: () => {
      toast.success("用户创建成功");
      setDialogOpen(false);
      queryClient.invalidateQueries({ queryKey: ["users"] });
    },
    onError: () => toast.error("创建用户失败"),
  });

  const deleteMutation = useMutation({
    mutationFn: async (userId: string) => {
      await apiClient.delete(`/iam/users/${userId}`);
    },
    onSuccess: () => {
      toast.success("用户已删除");
      queryClient.invalidateQueries({ queryKey: ["users"] });
    },
    onError: () => toast.error("删除用户失败"),
  });

  const columns: ProColumn<UserRecord>[] = [
    { key: "username", title: "用户名", width: 120 },
    { key: "nickname", title: "昵称", width: 140 },
    { key: "email", title: "邮箱", width: 200 },
    {
      key: "status",
      title: "状态",
      width: 100,
      render: (row) => (
        <Badge variant={row.status === 1 ? "success" : "destructive"}>
          {row.status === 1 ? "启用" : "禁用"}
        </Badge>
      ),
    },
    { key: "created_at", title: "创建时间", width: 160 },
  ];

  return (
    <PageCard
      title="用户管理"
      actions={
        <Button size="sm" onClick={() => setDialogOpen(true)}>
          新建用户
        </Button>
      }
    >
      <ProTable
        columns={columns}
        data={data?.list ?? []}
        loading={isLoading}
        error={error ? "加载用户失败" : undefined}
        total={data?.total ?? 0}
        page={page}
        pageSize={pageSize}
        rowKey={(row) => row.user_id}
        onPageChange={setPage}
        onPageSizeChange={(size) => {
          setPageSize(size);
          setPage(1);
        }}
        onRefresh={() => refetch()}
        actions={(row) => (
          <>
            <Button size="sm" variant="outline" disabled onClick={() => stubAction("编辑")}>
              编辑
            </Button>
            <Button size="sm" variant="outline" disabled onClick={() => stubAction("分配角色")}>
              分配角色
            </Button>
            <AlertDialog>
              <AlertDialogTrigger asChild>
                <Button size="sm" variant="destructive">
                  删除
                </Button>
              </AlertDialogTrigger>
              <AlertDialogContent>
                <AlertDialogHeader>
                  <AlertDialogTitle>确定要删除该用户吗？</AlertDialogTitle>
                  <AlertDialogDescription>此操作不可撤销。</AlertDialogDescription>
                </AlertDialogHeader>
                <AlertDialogFooter>
                  <AlertDialogCancel>取消</AlertDialogCancel>
                  <AlertDialogAction onClick={() => deleteMutation.mutate(row.user_id)}>
                    删除
                  </AlertDialogAction>
                </AlertDialogFooter>
              </AlertDialogContent>
            </AlertDialog>
          </>
        )}
      />

      <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>新建用户</DialogTitle>
          </DialogHeader>
          <form
            className="space-y-3"
            onSubmit={(e) => {
              e.preventDefault();
              createMutation.mutate(form);
            }}
          >
            {(
              [
                ["username", "用户名", true],
                ["password", "密码", true],
                ["nickname", "昵称", false],
                ["email", "邮箱", false],
                ["phone", "手机号", false],
              ] as const
            ).map(([key, label, required]) => (
              <div key={key} className="space-y-1">
                <Label>{label}</Label>
                <Input
                  required={required}
                  value={form[key]}
                  onChange={(e) => setForm((prev) => ({ ...prev, [key]: e.target.value }))}
                />
              </div>
            ))}
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
