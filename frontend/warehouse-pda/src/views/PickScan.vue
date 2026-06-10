<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import {
  useWarehouseStore,
  getErrorMessage,
  isDuplicateError,
  isNetworkError,
} from "@/stores/warehouse";

const route = useRoute();
const router = useRouter();
const store = useWarehouseStore();

const outboundId = (route.query.outbound_id as string) || "";
const scanInput = ref("");
const scanQty = ref(1);
const scanning = ref(false);
const scanInputRef = ref<HTMLInputElement | null>(null);

// Auto-focus scan input
onMounted(async () => {
  if (outboundId) {
    await store.fetchPickTasks(outboundId);
  }
  scanInputRef.value?.focus();
});

// Re-focus after each scan
watch(scanning, (val) => {
  if (!val) {
    setTimeout(() => scanInputRef.value?.focus(), 100);
  }
});

const pendingTasks = computed(() =>
  store.pickTasks.filter((t) => t.status !== "done" && t.picked_quantity < t.quantity)
);

const recentPicks = computed(() =>
  store.scanHistory.filter((r) => r.type === "pick").slice(0, 10)
);

function findTaskByScan(value: string) {
  const v = value.trim().toUpperCase();
  if (!v) return null;
  return pendingTasks.value.find(
    (t) =>
      t.sku_code?.toUpperCase() === v ||
      t.location_code?.toUpperCase() === v ||
      t.id?.toUpperCase() === v
  ) || null;
}

async function handleScan() {
  const value = scanInput.value.trim();
  if (!value || scanning.value) return;

  const task = findTaskByScan(value);
  if (!task) {
    ElMessage.warning(`未找到匹配的待拣任务: ${value}`);
    scanInput.value = "";
    scanInputRef.value?.focus();
    return;
  }

  scanning.value = true;
  const label = `${task.sku_name || task.sku_code} @ ${task.location_code || "-"}`;
  try {
    await store.pickScan(task.id, scanQty.value, label);
    ElMessage.success(`✓ ${label} 拣货完成`);
    scanInput.value = "";
    await store.fetchPickTasks(outboundId);
  } catch (e: unknown) {
    if (isDuplicateError(e)) {
      ElMessage.info(`已拣货（重复扫描）: ${label}`);
      scanInput.value = "";
      await store.fetchPickTasks(outboundId);
    } else if (isNetworkError(e)) {
      ElMessage.error("网络不可用，请检查连接后重试");
    } else {
      ElMessage.error(getErrorMessage(e, "拣货失败"));
    }
  } finally {
    scanning.value = false;
  }
}

async function handleQuickPick(taskId: string, label: string) {
  if (scanning.value) return;
  scanning.value = true;
  try {
    await store.pickScan(taskId, scanQty.value, label);
    ElMessage.success(`✓ ${label} 拣货完成`);
    await store.fetchPickTasks(outboundId);
  } catch (e: unknown) {
    if (isDuplicateError(e)) {
      ElMessage.info(`已拣货（重复扫描）: ${label}`);
      await store.fetchPickTasks(outboundId);
    } else if (isNetworkError(e)) {
      ElMessage.error("网络不可用，请检查连接后重试");
    } else {
      ElMessage.error(getErrorMessage(e, "拣货失败"));
    }
  } finally {
    scanning.value = false;
  }
}
</script>

<template>
  <div style="padding: 16px; max-width: 480px; margin: 0 auto">
    <!-- Header -->
    <div style="display: flex; align-items: center; justify-content: space-between; margin-bottom: 12px">
      <el-page-header @back="router.back()" content="拣货扫码" />
      <el-tag :type="store.online ? 'success' : 'danger'" size="small" effect="dark">
        {{ store.online ? "在线" : "离线" }}
      </el-tag>
    </div>

    <el-empty v-if="!outboundId" description="请从拣货列表选择出库单" />

    <template v-else>
      <!-- Scan Input -->
      <el-card style="margin-bottom: 12px">
        <div style="font-size: 13px; color: #909399; margin-bottom: 8px">
          扫描条码 / 输入 SKU 编码 · 库位号 · 任务 ID
        </div>
        <el-input
          ref="scanInputRef"
          v-model="scanInput"
          placeholder="扫描或输入…"
          size="large"
          clearable
          :disabled="scanning"
          @keyup.enter="handleScan"
        >
          <template #append>
            <el-button
              :loading="scanning"
              :disabled="!scanInput.trim()"
              @click="handleScan"
            >
              确认
            </el-button>
          </template>
        </el-input>
        <div style="margin-top: 8px; display: flex; align-items: center; gap: 8px">
          <span style="font-size: 13px; color: #909399">数量:</span>
          <el-input-number v-model="scanQty" :min="1" :max="9999" size="small" />
        </div>
      </el-card>

      <!-- Pending Tasks -->
      <div v-if="pendingTasks.length === 0 && store.pickTasks.length > 0" style="text-align: center; padding: 24px; color: #67c23a">
        全部拣货任务已完成 ✓
      </div>

      <el-card
        v-for="task in pendingTasks"
        :key="task.id"
        style="margin-bottom: 8px"
        :body-style="{ padding: '12px 16px' }"
      >
        <div style="display: flex; justify-content: space-between; align-items: center">
          <div>
            <div style="font-weight: 500">{{ task.sku_name || task.sku_code }}</div>
            <div style="font-size: 12px; color: #909399">
              {{ task.sku_code }} · 库位 {{ task.location_code || "-" }}
              · {{ task.picked_quantity || 0 }}/{{ task.quantity }}
            </div>
          </div>
          <el-button
            type="primary"
            size="small"
            :loading="scanning"
            @click="handleQuickPick(task.id, task.sku_name || task.sku_code)"
          >
            拣货
          </el-button>
        </div>
      </el-card>

      <!-- Recent Scan History -->
      <div v-if="recentPicks.length > 0" style="margin-top: 16px">
        <div style="font-size: 13px; color: #909399; margin-bottom: 8px">
          最近扫描记录
        </div>
        <div
          v-for="r in recentPicks"
          :key="r.id"
          style="font-size: 12px; padding: 4px 0; border-bottom: 1px solid #f0f0f0; display: flex; justify-content: space-between"
        >
          <span>
            <span :style="{ color: r.success ? '#67c23a' : '#f56c6c' }">
              {{ r.success ? "✓" : "✗" }}
            </span>
            {{ r.targetLabel }}
          </span>
          <span style="color: #c0c4cc">{{ new Date(r.time).toLocaleTimeString() }}</span>
        </div>
      </div>
    </template>
  </div>
</template>
