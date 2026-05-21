<script setup lang="ts">
import { ref } from "vue";

const tableData = ref([
  { id: "1", name: "美国站店铺A", platform_code: "amazon_us", site: "Amazon.com", auth_status: "authorized", status: "active", created_at: "2026-01-01" },
  { id: "2", name: "英国站店铺B", platform_code: "amazon_uk", site: "Amazon.co.uk", auth_status: "authorized", status: "active", created_at: "2026-02-15" },
]);

const dialogVisible = ref(false);
const form = ref({ name: "", platform_code: "amazon_us", site: "", store_code: "" });
</script>

<template>
  <div>
    <el-card>
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center">
          <span>店铺授权管理</span>
          <el-button type="primary" @click="dialogVisible = true">添加店铺</el-button>
        </div>
      </template>

      <el-table :data="tableData" stripe>
        <el-table-column prop="name" label="店铺名称" min-width="180" />
        <el-table-column prop="platform_code" label="平台" width="150" />
        <el-table-column prop="site" label="站点" width="150" />
        <el-table-column prop="auth_status" label="授权状态" width="120">
          <template #default="{ row }">
            <el-tag :type="row.auth_status === 'authorized' ? 'success' : 'warning'" size="small">
              {{ row.auth_status === 'authorized' ? '已授权' : '未授权' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="店铺状态" width="120">
          <template #default>
            <el-tag type="success" size="small">运营中</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="280" fixed="right">
          <template #default>
            <el-button type="primary" size="small">授权</el-button>
            <el-button type="success" size="small">同步订单</el-button>
            <el-button type="warning" size="small">编辑</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" title="添加店铺" width="500px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="店铺名称" required>
          <el-input v-model="form.name" placeholder="如: 美国站店铺A" />
        </el-form-item>
        <el-form-item label="平台" required>
          <el-select v-model="form.platform_code" style="width: 100%">
            <el-option label="Amazon 美国" value="amazon_us" />
            <el-option label="Amazon 英国" value="amazon_uk" />
            <el-option label="Amazon 日本" value="amazon_jp" />
          </el-select>
        </el-form-item>
        <el-form-item label="站点" required>
          <el-input v-model="form.site" placeholder="如: Amazon.com" />
        </el-form-item>
        <el-form-item label="店铺编码">
          <el-input v-model="form.store_code" placeholder="店铺唯一标识" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="dialogVisible = false">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>
