<script setup lang="ts">
import { ref } from "vue";

const tableData = ref([
  { id: "1", username: "admin", nickname: "系统管理员", email: "admin@erp.com", status: "active", created_at: "2026-01-01" },
  { id: "2", username: "operator", nickname: "运营人员", email: "op@erp.com", status: "active", created_at: "2026-01-15" },
]);

const columns = [
  { prop: "username", label: "用户名", width: 120 },
  { prop: "nickname", label: "昵称", width: 140 },
  { prop: "email", label: "邮箱", width: 200 },
  { prop: "status", label: "状态", width: 100 },
  { prop: "created_at", label: "创建时间", width: 160 },
  { prop: "actions", label: "操作", fixed: "right" },
];

const dialogVisible = ref(false);
const formData = ref({ username: "", password: "", nickname: "", email: "", phone: "" });

const openCreateDialog = () => {
  formData.value = { username: "", password: "", nickname: "", email: "", phone: "" };
  dialogVisible.value = true;
};

const handleSave = () => {
  dialogVisible.value = false;
};
</script>

<template>
  <div>
    <el-card>
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center">
          <span>用户管理</span>
          <el-button type="primary" @click="openCreateDialog">新建用户</el-button>
        </div>
      </template>

      <el-table :data="tableData" stripe>
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column prop="nickname" label="昵称" width="140" />
        <el-table-column prop="email" label="邮箱" width="200" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'danger'" size="small">
              {{ row.status === 'active' ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="160" />
        <el-table-column label="操作" width="220" fixed="right">
          <template #default>
            <el-button type="primary" size="small">编辑</el-button>
            <el-button type="warning" size="small">分配角色</el-button>
            <el-button type="danger" size="small">禁用</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        style="margin-top: 16px; justify-content: flex-end"
        background
        layout="total, prev, pager, next"
        :total="2"
        :page-size="20"
      />
    </el-card>

    <el-dialog v-model="dialogVisible" title="新建用户" width="500px">
      <el-form :model="formData" label-width="80px">
        <el-form-item label="用户名" required>
          <el-input v-model="formData.username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="密码" required>
          <el-input v-model="formData.password" type="password" placeholder="请输入密码" />
        </el-form-item>
        <el-form-item label="昵称">
          <el-input v-model="formData.nickname" placeholder="请输入昵称" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="formData.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="手机号">
          <el-input v-model="formData.phone" placeholder="请输入手机号" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>
