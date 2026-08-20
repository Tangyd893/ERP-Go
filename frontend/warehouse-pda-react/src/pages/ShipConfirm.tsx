import { useMemo, useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { Button, Card, CardContent, Input, Skeleton } from "@erp/shared-react";
import PageHeader from "@/components/PageHeader";
import { outboundQueryKey, useOutbounds } from "@/hooks/useOutbounds";
import {
  confirmShip,
  filterOutbounds,
  packScan,
  STATUS_FILTERS,
  weighScan,
} from "@/lib/warehouse-api";
import { getErrorMessage, isDuplicateError } from "@/stores/warehouse";

export default function ShipConfirm() {
  const queryClient = useQueryClient();
  const [trackingNo, setTrackingNo] = useState("");
  const [shippingId, setShippingId] = useState<string | null>(null);

  const { data: outbounds = [], isLoading, refetch } = useOutbounds();

  const shippableOutbounds = useMemo(
    () => filterOutbounds(outbounds, STATUS_FILTERS.ship),
    [outbounds]
  );

  async function handleShip(outboundId: string) {
    setShippingId(outboundId);
    try {
      try {
        await packScan(outboundId);
      } catch (e) {
        if (!isDuplicateError(e)) throw e;
      }
      try {
        await weighScan(outboundId, 1);
      } catch (e) {
        if (!isDuplicateError(e)) throw e;
      }
      await confirmShip(outboundId, trackingNo, "default");
      toast.success("出库确认成功");
      await refetch();
      await queryClient.invalidateQueries({ queryKey: outboundQueryKey });
    } catch (e: unknown) {
      toast.error(getErrorMessage(e, "出库确认失败"));
    } finally {
      setShippingId(null);
    }
  }

  return (
    <div>
      <PageHeader title="出库确认" backTo="/" />

      <Input
        value={trackingNo}
        placeholder="物流单号（可选）"
        className="mb-3 min-h-touch"
        onChange={(e) => setTrackingNo(e.target.value)}
      />

      {isLoading ? (
        <div className="space-y-3 py-2">
          <Skeleton className="h-16 w-full" />
          <Skeleton className="h-16 w-full" />
        </div>
      ) : shippableOutbounds.length === 0 ? (
        <p className="py-12 text-center text-sm text-pda-text-secondary">暂无待出库单</p>
      ) : (
        shippableOutbounds.map((ob) => (
          <Card key={ob.id} className="mb-3">
            <CardContent className="py-3">
              <div className="text-[15px] font-semibold">{ob.order_no}</div>
              <div className="mb-2 mt-1 truncate text-[13px] text-pda-text-secondary">
                {ob.id} · {ob.status}
              </div>
              <Button
                className="min-h-touch w-full bg-pda-success hover:bg-pda-success/90"
                disabled={shippingId === ob.id}
                onClick={() => void handleShip(ob.id)}
              >
                {shippingId === ob.id ? "处理中…" : "确认出库"}
              </Button>
            </CardContent>
          </Card>
        ))
      )}
    </div>
  );
}
