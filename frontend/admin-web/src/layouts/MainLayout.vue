<script setup lang="ts">
import { ref, watch } from "vue";
import { useRouter, useRoute } from "vue-router";

const router = useRouter();
const route = useRoute();
const activeMenu = ref(route.path);

watch(
  () => route.path,
  (newPath) => {
    activeMenu.value = newPath;
  }
);

const menuItems = [
  { path: "/dashboard", title: "首页", icon: "DataAnalysis" },
  {
    title: "系统管理",
    icon: "Setting",
    children: [
      { path: "/system/users", title: "用户管理", icon: "User" },
      { path: "/system/roles", title: "角色管理", icon: "Avatar" },
      { path: "/system/org", title: "组织管理", icon: "OfficeBuilding" },
      { path: "/system/audit", title: "操作审计", icon: "DocumentChecked" },
    ],
  },
  {
    title: "商品管理",
    icon: "Goods",
    children: [
      { path: "/product/list", title: "商品列表", icon: "List" },
      { path: "/product/sku-mapping", title: "SKU 映射", icon: "Connection" },
    ],
  },
  {
    title: "渠道管理",
    icon: "Connection",
    children: [
      { path: "/channel/stores", title: "店铺授权", icon: "Shop" },
      { path: "/channel/import", title: "订单导入", icon: "Upload" },
    ],
  },
  {
    title: "订单管理",
    icon: "Document",
    children: [
      { path: "/order/list", title: "订单列表", icon: "Tickets" },
      { path: "/order/audit", title: "订单审核", icon: "Checked" },
    ],
  },
  { path: "/inventory", title: "库存管理", icon: "Box" },
  { path: "/warehouse/outbound", title: "仓储管理", icon: "House" },
  {
    title: "物流管理",
    icon: "Van",
    children: [
      { path: "/transport", title: "物流配置", icon: "Van" },
    ],
  },
  { path: "/notification", title: "通知中心", icon: "Bell" },
  {
    title: "采购财务",
    icon: "Money",
    children: [
      { path: "/purchase", title: "采购管理", icon: "ShoppingCart" },
      { path: "/finance", title: "财务管理", icon: "Money" },
    ],
  },
];

const handleMenuSelect = (index: string) => {
  if (index.startsWith("/")) {
    router.push(index);
  }
};
</script>

<template>
  <el-container style="height: 100vh">
    <el-aside width="220px" style="background-color: #304156; overflow-y: auto">
      <div style="height: 60px; display: flex; align-items: center; justify-content: center; color: #fff; font-size: 18px; font-weight: bold; border-bottom: 1px solid #4a5a6a">
        ERP-Go
      </div>
      <el-menu
        :default-active="activeMenu"
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409EFF"
        @select="handleMenuSelect"
      >
        <template v-for="item in menuItems" :key="item.path || item.title">
          <el-sub-menu v-if="item.children" :index="item.title">
            <template #title>
              <el-icon>
                <component :is="item.icon" />
              </el-icon>
              <span>{{ item.title }}</span>
            </template>
            <el-menu-item v-for="child in item.children" :key="child.path" :index="child.path">
              <el-icon>
                <component :is="child.icon" />
              </el-icon>
              <span>{{ child.title }}</span>
            </el-menu-item>
          </el-sub-menu>
          <el-menu-item v-else :index="item.path">
            <el-icon>
              <component :is="item.icon" />
            </el-icon>
            <span>{{ item.title }}</span>
          </el-menu-item>
        </template>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header style="height: 60px; border-bottom: 1px solid #e6e6e6; display: flex; align-items: center; justify-content: space-between; padding: 0 20px">
        <el-breadcrumb>
          <el-breadcrumb-item :to="{ path: '/dashboard' }">首页</el-breadcrumb-item>
          <el-breadcrumb-item v-if="route.meta.title">{{ route.meta.title }}</el-breadcrumb-item>
        </el-breadcrumb>
        <el-dropdown>
          <span style="cursor: pointer">
            管理员 <el-icon><ArrowDown /></el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item>个人设置</el-dropdown-item>
              <el-dropdown-item divided @click="router.push('/login')">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </el-header>
      <el-main style="background: #f0f2f5">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<style scoped>
.el-menu {
  border-right: none;
}
.el-breadcrumb {
  font-size: 14px;
}
</style>
