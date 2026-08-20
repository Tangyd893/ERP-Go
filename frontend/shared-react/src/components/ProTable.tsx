import type { ReactNode } from "react";
import { Button } from "./ui/button";
import { Skeleton } from "./ui/skeleton";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "./ui/table";

export interface ProColumn<T> {
  key: string;
  title: string;
  width?: string | number;
  align?: "left" | "center" | "right";
  render?: (row: T) => ReactNode;
  accessor?: (row: T) => ReactNode;
}

export interface ProTableProps<T> {
  columns: ProColumn<T>[];
  data: T[];
  loading?: boolean;
  error?: string;
  total?: number;
  page?: number;
  pageSize?: number;
  rowKey?: (row: T, index: number) => string;
  toolbar?: ReactNode;
  emptyAction?: ReactNode;
  emptyTitle?: string;
  onPageChange?: (page: number) => void;
  onPageSizeChange?: (size: number) => void;
  onRefresh?: () => void;
  actions?: (row: T) => ReactNode;
}

export function ProTable<T extends object>({
  columns,
  data,
  loading,
  error,
  total = 0,
  page = 1,
  pageSize = 20,
  rowKey,
  toolbar,
  emptyAction,
  emptyTitle = "暂无数据",
  onPageChange,
  onPageSizeChange,
  onRefresh,
  actions,
}: ProTableProps<T>) {
  const totalPages = Math.max(1, Math.ceil(total / pageSize));

  return (
    <div className="space-y-3">
      {toolbar ? <div className="flex flex-wrap items-center gap-2">{toolbar}</div> : null}

      {error ? (
        <div className="flex items-center justify-between rounded-md border border-destructive/30 bg-destructive/5 px-3 py-2 text-sm text-destructive">
          <span>{error}</span>
          {onRefresh ? (
            <Button size="sm" variant="outline" onClick={onRefresh}>
              重试
            </Button>
          ) : null}
        </div>
      ) : null}

      {loading ? (
        <div className="space-y-2 rounded-md border p-3">
          <Skeleton className="h-8 w-full" />
          <Skeleton className="h-8 w-full" />
          <Skeleton className="h-8 w-3/4" />
        </div>
      ) : null}

      {!loading && !error && data.length === 0 ? (
        <div className="flex flex-col items-center justify-center gap-3 rounded-md border border-dashed px-6 py-12 text-center">
          <p className="text-sm text-muted-foreground text-pretty">{emptyTitle}</p>
          {emptyAction}
        </div>
      ) : null}

      {!loading && !error && data.length > 0 ? (
        <div className="rounded-md border bg-card">
          <Table>
            <TableHeader>
              <TableRow>
                {columns.map((col) => (
                  <TableHead
                    key={col.key}
                    style={col.width ? { width: col.width } : undefined}
                    className={
                      col.align === "right"
                        ? "text-right"
                        : col.align === "center"
                          ? "text-center"
                          : undefined
                    }
                  >
                    {col.title}
                  </TableHead>
                ))}
                {actions ? <TableHead className="text-right">操作</TableHead> : null}
              </TableRow>
            </TableHeader>
            <TableBody>
              {data.map((row, index) => (
                <TableRow key={rowKey?.(row, index) ?? String(index)}>
                  {columns.map((col) => (
                    <TableCell
                      key={col.key}
                      className={
                        col.align === "right"
                          ? "text-right"
                          : col.align === "center"
                            ? "text-center"
                            : undefined
                      }
                    >
                      {col.render
                        ? col.render(row)
                        : col.accessor
                          ? col.accessor(row)
                          : String((row as Record<string, unknown>)[col.key] ?? "")}
                    </TableCell>
                  ))}
                  {actions ? (
                    <TableCell className="text-right">
                      <div className="flex justify-end gap-1">{actions(row)}</div>
                    </TableCell>
                  ) : null}
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      ) : null}

      {total > 0 ? (
        <div className="flex flex-wrap items-center justify-between gap-2 text-sm text-muted-foreground">
          <span className="tabular-nums">共 {total} 条</span>
          <div className="flex items-center gap-2">
            <select
              className="h-8 rounded-md border bg-card px-2 text-sm"
              value={pageSize}
              onChange={(e) => onPageSizeChange?.(Number(e.target.value))}
              aria-label="每页条数"
            >
              {[10, 20, 50, 100].map((n) => (
                <option key={n} value={n}>
                  {n} / 页
                </option>
              ))}
            </select>
            <Button
              size="sm"
              variant="outline"
              disabled={page <= 1}
              onClick={() => onPageChange?.(page - 1)}
            >
              上一页
            </Button>
            <span className="tabular-nums">
              {page} / {totalPages}
            </span>
            <Button
              size="sm"
              variant="outline"
              disabled={page >= totalPages}
              onClick={() => onPageChange?.(page + 1)}
            >
              下一页
            </Button>
          </div>
        </div>
      ) : null}
    </div>
  );
}
