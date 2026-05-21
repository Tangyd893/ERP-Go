<script setup lang="ts">
import { ref } from "vue";

const orgTree = ref([
  {
    id: "1",
    name: "总公司",
    code: "headquarters",
    children: [
      { id: "2", name: "销售部", code: "sales" },
      { id: "3", name: "仓储部", code: "warehouse" },
      { id: "4", name: "财务部", code: "finance" },
    ],
  },
]);

const defaultExpandedKeys = ref(["1"]);

const treeProps = {
  children: "children",
  label: "name",
};
</script>

<template>
  <div>
    <el-row :gutter="16">
      <el-col :span="8">
        <el-card>
          <template #header>
            <span>组织架构</span>
          </template>
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
