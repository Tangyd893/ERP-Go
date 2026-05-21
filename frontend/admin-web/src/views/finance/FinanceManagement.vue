<script setup lang="ts">
import { ref } from "vue";
import { ProTable } from "@erp/shared";

const settlements = ref([{ id: "1", store: "美国站店铺A", period: "2026-05-01 ~ 2026-05-15", total_sales: 12500, commission: 1500, net_amount: 11000, status: "matched", currency: "USD" }]);
const profits = ref([{ order_no: "AMZ-20260520-004", sale_amount: 12.99, purchase_cost: 5.0, shipping_cost: 3.5, commission: 1.95, total_cost: 10.45, gross_profit: 2.54, profit_margin: 19.6 }]);

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
