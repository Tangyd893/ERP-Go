<script setup lang="ts">
import { ref, onMounted } from "vue";
import { ElMessage } from "element-plus";
import { useUserStore } from "@/stores/user";
import { ProForm } from "@erp/shared";

const userStore = useUserStore();

const dialogVisible = ref(false);
const newUserForm = ref({ username: "", password: "", nickname: "", email: "", phone: "" });

const formFields = [
  { prop: "username", label: "用户名", type: "input" as const, placeholder: "请输入用户名", required: true },
  { prop: "password", label: "密码", type: "input" as const, placeholder: "请输入密码", required: true },
  { prop: "nickname", label: "昵称", type: "input" as const, placeholder: "请输入昵称" },
  { prop: "email", label: "邮箱", type: "input" as const, placeholder: "请输入邮箱" },
  { prop: "phone", label: "手机号", type: "input" as const, placeholder: "请输入手机号" },
];

const currentPage = ref(1);
const pageSize = ref(20);

onMounted(() => {
  userStore.fetchUsers(1, 20);
});

function openCreateDialog() {
  newUserForm.value = { username: "", password: "", nickname: "", email: "", phone: "" };
  dialogVisible.value = true;
}

async function handleCreateUser(data: Record<string, unknown>) {
  try {
    await userStore.createUser(data);
    ElMessage.success("用户创建成功");
    dialogVisible.value = false;
    userStore.fetchUsers(currentPage.value, pageSize.value);
  } catch {
    ElMessage.error("创建用户失败");
  }
}

function handlePageChange(page: number) {
  currentPage.value = page;
  userStore.fetchUsers(page, pageSize.value);
}

function handleDeleteUser(userId: string) {
  userStore.deleteUser(userId);
}
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

      <el-table :data="userStore.users" stripe v-loading="userStore.loading">
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column prop="nickname" label="昵称" width="140" />
        <el-table-column prop="email" label="邮箱" width="200" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="160" />
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small">编辑</el-button>
            <el-button type="warning" size="small">分配角色</el-button>
            <el-popconfirm title="确定要删除该用户吗？" @confirm="handleDeleteUser(row.user_id)">
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
        :total="userStore.total"
        :page-size="pageSize"
        :current-page="currentPage"
        @current-change="handlePageChange"
      />
    </el-card>

    <ProForm
      v-model="dialogVisible"
      title="新建用户"
      :form-data="newUserForm"
      :fields="formFields"
      @submit="handleCreateUser"
    />
  </div>
</template>
