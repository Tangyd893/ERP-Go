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

interface RoleRecord {
  role_id: string;
  name: string;
  code: string;
  description: string;
  status: number;
  created_at: string;
}

const allPermissions = [
  { code: "user:read", name: "查看用户", checked: true },
  { code: "user:create", name: "创建用户", checked: true },
  { code: "user:update", name: "编辑用户", checked: true },
  { code: "user:disable", name: "禁用用户", checked: true },
  { code: "role:read", name: "查看角色", checked: true },
  { code: "role:create", name: "创建角色", checked: true },
  { code: "role:update", name: "编辑角色", checked: false },
  { code: "role:delete", name: "删除角色", checked: false },
];

export default function RoleManagement() {
  const queryClient = useQueryClient();
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [dialogOpen, setDialogOpen] = useState(false);
  const [permOpen, setPermOpen] = useState(false);
  const [currentRole, setCurrentRole] = useState<RoleRecord | null>(null);
  const [form, setForm] = useState({ name: "", code: "", description: "" });
  const [permissions, setPermissions] = useState(allPermissions);

  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ["roles", page, pageSize],
    queryFn: async () => {
      const res = await apiClient.get<ApiResponse<PageData<RoleRecord>>>("/iam/roles", {
        params: { page, page_size: pageSize },
      });
      return res.data.data;
    },
  });

  const createMutation = useMutation({
    mutationFn: async (payload: typeof form) => {
      const res = await apiClient.post<ApiResponse<RoleRecord>>("/iam/roles", payload);
      return res.data.data;
    },
    onSuccess: () => {
      toast.success("角色创建成功");
      setDialogOpen(false);
      queryClient.invalidateQueries({ queryKey: ["roles"] });
    },
    onError: () => toast.error("创建角色失败"),
  });

  const deleteMutation = useMutation({
    mutationFn: async (id: string) => {
      await apiClient.delete(`/iam/roles/${id}`);
    },
    onSuccess: () => {
      toast.success("角色已删除");
      queryClient.invalidateQueries({ queryKey: ["roles"] });
    },
    onError: () => toast.error("删除角色失败"),
  });

  const columns: ProColumn<RoleRecord>[] = [
    { key: "name", title: "角色名", width: 150 },
    { key: "code", title: "角色编码", width: 180 },
    { key: "description", title: "描述" },
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
      title="角色管理"
      actions={
        <Button size="sm" onClick={() => setDialogOpen(true)}>
          新建角色
        </Button>
      }
    >
      <ProTable
        columns={columns}
        data={data?.list ?? []}
        loading={isLoading}
        error={error ? "加载角色失败" : undefined}
        total={data?.total ?? 0}
        page={page}
        pageSize={pageSize}
        rowKey={(row) => row.role_id}
        onPageChange={setPage}
        onPageSizeChange={(size) => {
          setPageSize(size);
          setPage(1);
        }}
        onRefresh={() => refetch()}
        actions={(row) => (
          <>
            <Button
              size="sm"
              variant="outline"
              onClick={() => {
                setCurrentRole(row);
                setPermOpen(true);
              }}
            >
              分配权限
            </Button>
            <Button size="sm" variant="outline" disabled onClick={() => stubAction("编辑")}>
              编辑
            </Button>
            <AlertDialog>
              <AlertDialogTrigger asChild>
                <Button size="sm" variant="destructive">
                  删除
                </Button>
              </AlertDialogTrigger>
              <AlertDialogContent>
                <AlertDialogHeader>
                  <AlertDialogTitle>确定要删除该角色吗？</AlertDialogTitle>
                  <AlertDialogDescription>此操作不可撤销。</AlertDialogDescription>
                </AlertDialogHeader>
                <AlertDialogFooter>
                  <AlertDialogCancel>取消</AlertDialogCancel>
                  <AlertDialogAction onClick={() => deleteMutation.mutate(row.role_id)}>
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
            <DialogTitle>新建角色</DialogTitle>
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
                ["name", "角色名", true],
                ["code", "编码", true],
                ["description", "描述", false],
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

      <Dialog open={permOpen} onOpenChange={setPermOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>分配权限</DialogTitle>
          </DialogHeader>
          {currentRole ? (
            <div className="space-y-3">
              <p className="text-sm">
                角色: <Badge>{currentRole.name}</Badge>
              </p>
              <div className="max-h-64 space-y-2 overflow-y-auto">
                {permissions.map((perm) => (
                  <label key={perm.code} className="flex items-center gap-2 text-sm">
                    <input
                      type="checkbox"
                      checked={perm.checked}
                      onChange={(e) =>
                        setPermissions((prev) =>
                          prev.map((p) =>
                            p.code === perm.code ? { ...p, checked: e.target.checked } : p
                          )
                        )
                      }
                    />
                    <span>
                      {perm.name}{" "}
                      <span className="text-muted-foreground">({perm.code})</span>
                    </span>
                  </label>
                ))}
              </div>
              <div className="flex justify-end gap-2">
                <Button variant="outline" onClick={() => setPermOpen(false)}>
                  取消
                </Button>
                <Button onClick={() => setPermOpen(false)}>保存</Button>
              </div>
            </div>
          ) : null}
        </DialogContent>
      </Dialog>
    </PageCard>
  );
}
