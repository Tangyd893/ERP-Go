<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import {
  useWarehouseStore,
  isDuplicateError,
  isNetworkError,
  getErrorMessage,
} from "@/stores/warehouse";
import ScanInput from "@/components/ScanInput.vue";
import ScanFeedback from "@/components/ScanFeedback.vue";

const router = useRouter();
const store = useWarehouseStore();
const packingWeight = ref<Record<string, number>>({});
const scanning = ref(false);
const scanValue = ref("");
const scanInputRef = ref<InstanceType<typeof ScanInput> | null>(null);
const showHistory = ref(false);

function vibrate(pattern: number | number[] = 50) {
  try {
    if (typeof navigator !== "undefined" && "vibrate" in navigator) {
      (navigator as Navigator).vibrate(pattern);
    }
  } catch { /* 不支持静默忽略 */ }
}

onMounted(() => {
  store.fetchOutbounds();
  scanInputRef.value?.focus();
});

const packedOutbounds = computed(() =>
  store.outbounds.filter((o) => o.status === "checked")
);

function findOutbound(value: string) {
  const v = value.trim().toUpperCase();
  if (!v) return null;
  return (
    packedOutbounds.value.find(
      (o) => o.order_no?.toUpperCase() === v || o.id?.toUpperCase() === v
    ) || null
  );
}

async function handleScan(value: string) {
  if (scanning.value) return;

  const ob = findOutbound(value);
  if (!ob) {
    ElMessage.warning(`未找到待打包出库单: ${value}`);
    vibrate([80, 100, 80, 100, 200]);
    scanInputRef.value?.focus();
    return;
  }

  scanning.value = true;
  const weight = packingWeight.value[ob.id] || 0;
  try {
    await store.pack(ob.id, weight);
    ElMessage.success(`✓ ${ob.order_no || ob.id} 打包完成`);
    vibrate(80);
    packingWeight.value[ob.id] = 0;
    await store.fetchOutbounds();
  } catch (e: unknown) {
    if (isDuplicateError(e)) {
      ElMessage.info(`已打包（重复操作）: ${ob.order_no || ob.id}`);
      vibrate([30, 50, 30]);
      await store.fetchOutbounds();
    } else if (isNetworkError(e)) {
      ElMessage.error("网络不可用");
    } else {
      ElMessage.error(getErrorMessage(e, "打包失败"));
    }
  } finally {
    scanning.value = false;
    scanInputRef.value?.focus();
  }
}

async function handlePack(outboundId: string) {
  if (scanning.value) return;
  scanning.value = true;
  try {
    const weight = packingWeight.value[outboundId] || 0;
    await store.pack(outboundId, weight);
    ElMessage.success("打包完成");
    vibrate(80);
    packingWeight.value[outboundId] = 0;
    await store.fetchOutbounds();
  } catch (e: unknown) {
    if (isDuplicateError(e)) {
      ElMessage.info("已打包（重复操作）");
      vibrate([30, 50, 30]);
      await store.fetchOutbounds();
    } else {
      ElMessage.error(getErrorMessage(e, "打包失败"));
    }
  } finally {
    scanning.value = false;
  }
}
</script>

<template>
  <div class="pack-scan">
    <div class="pack-scan__header">
      <el-page-header @back="router.push('/')" content="打包" />
      <div class="pack-scan__header-right">
        <el-button size="small" text @click="showHistory = true">记录</el-button>
        <el-tag :type="store.online ? 'success' : 'danger'" size="small" effect="dark">
          {{ store.online ? "在线" : "离线" }}
        </el-tag>
      </div>
    </div>

    <ScanInput
      ref="scanInputRef"
      v-model="scanValue"
      hint="扫描出库单条码 / 输入订单号"
      :disabled="scanning"
      @scan="handleScan"
    />

    <el-empty v-if="packedOutbounds.length === 0" description="暂无待打包出库单" />

    <el-card
      v-for="ob in packedOutbounds"
      :key="ob.id"
      class="pack-scan__card"
      :body-style="{ padding: '12px 16px' }"
    >
      <div class="pack-scan__name">{{ ob.order_no }}</div>
      <div class="pack-scan__id">{{ ob.id }}</div>
      <div class="pack-scan__actions">
        <el-input-number
          v-model="packingWeight[ob.id]"
          :min="0"
          :step="1"
          placeholder="重量(g)"
          size="small"
        />
        <el-button
          type="primary"
          size="small"
          :loading="scanning"
          class="pack-scan__btn"
          @click="handlePack(ob.id)"
        >
          确认打包
        </el-button>
      </div>
    </el-card>

    <ScanFeedback v-if="showHistory" type="pack" @close="showHistory = false" />
  </div>
</template>

<style scoped>
.pack-scan__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.pack-scan__header-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.pack-scan__card {
  margin-bottom: 8px;
}

.pack-scan__name {
  font-weight: 500;
  font-size: 15px;
}

.pack-scan__id {
  color: var(--pda-text-secondary, #909399);
  margin: 8px 0;
  font-size: 13px;
}

.pack-scan__actions {
  display: flex;
  gap: 8px;
  align-items: center;
}

.pack-scan__btn {
  min-height: var(--pda-touch-min, 44px);
}
</style>
