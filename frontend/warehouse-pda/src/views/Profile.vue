<script setup lang="ts">
import { useRouter } from "vue-router";
import { useAuthStore } from "@/stores/auth";

const router = useRouter();
const auth = useAuthStore();

function handleLogout() {
  auth.logout();
  router.push("/login");
}
</script>

<template>
  <div class="profile-page">
    <el-card>
      <template #header>我的</template>
      <el-descriptions :column="1" border>
        <el-descriptions-item label="用户名">{{ auth.username || "-" }}</el-descriptions-item>
        <el-descriptions-item label="租户">{{ auth.tenantId || "-" }}</el-descriptions-item>
      </el-descriptions>
      <div class="profile-logout">
        <el-button type="danger" @click="handleLogout">退出登录</el-button>
      </div>
    </el-card>
  </div>
</template>

<style scoped>
.profile-page {
  padding: 16px;
}

.profile-logout {
  margin-top: 24px;
  text-align: center;
}
</style>
