<script setup lang="ts">
import { ref } from "vue";

const tableData = ref([
  { id: "1", code: "TSHIRT-001", name: "纯棉T恤-白色-M", spu_name: "纯棉T恤", barcode: "6901234567890", weight: 0.2, sale_price: 15.99, currency: "USD", status: "active" },
  { id: "2", code: "TSHIRT-002", name: "纯棉T恤-黑色-L", spu_name: "纯棉T恤", barcode: "6901234567891", weight: 0.22, sale_price: 15.99, currency: "USD", status: "active" },
  { id: "3", code: "MUG-001", name: "陶瓷马克杯-350ml", spu_name: "陶瓷马克杯", barcode: "6901234567892", weight: 0.35, sale_price: 12.99, currency: "USD", status: "active" },
]);
</script>

<template>
  <div>
    <el-card>
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center">
          <span>商品列表 (SKU)</span>
          <div>
            <el-button type="primary">新建 SPU</el-button>
            <el-button type="success">新建 SKU</el-button>
          </div>
        </div>
      </template>

      <el-form inline style="margin-bottom: 16px">
        <el-form-item label="SKU编码">
          <el-input placeholder="搜索编码" clearable style="width: 180px" />
        </el-form-item>
        <el-form-item label="商品名称">
          <el-input placeholder="搜索名称" clearable style="width: 180px" />
        </el-form-item>
        <el-form-item label="条码">
          <el-input placeholder="条码搜索" clearable style="width: 180px" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary">查询</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="tableData" stripe>
        <el-table-column prop="code" label="SKU编码" width="150" />
        <el-table-column prop="name" label="SKU名称" min-width="180" />
        <el-table-column prop="spu_name" label="所属SPU" width="120" />
        <el-table-column prop="barcode" label="条码" width="150" />
        <el-table-column prop="weight" label="重量(kg)" width="100" align="right" />
        <el-table-column prop="sale_price" label="售价" width="120" align="right">
          <template #default="{ row }">
            {{ row.currency }} {{ row.sale_price.toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'danger'" size="small">
              {{ row.status === 'active' ? '启售' : '停售' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default>
            <el-button type="primary" size="small">详情</el-button>
            <el-button type="warning" size="small">编辑</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination style="margin-top: 16px; justify-content: flex-end" background layout="total, prev, pager, next" :total="3" :page-size="20" />
    </el-card>
  </div>
</template>
