<script setup lang="ts">
import { ref } from "vue";

const tableData = ref([
  { id: "1", sku_code: "TSHIRT-001", sku_name: "纯棉T恤-白色-M", platform_code: "amazon_us", platform_sku: "B0XXX001", asin: "B0XXX001", fnsku: "X001AAA", store: "美国站店铺A" },
  { id: "2", sku_code: "TSHIRT-002", sku_name: "纯棉T恤-黑色-L", platform_code: "amazon_us", platform_sku: "B0XXX002", asin: "B0XXX002", fnsku: "X002BBB", store: "美国站店铺A" },
  { id: "3", sku_code: "MUG-001", sku_name: "陶瓷马克杯-350ml", platform_code: "amazon_uk", platform_sku: "B0XXX003", asin: "B0XXX003", fnsku: "X003CCC", store: "英国站店铺B" },
]);

const dialogVisible = ref(false);
const form = ref({ sku_id: "", store_id: "", platform_sku: "", asin: "", fnsku: "" });

const openCreateDialog = () => {
  form.value = { sku_id: "", store_id: "", platform_sku: "", asin: "", fnsku: "" };
  dialogVisible.value = true;
};
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

      <el-table :data="tableData" stripe>
        <el-table-column prop="sku_code" label="内部SKU" width="140" />
        <el-table-column prop="sku_name" label="SKU名称" min-width="180" />
        <el-table-column prop="store" label="店铺" width="150" />
        <el-table-column prop="platform_code" label="平台" width="120" />
        <el-table-column prop="platform_sku" label="平台SKU" width="140" />
        <el-table-column prop="asin" label="ASIN" width="120" />
        <el-table-column prop="fnsku" label="FNSKU" width="120" />
        <el-table-column label="操作" width="160" fixed="right">
          <template #default>
            <el-button type="primary" size="small">编辑</el-button>
            <el-button type="danger" size="small">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
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
        <el-button type="primary" @click="dialogVisible = false">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>
