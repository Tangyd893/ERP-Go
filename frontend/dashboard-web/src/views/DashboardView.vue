<script setup lang="ts">
import { onMounted, onUnmounted, ref, watch, nextTick } from "vue";
import { useDashboardStore } from "@/stores/dashboard";
import * as echarts from "echarts";

const store = useDashboardStore();
const trendChartRef = ref<HTMLDivElement | null>(null);
const rateChartRef = ref<HTMLDivElement | null>(null);

let trendChart: echarts.ECharts | null = null;
let rateChart: echarts.ECharts | null = null;

onMounted(async () => {
  await store.fetchAll();
  await nextTick();
  if (!store.loading && !store.error && !store.authRequired) {
    initCharts();
  }
  // T-638: 30 秒轮询
  store.startPolling(30000);
});

onUnmounted(() => {
  store.stopPolling();
  trendChart?.dispose();
  rateChart?.dispose();
});

watch(
  () => [store.loading, store.error, store.authRequired],
  async ([loading, error, auth]) => {
    if (!loading && !error && !auth) {
      await nextTick();
      initCharts();
    }
  }
);

function initCharts() {
  initTrendChart();
  initRateChart();
}

function initTrendChart() {
  if (!trendChartRef.value) return;
  if (!trendChart) {
    trendChart = echarts.init(trendChartRef.value);
    window.addEventListener("resize", () => trendChart?.resize());
  }

  const dates = store.trend.value.map((t) => t.date);
  const orders = store.trend.value.map((t) => t.order_count);
  const sales = store.trend.value.map((t) => t.sales_amount);

  // fallback 硬编码数据（API 无返回时）
  const useDates = dates.length > 0 ? dates : ["周一", "周二", "周三", "周四", "周五", "周六", "周日"];
  const useOrders = orders.length > 0 ? orders : [45, 52, 38, 65, 48, 32, 28];
  const useSales = sales.length > 0 ? sales : [2100, 3200, 1850, 4200, 2800, 1500, 1200];

  trendChart.setOption({
    tooltip: { trigger: "axis" },
    legend: { data: ["订单量", "销售额(元)"] },
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
        itemStyle: { color: "#409EFF" },
      },
      {
        name: "销售额(元)",
        type: "line",
        yAxisIndex: 1,
        data: useSales,
        itemStyle: { color: "#67C23A" },
      },
    ],
  });
}

function initRateChart() {
  if (!rateChartRef.value) return;
  if (!rateChart) {
    rateChart = echarts.init(rateChartRef.value);
    window.addEventListener("resize", () => rateChart?.resize());
  }

  const t = store.timeliness.value;
  const total = t.within_24h + t.within_48h + t.overdue;
  const useData =
    total > 0
      ? [
          { value: t.within_24h, name: "24h内出库", itemStyle: { color: "#67C23A" } },
          { value: t.within_48h, name: "48h内出库", itemStyle: { color: "#E6A23C" } },
          { value: t.overdue, name: "超时出库", itemStyle: { color: "#F56C6C" } },
        ]
      : [
          { value: 152, name: "24h内出库", itemStyle: { color: "#67C23A" } },
          { value: 28, name: "48h内出库", itemStyle: { color: "#E6A23C" } },
          { value: 6, name: "超时出库", itemStyle: { color: "#F56C6C" } },
        ];

  rateChart.setOption({
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
}
</script>

<template>
  <div style="padding: 24px; background: #f0f2f5; min-height: 100vh">
    <h2 style="margin-bottom: 24px">ERP-Go 经营看板</h2>

    <!-- 未登录引导 -->
    <el-alert
      v-if="store.authRequired"
      title="需要登录"
      type="warning"
      :closable="false"
      show-icon
      style="margin-bottom: 24px"
    >
      <template #default>
        <p style="margin: 4px 0">
          看板数据需要有效的访问令牌。请先登录管理后台获取令牌，或启用 Demo 模式查看演示数据。
        </p>
        <el-button
          type="primary"
          size="small"
          style="margin-top: 8px"
          @click="store.applyDemo()"
        >
          查看 Demo 数据
        </el-button>
      </template>
    </el-alert>

    <!-- 错误提示 -->
    <el-alert
      v-if="store.error && !store.authRequired"
      :title="store.error"
      type="error"
      :closable="false"
      show-icon
      style="margin-bottom: 16px"
    />

    <!-- KPI Cards（skeleton 加载态） -->
    <el-row :gutter="16">
      <el-col :xs="24" :sm="12" :md="6">
        <el-card shadow="hover">
          <el-skeleton v-if="store.loading" :rows="2" animated />
          <div v-else style="text-align: center">
            <div style="color: #909399; font-size: 14px">订单总数</div>
            <div style="font-size: 28px; font-weight: bold; color: #409EFF">
              {{ store.orderCount.toLocaleString() }}
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :md="6">
        <el-card shadow="hover">
          <el-skeleton v-if="store.loading" :rows="2" animated />
          <div v-else style="text-align: center">
            <div style="color: #909399; font-size: 14px">销售额</div>
            <div style="font-size: 28px; font-weight: bold; color: #67C23A">
              ¥{{ store.salesAmount.toLocaleString() }}
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :md="6">
        <el-card shadow="hover">
          <el-skeleton v-if="store.loading" :rows="2" animated />
          <div v-else style="text-align: center">
            <div style="color: #909399; font-size: 14px">出库量</div>
            <div style="font-size: 28px; font-weight: bold; color: #E6A23C">
              {{ store.outboundCount.toLocaleString() }}
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :md="6">
        <el-card shadow="hover">
          <el-skeleton v-if="store.loading" :rows="2" animated />
          <div v-else style="text-align: center">
            <div style="color: #909399; font-size: 14px">库存 SKU</div>
            <div style="font-size: 28px; font-weight: bold; color: #F56C6C">
              {{ store.skuCount.toLocaleString() }}
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- Charts -->
    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :xs="24" :md="16">
        <el-card>
          <template #header>销售趋势（近 7 天）</template>
          <el-skeleton v-if="store.loading" :rows="8" animated style="height: 320px" />
          <div v-else ref="trendChartRef" style="height: 320px" />
        </el-card>
      </el-col>
      <el-col :xs="24" :md="8">
        <el-card>
          <template #header>出库及时率</template>
          <el-skeleton v-if="store.loading" :rows="8" animated style="height: 320px" />
          <div v-else ref="rateChartRef" style="height: 320px" />
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>