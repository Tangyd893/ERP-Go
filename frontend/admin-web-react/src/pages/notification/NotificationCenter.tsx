import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { toast } from "sonner";
import {
  Badge,
  Button,
  Skeleton,
  apiClient,
  isDemo,
  type ApiResponse,
  type PageData,
} from "@erp/shared-react";
import { PageCard } from "@/components/PageCard";

interface Notification {
  id: string;
  title: string;
  content: string;
  type: string;
  read: boolean;
  created_at: string;
}

function getTagVariant(type: string): "secondary" | "warning" | "success" {
  if (type === "info") return "secondary";
  if (type === "warning") return "warning";
  return "success";
}

function getTypeLabel(type: string): string {
  if (type === "info") return "提醒";
  if (type === "warning") return "告警";
  return "完成";
}

export default function NotificationCenter() {
  const [page, setPage] = useState(1);
  const pageSize = 10;
  const [localRead, setLocalRead] = useState<Record<string, boolean>>({});

  const { data, isLoading, isError, refetch } = useQuery({
    queryKey: ["notifications", page],
    queryFn: async () => {
      if (isDemo()) {
        return { list: [] as Notification[], total: 0 };
      }
      const res = await apiClient.get<ApiResponse<PageData<Notification>>>("/notification/list", {
        params: { page, page_size: pageSize },
      });
      return {
        list: res.data.data?.list ?? [],
        total: res.data.data?.total ?? 0,
      };
    },
    meta: {
      onError: () => toast.warning("通知服务暂不可用"),
    },
  });

  const notifications = (data?.list ?? []).map((n) => ({
    ...n,
    read: localRead[n.id] ?? n.read,
  }));
  const unreadCount = notifications.filter((n) => !n.read).length;

  function markAllRead() {
    const next: Record<string, boolean> = {};
    notifications.forEach((n) => {
      next[n.id] = true;
    });
    setLocalRead((prev) => ({ ...prev, ...next }));
  }

  return (
    <PageCard
      title={
        <span className="inline-flex items-center gap-2">
          通知中心
          {unreadCount > 0 ? (
            <Badge variant="destructive" className="tabular-nums">
              {unreadCount}
            </Badge>
          ) : null}
        </span>
      }
      actions={
        <Button size="sm" variant="outline" onClick={markAllRead}>
          全部已读
        </Button>
      }
    >
      {isLoading ? (
        <div className="space-y-3">
          <Skeleton className="h-12 w-full" />
          <Skeleton className="h-12 w-full" />
        </div>
      ) : isError ? (
        <p className="py-8 text-center text-sm text-muted-foreground">
          通知服务暂不可用，请稍后重试
        </p>
      ) : !isDemo() && notifications.length === 0 ? (
        <p className="py-8 text-center text-sm text-muted-foreground">暂无通知</p>
      ) : (
        <div className="divide-y">
          {notifications.map((n) => (
            <div
              key={n.id}
              className="flex items-center justify-between gap-4 py-3 text-sm"
            >
              <div className="min-w-0">
                <Badge variant={getTagVariant(n.type)} className="mr-2">
                  {getTypeLabel(n.type)}
                </Badge>
                <span className={n.read ? "font-normal" : "font-semibold"}>{n.title}</span>
                <span className="ml-2 text-muted-foreground">{n.content}</span>
              </div>
              <span className="shrink-0 text-xs text-muted-foreground">{n.created_at}</span>
            </div>
          ))}
        </div>
      )}

      {(data?.total ?? 0) > pageSize ? (
        <div className="mt-4 flex justify-end gap-2">
          <Button
            size="sm"
            variant="outline"
            disabled={page <= 1}
            onClick={() => setPage((p) => p - 1)}
          >
            上一页
          </Button>
          <Button
            size="sm"
            variant="outline"
            disabled={page * pageSize >= (data?.total ?? 0)}
            onClick={() => setPage((p) => p + 1)}
          >
            下一页
          </Button>
        </div>
      ) : null}

      {isError ? (
        <div className="mt-4 flex justify-center">
          <Button size="sm" variant="outline" onClick={() => refetch()}>
            重试
          </Button>
        </div>
      ) : null}
    </PageCard>
  );
}
