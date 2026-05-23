<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { ProTable } from "@erp/shared";
import { useProductStore } from "@/stores/product";

const productStore = useProductStore();

const displayData = computed(() => productStore.skuMappings);

const columns = [
  { prop: "sku_code", label: "内部SKU", width: 140 },
  { prop: "sku_name", label: "SKU名称", minWidth: 180 },
  { prop: "store", label: "店铺", width: 150 },
  { prop: "platform_code", label: "平台", width: 120 },
  { prop: "platform_sku", label: "平台SKU", width: 140 },
  { prop: "asin", label: "ASIN", width: 120 },
  { prop: "fnsku", label: "FNSKU", width: 120 },
  { prop: "actions", label: "操作", width: 160, fixed: "right" as const },
];

const dialogVisible = ref(false);
const form = ref({ sku_id: "", store_id: "", platform_sku: "", asin: "", fnsku: "" });

const openCreateDialog = () => {
  form.value = { sku_id: "", store_id: "", platform_sku: "", asin: "", fnsku: "" };
  dialogVisible.value = true;
};

const handleSubmit = () => {
  productStore.createSKUMapping(form.value).then(() => {
    dialogVisible.value = false;
    productStore.fetchSKUMappings();
  });
};

onMounted(() => {
  productStore.fetchSKUMappings();
});
</script>

<template>
  <div>
    <el-card>
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center">
          <span>SKU 映射管理</span>
          <el-button type="primary" @click="openCreateDialog">新建映射</el-button>
        </div>
      </template>

      <ProTable
        :columns="columns"
        :data="displayData"
        :loading="productStore.loading"
        :total="displayData.length"
        @page-change="() => productStore.fetchSKUMappings()"
      >
        <template #actions>
          <el-button type="primary" size="small">编辑</el-button>
          <el-button type="danger" size="small">删除</el-button>
        </template>
      </ProTable>
    </el-card>

    <el-dialog v-model="dialogVisible" title="新建 SKU 映射" width="500px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="内部SKU" required>
          <el-select v-model="form.sku_id" placeholder="选择SKU" style="width: 100%">
            <el-option label="TSHIRT-001" value="1" />
            <el-option label="MUG-001" value="3" />
          </el-select>
        </el-form-item>
        <el-form-item label="店铺" required>
          <el-select v-model="form.store_id" placeholder="选择店铺" style="width: 100%">
            <el-option label="美国站店铺A" value="store-1" />
          </el-select>
        </el-form-item>
        <el-form-item label="平台SKU">
          <el-input v-model="form.platform_sku" placeholder="平台上的SKU编码" />
        </el-form-item>
        <el-form-item label="ASIN">
          <el-input v-model="form.asin" placeholder="Amazon ASIN" />
        </el-form-item>
        <el-form-item label="FNSKU">
          <el-input v-model="form.fnsku" placeholder="Amazon FNSKU" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>
