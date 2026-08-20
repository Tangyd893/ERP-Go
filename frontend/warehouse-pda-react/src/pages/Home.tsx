import { useMemo } from "react";
import { useNavigate } from "react-router-dom";
import { Button, Card, CardContent, Skeleton } from "@erp/shared-react";
import { useOutbounds } from "@/hooks/useOutbounds";
import { filterOutbounds, STATUS_FILTERS } from "@/lib/warehouse-api";
import { ErrorState } from "@erp/shared-react";

const cards = [
  { key: "pick", label: "拣货任务", path: "/pick", color: "text-pda-primary", filter: STATUS_FILTERS.pick },
  { key: "check", label: "复核任务", path: "/check", color: "text-pda-warning", filter: STATUS_FILTERS.check },
  { key: "pack", label: "打包", path: "/pack", color: "text-pda-text-secondary", filter: STATUS_FILTERS.pack },
  { key: "weigh", label: "称重", path: "/weigh", color: "text-pda-danger", filter: STATUS_FILTERS.weigh },
  {
    key: "ship",
    label: "待出库",
    path: "/ship",
    color: "text-pda-success",
    filter: null as null,
  },
] as const;

export default function Home() {
  const navigate = useNavigate();
  const { data: outbounds = [], isLoading, isError, refetch } = useOutbounds();

  const counts = useMemo(() => {
    const pendingShip = outbounds.filter(
      (o) => o.status !== "shipped" && o.status !== "cancelled"
    ).length;
    return {
      pick: filterOutbounds(outbounds, STATUS_FILTERS.pick).length,
      check: filterOutbounds(outbounds, STATUS_FILTERS.check).length,
      pack: filterOutbounds(outbounds, STATUS_FILTERS.pack).length,
      weigh: filterOutbounds(outbounds, STATUS_FILTERS.weigh).length,
      ship: pendingShip,
    };
  }, [outbounds]);

  if (isLoading) {
    return (
      <div>
        <h3 className="mb-3 text-lg font-semibold">WMS 仓库作业</h3>
        <div className="space-y-3">
          <Skeleton className="h-20 w-full" />
          <Skeleton className="h-20 w-full" />
          <Skeleton className="h-20 w-full" />
        </div>
      </div>
    );
  }

  if (isError) {
    return (
      <div>
        <h3 className="mb-3 text-lg font-semibold">WMS 仓库作业</h3>
        <ErrorState
          description="加载出库任务失败，请下拉刷新重试"
          onRetry={() => void refetch()}
        />
      </div>
    );
  }

  if (outbounds.length === 0) {
    return (
      <div>
        <h3 className="mb-3 text-lg font-semibold">WMS 仓库作业</h3>
        <div className="py-12 text-center">
          <p className="mb-4 text-sm text-pda-text-secondary">暂无出库任务</p>
          <Button onClick={() => void refetch()}>刷新</Button>
        </div>
      </div>
    );
  }

  return (
    <div>
      <h3 className="mb-3 text-lg font-semibold">WMS 仓库作业</h3>
      <div className="space-y-3">
        {cards.map((card) => (
          <Card
            key={card.key}
            className="pda-touch-target cursor-pointer transition-shadow hover:shadow-md"
            onClick={() => navigate(card.path)}
          >
            <CardContent className="flex flex-col items-center justify-center py-4">
              <div className="text-base">{card.label}</div>
              <div className={`mt-1 text-[28px] font-bold ${card.color}`}>
                {counts[card.key]}
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
}
