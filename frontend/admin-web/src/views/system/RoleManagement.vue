<script setup lang="ts">
import { ref } from "vue";

const tableData = ref([
  { id: "1", name: "超级管理员", code: "super_admin", description: "系统超级管理员", status: "active", created_at: "2026-01-01" },
  { id: "2", name: "运营专员", code: "operator", description: "负责订单和仓库日常操作", status: "active", created_at: "2026-01-15" },
]);

const dialogVisible = ref(false);
const formData = ref({ name: "", code: "", description: "" });

const openCreateDialog = () => {
  formData.value = { name: "", code: "", description: "" };
  dialogVisible.value = true;
};

const handleSave = () => {
  dialogVisible.value = false;
};

const permDialogVisible = ref(false);
const currentRole = ref<{ id: string; name: string } | null>(null);

const openPermDialog = (row: { id: string; name: string }) => {
  currentRole.value = row;
  permDialogVisible.value = true;
};

const allPermissions = ref([
  { code: "user:read", name: "查看用户", checked: true },
  { code: "user:create", name: "创建用户", checked: true },
  { code: "user:update", name: "编辑用户", checked: true },
  { code: "user:disable", name: "禁用用户", checked: true },
  { code: "role:read", name: "查看角色", checked: true },
  { code: "role:create", name: "创建角色", checked: true },
  { code: "role:update", name: "编辑角色", checked: false },
  { code: "role:delete", name: "删除角色", checked: false },
]);
</script>

<template>
  <div>
    <el-card>
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center">
          <span>角色管理</span>
          <el-button type="primary" @click="openCreateDialog">新建角色</el-button>
        </div>
      </template>

      <el-table :data="tableData" stripe>
        <el-table-column prop="name" label="角色名" width="150" />
        <el-table-column prop="code" label="角色编码" width="180" />
        <el-table-column prop="description" label="描述" min-width="200" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'danger'" size="small">
              {{ row.status === 'active' ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="160" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="openPermDialog(row)">分配权限</el-button>
            <el-button type="warning" size="small">编辑</el-button>
            <el-button type="danger" size="small">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" title="新建角色" width="500px">
      <el-form :model="formData" label-width="80px">
        <el-form-item label="角色名" required>
          <el-input v-model="formData.name" placeholder="请输入角色名称" />
        </el-form-item>
        <el-form-item label="编码" required>
          <el-input v-model="formData.code" placeholder="请输入角色编码" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="formData.description" type="textarea" placeholder="请输入角色描述" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="permDialogVisible" title="分配权限" width="500px">
      <template v-if="currentRole">
        <p style="margin-bottom: 16px">角色: <el-tag>{{ currentRole.name }}</el-tag></p>
        <el-checkbox-group>
          <div v-for="perm in allPermissions" :key="perm.code" style="margin-bottom: 8px">
            <el-checkbox v-model="perm.checked" :label="perm.code">
              {{ perm.name }} <span style="color: #909399; font-size: 12px">({{ perm.code }})</span>
            </el-checkbox>
          </div>
        </el-checkbox-group>
      </template>
      <template #footer>
        <el-button @click="permDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="permDialogVisible = false">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>
