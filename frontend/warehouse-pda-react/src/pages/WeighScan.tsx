import { useMemo, useRef, useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { Badge, Button, Card, CardContent, Input } from "@erp/shared-react";
import PageHeader from "@/components/PageHeader";
import ScanInput, { type ScanInputHandle } from "@/components/ScanInput";
import ScanFeedback from "@/components/ScanFeedback";
import { outboundQueryKey, useOutbounds } from "@/hooks/useOutbounds";
import { filterOutbounds, STATUS_FILTERS, weighScan } from "@/lib/warehouse-api";
import { scanFeedback } from "@/lib/vibrate";
import { getErrorMessage, isDuplicateError, isNetworkError, useWarehouseStore } from "@/stores/warehouse";

export default function WeighScan() {
  const queryClient = useQueryClient();
  const scanInputRef = useRef<ScanInputHandle>(null);
  const [weighWeight, setWeighWeight] = useState<Record<string, number>>({});
  const [scanning, setScanning] = useState(false);
  const [showHistory, setShowHistory] = useState(false);

  const online = useWarehouseStore((s) => s.online);
  const { data: outbounds = [], refetch } = useOutbounds();

  const weighedOutbounds = useMemo(
    () => filterOutbounds(outbounds, STATUS_FILTERS.weigh),
    [outbounds]
  );

  function findOutbound(value: string) {
    const v = value.trim().toUpperCase();
    if (!v) return null;
    return (
      weighedOutbounds.find(
        (o) => o.order_no?.toUpperCase() === v || o.id?.toUpperCase() === v
      ) ?? null
    );
  }

  async function runWeigh(outboundId: string, label: string) {
    if (scanning) return;
    const weight = weighWeight[outboundId];
    if (!weight || weight <= 0) {
      toast.warning("请输入有效重量");
      scanInputRef.current?.focus();
      return;
    }
    setScanning(true);
    try {
      await weighScan(outboundId, weight);
      toast.success(`✓ ${label} 称重完成`);
      scanFeedback(true);
      setWeighWeight((prev) => ({ ...prev, [outboundId]: 0 }));
      await refetch();
      await queryClient.invalidateQueries({ queryKey: outboundQueryKey });
    } catch (e: unknown) {
      if (isDuplicateError(e)) {
        toast.info(`已称重（重复操作）: ${label}`);
        scanFeedback(true, true);
        await refetch();
        await queryClient.invalidateQueries({ queryKey: outboundQueryKey });
      } else if (isNetworkError(e)) {
        toast.error("网络不可用");
        scanFeedback(false);
      } else {
        toast.error(getErrorMessage(e, "称重失败"));
        scanFeedback(false);
      }
    } finally {
      setScanning(false);
      scanInputRef.current?.focus();
    }
  }

  async function handleScan(value: string) {
    const ob = findOutbound(value);
    if (!ob) {
      toast.warning(`未找到待称重出库单: ${value}`);
      scanFeedback(false);
      scanInputRef.current?.focus();
      return;
    }
    await runWeigh(ob.id, ob.order_no || ob.id);
  }

  return (
    <div>
      <PageHeader
        title="称重"
        backTo="/"
        right={
          <>
            <button
              type="button"
              className="min-h-touch px-2 text-sm text-primary"
              onClick={() => setShowHistory(true)}
            >
              记录
            </button>
            <Badge variant={online ? "success" : "destructive"} className="text-xs">
              {online ? "在线" : "离线"}
            </Badge>
          </>
        }
      />

      <ScanInput
        ref={scanInputRef}
        hint="扫描出库单条码 / 输入订单号"
        disabled={scanning}
        onScan={(v) => void handleScan(v)}
      />

      {weighedOutbounds.length === 0 ? (
        <p className="py-12 text-center text-sm text-pda-text-secondary">暂无待称重出库单</p>
      ) : (
        weighedOutbounds.map((ob) => (
          <Card key={ob.id} className="mb-2">
            <CardContent className="py-3">
              <div className="text-[15px] font-medium">{ob.order_no}</div>
              <div className="my-2 text-[13px] text-pda-text-secondary">{ob.id}</div>
              <div className="flex items-center gap-2">
                <Input
                  type="number"
                  min={1}
                  step={1}
                  placeholder="重量(g)"
                  value={weighWeight[ob.id] ?? ""}
                  className="h-9 w-28"
                  onChange={(e) =>
                    setWeighWeight((prev) => ({
                      ...prev,
                      [ob.id]: Number(e.target.value) || 0,
                    }))
                  }
                />
                <Button
                  variant="secondary"
                  size="sm"
                  className="min-h-touch bg-pda-warning text-white hover:bg-pda-warning/90"
                  disabled={scanning}
                  onClick={() => void runWeigh(ob.id, ob.order_no || ob.id)}
                >
                  确认称重
                </Button>
              </div>
            </CardContent>
          </Card>
        ))
      )}

      {showHistory ? (
        <ScanFeedback type="weigh" onClose={() => setShowHistory(false)} />
      ) : null}
    </div>
  );
}
