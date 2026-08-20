import { useMemo, useRef, useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { Badge, Button, Card, CardContent, Input } from "@erp/shared-react";
import PageHeader from "@/components/PageHeader";
import ScanInput, { type ScanInputHandle } from "@/components/ScanInput";
import ScanFeedback from "@/components/ScanFeedback";
import { outboundQueryKey, useOutbounds } from "@/hooks/useOutbounds";
import { filterOutbounds, packScan, STATUS_FILTERS } from "@/lib/warehouse-api";
import { scanFeedback } from "@/lib/vibrate";
import { getErrorMessage, isDuplicateError, isNetworkError, useWarehouseStore } from "@/stores/warehouse";

export default function PackScan() {
  const queryClient = useQueryClient();
  const scanInputRef = useRef<ScanInputHandle>(null);
  const [packingWeight, setPackingWeight] = useState<Record<string, number>>({});
  const [scanning, setScanning] = useState(false);
  const [showHistory, setShowHistory] = useState(false);

  const online = useWarehouseStore((s) => s.online);
  const { data: outbounds = [], refetch } = useOutbounds();

  const packedOutbounds = useMemo(
    () => filterOutbounds(outbounds, STATUS_FILTERS.pack),
    [outbounds]
  );

  function findOutbound(value: string) {
    const v = value.trim().toUpperCase();
    if (!v) return null;
    return (
      packedOutbounds.find(
        (o) => o.order_no?.toUpperCase() === v || o.id?.toUpperCase() === v
      ) ?? null
    );
  }

  async function runPack(outboundId: string, label: string) {
    if (scanning) return;
    setScanning(true);
    const weight = packingWeight[outboundId] || 0;
    try {
      await packScan(outboundId, weight);
      toast.success(`✓ ${label} 打包完成`);
      scanFeedback(true);
      setPackingWeight((prev) => ({ ...prev, [outboundId]: 0 }));
      await refetch();
      await queryClient.invalidateQueries({ queryKey: outboundQueryKey });
    } catch (e: unknown) {
      if (isDuplicateError(e)) {
        toast.info(`已打包（重复操作）: ${label}`);
        scanFeedback(true, true);
        await refetch();
        await queryClient.invalidateQueries({ queryKey: outboundQueryKey });
      } else if (isNetworkError(e)) {
        toast.error("网络不可用");
        scanFeedback(false);
      } else {
        toast.error(getErrorMessage(e, "打包失败"));
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
      toast.warning(`未找到待打包出库单: ${value}`);
      scanFeedback(false);
      scanInputRef.current?.focus();
      return;
    }
    await runPack(ob.id, ob.order_no || ob.id);
  }

  return (
    <div>
      <PageHeader
        title="打包"
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

      {packedOutbounds.length === 0 ? (
        <p className="py-12 text-center text-sm text-pda-text-secondary">暂无待打包出库单</p>
      ) : (
        packedOutbounds.map((ob) => (
          <Card key={ob.id} className="mb-2">
            <CardContent className="py-3">
              <div className="text-[15px] font-medium">{ob.order_no}</div>
              <div className="my-2 text-[13px] text-pda-text-secondary">{ob.id}</div>
              <div className="flex items-center gap-2">
                <Input
                  type="number"
                  min={0}
                  step={1}
                  placeholder="重量(g)"
                  value={packingWeight[ob.id] ?? ""}
                  className="h-9 w-28"
                  onChange={(e) =>
                    setPackingWeight((prev) => ({
                      ...prev,
                      [ob.id]: Number(e.target.value) || 0,
                    }))
                  }
                />
                <Button
                  size="sm"
                  className="min-h-touch"
                  disabled={scanning}
                  onClick={() => void runPack(ob.id, ob.order_no || ob.id)}
                >
                  确认打包
                </Button>
              </div>
            </CardContent>
          </Card>
        ))
      )}

      {showHistory ? (
        <ScanFeedback type="pack" onClose={() => setShowHistory(false)} />
      ) : null}
    </div>
  );
}
