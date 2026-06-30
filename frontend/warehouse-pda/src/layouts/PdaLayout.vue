<script setup lang="ts">
import { computed } from "vue";
import { useRouter, useRoute } from "vue-router";
import { useAuthStore } from "@/stores/auth";
import {
  HomeFilled,
  List,
  Checked,
  Box,
  User,
} from "@element-plus/icons-vue";

const router = useRouter();
const route = useRoute();
const auth = useAuthStore();

const tabs = [
  { path: "/", label: "首页", icon: HomeFilled },
  { path: "/pick", label: "拣货", icon: List },
  { path: "/check", label: "复核", icon: Checked },
  { path: "/pack", label: "打包", icon: Box },
  { path: "/profile", label: "我的", icon: User },
];

const activeTab = computed(() => {
  // 子路由如 /pick/scan 也匹配 /pick
  for (const tab of tabs) {
    if (tab.path === "/") {
      if (route.path === "/") return "/";
    } else if (route.path.startsWith(tab.path)) {
      return tab.path;
    }
  }
  return "/";
});

function onTabClick(path: string) {
  if (route.path !== path) {
    router.push(path);
  }
}

function handleLogout() {
  auth.logout();
  router.push("/login");
}
</script>

<template>
  <div class="pda-layout">
    <!-- 顶栏 -->
    <header class="pda-header">
      <span class="pda-header__title">WMS PDA</span>
      <div class="pda-header__right">
        <span class="pda-header__user">{{ auth.username || auth.tenantId || "未登录" }}</span>
        <el-button text size="small" @click="handleLogout">退出</el-button>
      </div>
    </header>

    <!-- 内容区 -->
    <main class="pda-main">
      <router-view />
    </main>

    <!-- 底栏 Tab -->
    <footer class="pda-tabbar">
      <div
        v-for="tab in tabs"
        :key="tab.path"
        class="pda-tabbar__item"
        :class="{ 'is-active': activeTab === tab.path }"
        @click="onTabClick(tab.path)"
      >
        <component :is="tab.icon" class="pda-tabbar__icon" />
        <span class="pda-tabbar__label">{{ tab.label }}</span>
      </div>
    </footer>
  </div>
</template>

<style scoped>
.pda-layout {
  display: flex;
  flex-direction: column;
  height: 100vh;
  height: 100dvh;
  background: var(--pda-bg, #f5f5f5);
}

/* ---- 顶栏 ---- */
.pda-header {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 48px;
  padding: 0 12px;
  padding-top: var(--pda-safe-top, 0px);
  background: #fff;
  border-bottom: 1px solid #ebeef5;
  z-index: 100;
}

.pda-header__title {
  font-size: 16px;
  font-weight: 600;
  color: var(--pda-text, #303133);
}

.pda-header__right {
  display: flex;
  align-items: center;
  gap: 4px;
}

.pda-header__user {
  font-size: 12px;
  color: var(--pda-text-secondary, #909399);
  max-width: 120px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* ---- 内容 ---- */
.pda-main {
  flex: 1;
  overflow-y: auto;
  -webkit-overflow-scrolling: touch;
  padding: 12px;
}

/* ---- 底栏 ---- */
.pda-tabbar {
  flex-shrink: 0;
  display: flex;
  height: 56px;
  padding-bottom: var(--pda-safe-bottom, 0px);
  background: #fff;
  border-top: 1px solid #ebeef5;
  z-index: 100;
}

.pda-tabbar__item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 2px;
  color: var(--pda-text-secondary, #909399);
  cursor: pointer;
  transition: color 0.2s;
  min-height: var(--pda-touch-min, 44px);
}

.pda-tabbar__item.is-active {
  color: var(--pda-primary, #409EFF);
}

.pda-tabbar__icon {
  font-size: 20px;
}

.pda-tabbar__label {
  font-size: 10px;
  line-height: 1;
}
</style>
