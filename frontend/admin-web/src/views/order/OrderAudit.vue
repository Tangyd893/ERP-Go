<script setup lang="ts">
import { computed, onMounted } from "vue";
import { ElMessage } from "element-plus";
import { ProTable } from "@erp/shared";
import { useOrderStore } from "@/stores/order";

const orderStore = useOrderStore();

const displayData = computed(() =>
  orderStore.orders.filter((o) => o.status === "pending")
);

const columns = [
  { prop: "platform_order_no", label: "平台订单号", width: 200 },
  { prop: "store", label: "店铺", width: 150 },
  { prop: "buyer_name", label: "买家", width: 120 },
  { prop: "total_amount", label: "金额", width: 120, align: "right" as const },
  { prop: "created_at", label: "创建时间", width: 170 },
  { prop: "actions", label: "操作", width: 280, fixed: "right" as const },
];

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const handleApprove = (order: any) => {
  const orderId = order.order_id || order.id;
  orderStore.auditOrder(orderId, true).then(() => {
    order.status = "approved";
    ElMessage.success(`订单 ${order.platform_order_no || order.order_no} 审核通过`);
  });
};

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const handleReject = (order: any) => {
  const orderId = order.order_id || order.id;
  orderStore.auditOrder(orderId, false).then(() => {
    order.status = "rejected";
    ElMessage.error(`订单 ${order.platform_order_no || order.order_no} 已驳回`);
  });
};

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const handleDetail = (order: any) => {
  ElMessage.info(`查看订单 ${order.platform_order_no || order.order_no} 详情`);
};

onMounted(() => {
  orderStore.fetchOrders({ status: "pending" });
});
</script>

<template>
  <div>
    <el-card>
      <template #header>
        <span>订单审核</span>
      </template>

      <ProTable
        :columns="columns"
        :data="displayData"
        :loading="orderStore.loading"
        :total="displayData.length"
        @page-change="(page: number) => orderStore.fetchOrders({ status: 'pending', page })"
      >
        <template #total_amount="{ row }">
          {{ row.currency }} {{ row.total_amount.toFixed(2) }}
        </template>
        <template #actions="{ row }">
          <el-button type="success" size="small" @click="handleApprove(row)">审核通过</el-button>
          <el-button type="danger" size="small" @click="handleReject(row)">驳回</el-button>
          <el-button size="small" @click="handleDetail(row)">查看详情</el-button>
        </template>
      </ProTable>
    </el-card>
  </div>
</template>
