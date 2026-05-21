<script setup lang="ts">
import { ref } from "vue";
const suppliers = ref([{ id:"1", name:"纺织品供应商A", code:"SUP-001", contact:"张三", phone:"13800138001", status:"active" }]);
const orders = ref([{ id:"1", order_no:"PO-20260521-001", supplier_name:"纺织品供应商A", status:"approved", total_amount:5000, currency:"USD", created_at:"2026-05-21" }]);
const statusMap: Record<string,{type:string;label:string}> = { draft:{type:"info",label:"草稿"}, pending:{type:"warning",label:"待审核"}, approved:{type:"",label:"已审核"}, ordered:{type:"success",label:"已下单"}, completed:{type:"success",label:"已完成"} };
</script>
<template>
  <div>
    <el-card style="margin-bottom:16px"><template #header><div style="display:flex;justify-content:space-between"><span>供应商</span><el-button type="primary" size="small">添加供应商</el-button></div></template>
      <el-table :data="suppliers" stripe size="small">
        <el-table-column prop="name" label="名称" width="180" /><el-table-column prop="code" label="编码" width="120" />
        <el-table-column prop="contact" label="联系人" /><el-table-column prop="phone" label="电话" width="150" />
        <el-table-column label="操作" width="160"><template #default><el-button type="primary" size="small">编辑</el-button><el-button type="danger" size="small">禁用</el-button></template></el-table-column>
      </el-table>
    </el-card>
    <el-card><template #header><div style="display:flex;justify-content:space-between"><span>采购单</span><el-button type="primary" size="small">新建采购单</el-button></div></template>
      <el-table :data="orders" stripe>
        <el-table-column prop="order_no" label="采购单号" width="200" /><el-table-column prop="supplier_name" label="供应商" width="150" />
        <el-table-column label="状态" width="120"><template #default="{row}"><el-tag :type="statusMap[row.status]?.type||'info'" size="small">{{statusMap[row.status]?.label||row.status}}</el-tag></template></el-table-column>
        <el-table-column label="金额" width="120"><template #default="{row}">{{row.currency}} {{row.total_amount}}</template></el-table-column>
        <el-table-column prop="created_at" label="创建时间" /><el-table-column label="操作" width="160"><template #default><el-button type="primary" size="small">审核</el-button><el-button type="success" size="small">收货</el-button></template></el-table-column>
      </el-table>
    </el-card>
  </div>
</template>
