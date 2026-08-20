import { useMemo, useRef, useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { Badge, Button, Card, CardContent } from "@erp/shared-react";
import PageHeader from "@/components/PageHeader";
import ScanInput, { type ScanInputHandle } from "@/components/ScanInput";
import ScanFeedback from "@/components/ScanFeedback";
import { outboundQueryKey, useOutbounds } from "@/hooks/useOutbounds";
import { checkScan, filterOutbounds, STATUS_FILTERS } from "@/lib/warehouse-api";
import { scanFeedback } from "@/lib/vibrate";
import { getErrorMessage, isDuplicateError, isNetworkError, useWarehouseStore } from "@/stores/warehouse";

export default function CheckScan() {
  const queryClient = useQueryClient();
  const scanInputRef = useRef<ScanInputHandle>(null);
  const [scanning, setScanning] = useState(false);
  const [showHistory, setShowHistory] = useState(false);

  const online = useWarehouseStore((s) => s.online);
  const checkHistoryCount = useWarehouseStore((s) =>
    s.scanHistory.filter((r) => r.type === "check").length
  );

  const { data: outbounds = [], refetch } = useOutbounds();
  const checkingOutbounds = useMemo(
    () => filterOutbounds(outbounds, STATUS_FILTERS.check),
    [outbounds]
  );

  function findOutbound(value: string) {
    const v = value.trim().toUpperCase();
    if (!v) return null;
    return (
      checkingOutbounds.find(
        (o) => o.order_no?.toUpperCase() === v || o.id?.toUpperCase() === v
      ) ?? null
    );
  }

  async function runCheck(outboundId: string, label: string) {
    if (scanning) return;
    setScanning(true);
    try {
      await checkScan(outboundId, "", 1, label);
      toast.success(`✓ ${label} 复核完成`);
      scanFeedback(true);
      await refetch();
      await queryClient.invalidateQueries({ queryKey: outboundQueryKey });
    } catch (e: unknown) {
      if (isDuplicateError(e)) {
        toast.info(`已复核（重复扫描）: ${label}`);
        scanFeedback(true, true);
        await refetch();
        await queryClient.invalidateQueries({ queryKey: outboundQueryKey });
      } else if (isNetworkError(e)) {
        toast.error("网络不可用，请检查连接后重试");
        scanFeedback(false);
      } else {
        toast.error(getErrorMessage(e, "复核失败"));
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
      toast.warning(`未找到匹配的待复核出库单: ${value}`);
      scanFeedback(false);
      scanInputRef.current?.focus();
      return;
    }
    await runCheck(ob.id, ob.order_no || ob.id);
  }

  return (
    <div>
      <PageHeader
        title="复核扫码"
        backTo="/"
        right={
          <>
            <button
              type="button"
              className="relative min-h-touch px-2 text-sm text-primary"
              onClick={() => setShowHistory(true)}
            >
              记录
              {checkHistoryCount > 0 ? (
                <span className="absolute -right-0.5 -top-0.5 flex h-4 min-w-4 items-center justify-center rounded-full bg-primary px-1 text-[10px] text-primary-foreground">
                  {checkHistoryCount}
                </span>
              ) : null}
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

      {checkingOutbounds.length === 0 ? (
        <p className="py-12 text-center text-sm text-pda-text-secondary">暂无待复核出库单</p>
      ) : (
        checkingOutbounds.map((ob) => (
          <Card key={ob.id} className="mb-2">
            <CardContent className="flex items-center justify-between gap-2 py-3">
              <div className="min-w-0 flex-1">
                <div className="text-[15px] font-medium">{ob.order_no}</div>
                <div className="mt-0.5 text-xs text-pda-text-secondary">{ob.id}</div>
              </div>
              <Button
                variant="secondary"
                size="sm"
                className="min-h-touch shrink-0 bg-pda-warning text-white hover:bg-pda-warning/90"
                disabled={scanning}
                onClick={() => void runCheck(ob.id, ob.order_no || ob.id)}
              >
                复核确认
              </Button>
            </CardContent>
          </Card>
        ))
      )}

      {showHistory ? (
        <ScanFeedback type="check" onClose={() => setShowHistory(false)} />
      ) : null}
    </div>
  );
}
