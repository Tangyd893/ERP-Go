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

const router = useRouter();
const store = useWarehouseStore();

const scanInput = ref("");
const scanning = ref(false);
const scanInputRef = ref<HTMLInputElement | null>(null);

onMounted(async () => {
  await store.fetchOutbounds();
  scanInputRef.value?.focus();
});

const checkingOutbounds = computed(() =>
  store.outbounds.filter((o) => o.status === "picked" || o.status === "picking")
);

const recentChecks = computed(() =>
  store.scanHistory.filter((r) => r.type === "check").slice(0, 10)
);

function findOutbound(value: string) {
  const v = value.trim().toUpperCase();
  if (!v) return null;
  return (
    checkingOutbounds.value.find(
      (o) =>
        o.order_no?.toUpperCase() === v ||
        o.id?.toUpperCase() === v
    ) || null
  );
}

async function handleCheckClick(outboundId: string, label: string) {
  if (scanning.value) return;
  scanning.value = true;
  try {
    await store.checkScan(outboundId, "", 1, label);
    ElMessage.success(`✓ ${label} 复核完成`);
    await store.fetchOutbounds();
  } catch (e: unknown) {
    if (isDuplicateError(e)) {
      ElMessage.info(`已复核（重复扫描）: ${label}`);
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

async function handleScan() {
  const value = scanInput.value.trim();
  if (!value || scanning.value) return;

  const ob = findOutbound(value);
  if (!ob) {
    ElMessage.warning(`未找到匹配的待复核出库单: ${value}`);
    scanInput.value = "";
    scanInputRef.value?.focus();
    return;
  }

  scanning.value = true;
  const label = `${ob.order_no || ob.id}`;
  try {
    await store.checkScan(ob.id, "", 1, label);
    ElMessage.success(`✓ ${label} 复核完成`);
    scanInput.value = "";
    await store.fetchOutbounds();
  } catch (e: unknown) {
    if (isDuplicateError(e)) {
      ElMessage.info(`已复核（重复扫描）: ${label}`);
      scanInput.value = "";
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
  <div style="padding: 16px; max-width: 480px; margin: 0 auto">
    <!-- Header -->
    <div style="display: flex; align-items: center; justify-content: space-between; margin-bottom: 12px">
      <el-page-header @back="router.push('/')" content="复核扫码" />
      <el-tag :type="store.online ? 'success' : 'danger'" size="small" effect="dark">
        {{ store.online ? "在线" : "离线" }}
      </el-tag>
    </div>

    <!-- Scan Input -->
    <el-card style="margin-bottom: 12px">
      <div style="font-size: 13px; color: #909399; margin-bottom: 8px">
        扫描出库单条码 / 输入订单号
      </div>
      <el-input
        ref="scanInputRef"
        v-model="scanInput"
        placeholder="扫描或输入订单号…"
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
    </el-card>

    <!-- Pending Check Tasks -->
    <el-empty
      v-if="checkingOutbounds.length === 0"
      description="暂无待复核出库单"
    />

    <el-card
      v-for="ob in checkingOutbounds"
      :key="ob.id"
      style="margin-bottom: 8px"
      :body-style="{ padding: '12px 16px' }"
    >
      <div style="display: flex; justify-content: space-between; align-items: center">
        <div>
          <div style="font-weight: 500">{{ ob.order_no }}</div>
          <div style="font-size: 12px; color: #909399">{{ ob.id }}</div>
        </div>
        <el-button
          type="warning"
          size="small"
          :loading="scanning"
          @click="handleCheckClick(ob.id, ob.order_no || ob.id)"
        >
          复核确认
        </el-button>
      </div>
    </el-card>

    <!-- Recent Scan History -->
    <div v-if="recentChecks.length > 0" style="margin-top: 16px">
      <div style="font-size: 13px; color: #909399; margin-bottom: 8px">
        最近复核记录
      </div>
      <div
        v-for="r in recentChecks"
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
  </div>
</template>
