<script setup lang="ts">
import { ref, onMounted } from "vue";
import { ProTable, apiClient } from "@erp/shared";

interface Settlement {
  id: string; store: string; period: string; total_sales: number;
  commission: number; net_amount: number; status: string; currency: string;
}
interface Profit {
  order_no: string; sale_amount: number; purchase_cost: number;
  shipping_cost: number; commission: number; total_cost: number;
  gross_profit: number; profit_margin: number;
}

const settlements = ref<Settlement[]>([]);
const profits = ref<Profit[]>([]);
const loading = ref(false);
const error = ref("");

const settlementColumns = [
  { prop: "store", label: "店铺", width: 150 },
  { prop: "period", label: "结算周期", width: 200 },
  { prop: "total_sales", label: "销售额", width: 130 },
  { prop: "commission", label: "佣金", width: 120 },
  { prop: "net_amount", label: "净收入", width: 120 },
  { prop: "status", label: "状态", width: 100 },
];

const profitColumns = [
  { prop: "order_no", label: "订单号", width: 200 },
  { prop: "sale_amount", label: "销售额", width: 100 },
  { prop: "purchase_cost", label: "采购成本", width: 100 },
  { prop: "shipping_cost", label: "物流费", width: 100 },
  { prop: "commission", label: "佣金", width: 80 },
  { prop: "total_cost", label: "总成本", width: 100 },
  { prop: "gross_profit", label: "毛利", width: 80 },
  { prop: "profit_margin", label: "利润率", width: 80 },
];

onMounted(async () => {
  loading.value = true;
  try {
    const [settleRes, profitRes] = await Promise.allSettled([
      apiClient.get("/finance/settlements", { params: { page: 1, page_size: 100 } }),
      apiClient.get("/finance/profit", { params: { page: 1, page_size: 100 } }),
    ]);
    if (settleRes.status === "fulfilled") {
      settlements.value = settleRes.value.data?.data?.list ?? [];
    }
    if (profitRes.status === "fulfilled") {
      profits.value = profitRes.value.data?.data?.list ?? [];
    }
  } catch {
    error.value = "加载财务数据失败";
  } finally {
    loading.value = false;
  }
});
</script>

<template>
  <div>
    <el-card style="margin-bottom: 16px">
      <template #header>
        <div style="display: flex; justify-content: space-between">
          <span>平台结算</span>
          <el-button type="primary" size="small">导入结算报告</el-button>
        </div>
      </template>
      <ProTable
        :columns="settlementColumns"
        :data="settlements"
        :loading="loading"
        :error="error"
        :total="settlements.length"
        :page-size="settlements.length"
      >
        <template #total_sales="{ row }">{{ row.currency }} {{ row.total_sales }}</template>
        <template #commission="{ row }">{{ row.currency }} {{ row.commission }}</template>
        <template #net_amount="{ row }">
          <span style="color: #67C23A; font-weight: bold">{{ row.currency }} {{ row.net_amount }}</span>
        </template>
        <template #status>
          <el-tag type="success" size="small">已匹配</el-tag>
        </template>
      </ProTable>
    </el-card>

    <el-card>
      <template #header><span>订单利润分析</span></template>
      <ProTable
        :columns="profitColumns"
        :data="profits"
        :loading="loading"
        :total="profits.length"
        :page-size="profits.length"
      >
        <template #gross_profit="{ row }">
          <span :style="{ color: row.gross_profit > 0 ? '#67C23A' : '#F56C6C' }">{{ row.gross_profit.toFixed(2) }}</span>
        </template>
        <template #profit_margin="{ row }">{{ row.profit_margin.toFixed(1) }}%</template>
      </ProTable>
    </el-card>
  </div>
</template>
