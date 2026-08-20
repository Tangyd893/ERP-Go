import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import {
  Badge,
  Button,
  Input,
  Label,
  ProTable,
  type ProColumn,
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
  apiClient,
  type ApiResponse,
  type PageData,
} from "@erp/shared-react";
import { PageCard } from "@/components/PageCard";

interface AuditRecord {
  log_id: string;
  username: string;
  action: string;
  resource_type: string;
  detail: string;
  ip: string;
  created_at: string;
}

const actionLabels: Record<string, { variant: "success" | "secondary" | "default" | "warning" | "destructive"; label: string }> = {
  login: { variant: "success", label: "登录" },
  logout: { variant: "secondary", label: "登出" },
  create: { variant: "default", label: "创建" },
  update: { variant: "warning", label: "更新" },
  delete: { variant: "destructive", label: "删除" },
  export: { variant: "warning", label: "导出" },
  permission_change: { variant: "destructive", label: "权限变更" },
};

export default function AuditLog() {
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [username, setUsername] = useState("");
  const [action, setAction] = useState("");
  const [startDate, setStartDate] = useState("");
  const [endDate, setEndDate] = useState("");

  const filters = { username, action, startDate, endDate };

  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ["audit-logs", page, pageSize, filters],
    queryFn: async () => {
      const res = await apiClient.get<ApiResponse<PageData<AuditRecord>>>("/iam/audit-logs", {
        params: {
          page,
          page_size: pageSize,
          username: username || undefined,
          action: action || undefined,
          start_date: startDate || undefined,
          end_date: endDate || undefined,
        },
      });
      return res.data.data;
    },
  });

  const columns: ProColumn<AuditRecord>[] = [
    { key: "username", title: "用户", width: 120 },
    {
      key: "action",
      title: "操作",
      width: 120,
      render: (row) => {
        const meta = actionLabels[row.action];
        return (
          <Badge variant={meta?.variant ?? "secondary"}>{meta?.label ?? row.action}</Badge>
        );
      },
    },
    { key: "resource_type", title: "资源类型", width: 120 },
    { key: "detail", title: "操作详情" },
    { key: "ip", title: "IP 地址", width: 140 },
    { key: "created_at", title: "操作时间", width: 180 },
  ];

  return (
    <PageCard title="操作审计日志">
      <div className="mb-4 grid gap-3 md:grid-cols-2 xl:grid-cols-4">
        <div className="space-y-1">
          <Label>用户</Label>
          <Input
            placeholder="搜索用户名"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
          />
        </div>
        <div className="space-y-1">
          <Label>操作类型</Label>
          <Select value={action || "all"} onValueChange={(v) => setAction(v === "all" ? "" : v)}>
            <SelectTrigger>
              <SelectValue placeholder="全部" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">全部</SelectItem>
              {Object.entries(actionLabels).map(([key, meta]) => (
                <SelectItem key={key} value={key}>
                  {meta.label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
        <div className="space-y-1">
          <Label>开始日期</Label>
          <Input type="date" value={startDate} onChange={(e) => setStartDate(e.target.value)} />
        </div>
        <div className="space-y-1">
          <Label>结束日期</Label>
          <Input type="date" value={endDate} onChange={(e) => setEndDate(e.target.value)} />
        </div>
      </div>
      <div className="mb-4 flex gap-2">
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
            setUsername("");
            setAction("");
            setStartDate("");
            setEndDate("");
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
        error={error ? "加载审计日志失败" : undefined}
        total={data?.total ?? 0}
        page={page}
        pageSize={pageSize}
        rowKey={(row) => row.log_id}
        onPageChange={setPage}
        onPageSizeChange={(size) => {
          setPageSize(size);
          setPage(1);
        }}
        onRefresh={() => refetch()}
      />
    </PageCard>
  );
}
