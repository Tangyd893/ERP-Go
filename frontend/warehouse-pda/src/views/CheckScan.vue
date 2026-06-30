<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import {
  useWarehouseStore,
  getErrorMessage,
  isDuplicateError,
  isNetworkError,
} from "@/stores/warehouse";
import ScanInput from "@/components/ScanInput.vue";
import ScanFeedback from "@/components/ScanFeedback.vue";

const router = useRouter();
const store = useWarehouseStore();

const scanValue = ref("");
const scanning = ref(false);
const scanInputRef = ref<InstanceType<typeof ScanInput> | null>(null);
const showHistory = ref(false);

function vibrate(pattern: number | number[] = 50) {
  try {
    if (typeof navigator !== "undefined" && "vibrate" in navigator) {
      (navigator as Navigator).vibrate(pattern);
    }
  } catch { /* 不支持静默忽略 */ }
}

onMounted(async () => {
  await store.fetchOutbounds();
  scanInputRef.value?.focus();
});

const checkingOutbounds = computed(() =>
  store.outbounds.filter((o) => o.status === "picked" || o.status === "picking")
);

function findOutbound(value: string) {
  const v = value.trim().toUpperCase();
  if (!v) return null;
  return (
    checkingOutbounds.value.find(
      (o) => o.order_no?.toUpperCase() === v || o.id?.toUpperCase() === v
    ) || null
  );
}

async function handleCheckClick(outboundId: string, label: string) {
  if (scanning.value) return;
  scanning.value = true;
  try {
    await store.checkScan(outboundId, "", 1, label);
    ElMessage.success(`✓ ${label} 复核完成`);
    vibrate(80);
    await store.fetchOutbounds();
  } catch (e: unknown) {
    if (isDuplicateError(e)) {
      ElMessage.info(`已复核（重复扫描）: ${label}`);
      vibrate([30, 50, 30]);
      await store.fetchOutbounds();
    } else if (isNetworkError(e)) {
      ElMessage.error("网络不可用，请检查连接后重试");
    } else {
      ElMessage.error(getErrorMessage(e, "复核失败"));
    }
  } finally {
    scanning.value = false;
  }
}

async function handleScan(value: string) {
  if (scanning.value) return;

  const ob = findOutbound(value);
  if (!ob) {
    ElMessage.warning(`未找到匹配的待复核出库单: ${value}`);
    vibrate([80, 100, 80, 100, 200]);
    scanInputRef.value?.focus();
    return;
  }

  scanning.value = true;
  const label = `${ob.order_no || ob.id}`;
  try {
    await store.checkScan(ob.id, "", 1, label);
    ElMessage.success(`✓ ${label} 复核完成`);
    vibrate(80);
    await store.fetchOutbounds();
  } catch (e: unknown) {
    if (isDuplicateError(e)) {
      ElMessage.info(`已复核（重复扫描）: ${label}`);
      vibrate([30, 50, 30]);
      await store.fetchOutbounds();
    } else if (isNetworkError(e)) {
      ElMessage.error("网络不可用，请检查连接后重试");
    } else {
      ElMessage.error(getErrorMessage(e, "复核失败"));
    }
  } finally {
    scanning.value = false;
    scanInputRef.value?.focus();
  }
}
</script>

<template>
  <div class="check-scan">
    <!-- Header -->
    <div class="check-scan__header">
      <el-page-header @back="router.push('/')" content="复核扫码" />
      <div class="check-scan__header-right">
        <el-badge :value="store.scanHistory.filter(r => r.type === 'check').length" :hidden="store.scanHistory.length === 0">
          <el-button size="small" text @click="showHistory = true">记录</el-button>
        </el-badge>
        <el-tag :type="store.online ? 'success' : 'danger'" size="small" effect="dark">
          {{ store.online ? "在线" : "离线" }}
        </el-tag>
      </div>
    </div>

    <!-- Scan Input -->
    <ScanInput
      ref="scanInputRef"
      v-model="scanValue"
      hint="扫描出库单条码 / 输入订单号"
      :disabled="scanning"
      @scan="handleScan"
    />

    <!-- Empty -->
    <el-empty v-if="checkingOutbounds.length === 0" description="暂无待复核出库单" />

    <!-- Task Cards -->
    <el-card
      v-for="ob in checkingOutbounds"
      :key="ob.id"
      class="check-scan__card"
      :body-style="{ padding: '12px 16px' }"
    >
      <div class="check-scan__row">
        <div class="check-scan__info">
          <div class="check-scan__name">{{ ob.order_no }}</div>
          <div class="check-scan__id">{{ ob.id }}</div>
        </div>
        <el-button
          type="warning"
          size="small"
          :loading="scanning"
          class="check-scan__btn"
          @click="handleCheckClick(ob.id, ob.order_no || ob.id)"
        >
          复核确认
        </el-button>
      </div>
    </el-card>

    <!-- History Drawer -->
    <ScanFeedback
      v-if="showHistory"
      type="check"
      @close="showHistory = false"
    />
  </div>
</template>

<style scoped>
.check-scan__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.check-scan__header-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.check-scan__card {
  margin-bottom: 8px;
}

.check-scan__row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.check-scan__info {
  flex: 1;
  min-width: 0;
}

.check-scan__name {
  font-weight: 500;
  font-size: 15px;
}

.check-scan__id {
  font-size: 12px;
  color: var(--pda-text-secondary, #909399);
  margin-top: 2px;
}

.check-scan__btn {
  flex-shrink: 0;
  min-height: var(--pda-touch-min, 44px);
}
</style>
