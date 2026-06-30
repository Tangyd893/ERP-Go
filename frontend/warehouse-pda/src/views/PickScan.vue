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
import ScanInput from "@/components/ScanInput.vue";
import ScanFeedback from "@/components/ScanFeedback.vue";

const route = useRoute();
const router = useRouter();
const store = useWarehouseStore();

const outboundId = (route.query.outbound_id as string) || "";
const scanValue = ref("");
const scanQty = ref(1);
const scanning = ref(false);
const scanInputRef = ref<InstanceType<typeof ScanInput> | null>(null);
const showHistory = ref(false);

/** 触发振动反馈 */
function vibrate(pattern: number | number[] = 50) {
  try {
    if (typeof navigator !== "undefined" && "vibrate" in navigator) {
      (navigator as Navigator).vibrate(pattern);
    }
  } catch { /* 不支持静默忽略 */ }
}

onMounted(async () => {
  if (outboundId) {
    await store.fetchPickTasks(outboundId);
  }
  scanInputRef.value?.focus();
});

watch(scanning, (val) => {
  if (!val) {
    setTimeout(() => scanInputRef.value?.focus(), 100);
  }
});

const pendingTasks = computed(() =>
  store.pickTasks.filter((t) => t.status !== "done" && t.picked_quantity < t.quantity)
);

function findTaskByScan(value: string) {
  const v = value.trim().toUpperCase();
  if (!v) return null;
  return (
    pendingTasks.value.find(
      (t) =>
        t.sku_code?.toUpperCase() === v ||
        t.location_code?.toUpperCase() === v ||
        t.id?.toUpperCase() === v
    ) || null
  );
}

async function handleScan(value: string) {
  if (scanning.value) return;

  const task = findTaskByScan(value);
  if (!task) {
    ElMessage.warning(`未找到匹配的待拣任务: ${value}`);
    vibrate([80, 100, 80, 100, 200]);
    scanInputRef.value?.focus();
    return;
  }

  scanning.value = true;
  const label = `${task.sku_name || task.sku_code} @ ${task.location_code || "-"}`;
  try {
    await store.pickScan(task.id, scanQty.value, label);
    ElMessage.success(`✓ ${label} 拣货完成`);
    vibrate(80);
    await store.fetchPickTasks(outboundId);
  } catch (e: unknown) {
    if (isDuplicateError(e)) {
      ElMessage.info(`已拣货（重复扫描）: ${label}`);
      vibrate([30, 50, 30]);
      await store.fetchPickTasks(outboundId);
    } else if (isNetworkError(e)) {
      ElMessage.error("网络不可用，请检查连接后重试");
      vibrate([80, 100, 80, 100, 200]);
    } else {
      ElMessage.error(getErrorMessage(e, "拣货失败"));
      vibrate([80, 100, 80, 100, 200]);
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
    vibrate(80);
    await store.fetchPickTasks(outboundId);
  } catch (e: unknown) {
    if (isDuplicateError(e)) {
      ElMessage.info(`已拣货（重复扫描）: ${label}`);
      vibrate([30, 50, 30]);
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
  <div class="pick-scan">
    <!-- Header -->
    <div class="pick-scan__header">
      <el-page-header @back="router.back()" content="拣货扫码" />
      <div class="pick-scan__header-right">
        <el-badge :value="store.scanHistory.filter(r => r.type === 'pick').length" :hidden="store.scanHistory.length === 0">
          <el-button size="small" text @click="showHistory = true">记录</el-button>
        </el-badge>
        <el-tag :type="store.online ? 'success' : 'danger'" size="small" effect="dark">
          {{ store.online ? "在线" : "离线" }}
        </el-tag>
      </div>
    </div>

    <el-empty v-if="!outboundId" description="请从拣货列表选择出库单" />

    <template v-else>
      <!-- Scan Input -->
      <ScanInput
        ref="scanInputRef"
        v-model="scanValue"
        hint="扫描条码 / 输入 SKU 编码 · 库位号 · 任务 ID"
        :disabled="scanning"
        @scan="handleScan"
      />

      <!-- Quantity selector -->
      <div class="pick-scan__qty">
        <span class="pick-scan__qty-label">数量:</span>
        <el-input-number
          v-model="scanQty"
          :min="1"
          :max="9999"
          size="small"
          class="pick-scan__qty-input"
        />
      </div>

      <!-- All done -->
      <div v-if="pendingTasks.length === 0 && store.pickTasks.length > 0" class="pick-scan__all-done">
        全部拣货任务已完成 ✓
      </div>

      <!-- Pending Tasks -->
      <el-card
        v-for="task in pendingTasks"
        :key="task.id"
        class="pick-scan__task-card"
        :body-style="{ padding: '12px 16px' }"
      >
        <div class="pick-scan__task-row">
          <div class="pick-scan__task-info">
            <div class="pick-scan__task-name">{{ task.sku_name || task.sku_code }}</div>
            <div class="pick-scan__task-meta">
              {{ task.sku_code }} · 库位 {{ task.location_code || "-" }}
              · {{ task.picked_quantity || 0 }}/{{ task.quantity }}
            </div>
          </div>
          <el-button
            type="primary"
            size="small"
            :loading="scanning"
            class="pick-scan__task-btn"
            @click="handleQuickPick(task.id, task.sku_name || task.sku_code)"
          >
            拣货
          </el-button>
        </div>
      </el-card>
    </template>

    <!-- Scan History Drawer -->
    <ScanFeedback
      v-if="showHistory"
      type="pick"
      @close="showHistory = false"
    />
  </div>
</template>

<style scoped>
.pick-scan {
  /* 移除内边距 — PdaLayout 的 pda-main 已提供 */
}

.pick-scan__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.pick-scan__header-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.pick-scan__qty {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
  padding: 0 4px;
}

.pick-scan__qty-label {
  font-size: 13px;
  color: var(--pda-text-secondary, #909399);
}

.pick-scan__qty-input {
  width: 120px;
}

.pick-scan__all-done {
  text-align: center;
  padding: 24px;
  color: var(--pda-success, #67c23a);
  font-weight: 500;
}

.pick-scan__task-card {
  margin-bottom: 8px;
}

.pick-scan__task-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.pick-scan__task-info {
  flex: 1;
  min-width: 0;
}

.pick-scan__task-name {
  font-weight: 500;
  font-size: 15px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.pick-scan__task-meta {
  font-size: 12px;
  color: var(--pda-text-secondary, #909399);
  margin-top: 2px;
}

.pick-scan__task-btn {
  flex-shrink: 0;
  min-height: var(--pda-touch-min, 44px);
}
</style>
