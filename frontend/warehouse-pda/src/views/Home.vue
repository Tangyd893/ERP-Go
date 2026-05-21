<script setup lang="ts">
</script>

<template>
  <div style="padding: 16px">
    <h3 style="margin-bottom: 16px">WMS 仓库作业</h3>
    <el-space direction="vertical" style="width: 100%" :size="12">
      <el-card shadow="hover">
        <div style="text-align: center; font-size: 18px">拣货任务</div>
        <div style="text-align: center; color: #409EFF; font-size: 24px; font-weight: bold">3</div>
      </el-card>
      <el-card shadow="hover">
        <div style="text-align: center; font-size: 18px">复核任务</div>
        <div style="text-align: center; color: #E6A23C; font-size: 24px; font-weight: bold">2</div>
      </el-card>
      <el-card shadow="hover">
        <div style="text-align: center; font-size: 18px">打包任务</div>
        <div style="text-align: center; color: #67C23A; font-size: 24px; font-weight: bold">5</div>
      </el-card>
    </el-space>
    <div style="margin-top: 24px">
      <el-button type="primary" size="large" style="width: 100%">扫描拣货</el-button>
    </div>
  </div>
</template>
