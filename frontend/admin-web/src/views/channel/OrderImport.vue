<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { ProTable } from "@erp/shared";
import { useChannelStore } from "@/stores/channel";

const mockData = ref([
  { id: "1", import_type: "csv", store_name: "美国站店铺A", status: "completed", total_rows: 150, success_rows: 148, failed_rows: 2, created_at: "2026-05-21 10:00" },
  { id: "2", import_type: "api", store_name: "英国站店铺B", status: "processing", total_rows: 0, success_rows: 0, failed_rows: 0, created_at: "2026-05-21 11:00" },
]);

const channelStore = useChannelStore();

const displayData = computed(() =>
  channelStore.importTasks.length > 0 ? channelStore.importTasks : mockData.value
);

const columns = [
  { prop: "import_type", label: "导入方式", width: 100 },
  { prop: "store_name", label: "店铺", width: 150 },
  { prop: "status", label: "状态", width: 120 },
  { prop: "total_rows", label: "总行数", width: 100, align: "right" as const },
  { prop: "success_rows", label: "成功", width: 100, align: "right" as const },
  { prop: "failed_rows", label: "失败", width: 100, align: "right" as const },
  { prop: "created_at", label: "创建时间", width: 180 },
  { prop: "actions", label: "操作", width: 160, fixed: "right" as const },
];

const uploadDialogVisible = ref(false);

onMounted(() => {
  channelStore.fetchImportTasks(1, 20);
});
</script>

<template>
  <div>
    <el-card>
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center">
          <span>订单导入任务</span>
          <div>
            <el-button type="primary" @click="uploadDialogVisible = true">上传 CSV 导入</el-button>
            <el-button type="success">API 同步</el-button>
          </div>
        </div>
      </template>

      <ProTable
        :columns="columns"
        :data="displayData"
        :loading="channelStore.loading"
        :total="displayData.length"
        @page-change="(page: number) => channelStore.fetchImportTasks(page, 20)"
      >
        <template #import_type="{ row }">
          <el-tag :type="row.import_type === 'csv' ? 'warning' : 'primary'" size="small">
            {{ row.import_type.toUpperCase() }}
          </el-tag>
        </template>
        <template #status="{ row }">
          <el-tag :type="row.status === 'completed' ? 'success' : 'warning'" size="small">
            {{ row.status === 'completed' ? '已完成' : '处理中' }}
          </el-tag>
        </template>
        <template #success_rows="{ row }">
          <span style="color: #67C23A">{{ row.success_rows }}</span>
        </template>
        <template #failed_rows="{ row }">
          <span :style="{ color: row.failed_rows > 0 ? '#F56C6C' : '#909399' }">{{ row.failed_rows }}</span>
        </template>
        <template #actions="{ row }">
          <el-button type="primary" size="small" :disabled="row.status !== 'completed'">查看结果</el-button>
        </template>
      </ProTable>
    </el-card>

    <el-dialog v-model="uploadDialogVisible" title="上传 CSV 订单导入" width="500px">
      <el-form label-width="100px">
        <el-form-item label="店铺" required>
          <el-select placeholder="选择店铺" style="width: 100%">
            <el-option label="美国站店铺A" value="store-1" />
            <el-option label="英国站店铺B" value="store-2" />
          </el-select>
        </el-form-item>
        <el-form-item label="CSV文件" required>
          <el-upload
            drag
            :auto-upload="false"
            accept=".csv"
          >
            <el-icon style="font-size: 48px; color: #c0c4cc"><UploadFilled /></el-icon>
            <div style="color: #606266; margin-top: 8px">将 CSV 文件拖到此处，或点击上传</div>
            <template #tip>
              <div style="margin-top: 8px; font-size: 12px; color: #909399">
                文件需包含: OrderID, SKU, Quantity, BuyerName, Address, Phone 等字段
              </div>
            </template>
          </el-upload>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="uploadDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="uploadDialogVisible = false">开始导入</el-button>
      </template>
    </el-dialog>
  </div>
</template>
