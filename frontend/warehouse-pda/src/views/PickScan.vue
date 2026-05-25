<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import { useWarehouseStore, getErrorMessage } from "@/stores/warehouse";

const route = useRoute();
const router = useRouter();
const store = useWarehouseStore();
const scanQty = ref(1);
const scanning = ref(false);

const outboundId = route.query.outbound_id as string;

onMounted(async () => {
  if (outboundId) {
    await store.fetchPickTasks(outboundId);
  }
});

async function handlePick(taskId: string, maxQty: number) {
  if (scanning.value) return;
  scanning.value = true;
  try {
    await store.pickScan(taskId, scanQty.value || maxQty);
    ElMessage.success("拣货成功");
    await store.fetchPickTasks(outboundId);
  } catch (e: unknown) {
    ElMessage.error(getErrorMessage(e, "拣货失败，请重试"));
  } finally {
    scanning.value = false;
  }
}
</script>

<template>
  <div style="padding: 16px">
    <el-page-header @back="router.back()" content="拣货扫码" />
    <el-empty v-if="!outboundId" description="请从拣货列表选择出库单" />
    <el-card v-for="task in store.pickTasks" :key="task.id" style="margin-top: 12px">
      <div>{{ task.sku_name }} ({{ task.sku_code }})</div>
      <div style="color: #909399; margin: 8px 0">
        库位 {{ task.location_code || "-" }} · 待拣 {{ task.quantity - (task.picked_quantity || 0) }}
      </div>
      <el-button type="primary" :loading="scanning" @click="handlePick(task.id, task.quantity)">
        确认拣货
      </el-button>
    </el-card>
  </div>
</template>
