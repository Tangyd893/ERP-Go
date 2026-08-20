import { useMemo } from "react";
import { useNavigate } from "react-router-dom";
import { Card, CardContent, Skeleton } from "@erp/shared-react";
import PageHeader from "@/components/PageHeader";
import { useOutbounds } from "@/hooks/useOutbounds";
import { filterOutbounds, STATUS_FILTERS } from "@/lib/warehouse-api";

export default function PickList() {
  const navigate = useNavigate();
  const { data: outbounds = [], isLoading } = useOutbounds();

  const pickableOutbounds = useMemo(
    () => filterOutbounds(outbounds, STATUS_FILTERS.pick),
    [outbounds]
  );

  function openPick(outboundId: string) {
    navigate({ pathname: "/pick/scan", search: `?outbound_id=${encodeURIComponent(outboundId)}` });
  }

  return (
    <div>
      <PageHeader title="拣货任务" backTo="/" />

      {isLoading ? (
        <div className="space-y-3 py-2">
          <Skeleton className="h-16 w-full" />
          <Skeleton className="h-16 w-full" />
        </div>
      ) : pickableOutbounds.length === 0 ? (
        <p className="py-12 text-center text-sm text-pda-text-secondary">暂无拣货任务</p>
      ) : (
        pickableOutbounds.map((ob) => (
          <Card
            key={ob.id}
            className="mb-3 cursor-pointer transition-shadow hover:shadow-md"
            onClick={() => openPick(ob.id)}
          >
            <CardContent className="py-3">
              <div className="text-[15px] font-semibold">{ob.order_no}</div>
              <div className="mt-1 truncate text-[13px] text-pda-text-secondary">
                {ob.id} · {ob.status}
              </div>
            </CardContent>
          </Card>
        ))
      )}
    </div>
  );
}
