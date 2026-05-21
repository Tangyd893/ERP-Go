<script setup lang="ts">
import { ref } from "vue";

const notifications = ref([
  { id: "1", title: "订单审核提醒", content: "有3条订单待审核", type: "info", read: false, created_at: "2026-05-21 10:00" },
  { id: "2", title: "库存告警", content: "SKU TSHIRT-002 库存低于安全值", type: "warning", read: false, created_at: "2026-05-21 09:30" },
  { id: "3", title: "出库完成", content: "订单 AMZ-20260520-004 已出库", type: "success", read: true, created_at: "2026-05-20 17:00" },
]);
</script>

<template>
  <div>
    <el-card>
      <template #header>
        <div style="display: flex; justify-content: space-between">
          <span>通知中心 <el-badge :value="2" style="margin-left: 8px" /></span>
          <el-button size="small">全部已读</el-button>
        </div>
      </template>
      <div v-for="n in notifications" :key="n.id" style="padding: 12px 0; border-bottom: 1px solid #f0f0f0; display: flex; align-items: center; justify-content: space-between">
        <div>
          <el-tag :type="n.type" size="small" style="margin-right: 8px">{{ n.type === 'info' ? '提醒' : n.type === 'warning' ? '告警' : '完成' }}</el-tag>
          <span :style="{ fontWeight: n.read ? 'normal' : 'bold' }">{{ n.title }}</span>
          <span style="color: #909399; margin-left: 8px; font-size: 13px">{{ n.content }}</span>
        </div>
        <span style="color: #c0c4cc; font-size: 12px">{{ n.created_at }}</span>
      </div>
      <el-pagination style="margin-top: 16px; justify-content: flex-end" background layout="total, prev, pager, next" :total="3" :page-size="10" size="small" />
    </el-card>
  </div>
</template>
