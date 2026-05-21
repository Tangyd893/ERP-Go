<script setup lang="ts">
import { ref } from "vue";

const activeMenu = ref("dashboard");

const handleMenuSelect = (index: string) => {
  activeMenu.value = index;
};
</script>

<template>
  <el-container style="height: 100vh">
    <el-aside width="220px" style="background-color: #304156">
      <div style="height: 60px; display: flex; align-items: center; justify-content: center; color: #fff; font-size: 18px; font-weight: bold">
        ERP-Go
      </div>
      <el-menu
        :default-active="activeMenu"
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409EFF"
        @select="handleMenuSelect"
      >
        <el-menu-item index="dashboard">
          <el-icon><el-icon-data-analysis /></el-icon>
          <span>首页</span>
        </el-menu-item>
        <el-sub-menu index="product">
          <template #title>
            <el-icon><el-icon-goods /></el-icon>
            <span>商品管理</span>
          </template>
          <el-menu-item index="product-list">商品列表</el-menu-item>
          <el-menu-item index="sku-mapping">SKU 映射</el-menu-item>
        </el-sub-menu>
        <el-sub-menu index="order">
          <template #title>
            <el-icon><el-icon-document /></el-icon>
            <span>订单管理</span>
          </template>
          <el-menu-item index="order-list">订单列表</el-menu-item>
          <el-menu-item index="order-audit">订单审核</el-menu-item>
        </el-sub-menu>
        <el-menu-item index="inventory">
          <el-icon><el-icon-box /></el-icon>
          <span>库存管理</span>
        </el-menu-item>
        <el-menu-item index="warehouse">
          <el-icon><el-icon-house /></el-icon>
          <span>仓储管理</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header style="height: 60px; border-bottom: 1px solid #e6e6e6; display: flex; align-items: center; justify-content: flex-end; padding: 0 20px">
        <el-dropdown>
          <span style="cursor: pointer">
            管理员 <el-icon><el-icon-arrow-down /></el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item>个人设置</el-dropdown-item>
              <el-dropdown-item divided>退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </el-header>
      <el-main>
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<style scoped>
.el-menu {
  border-right: none;
}
</style>
