import { useEffect, useMemo, useRef, useState } from "react";
import { useSearchParams } from "react-router-dom";
import { useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { Badge, Button, Card, CardContent, Input } from "@erp/shared-react";
import PageHeader from "@/components/PageHeader";
import ScanInput, { type ScanInputHandle } from "@/components/ScanInput";
import ScanFeedback from "@/components/ScanFeedback";
import { outboundQueryKey, usePickTasks } from "@/hooks/useOutbounds";
import { pickScan } from "@/lib/warehouse-api";
import { scanFeedback } from "@/lib/vibrate";
import { getErrorMessage, isDuplicateError, isNetworkError, useWarehouseStore } from "@/stores/warehouse";

export default function PickScan() {
  const [searchParams] = useSearchParams();
  const outboundId = searchParams.get("outbound_id") || "";
  const queryClient = useQueryClient();
  const scanInputRef = useRef<ScanInputHandle>(null);
  const [scanQty, setScanQty] = useState(1);
  const [scanning, setScanning] = useState(false);
  const [showHistory, setShowHistory] = useState(false);

  const online = useWarehouseStore((s) => s.online);
  const pickHistoryCount = useWarehouseStore((s) =>
    s.scanHistory.filter((r) => r.type === "pick").length
  );

  const { data: pickTasks = [], refetch } = usePickTasks(outboundId);

  const pendingTasks = useMemo(
    () => pickTasks.filter((t) => t.status !== "done" && t.picked_quantity < t.quantity),
    [pickTasks]
  );

  useEffect(() => {
    scanInputRef.current?.focus();
  }, [scanning]);

  function findTaskByScan(value: string) {
    const v = value.trim().toUpperCase();
    if (!v) return null;
    return (
      pendingTasks.find(
        (t) =>
          t.sku_code?.toUpperCase() === v ||
          t.location_code?.toUpperCase() === v ||
          t.id?.toUpperCase() === v
      ) ?? null
    );
  }

  async function runPick(taskId: string, qty: number, label: string) {
    if (scanning) return;
    setScanning(true);
    try {
      await pickScan(taskId, qty, label);
      toast.success(`✓ ${label} 拣货完成`);
      scanFeedback(true);
      await refetch();
      await queryClient.invalidateQueries({ queryKey: outboundQueryKey });
    } catch (e: unknown) {
      if (isDuplicateError(e)) {
        toast.info(`已拣货（重复扫描）: ${label}`);
        scanFeedback(true, true);
        await refetch();
        await queryClient.invalidateQueries({ queryKey: outboundQueryKey });
      } else if (isNetworkError(e)) {
        toast.error("网络不可用，请检查连接后重试");
        scanFeedback(false);
      } else {
        toast.error(getErrorMessage(e, "拣货失败"));
        scanFeedback(false);
      }
    } finally {
      setScanning(false);
    }
  }

  async function handleScan(value: string) {
    const task = findTaskByScan(value);
    if (!task) {
      toast.warning(`未找到匹配的待拣任务: ${value}`);
      scanFeedback(false);
      scanInputRef.current?.focus();
      return;
    }
    const label = `${task.sku_name || task.sku_code} @ ${task.location_code || "-"}`;
    await runPick(task.id, scanQty, label);
  }

  return (
    <div>
      <PageHeader
        title="拣货扫码"
        right={
          <>
            <button
              type="button"
              className="relative min-h-touch px-2 text-sm text-primary"
              onClick={() => setShowHistory(true)}
            >
              记录
              {pickHistoryCount > 0 ? (
                <span className="absolute -right-0.5 -top-0.5 flex h-4 min-w-4 items-center justify-center rounded-full bg-primary px-1 text-[10px] text-primary-foreground">
                  {pickHistoryCount}
                </span>
              ) : null}
            </button>
            <Badge variant={online ? "success" : "destructive"} className="text-xs">
              {online ? "在线" : "离线"}
            </Badge>
          </>
        }
      />

      {!outboundId ? (
        <p className="py-12 text-center text-sm text-pda-text-secondary">请从拣货列表选择出库单</p>
      ) : (
        <>
          <ScanInput
            ref={scanInputRef}
            hint="扫描条码 / 输入 SKU 编码 · 库位号 · 任务 ID"
            disabled={scanning}
            onScan={(v) => void handleScan(v)}
          />

          <div className="mb-3 flex items-center gap-2 px-1">
            <span className="text-[13px] text-pda-text-secondary">数量:</span>
            <Input
              type="number"
              min={1}
              max={9999}
              value={scanQty}
              className="h-9 w-28"
              onChange={(e) => setScanQty(Math.max(1, Number(e.target.value) || 1))}
            />
          </div>

          {pendingTasks.length === 0 && pickTasks.length > 0 ? (
            <p className="py-6 text-center font-medium text-pda-success">全部拣货任务已完成 ✓</p>
          ) : null}

          {pendingTasks.map((task) => {
            const label = task.sku_name || task.sku_code;
            return (
              <Card key={task.id} className="mb-2">
                <CardContent className="flex items-center justify-between gap-2 py-3">
                  <div className="min-w-0 flex-1">
                    <div className="truncate text-[15px] font-medium">{label}</div>
                    <div className="mt-0.5 text-xs text-pda-text-secondary">
                      {task.sku_code} · 库位 {task.location_code || "-"} ·{" "}
                      {task.picked_quantity || 0}/{task.quantity}
                    </div>
                  </div>
                  <Button
                    size="sm"
                    className="min-h-touch shrink-0"
                    disabled={scanning}
                    onClick={() => void runPick(task.id, scanQty, label)}
                  >
                    拣货
                  </Button>
                </CardContent>
              </Card>
            );
          })}
        </>
      )}

      {showHistory ? (
        <ScanFeedback type="pick" onClose={() => setShowHistory(false)} />
      ) : null}
    </div>
  );
}
