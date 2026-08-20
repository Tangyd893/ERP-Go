import { useMemo } from "react";
import { X } from "lucide-react";
import { useWarehouseStore } from "@/stores/warehouse";
import { cn } from "@erp/shared-react";

interface ScanFeedbackProps {
  type?: "pick" | "check" | "pack" | "weigh";
  limit?: number;
  onClose: () => void;
}

export default function ScanFeedback({ type, limit = 50, onClose }: ScanFeedbackProps) {
  const scanHistory = useWarehouseStore((s) => s.scanHistory);

  const records = useMemo(() => {
    let list = scanHistory;
    if (type) {
      list = list.filter((r) => r.type === type);
    }
    return list.slice(0, limit);
  }, [scanHistory, type, limit]);

  return (
    <div className="fixed inset-0 z-[var(--z-modal,60)] flex flex-col justify-end">
      <button
        type="button"
        aria-label="关闭"
        className="absolute inset-0 bg-black/40"
        onClick={onClose}
      />
      <div className="relative max-h-[60vh] rounded-t-xl bg-white p-4 pb-[calc(1rem+var(--pda-safe-bottom))] shadow-xl">
        <div className="mb-3 flex items-center justify-between">
          <h3 className="text-base font-semibold text-pda-text">扫码历史</h3>
          <button
            type="button"
            className="pda-touch-target flex items-center justify-center rounded-md text-pda-text-secondary"
            onClick={onClose}
          >
            <X className="h-5 w-5" />
          </button>
        </div>
        {records.length === 0 ? (
          <p className="py-6 text-center text-sm text-pda-text-secondary">暂无扫码记录</p>
        ) : (
          <ul className="max-h-[45vh] overflow-y-auto">
            {records.map((r) => (
              <li
                key={r.id}
                className="flex items-center justify-between border-b border-gray-100 py-2.5 text-[13px]"
              >
                <span
                  className={cn(
                    "w-6 shrink-0 text-center font-bold",
                    r.success ? "text-pda-success" : "text-pda-danger"
                  )}
                >
                  {r.success ? "✓" : "✗"}
                </span>
                <span className="mx-2 min-w-0 flex-1 truncate">{r.targetLabel}</span>
                <span className="shrink-0 text-[11px] text-pda-text-secondary">
                  {new Date(r.time).toLocaleTimeString()}
                </span>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
}
