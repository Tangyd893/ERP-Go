<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { ProTable } from "@erp/shared";
import { useProductStore } from "@/stores/product";

const mockData = ref([
  { id: "1", code: "TSHIRT-001", name: "纯棉T恤-白色-M", spu_name: "纯棉T恤", barcode: "6901234567890", weight: 0.2, sale_price: 15.99, currency: "USD", status: "active" },
  { id: "2", code: "TSHIRT-002", name: "纯棉T恤-黑色-L", spu_name: "纯棉T恤", barcode: "6901234567891", weight: 0.22, sale_price: 15.99, currency: "USD", status: "active" },
  { id: "3", code: "MUG-001", name: "陶瓷马克杯-350ml", spu_name: "陶瓷马克杯", barcode: "6901234567892", weight: 0.35, sale_price: 12.99, currency: "USD", status: "active" },
]);

const productStore = useProductStore();

const displayData = computed(() =>
  productStore.products.length > 0 ? productStore.products : mockData.value
);
const displayTotal = computed(() =>
  productStore.products.length > 0 ? productStore.productTotal : mockData.value.length
);

const columns = [
  { prop: "code", label: "SKU编码", width: 150 },
  { prop: "name", label: "SKU名称", minWidth: 180 },
  { prop: "spu_name", label: "所属SPU", width: 120 },
  { prop: "barcode", label: "条码", width: 150 },
  { prop: "weight", label: "重量(kg)", width: 100, align: "right" as const },
  { prop: "sale_price", label: "售价", width: 120, align: "right" as const },
  { prop: "status", label: "状态", width: 100 },
  { prop: "actions", label: "操作", width: 200, fixed: "right" as const },
];

onMounted(() => {
  productStore.fetchProducts(1, 20);
});
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

      <ProTable
        :columns="columns"
        :data="displayData"
        :loading="productStore.loading"
        :total="displayTotal"
        @page-change="(page: number) => productStore.fetchProducts(page, 20)"
      >
        <template #sale_price="{ row }">
          {{ row.currency }} {{ row.sale_price.toFixed(2) }}
        </template>
        <template #status="{ row }">
          <el-tag :type="row.status === 'active' ? 'success' : 'danger'" size="small">
            {{ row.status === 'active' ? '启售' : '停售' }}
          </el-tag>
        </template>
        <template #actions>
          <el-button type="primary" size="small">详情</el-button>
          <el-button type="warning" size="small">编辑</el-button>
        </template>
      </ProTable>
    </el-card>
  </div>
</template>
