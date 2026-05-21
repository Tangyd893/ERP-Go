<script setup lang="ts">
</script>

<template>
  <div>
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>今日订单</template>
          <div style="font-size: 32px; font-weight: bold; color: #409EFF">128</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>待审核</template>
          <div style="font-size: 32px; font-weight: bold; color: #E6A23C">35</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>待出库</template>
          <div style="font-size: 32px; font-weight: bold; color: #67C23A">52</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>库存告警</template>
          <div style="font-size: 32px; font-weight: bold; color: #F56C6C">8</div>
        </el-card>
      </el-col>
    </el-row>
    <el-card style="margin-top: 20px">
      <template #header>欢迎使用 ERP-Go 跨境电商 ERP 系统</template>
      <p>系统当前处于 MVP 阶段，功能正在持续构建中。</p>
      <p>首期功能：商品管理、订单管理、库存管理、仓储出库、物流发货。</p>
    </el-card>
  </div>
</template>
