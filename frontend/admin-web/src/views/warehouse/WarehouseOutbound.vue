<script setup lang="ts">
import { ref } from "vue";
import { ProTable } from "@erp/shared";

const outbounds = ref([
  { id: "1", order_no: "AMZ-20260521-001", status: "packed", items: 2, created_at: "2026-05-21 10:02" },
  { id: "2", order_no: "AMZ-20260520-004", status: "shipped", items: 1, created_at: "2026-05-20 16:30" },
]);

const columns = [
  { prop: "order_no", label: "关联订单", width: 200 },
  { prop: "status", label: "出库状态", width: 130 },
  { prop: "items", label: "商品数", width: 80 },
  { prop: "created_at", label: "创建时间", width: 180 },
  { prop: "actions", label: "操作", width: 200 },
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
</script>

<template>
  <div>
    <el-card style="margin-bottom: 16px">
      <template #header>
        <div style="display: flex; justify-content: space-between">
          <span>出库单列表</span>
          <el-button type="primary" size="small">创建波次</el-button>
        </div>
      </template>
      <ProTable
        :columns="columns"
        :data="outbounds"
        :total="outbounds.length"
      >
        <template #status="{ row }">
          <el-tag :type="statusLabels[row.status]?.type || 'info'" size="small">
            {{ statusLabels[row.status]?.label || row.status }}
          </el-tag>
        </template>
        <template #actions="{ row }">
          <el-button type="primary" size="small" :disabled="row.status === 'shipped'">处理</el-button>
          <el-button type="danger" size="small" :disabled="row.status === 'shipped'">异常</el-button>
        </template>
      </ProTable>
    </el-card>
    <el-row :gutter="16">
      <el-col :span="8">
        <el-card><template #header>拣货任务</template><div style="text-align: center; padding: 20px; font-size: 32px; color: #409EFF">3</div></el-card>
      </el-col>
      <el-col :span="8">
        <el-card><template #header>复核任务</template><div style="text-align: center; padding: 20px; font-size: 32px; color: #E6A23C">2</div></el-card>
      </el-col>
      <el-col :span="8">
        <el-card><template #header>待打包</template><div style="text-align: center; padding: 20px; font-size: 32px; color: #67C23A">5</div></el-card>
      </el-col>
    </el-row>
  </div>
</template>
