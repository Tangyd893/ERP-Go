import { useQuery } from "@tanstack/react-query";
import {
  ProTable,
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
  type ProColumn,
  apiClient,
  type ApiResponse,
  type PageData,
} from "@erp/shared-react";
import { PageCard } from "@/components/PageCard";

interface InventoryBalance {
  sku_id: string;
  sku_code: string;
  sku_name: string;
  warehouse_name: string;
  qty: number;
  locked_qty: number;
  available_qty: number;
}

interface InventoryJournal {
  journal_id: string;
  sku_code?: string;
  change_type?: string;
  change_qty: number;
  before_total?: number;
  after_total?: number;
  before_avail?: number;
  after_avail?: number;
  created_at: string;
}

const changeLabels: Record<string, { label: string; color: string }> = {
  lock: { label: "锁定", color: "text-amber-600" },
  release: { label: "释放", color: "text-emerald-600" },
  deduct: { label: "扣减", color: "text-destructive" },
  increase: { label: "入库", color: "text-primary" },
};

export default function InventoryList() {
  const balancesQuery = useQuery({
    queryKey: ["inventory-balances"],
    queryFn: async () => {
      const res = await apiClient.get<ApiResponse<PageData<InventoryBalance>>>(
        "/inventory/balances"
      );
      return res.data.data;
    },
  });

  const journalsQuery = useQuery({
    queryKey: ["inventory-journals"],
    queryFn: async () => {
      const res = await apiClient.get<ApiResponse<PageData<InventoryJournal>>>(
        "/inventory/journals"
      );
      return res.data.data;
    },
  });

  const balanceColumns: ProColumn<InventoryBalance>[] = [
    { key: "sku_code", title: "SKU编码", width: 150 },
    { key: "sku_name", title: "SKU名称" },
    { key: "warehouse_name", title: "仓库", width: 120 },
    {
      key: "qty",
      title: "总库存",
      width: 100,
      align: "right",
      render: (row) => <span className="font-semibold">{row.qty}</span>,
    },
    {
      key: "locked_qty",
      title: "已锁定",
      width: 100,
      align: "right",
      render: (row) => (
        <span className={row.locked_qty > 0 ? "text-amber-600" : "text-muted-foreground"}>
          {row.locked_qty}
        </span>
      ),
    },
    {
      key: "available_qty",
      title: "可用库存",
      width: 100,
      align: "right",
      render: (row) => (
        <span
          className={`font-semibold ${row.available_qty < 20 ? "text-destructive" : "text-emerald-600"}`}
        >
          {row.available_qty}
        </span>
      ),
    },
  ];

  const journalColumns: ProColumn<InventoryJournal>[] = [
    { key: "sku_code", title: "SKU", width: 140 },
    {
      key: "change_type",
      title: "类型",
      width: 80,
      render: (row) => {
        const meta = changeLabels[row.change_type ?? ""];
        return (
          <span className={`font-semibold ${meta?.color ?? ""}`}>
            {meta?.label ?? row.change_type}
          </span>
        );
      },
    },
    { key: "change_qty", title: "数量", width: 80, align: "center" },
    { key: "before_total", title: "变动前总数", width: 100, align: "right" },
    { key: "after_total", title: "变动后总数", width: 100, align: "right" },
    { key: "before_avail", title: "变动前可用", width: 110, align: "right" },
    { key: "after_avail", title: "变动后可用", width: 110, align: "right" },
    { key: "created_at", title: "时间", width: 170 },
  ];

  const balances = balancesQuery.data?.list ?? [];
  const journals = journalsQuery.data?.list ?? [];

  return (
    <div className="space-y-4">
      <PageCard title="库存余额">
        <ProTable
          columns={balanceColumns}
          data={balances}
          loading={balancesQuery.isLoading}
          error={balancesQuery.error ? "加载库存余额失败" : undefined}
          total={balances.length}
          rowKey={(row) => `${row.sku_id}-${row.warehouse_name}`}
          onRefresh={() => balancesQuery.refetch()}
        />
      </PageCard>

      <PageCard
        title="库存流水"
        actions={
          <Select disabled>
            <SelectTrigger className="h-8 w-36">
              <SelectValue placeholder="变动类型" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="lock">锁定</SelectItem>
              <SelectItem value="release">释放</SelectItem>
              <SelectItem value="deduct">扣减</SelectItem>
              <SelectItem value="increase">入库</SelectItem>
            </SelectContent>
          </Select>
        }
      >
        <ProTable
          columns={journalColumns}
          data={journals}
          loading={journalsQuery.isLoading}
          error={journalsQuery.error ? "加载库存流水失败" : undefined}
          total={journals.length}
          rowKey={(row) => row.journal_id}
          onRefresh={() => journalsQuery.refetch()}
        />
      </PageCard>
    </div>
  );
}
