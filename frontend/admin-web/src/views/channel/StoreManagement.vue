<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { ProTable, ProForm } from "@erp/shared";
import { useChannelStore } from "@/stores/channel";

const mockData = ref([
  { id: "1", name: "美国站店铺A", platform_code: "amazon_us", site: "Amazon.com", auth_status: "authorized", status: "active", created_at: "2026-01-01" },
  { id: "2", name: "英国站店铺B", platform_code: "amazon_uk", site: "Amazon.co.uk", auth_status: "authorized", status: "active", created_at: "2026-02-15" },
]);

const channelStore = useChannelStore();

const displayData = computed(() =>
  channelStore.stores.length > 0 ? channelStore.stores : mockData.value
);

const columns = [
  { prop: "name", label: "店铺名称", minWidth: 180 },
  { prop: "platform_code", label: "平台", width: 150 },
  { prop: "site", label: "站点", width: 150 },
  { prop: "auth_status", label: "授权状态", width: 120 },
  { prop: "status", label: "店铺状态", width: 120 },
  { prop: "actions", label: "操作", width: 280, fixed: "right" as const },
];

const dialogVisible = ref(false);
const formData = ref({ name: "", platform_code: "amazon_us", site: "", store_code: "" });

const proFormFields = [
  { label: "店铺名称", prop: "name", type: "input" as const, placeholder: "如: 美国站店铺A", required: true },
  {
    label: "平台", prop: "platform_code", type: "select" as const, required: true,
    options: [
      { label: "Amazon 美国", value: "amazon_us" },
      { label: "Amazon 英国", value: "amazon_uk" },
      { label: "Amazon 日本", value: "amazon_jp" },
    ],
  },
  { label: "站点", prop: "site", type: "input" as const, placeholder: "如: Amazon.com", required: true },
  { label: "店铺编码", prop: "store_code", type: "input" as const, placeholder: "店铺唯一标识" },
];

const handleSubmit = (data: Record<string, unknown>) => {
  channelStore.createStore(data as Partial<{ name: string; platform_code: string; site: string; store_code: string }>).then(() => {
    dialogVisible.value = false;
    channelStore.fetchStores();
  });
};

onMounted(() => {
  channelStore.fetchStores();
});
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

      <ProTable
        :columns="columns"
        :data="displayData"
        :loading="channelStore.loading"
        :total="displayData.length"
        @page-change="() => channelStore.fetchStores()"
      >
        <template #auth_status="{ row }">
          <el-tag :type="row.auth_status === 'authorized' ? 'success' : 'warning'" size="small">
            {{ row.auth_status === 'authorized' ? '已授权' : '未授权' }}
          </el-tag>
        </template>
        <template #status>
          <el-tag type="success" size="small">运营中</el-tag>
        </template>
        <template #actions>
          <el-button type="primary" size="small">授权</el-button>
          <el-button type="success" size="small">同步订单</el-button>
          <el-button type="warning" size="small">编辑</el-button>
        </template>
      </ProTable>
    </el-card>

    <ProForm
      v-model="dialogVisible"
      title="添加店铺"
      :form-data="formData"
      :fields="proFormFields"
      @submit="handleSubmit"
    />
  </div>
</template>
