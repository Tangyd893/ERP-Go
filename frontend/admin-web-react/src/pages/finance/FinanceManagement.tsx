import { useQuery } from "@tanstack/react-query";
import {
  Badge,
  Button,
  ProTable,
  type ProColumn,
  apiClient,
} from "@erp/shared-react";
import { PageCard } from "@/components/PageCard";
import { stubAction, unwrapList } from "@/lib/utils";

interface Settlement {
  id: string;
  store: string;
  period: string;
  total_sales: number;
  commission: number;
  net_amount: number;
  status: string;
  currency: string;
}

interface Profit {
  order_no: string;
  sale_amount: number;
  purchase_cost: number;
  shipping_cost: number;
  commission: number;
  total_cost: number;
  gross_profit: number;
  profit_margin: number;
}

export default function FinanceManagement() {
  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ["finance"],
    queryFn: async () => {
      const [settleRes, profitRes] = await Promise.allSettled([
        apiClient.get("/finance/settlements", { params: { page: 1, page_size: 100 } }),
        apiClient.get("/finance/profit", { params: { page: 1, page_size: 100 } }),
      ]);
      return {
        settlements:
          settleRes.status === "fulfilled"
            ? unwrapList<Settlement>(settleRes.value.data?.data ?? settleRes.value.data)
            : [],
        profits:
          profitRes.status === "fulfilled"
            ? unwrapList<Profit>(profitRes.value.data?.data ?? profitRes.value.data)
            : [],
      };
    },
  });

  const settlements = data?.settlements ?? [];
  const profits = data?.profits ?? [];

  const settlementColumns: ProColumn<Settlement>[] = [
    { key: "store", title: "店铺", width: 150 },
    { key: "period", title: "结算周期", width: 200 },
    {
      key: "total_sales",
      title: "销售额",
      width: 130,
      render: (row) => `${row.currency} ${row.total_sales}`,
    },
    {
      key: "commission",
      title: "佣金",
      width: 120,
      render: (row) => `${row.currency} ${row.commission}`,
    },
    {
      key: "net_amount",
      title: "净收入",
      width: 120,
      render: (row) => (
        <span className="font-semibold text-emerald-600">
          {row.currency} {row.net_amount}
        </span>
      ),
    },
    {
      key: "status",
      title: "状态",
      width: 100,
      render: () => <Badge variant="success">已匹配</Badge>,
    },
  ];

  const profitColumns: ProColumn<Profit>[] = [
    { key: "order_no", title: "订单号", width: 200 },
    { key: "sale_amount", title: "销售额", width: 100 },
    { key: "purchase_cost", title: "采购成本", width: 100 },
    { key: "shipping_cost", title: "物流费", width: 100 },
    { key: "commission", title: "佣金", width: 80 },
    { key: "total_cost", title: "总成本", width: 100 },
    {
      key: "gross_profit",
      title: "毛利",
      width: 80,
      render: (row) => (
        <span className={row.gross_profit > 0 ? "text-emerald-600" : "text-destructive"}>
          {row.gross_profit.toFixed(2)}
        </span>
      ),
    },
    {
      key: "profit_margin",
      title: "利润率",
      width: 80,
      render: (row) => `${row.profit_margin.toFixed(1)}%`,
    },
  ];

  return (
    <div className="space-y-4">
      <PageCard
        title="平台结算"
        actions={
          <Button size="sm" disabled onClick={() => stubAction("导入结算报告")}>
            导入结算报告
          </Button>
        }
      >
        <ProTable
          columns={settlementColumns}
          data={settlements}
          loading={isLoading}
          error={error ? "加载财务数据失败" : undefined}
          total={settlements.length}
          pageSize={settlements.length || 10}
          rowKey={(row) => row.id}
          onRefresh={() => refetch()}
        />
      </PageCard>

      <PageCard title="订单利润分析">
        <ProTable
          columns={profitColumns}
          data={profits}
          loading={isLoading}
          total={profits.length}
          pageSize={profits.length || 10}
          rowKey={(row) => row.order_no}
          onRefresh={() => refetch()}
        />
      </PageCard>
    </div>
  );
}
