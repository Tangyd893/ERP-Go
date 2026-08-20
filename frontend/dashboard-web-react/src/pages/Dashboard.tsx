import { useCallback, useEffect, useRef, useState, type RefObject } from "react";
import { useQuery } from "@tanstack/react-query";
import axios from "axios";
import * as echarts from "echarts";
import {
  apiClient,
  isDemo,
  Button,
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  Skeleton,
} from "@erp/shared-react";
import {
  type DashboardData,
  DEMO_DASHBOARD,
  FALLBACK_TREND,
  FALLBACK_TIMELINESS,
  CHART_COLORS,
} from "../lib/dashboard-data";

const POLL_INTERVAL_MS = 30_000;

interface DashboardQueryResult {
  data: DashboardData | null;
  authRequired: boolean;
  error: string | null;
}

async function loadDashboard(forceDemo: boolean): Promise<DashboardQueryResult> {
  if (forceDemo || isDemo()) {
    return { data: DEMO_DASHBOARD, authRequired: false, error: null };
  }

  try {
    const res = await apiClient.get("/report/dashboard");
    const payload: DashboardData = res.data?.data ?? res.data;
    return {
      data: {
        order_count: payload.order_count ?? 0,
        sales_amount: payload.sales_amount ?? 0,
        outbound_count: payload.outbound_count ?? 0,
        sku_count: payload.sku_count ?? 0,
        trend: payload.trend ?? [],
        timeliness: payload.timeliness ?? { within_24h: 0, within_48h: 0, overdue: 0 },
      },
      authRequired: false,
      error: null,
    };
  } catch (err: unknown) {
    if (axios.isAxiosError(err) && err.response?.status === 401) {
      return {
        data: null,
        authRequired: true,
        error: "请先登录管理后台获取访问令牌",
      };
    }
    return { data: null, authRequired: false, error: "数据加载失败" };
  }
}

function useTrendChart(
  containerRef: RefObject<HTMLDivElement | null>,
  trend: DashboardData["trend"] | undefined,
  enabled: boolean
) {
  const chartRef = useRef<echarts.ECharts | null>(null);

  useEffect(() => {
    if (!enabled || !containerRef.current) return;

    if (!chartRef.current) {
      chartRef.current = echarts.init(containerRef.current);
      const onResize = () => chartRef.current?.resize();
      window.addEventListener("resize", onResize);
    }

    const dates = trend?.map((t) => t.date) ?? [];
    const orders = trend?.map((t) => t.order_count) ?? [];
    const sales = trend?.map((t) => t.sales_amount) ?? [];

    const useDates = dates.length > 0 ? dates : FALLBACK_TREND.dates;
    const useOrders = orders.length > 0 ? orders : FALLBACK_TREND.orders;
    const useSales = sales.length > 0 ? sales : FALLBACK_TREND.sales;

    chartRef.current.setOption({
      tooltip: { trigger: "axis" },
      legend: { data: ["订单量", "销售额(元)"] },
      grid: { left: 48, right: 48, top: 40, bottom: 32 },
      xAxis: { type: "category", data: useDates },
      yAxis: [
        { type: "value", name: "单量" },
        { type: "value", name: "金额(元)" },
      ],
      series: [
        {
          name: "订单量",
          type: "bar",
          data: useOrders,
          itemStyle: { color: CHART_COLORS.orders },
        },
        {
          name: "销售额(元)",
          type: "line",
          yAxisIndex: 1,
          data: useSales,
          itemStyle: { color: CHART_COLORS.sales },
          smooth: true,
        },
      ],
    });
  }, [containerRef, trend, enabled]);

  useEffect(() => {
    return () => {
      chartRef.current?.dispose();
      chartRef.current = null;
    };
  }, []);
}

