<script setup lang="ts">
import { ref, onMounted } from "vue";
import { ElMessage } from "element-plus";
import { apiClient } from "@erp/shared";
import type { ApiResponse, PageData } from "@erp/shared";

const isDemo = import.meta.env.VITE_DEMO_MODE === "true";

interface Notification {
  id: string;
  title: string;
  content: string;
  type: string;
  read: boolean;
  created_at: string;
}

const notifications = ref<Notification[]>([]);
const total = ref(0);
const loading = ref(false);
const unreadCount = ref(0);
const fetchError = ref(false);

async function fetchNotifications(page: number = 1, pageSize: number = 10) {
  if (isDemo) {
    notifications.value = [];
    total.value = 0;
    unreadCount.value = 0;
    return;
  }
  loading.value = true;
  fetchError.value = false;
  try {
    const res = await apiClient.get<ApiResponse<PageData<Notification>>>(
      "/notification/list",
      { params: { page, page_size: pageSize } }
    );
    notifications.value = res.data.data?.list ?? [];
    total.value = res.data.data?.total ?? 0;
    unreadCount.value = notifications.value.filter((n) => !n.read).length;
  } catch {
    notifications.value = [];
    total.value = 0;
    unreadCount.value = 0;
    fetchError.value = true;
    ElMessage.warning("通知服务暂不可用");
  } finally {
    loading.value = false;
  }
}

function markAllRead() {
  notifications.value.forEach((n) => {
    n.read = true;
  });
  unreadCount.value = 0;
}

onMounted(() => {
  fetchNotifications();
});
</script>

<template>
  <div>
    <el-card>
      <template #header>
        <div style="display: flex; justify-content: space-between">
          <span>通知中心 <el-badge :value="unreadCount" style="margin-left: 8px" /></span>
          <el-button size="small" @click="markAllRead">全部已读</el-button>
        </div>
      </template>
      <div v-loading="loading">
        <el-empty v-if="fetchError" description="通知服务暂不可用，请稍后重试" />
        <el-empty v-else-if="!isDemo && notifications.length === 0 && !loading" description="暂无通知" />
        <div v-for="n in notifications" :key="n.id" style="padding: 12px 0; border-bottom: 1px solid #f0f0f0; display: flex; align-items: center; justify-content: space-between">
          <div>
            <el-tag :type="n.type as 'info' | 'warning' | 'success'" size="small" style="margin-right: 8px">{{ n.type === 'info' ? '提醒' : n.type === 'warning' ? '告警' : '完成' }}</el-tag>
            <span :style="{ fontWeight: n.read ? 'normal' : 'bold' }">{{ n.title }}</span>
            <span style="color: #909399; margin-left: 8px; font-size: 13px">{{ n.content }}</span>
          </div>
          <span style="color: #c0c4cc; font-size: 12px">{{ n.created_at }}</span>
        </div>
      </div>
      <el-pagination
        style="margin-top: 16px; justify-content: flex-end"
        background
        layout="total, prev, pager, next"
        :total="total"
        :page-size="10"
        size="small"
        @current-change="(page: number) => fetchNotifications(page)"
      />
    </el-card>
  </div>
</template>
