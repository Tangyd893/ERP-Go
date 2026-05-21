<script setup lang="ts">
import { ref } from "vue";
const settlements = ref([{ id:"1", store:"美国站店铺A", period:"2026-05-01 ~ 2026-05-15", total_sales:12500, commission:1500, net_amount:11000, status:"matched", currency:"USD" }]);
const profits = ref([{ order_no:"AMZ-20260520-004", sale_amount:12.99, purchase_cost:5.0, shipping_cost:3.5, commission:1.95, total_cost:10.45, gross_profit:2.54, profit_margin:19.6 }]);
</script>
<template>
  <div>
    <el-card style="margin-bottom:16px"><template #header><div style="display:flex;justify-content:space-between"><span>平台结算</span><el-button type="primary" size="small">导入结算报告</el-button></div></template>
      <el-table :data="settlements" stripe size="small">
        <el-table-column prop="store" label="店铺" width="150" /><el-table-column prop="period" label="结算周期" width="200" />
        <el-table-column label="销售额" width="130"><template #default="{row}">{{row.currency}} {{row.total_sales}}</template></el-table-column>
        <el-table-column label="佣金" width="120"><template #default="{row}">{{row.currency}} {{row.commission}}</template></el-table-column>
        <el-table-column label="净收入" width="120"><template #default="{row}"><span style="color:#67C23A;font-weight:bold">{{row.currency}} {{row.net_amount}}</span></template></el-table-column>
        <el-table-column label="状态" width="100"><template #default><el-tag type="success" size="small">已匹配</el-tag></template></el-table-column>
      </el-table>
    </el-card>
    <el-card><template #header><span>订单利润分析</span></template>
      <el-table :data="profits" stripe size="small">
        <el-table-column prop="order_no" label="订单号" width="200" /><el-table-column prop="sale_amount" label="销售额" width="100" />
        <el-table-column prop="purchase_cost" label="采购成本" width="100" /><el-table-column prop="shipping_cost" label="物流费" width="100" />
        <el-table-column prop="commission" label="佣金" width="80" /><el-table-column prop="total_cost" label="总成本" width="100" />
        <el-table-column prop="gross_profit" label="毛利" width="80"><template #default="{row}"><span :style="{color:row.gross_profit>0?'#67C23A':'#F56C6C'}">{{row.gross_profit.toFixed(2)}}</span></template></el-table-column>
        <el-table-column prop="profit_margin" label="利润率" width="80"><template #default="{row}">{{row.profit_margin.toFixed(1)}}%</template></el-table-column>
      </el-table>
    </el-card>
  </div>
</template>
