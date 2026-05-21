<script setup lang="ts">
import { ref } from "vue";

const tableData = ref([
  { id: "1", username: "admin", action: "login", resource_type: "user", resource_id: "1", detail: "登录成功", ip: "192.168.1.1", created_at: "2026-05-21 10:30:00" },
  { id: "2", username: "admin", action: "create", resource_type: "user", resource_id: "2", detail: "创建用户 operator", ip: "192.168.1.1", created_at: "2026-05-21 10:35:00" },
  { id: "3", username: "operator", action: "login", resource_type: "user", resource_id: "2", detail: "登录成功", ip: "192.168.1.2", created_at: "2026-05-21 11:00:00" },
]);

const actionLabels: Record<string, { type: string; label: string }> = {
  login: { type: "success", label: "登录" },
  logout: { type: "info", label: "登出" },
  create: { type: "primary", label: "创建" },
  update: { type: "warning", label: "更新" },
  delete: { type: "danger", label: "删除" },
  export: { type: "warning", label: "导出" },
  permission_change: { type: "danger", label: "权限变更" },
};
</script>

<template>
  <div>
    <el-card>
      <template #header>
        <span>操作审计日志</span>
      </template>

      <el-form inline>
        <el-form-item label="用户">
          <el-input placeholder="搜索用户名" clearable style="width: 200px" />
        </el-form-item>
        <el-form-item label="操作类型">
          <el-select placeholder="全部" clearable style="width: 150px">
            <el-option label="登录" value="login" />
            <el-option label="登出" value="logout" />
            <el-option label="创建" value="create" />
            <el-option label="更新" value="update" />
            <el-option label="删除" value="delete" />
            <el-option label="导出" value="export" />
            <el-option label="权限变更" value="permission_change" />
          </el-select>
        </el-form-item>
        <el-form-item label="时间范围">
          <el-date-picker type="daterange" range-separator="至" start-placeholder="开始日期" end-placeholder="结束日期" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary">查询</el-button>
          <el-button>重置</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="tableData" stripe>
        <el-table-column prop="username" label="用户" width="120" />
        <el-table-column prop="action" label="操作" width="120">
          <template #default="{ row }">
            <el-tag
              :type="actionLabels[row.action]?.type || 'info'"
              size="small"
            >
              {{ actionLabels[row.action]?.label || row.action }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="resource_type" label="资源类型" width="120" />
        <el-table-column prop="detail" label="操作详情" min-width="200" />
        <el-table-column prop="ip" label="IP 地址" width="140" />
        <el-table-column prop="created_at" label="操作时间" width="180" />
      </el-table>

      <el-pagination
        style="margin-top: 16px; justify-content: flex-end"
        background
        layout="total, prev, pager, next"
        :total="3"
        :page-size="20"
      />
    </el-card>
  </div>
</template>
