<script setup lang="ts">
import { computed, onMounted } from "vue";
import { ElMessage } from "element-plus";
import { ProTable } from "@erp/shared";
import { useWarehouseStore } from "@/stores/warehouse";

const warehouseStore = useWarehouseStore();

onMounted(() => {
  warehouseStore.fetchOutbounds();
});

const stats = computed(() => {
  const items = warehouseStore.outbounds;
  return {
    picking: items.filter((o) => o.status === "picking" || o.status === "created").length,
    checking: items.filter((o) => o.status === "picked" || o.status === "checking" || o.status === "checked").length,
    packing: items.filter((o) => o.status === "packed" || o.status === "weighed").length,
    total: items.length,
  };
});

const columns = [
  { prop: "order_no", label: "关联订单", width: 200 },
  { prop: "id", label: "出库单号", width: 180 },
  { prop: "status", label: "出库状态", width: 110 },
  { prop: "items", label: "商品数", width: 80, align: "center" as const },
  { prop: "created_at", label: "创建时间", width: 180 },
  { prop: "actions", label: "操作", width: 200, fixed: "right" as const },
];

const statusLabels: Record<string, { type: string; label: string }> = {
  created: { type: "info", label: "待波次" },
  picking: { type: "", label: "拣货中" },
  picked: { type: "warning", label: "已拣货" },
  checking: { type: "", label: "复核中" },
  checked: { type: "warning", label: "已复核" },
  packed: { type: "", label: "已打包" },
  weighed: { type: "", label: "已称重" },
  shipped: { type: "success", label: "已出库" },
  abnormal: { type: "danger", label: "异常" },
};

function handleCreateOutbound() {
  ElMessage.info("创建出库单功能：请从已审核订单触发出库");
}

function handleProcess(row: Record<string, unknown>) {
  const status = row.status as string;
  if (status === "created") {
    ElMessage.info("出库单已创建，等待波次分配");
  } else if (status === "picking" || status === "picked") {
    ElMessage.info("拣货进行中，查看 PDA 端进度");
  } else if (status === "checked" || status === "packed") {
    ElMessage.info("请继续打包/称重操作");
  } else {
    ElMessage.info(`当前状态: ${statusLabels[status]?.label || status}`);
  }
}

function handleAbnormal(row: Record<string, unknown>) {
  ElMessage.warning(`出库单 ${row.id} 标记异常 — 功能待接入`);
}
</script>

<template>
  <div>
    <!-- Stats Cards -->
    <el-row :gutter="16" style="margin-bottom: 16px">
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>全部出库单</template>
          <div style="text-align: center; font-size: 32px; color: #303133">
            {{ stats.total }}
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>拣货任务</template>
          <div style="text-align: center; font-size: 32px; color: #409eff">
            {{ stats.picking }}
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>复核任务</template>
          <div style="text-align: center; font-size: 32px; color: #e6a23c">
            {{ stats.checking }}
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>待打包</template>
          <div style="text-align: center; font-size: 32px; color: #67c23a">
            {{ stats.packing }}
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- Table -->
    <el-card>
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center">
          <span>出库单列表</span>
          <el-button type="primary" size="small" @click="handleCreateOutbound">
            创建出库单
          </el-button>
        </div>
      </template>

      <ProTable
        :columns="columns"
        :data="warehouseStore.outbounds"
        :total="warehouseStore.total"
        :loading="warehouseStore.loading"
        @page-change="(page: number) => warehouseStore.fetchOutbounds({ page })"
      >
        <template #status="{ row }">
          <el-tag
            :type="statusLabels[row.status]?.type || 'info'"
            size="small"
            disable-transitions
          >
            {{ statusLabels[row.status]?.label || row.status }}
          </el-tag>
        </template>
        <template #items="{ row }">
          {{ Array.isArray(row.items) ? row.items.length : (row.items || "-") }}
        </template>
        <template #actions="{ row }">
          <el-button
            type="primary"
            size="small"
            :disabled="row.status === 'shipped'"
            @click="handleProcess(row)"
          >
            处理
          </el-button>
          <el-button
            type="danger"
            size="small"
            :disabled="row.status === 'shipped' || row.status === 'abnormal'"
            @click="handleAbnormal(row)"
          >
            异常
          </el-button>
        </template>
      </ProTable>
    </el-card>
  </div>
</template>
