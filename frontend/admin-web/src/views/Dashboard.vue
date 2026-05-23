<script setup lang="ts">
import { onMounted, computed } from "vue";
import { useOrderStore } from "@/stores/order";
import { useInventoryStore } from "@/stores/inventory";
import { useWarehouseStore } from "@/stores/warehouse";

const orderStore = useOrderStore();
const inventoryStore = useInventoryStore();
const warehouseStore = useWarehouseStore();

onMounted(() => {
  orderStore.fetchOrders();
  inventoryStore.fetchBalances();
  warehouseStore.fetchOutbounds();
});

const todayOrders = computed(() => orderStore.total);
const pendingAudit = computed(() =>
  (Array.isArray(orderStore.orders) ? orderStore.orders : [])
    .filter((o: any) => o.status === "pending").length
);
const pendingOutbound = computed(() => warehouseStore.total);
const alertCount = computed(() =>
  (Array.isArray(inventoryStore.balances) ? inventoryStore.balances : [])
    .filter((b: any) => (b.quantity ?? 0) < (b.alert_threshold ?? 10)).length
);
</script>

<template>
  <div>
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>今日订单</template>
          <div style="font-size: 32px; font-weight: bold; color: #409EFF">{{ todayOrders }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>待审核</template>
          <div style="font-size: 32px; font-weight: bold; color: #E6A23C">{{ pendingAudit }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>待出库</template>
          <div style="font-size: 32px; font-weight: bold; color: #67C23A">{{ pendingOutbound }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>库存告警</template>
          <div style="font-size: 32px; font-weight: bold; color: #F56C6C">{{ alertCount }}</div>
        </el-card>
      </el-col>
    </el-row>
    <el-card style="margin-top: 20px">
      <template #header>欢迎使用 ERP-Go 跨境电商 ERP 系统</template>
      <p>系统当前处于 MVP 阶段，功能正在持续构建中。</p>
      <p>首期功能：商品管理、订单管理、库存管理、仓储出库、物流发货。</p>
    </el-card>
  </div>
</template>
