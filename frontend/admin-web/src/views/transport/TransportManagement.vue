<script setup lang="ts">
import { ref } from "vue";

const carriers = ref([
  { id: "1", name: "USPS", code: "usps", service: "Priority Mail", status: "active" },
  { id: "2", name: "FedEx", code: "fedex", service: "Ground", status: "active" },
  { id: "3", name: "UPS", code: "ups", service: "Standard", status: "active" },
]);

const shipments = ref([
  { id: "1", tracking_no: "1Z999AA10123456784", carrier: "USPS", status: "shipped", order_no: "AMZ-20260520-004", weight: 0.5, created_at: "2026-05-20 17:00" },
]);

const statusLabels: Record<string, { type: string; label: string }> = {
  pending: { type: "info", label: "待获取面单" },
  labeled: { type: "warning", label: "已获取面单" },
  shipped: { type: "success", label: "已发货" },
  in_transit: { type: "", label: "运输中" },
  delivered: { type: "success", label: "已签收" },
  failed: { type: "danger", label: "发货失败" },
};
</script>

<template>
  <div>
    <el-card style="margin-bottom: 16px">
      <template #header>
        <div style="display: flex; justify-content: space-between">
          <span>物流配置</span>
          <el-button type="primary" size="small">添加物流商</el-button>
        </div>
      </template>
      <el-table :data="carriers" stripe size="small">
        <el-table-column prop="name" label="物流商" width="120" />
        <el-table-column prop="code" label="编码" width="100" />
        <el-table-column prop="service" label="物流产品" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'danger'" size="small">
              {{ row.status === 'active' ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-card>
      <template #header><span>发运记录</span></template>
      <el-table :data="shipments" stripe>
        <el-table-column prop="tracking_no" label="运单号" width="200" />
        <el-table-column prop="carrier" label="物流商" width="100" />
        <el-table-column prop="order_no" label="订单号" width="200" />
        <el-table-column prop="weight" label="重量(kg)" width="100" />
        <el-table-column prop="status" label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="statusLabels[row.status]?.type || 'info'" size="small">
              {{ statusLabels[row.status]?.label || row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
      </el-table>
    </el-card>
  </div>
</template>