function useTimelinessChart(
  containerRef: RefObject<HTMLDivElement | null>,
  timeliness: DashboardData["timeliness"] | undefined,
  enabled: boolean
) {
  const chartRef = useRef<echarts.ECharts | null>(null);

  useEffect(() => {
    if (!enabled || !containerRef.current) return;

    if (!chartRef.current) {
      chartRef.current = echarts.init(containerRef.current);
      const onResize = () => chartRef.current?.resize();
      window.addEventListener("resize", onResize);
    }

    const t = timeliness ?? { within_24h: 0, within_48h: 0, overdue: 0 };
    const total = t.within_24h + t.within_48h + t.overdue;
    const useData =
      total > 0
        ? [
            { value: t.within_24h, name: "24h内出库", itemStyle: { color: "#67C23A" } },
            { value: t.within_48h, name: "48h内出库", itemStyle: { color: "#E6A23C" } },
            { value: t.overdue, name: "超时出库", itemStyle: { color: "#F56C6C" } },
          ]
        : FALLBACK_TIMELINESS.map((item) => ({
            value: item.value,
            name: item.name,
            itemStyle: { color: item.color },
          }));

    chartRef.current.setOption({
      tooltip: { trigger: "item" },
      legend: { orient: "vertical", left: "left" },
      series: [
        {
          name: "出库时效",
          type: "pie",
          radius: ["50%", "75%"],
          label: { show: false },
          data: useData,
        },
      ],
    });
  }, [containerRef, timeliness, enabled]);

  useEffect(() => {
    return () => {
      chartRef.current?.dispose();
      chartRef.current = null;
    };
  }, []);
}

const KPI_STATS = [
  { key: "order_count" as const, label: "订单总数", format: (v: number) => v.toLocaleString(), color: "text-primary" },
  { key: "sales_amount" as const, label: "销售额", format: (v: number) => `¥${v.toLocaleString()}`, color: "text-emerald-600" },
  { key: "outbound_count" as const, label: "出库量", format: (v: number) => v.toLocaleString(), color: "text-amber-600" },
  { key: "sku_count" as const, label: "库存 SKU", format: (v: number) => v.toLocaleString(), color: "text-destructive" },
];

export default function Dashboard() {
  const [forceDemo, setForceDemo] = useState(false);
  const trendChartRef = useRef<HTMLDivElement>(null);
  const rateChartRef = useRef<HTMLDivElement>(null);

  const { data: result, isLoading, refetch } = useQuery({
    queryKey: ["dashboard", forceDemo],
    queryFn: () => loadDashboard(forceDemo),
    refetchInterval: forceDemo || isDemo() ? false : POLL_INTERVAL_MS,
  });

  const dashboard = result?.data;
  const authRequired = result?.authRequired ?? false;
  const error = result?.error ?? null;
  const chartsReady = !isLoading && !error && !authRequired && !!dashboard;

  useTrendChart(trendChartRef, dashboard?.trend, chartsReady);
  useTimelinessChart(rateChartRef, dashboard?.timeliness, chartsReady);

  const applyDemo = useCallback(() => {
    setForceDemo(true);
    void refetch();
  }, [refetch]);

  return (
    <div className="dashboard-shell">
      <header className="mb-6">
        <h1 className="text-xl font-semibold tracking-tight text-foreground">ERP-Go 经营看板</h1>
      </header>

      {authRequired && (
        <div
          role="alert"
          className="mb-6 rounded-lg border border-amber-200 bg-amber-50 p-4 text-sm text-amber-900"
        >
          <p className="font-medium">需要登录</p>
          <p className="mt-1 text-amber-800">
            看板数据需要有效的访问令牌。请先登录管理后台获取令牌，或启用 Demo 模式查看演示数据。
          </p>
          <Button type="button" size="sm" className="mt-3" onClick={applyDemo}>
            查看 Demo 数据
          </Button>
        </div>
      )}

      {error && !authRequired && (
        <div
          role="alert"
          className="mb-4 rounded-lg border border-destructive/30 bg-destructive/5 p-4 text-sm text-destructive"
        >
          {error}
        </div>
      )}

      <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
        {KPI_STATS.map((stat) => (
          <Card key={stat.key}>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium text-muted-foreground">{stat.label}</CardTitle>
            </CardHeader>
            <CardContent className="text-center">
              {isLoading ? (
                <Skeleton className="mx-auto h-9 w-24" />
              ) : (
                <div className={`text-3xl font-bold tabular-nums ${stat.color}`}>
                  {stat.format(dashboard?.[stat.key] ?? 0)}
                </div>
              )}
            </CardContent>
          </Card>
        ))}
      </div>

      <div className="mt-4 grid gap-4 lg:grid-cols-3">
        <Card className="lg:col-span-2">
          <CardHeader>
            <CardTitle className="text-base">销售趋势（近 7 天）</CardTitle>
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <Skeleton className="h-80 w-full" />
            ) : (
              <div ref={trendChartRef} className="h-80 w-full" />
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="text-base">出库及时率</CardTitle>
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <Skeleton className="h-80 w-full" />
            ) : (
              <div ref={rateChartRef} className="h-80 w-full" />
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
