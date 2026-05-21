<script setup lang="ts">
import { ref, onMounted } from "vue";
import { ElMessage } from "element-plus";
import { useRoleStore } from "@/stores/role";
import { ProForm } from "@erp/shared";

const roleStore = useRoleStore();

const dialogVisible = ref(false);
const newRoleForm = ref({ name: "", code: "", description: "" });

const formFields = [
  { prop: "name", label: "角色名", type: "input" as const, placeholder: "请输入角色名称", required: true },
  { prop: "code", label: "编码", type: "input" as const, placeholder: "请输入角色编码", required: true },
  { prop: "description", label: "描述", type: "textarea" as const, placeholder: "请输入角色描述" },
];

const currentPage = ref(1);
const pageSize = ref(20);

onMounted(() => {
  roleStore.fetchRoles(1, 20);
});

function openCreateDialog() {
  newRoleForm.value = { name: "", code: "", description: "" };
  dialogVisible.value = true;
}

async function handleCreateRole(data: Record<string, unknown>) {
  try {
    await roleStore.createRole(data);
    ElMessage.success("角色创建成功");
    dialogVisible.value = false;
    roleStore.fetchRoles(currentPage.value, pageSize.value);
  } catch {
    ElMessage.error("创建角色失败");
  }
}

function handlePageChange(page: number) {
  currentPage.value = page;
  roleStore.fetchRoles(page, pageSize.value);
}

async function handleDeleteRole(id: string) {
  try {
    await roleStore.deleteRole(id);
    ElMessage.success("角色已删除");
    roleStore.fetchRoles(currentPage.value, pageSize.value);
  } catch {
    ElMessage.error("删除角色失败");
  }
}

const permDialogVisible = ref(false);
const currentRole = ref<{ role_id: string; name: string } | null>(null);

function openPermDialog(row: { role_id: string; name: string }) {
  currentRole.value = row;
  permDialogVisible.value = true;
}

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

      <el-table :data="roleStore.roles" stripe v-loading="roleStore.loading">
        <el-table-column prop="name" label="角色名" width="150" />
        <el-table-column prop="code" label="角色编码" width="180" />
        <el-table-column prop="description" label="描述" min-width="200" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="160" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="openPermDialog(row)">分配权限</el-button>
            <el-button type="warning" size="small">编辑</el-button>
            <el-popconfirm title="确定要删除该角色吗？" @confirm="handleDeleteRole(row.role_id)">
              <template #reference>
                <el-button type="danger" size="small">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        style="margin-top: 16px; justify-content: flex-end"
        background
        layout="total, prev, pager, next"
        :total="roleStore.total"
        :page-size="pageSize"
        :current-page="currentPage"
        @current-change="handlePageChange"
      />
    </el-card>

    <ProForm
      v-model="dialogVisible"
      title="新建角色"
      :form-data="newRoleForm"
      :fields="formFields"
      @submit="handleCreateRole"
    />

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
