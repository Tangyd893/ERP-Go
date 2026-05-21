<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { ElMessage } from "element-plus";
import { ProTable } from "@erp/shared";
import { useOrderStore } from "@/stores/order";

const mockOrders = ref([
  { id: "1", platform_order_no: "AMZ-20260521-100", store: "美国站店铺A", buyer_name: "张三", currency: "USD", total_amount: 99.99, status: "pending", created_at: "2026-05-21 10:30" },
  { id: "2", platform_order_no: "AMZ-20260521-101", store: "英国站店铺B", buyer_name: "李四", currency: "GBP", total_amount: 45.50, status: "pending", created_at: "2026-05-21 11:00" },
  { id: "3", platform_order_no: "AMZ-20260521-102", store: "美国站店铺A", buyer_name: "王五", currency: "USD", total_amount: 23.80, status: "pending", created_at: "2026-05-21 11:30" },
]);

const orderStore = useOrderStore();

const displayData = computed(() => {
  if (orderStore.orders.length > 0) {
    return orderStore.orders.filter((o) => o.status === "pending");
  }
  return mockOrders.value;
});

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
