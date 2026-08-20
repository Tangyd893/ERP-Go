import type { ReactNode } from "react";
import { Button } from "./ui/button";
import { AlertCircle } from "lucide-react";

export function ErrorState({
  title = "加载失败",
  description,
  onRetry,
  action,
}: {
  title?: string;
  description?: string;
  onRetry?: () => void;
  action?: ReactNode;
}) {
  return (
    <div className="flex flex-col items-center justify-center gap-3 rounded-md border border-dashed px-6 py-12 text-center">
      <AlertCircle className="size-8 text-muted-foreground" aria-hidden />
      <div className="space-y-1">
        <h3 className="text-base font-semibold text-balance">{title}</h3>
        {description ? (
          <p className="text-sm text-muted-foreground text-pretty">{description}</p>
        ) : null}
      </div>
      {onRetry ? (
        <Button variant="outline" onClick={onRetry}>
          重试
        </Button>
      ) : (
        action
      )}
    </div>
  );
}
