<script setup lang="ts">
import { ref } from "vue";
import { ProTable } from "@erp/shared";

const suppliers = ref([{ id: "1", name: "纺织品供应商A", code: "SUP-001", contact: "张三", phone: "13800138001", status: "active" }]);
const orders = ref([{ id: "1", order_no: "PO-20260521-001", supplier_name: "纺织品供应商A", status: "approved", total_amount: 5000, currency: "USD", created_at: "2026-05-21" }]);

const supplierColumns = [
  { prop: "name", label: "名称", width: 180 },
  { prop: "code", label: "编码", width: 120 },
  { prop: "contact", label: "联系人" },
  { prop: "phone", label: "电话", width: 150 },
  { prop: "actions", label: "操作", width: 160 },
];

const orderColumns = [
  { prop: "order_no", label: "采购单号", width: 200 },
  { prop: "supplier_name", label: "供应商", width: 150 },
  { prop: "status", label: "状态", width: 120 },
  { prop: "total_amount", label: "金额", width: 120 },
  { prop: "created_at", label: "创建时间" },
  { prop: "actions", label: "操作", width: 160 },
];

const statusMap: Record<string, { type: string; label: string }> = {
  draft: { type: "info", label: "草稿" },
  pending: { type: "warning", label: "待审核" },
  approved: { type: "", label: "已审核" },
  ordered: { type: "success", label: "已下单" },
  completed: { type: "success", label: "已完成" },
};
</script>

<template>
  <div>
    <el-card style="margin-bottom: 16px">
      <template #header>
        <div style="display: flex; justify-content: space-between"><span>供应商</span><el-button type="primary" size="small">添加供应商</el-button></div>
      </template>
      <ProTable
        :columns="supplierColumns"
        :data="suppliers"
        :total="suppliers.length"
        :page-size="suppliers.length"
      >
        <template #actions>
          <el-button type="primary" size="small">编辑</el-button>
          <el-button type="danger" size="small">禁用</el-button>
        </template>
      </ProTable>
    </el-card>

    <el-card>
      <template #header>
        <div style="display: flex; justify-content: space-between"><span>采购单</span><el-button type="primary" size="small">新建采购单</el-button></div>
      </template>
      <ProTable
        :columns="orderColumns"
        :data="orders"
        :total="orders.length"
        :page-size="orders.length"
      >
        <template #status="{ row }">
          <el-tag :type="statusMap[row.status]?.type || 'info'" size="small">{{ statusMap[row.status]?.label || row.status }}</el-tag>
        </template>
        <template #total_amount="{ row }">
          {{ row.currency }} {{ row.total_amount }}
        </template>
        <template #actions>
          <el-button type="primary" size="small">审核</el-button>
          <el-button type="success" size="small">收货</el-button>
        </template>
      </ProTable>
    </el-card>
  </div>
</template>
