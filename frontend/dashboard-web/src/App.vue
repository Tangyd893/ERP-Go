<script setup lang="ts">
</script>

<template>
  <div style="padding: 24px; background: #f0f2f5; min-height: 100vh">
    <h2 style="text-align: center; margin-bottom: 24px">ERP-Go 经营看板</h2>
    <el-row :gutter="16">
      <el-col :span="6">
        <el-card shadow="hover">
          <div style="text-align: center">
            <div style="color: #909399; font-size: 14px">今日销售额</div>
            <div style="font-size: 28px; font-weight: bold; color: #409EFF">$12,850</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <div style="text-align: center">
            <div style="color: #909399; font-size: 14px">订单数</div>
            <div style="font-size: 28px; font-weight: bold; color: #67C23A">328</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <div style="text-align: center">
            <div style="color: #909399; font-size: 14px">出库量</div>
            <div style="font-size: 28px; font-weight: bold; color: #E6A23C">186</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <div style="text-align: center">
            <div style="color: #909399; font-size: 14px">库存 SKU</div>
            <div style="font-size: 28px; font-weight: bold; color: #F56C6C">1,256</div>
          </div>
        </el-card>
      </el-col>
    </el-row>
    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :span="16">
        <el-card>
          <template #header>销售趋势</template>
          <div style="height: 300px; display: flex; align-items: center; justify-content: center; color: #909399">
            图表区域 - 待接入 ECharts
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card>
          <template #header>出库及时率</template>
          <div style="height: 300px; display: flex; align-items: center; justify-content: center; color: #909399">
            图表区域 - 待接入 ECharts
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<style scoped>
.el-row {
  margin-bottom: 16px;
}
</style>
