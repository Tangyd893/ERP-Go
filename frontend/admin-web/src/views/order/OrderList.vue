<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { ElMessage, ElMessageBox } from "element-plus";
import { ProTable } from "@erp/shared";
import { useOrderStore } from "@/stores/order";

const router = useRouter();
const orderStore = useOrderStore();

const searchForm = ref({
  keyword: "",
  status: "",
  store: "",
});

const columns = [
  { prop: "platform_order_no", label: "平台订单号", width: 200 },
  { prop: "store", label: "店铺", width: 150 },
  { prop: "buyer_name", label: "买家", width: 130 },
  { prop: "total_amount", label: "金额", width: 100, align: "right" as const },
  { prop: "items", label: "件数", width: 70, align: "center" as const },
  { prop: "outbound_no", label: "出库单号", width: 150 },
  { prop: "status", label: "状态", width: 100 },
  { prop: "shipped_at", label: "发货时间", width: 170 },
  { prop: "created_at", label: "创建时间", width: 170 },
  { prop: "actions", label: "操作", width: 260, fixed: "right" as const },
];

const statusLabels: Record<string, { type: string; label: string }> = {
  pending: { type: "info", label: "待审核" },
  approved: { type: "warning", label: "已审核" },
  locked: { type: "", label: "已锁定" },
  picking: { type: "", label: "拣货中" },
  packed: { type: "", label: "已打包" },
  shipped: { type: "success", label: "已发货" },
  abnormal: { type: "danger", label: "异常" },
  cancelled: { type: "info", label: "已取消" },
};

// Fetch on mount
onMounted(() => {
  orderStore.fetchOrders();
});

// Search
function handleSearch() {
  orderStore.fetchOrders({
    page: 1,
    keyword: searchForm.value.keyword || undefined,
    status: searchForm.value.status || undefined,
  });
}

function handleReset() {
  searchForm.value = { keyword: "", status: "", store: "" };
  orderStore.fetchOrders({ page: 1 });
}

// Audit
async function handleAudit(row: Record<string, unknown>) {
  try {
    await orderStore.auditOrder(row.id as string, true);
    ElMessage.success("审核通过");
    orderStore.fetchOrders();
  } catch {
    ElMessage.error("审核失败");
  }
}

// Mark abnormal
async function handleAbnormal(row: Record<string, unknown>) {
  try {
    const { value: reason } = await ElMessageBox.prompt("请输入异常原因", "标记异常", {
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      inputType: "textarea",
    });
    if (reason) {
      await orderStore.markAbnormal(row.id as string, reason);
      ElMessage.success("已标记为异常");
      orderStore.fetchOrders();
    }
  } catch {
    // user cancelled
  }
}

// Cancel
async function handleCancel(row: Record<string, unknown>) {
  try {
    await ElMessageBox.confirm("确认取消该订单？此操作不可撤销", "取消订单", {
      confirmButtonText: "确认取消",
      cancelButtonText: "返回",
      type: "warning",
    });
    await orderStore.cancelOrder(row.id as string, "手动取消");
    ElMessage.success("订单已取消");
    orderStore.fetchOrders();
  } catch {
    // user cancelled
  }
}

function handleImport() {
  router.push("/channel/import");
}
</script>

<template>
  <div>
    <el-card>
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center">
          <span>订单列表</span>
          <el-button type="primary" @click="handleImport">导入订单</el-button>
        </div>
      </template>

      <el-form inline @submit.prevent="handleSearch" style="margin-bottom: 16px">
        <el-form-item label="搜索">
          <el-input
            v-model="searchForm.keyword"
            placeholder="订单号 / 买家"
            clearable
            style="width: 220px"
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable style="width: 140px">
            <el-option
              v-for="(v, k) in statusLabels"
              :key="k"
              :label="v.label"
              :value="k"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">查询</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <ProTable
        :columns="columns"
        :data="orderStore.orders"
        :loading="orderStore.loading"
        :total="orderStore.total"
        @page-change="(page: number) => orderStore.fetchOrders({ page })"
      >
        <template #total_amount="{ row }">
          {{ (row as any).currency || "¥" }} {{ Number(row.total_amount || 0).toFixed(2) }}
        </template>
        <template #items="{ row }">
          {{ Array.isArray(row.items) ? row.items.length : (row.items || 0) }}
        </template>
        <template #status="{ row }">
          <el-tag
            :type="statusLabels[row.status]?.type || 'info'"
            size="small"
          >
            {{ statusLabels[row.status]?.label || row.status }}
          </el-tag>
        </template>
        <template #actions="{ row }">
          <el-button
            type="primary"
            size="small"
            :disabled="row.status !== 'pending'"
            @click="handleAudit(row)"
          >
            审核
          </el-button>
          <el-button
            type="warning"
            size="small"
            :disabled="row.status === 'shipped' || row.status === 'cancelled' || row.status === 'abnormal'"
            @click="handleAbnormal(row)"
          >
            异常
          </el-button>
          <el-button
            type="danger"
            size="small"
            :disabled="row.status === 'cancelled' || row.status === 'shipped'"
            @click="handleCancel(row)"
          >
            取消
          </el-button>
        </template>
      </ProTable>
    </el-card>
  </div>
</template>
