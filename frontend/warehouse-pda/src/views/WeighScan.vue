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
const weighWeight = ref<Record<string, number>>({});
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

const weighedOutbounds = computed(() =>
  store.outbounds.filter((o) => o.status === "packed")
);

function findOutbound(value: string) {
  const v = value.trim().toUpperCase();
  if (!v) return null;
  return (
    weighedOutbounds.value.find(
      (o) => o.order_no?.toUpperCase() === v || o.id?.toUpperCase() === v
    ) || null
  );
}

async function handleScan(value: string) {
  if (scanning.value) return;

  const ob = findOutbound(value);
  if (!ob) {
    ElMessage.warning(`未找到待称重出库单: ${value}`);
    vibrate([80, 100, 80, 100, 200]);
    scanInputRef.value?.focus();
    return;
  }

  const weight = weighWeight.value[ob.id];
  if (!weight || weight <= 0) {
    ElMessage.warning("请输入有效重量");
    scanInputRef.value?.focus();
    return;
  }

  scanning.value = true;
  try {
    await store.weigh(ob.id, weight);
    ElMessage.success(`✓ ${ob.order_no || ob.id} 称重完成`);
    vibrate(80);
    weighWeight.value[ob.id] = 0;
    await store.fetchOutbounds();
  } catch (e: unknown) {
    if (isDuplicateError(e)) {
      ElMessage.info(`已称重（重复操作）: ${ob.order_no || ob.id}`);
      vibrate([30, 50, 30]);
      await store.fetchOutbounds();
    } else if (isNetworkError(e)) {
      ElMessage.error("网络不可用");
    } else {
      ElMessage.error(getErrorMessage(e, "称重失败"));
    }
  } finally {
    scanning.value = false;
    scanInputRef.value?.focus();
  }
}

async function handleWeigh(outboundId: string) {
  if (scanning.value) return;
  const weight = weighWeight.value[outboundId];
  if (!weight || weight <= 0) {
    ElMessage.warning("请输入有效重量");
    return;
  }
  scanning.value = true;
  try {
    await store.weigh(outboundId, weight);
    ElMessage.success("称重完成");
    vibrate(80);
    weighWeight.value[outboundId] = 0;
    await store.fetchOutbounds();
  } catch (e: unknown) {
    if (isDuplicateError(e)) {
      ElMessage.info("已称重（重复操作）");
      vibrate([30, 50, 30]);
      await store.fetchOutbounds();
    } else {
      ElMessage.error(getErrorMessage(e, "称重失败"));
    }
  } finally {
    scanning.value = false;
  }
}
</script>

<template>
  <div class="weigh-scan">
    <div class="weigh-scan__header">
      <el-page-header @back="router.push('/')" content="称重" />
      <div class="weigh-scan__header-right">
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

    <el-empty v-if="weighedOutbounds.length === 0" description="暂无待称重出库单" />

    <el-card
      v-for="ob in weighedOutbounds"
      :key="ob.id"
      class="weigh-scan__card"
      :body-style="{ padding: '12px 16px' }"
    >
      <div class="weigh-scan__name">{{ ob.order_no }}</div>
      <div class="weigh-scan__id">{{ ob.id }}</div>
      <div class="weigh-scan__actions">
        <el-input-number
          v-model="weighWeight[ob.id]"
          :min="1"
          :step="1"
          placeholder="重量(g)"
          size="small"
        />
        <el-button
          type="warning"
          size="small"
          :loading="scanning"
          class="weigh-scan__btn"
          @click="handleWeigh(ob.id)"
        >
          确认称重
        </el-button>
      </div>
    </el-card>

    <ScanFeedback v-if="showHistory" type="weigh" @close="showHistory = false" />
  </div>
</template>

<style scoped>
.weigh-scan__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.weigh-scan__header-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.weigh-scan__card {
  margin-bottom: 8px;
}

.weigh-scan__name {
  font-weight: 500;
  font-size: 15px;
}

.weigh-scan__id {
  color: var(--pda-text-secondary, #909399);
  margin: 8px 0;
  font-size: 13px;
}

.weigh-scan__actions {
  display: flex;
  gap: 8px;
  align-items: center;
}

.weigh-scan__btn {
  min-height: var(--pda-touch-min, 44px);
}
</style>
