<script setup lang="ts">
import { ref, onMounted } from "vue";
import { ProTable, apiClient } from "@erp/shared";

interface Supplier {
  id: string; name: string; code: string; contact: string; phone: string; status: string;
}
interface PurchaseOrder {
  id: string; order_no: string; supplier_name: string; status: string;
  total_amount: number; currency: string; created_at: string;
}

const suppliers = ref<Supplier[]>([]);
const orders = ref<PurchaseOrder[]>([]);
const loading = ref(false);
const error = ref("");

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

onMounted(async () => {
  loading.value = true;
  try {
    const [supRes, ordRes] = await Promise.allSettled([
      apiClient.get("/purchase/suppliers", { params: { page: 1, page_size: 100 } }),
      apiClient.get("/purchase/orders", { params: { page: 1, page_size: 100 } }),
    ]);
    if (supRes.status === "fulfilled") {
      suppliers.value = supRes.value.data?.data?.list ?? [];
    }
    if (ordRes.status === "fulfilled") {
      orders.value = ordRes.value.data?.data?.list ?? [];
    }
  } catch {
    error.value = "加载采购数据失败";
  } finally {
    loading.value = false;
  }
});
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
        :loading="loading"
        :error="error"
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
        :loading="loading"
        :total="orders.length"
        :page-size="orders.length"
      >
        <template #status="{ row }">
          <el-tag :type="statusMap[row.status]?.type || 'info'" size="small">{{ statusMap[row.status]?.label || row.status }}</el-tag>
        </template>
        <template #total_amount="{ row }">{{ row.currency }} {{ row.total_amount }}</template>
        <template #actions>
          <el-button type="primary" size="small">审核</el-button>
          <el-button type="success" size="small">收货</el-button>
        </template>
      </ProTable>
    </el-card>
  </div>
</template>
