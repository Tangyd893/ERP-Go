<script setup lang="ts">
import { ref, onMounted } from "vue";
import { apiClient } from "@erp/shared";

interface OrgNode {
  id: string; name: string; code: string; children?: OrgNode[];
}

const orgTree = ref<OrgNode[]>([]);
const loading = ref(false);
const error = ref("");

const defaultExpandedKeys = ref<string[]>([]);

const treeProps = {
  children: "children",
  label: "name",
};

onMounted(async () => {
  loading.value = true;
  try {
    const res = await apiClient.get("/tenant/organizations");
    const list = res.data?.data?.list ?? res.data?.data ?? [];
    orgTree.value = list;
    if (list.length > 0) {
      defaultExpandedKeys.value = [list[0].id];
    }
  } catch {
    error.value = "加载组织架构失败";
  } finally {
    loading.value = false;
  }
});
</script>

<template>
  <div>
    <!-- 加载/错误态 -->
    <div v-if="loading" style="text-align: center; padding: 48px 0">
      <el-skeleton :rows="6" animated />
    </div>
    <el-alert
      v-else-if="error"
      :title="error"
      type="error"
      :closable="false"
      show-icon
      style="margin-bottom: 16px"
    />
    <el-empty v-else-if="orgTree.length === 0" description="暂无组织架构数据，请先创建组织">
      <el-button type="primary" disabled>新建组织</el-button>
    </el-empty>

    <!-- 正常 -->
    <el-row v-else :gutter="16">
      <el-col :span="8">
        <el-card>
          <template #header><span>组织架构</span></template>
          <el-tree
            :data="orgTree"
            :props="treeProps"
            :default-expanded-keys="defaultExpandedKeys"
            node-key="id"
            highlight-current
          >
            <template #default="{ data }">
              <span style="font-size: 14px">{{ data.name }}</span>
              <span style="margin-left: 8px; font-size: 12px; color: #909399">{{ data.code }}</span>
            </template>
          </el-tree>
        </el-card>
      </el-col>
      <el-col :span="16">
        <el-card>
          <template #header>
            <div style="display: flex; justify-content: space-between; align-items: center">
              <span>部门/岗位</span>
              <el-button type="primary" size="small">新建部门</el-button>
            </div>
          </template>
          <el-table :data="[]" empty-text="请在左侧选择组织节点">
            <el-table-column prop="name" label="名称" />
            <el-table-column prop="code" label="编码" />
            <el-table-column prop="manager" label="负责人" />
            <el-table-column label="操作">
              <template #default>
                <el-button type="primary" size="small">编辑</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>
