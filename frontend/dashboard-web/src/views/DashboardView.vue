<script setup lang="ts">
import { onMounted, ref, watch, nextTick } from "vue";
import { useDashboardStore } from "@/stores/dashboard";
import * as echarts from "echarts";

const store = useDashboardStore();
const trendChartRef = ref<HTMLDivElement | null>(null);
const rateChartRef = ref<HTMLDivElement | null>(null);

onMounted(async () => {
  await store.fetchAll();
  await nextTick();
  initCharts();
});

watch(() => store.loading, async (v) => {
  if (!v) {
    await nextTick();
    initCharts();
  }
});

function initCharts() {
  if (trendChartRef.value) {
    const chart = echarts.init(trendChartRef.value);
    chart.setOption({
      tooltip: { trigger: "axis" },
      legend: { data: ["订单量", "销售额(元)"] },
      xAxis: { type: "category", data: ["周一", "周二", "周三", "周四", "周五", "周六", "周日"] },
      yAxis: [
        { type: "value", name: "单量" },
        { type: "value", name: "金额(元)" },
      ],
      series: [
        {
          name: "订单量",
          type: "bar",
          data: [45, 52, 38, 65, 48, 32, 28],
          itemStyle: { color: "#409EFF" },
        },
        {
          name: "销售额(元)",
          type: "line",
          yAxisIndex: 1,
          data: [2100, 3200, 1850, 4200, 2800, 1500, 1200],
          itemStyle: { color: "#67C23A" },
        },
      ],
    });
    window.addEventListener("resize", () => chart.resize());
  }

  if (rateChartRef.value) {
    const chart = echarts.init(rateChartRef.value);
    chart.setOption({
      tooltip: { trigger: "item" },
      legend: { orient: "vertical", left: "left" },
      series: [
        {
          name: "出库时效",
          type: "pie",
          radius: ["50%", "75%"],
          label: { show: false },
          data: [
            { value: 152, name: "24h内出库", itemStyle: { color: "#67C23A" } },
            { value: 28, name: "48h内出库", itemStyle: { color: "#E6A23C" } },
            { value: 6, name: "超时出库", itemStyle: { color: "#F56C6C" } },
          ],
        },
      ],
    });
    window.addEventListener("resize", () => chart.resize());
  }
}
</script>

<template>
  <div style="padding: 24px; background: #f0f2f5; min-height: 100vh">
    <h2 style="margin-bottom: 24px">ERP-Go 经营看板</h2>

    <div v-loading="store.loading">
      <!-- KPI Cards -->
      <el-row :gutter="16">
        <el-col :xs="24" :sm="12" :md="6">
          <el-card shadow="hover">
            <div style="text-align: center">
              <div style="color: #909399; font-size: 14px">订单总数</div>
              <div style="font-size: 28px; font-weight: bold; color: #409EFF">
                {{ store.orderCount.toLocaleString() }}
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :xs="24" :sm="12" :md="6">
          <el-card shadow="hover">
            <div style="text-align: center">
              <div style="color: #909399; font-size: 14px">销售额</div>
              <div style="font-size: 28px; font-weight: bold; color: #67C23A">
                ¥{{ store.salesAmount.toLocaleString() }}
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :xs="24" :sm="12" :md="6">
          <el-card shadow="hover">
            <div style="text-align: center">
              <div style="color: #909399; font-size: 14px">出库量</div>
              <div style="font-size: 28px; font-weight: bold; color: #E6A23C">
                {{ store.outboundCount.toLocaleString() }}
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :xs="24" :sm="12" :md="6">
          <el-card shadow="hover">
            <div style="text-align: center">
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
            <div ref="trendChartRef" style="height: 320px" />
          </el-card>
        </el-col>
        <el-col :xs="24" :md="8">
          <el-card>
            <template #header>出库及时率</template>
            <div ref="rateChartRef" style="height: 320px" />
          </el-card>
        </el-col>
      </el-row>

      <!-- Error state -->
      <el-empty v-if="store.error && !store.loading" :description="store.error" />
    </div>
  </div>
</template>
