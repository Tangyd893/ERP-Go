<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import { useWarehouseStore, isDuplicateError, isNetworkError, getErrorMessage } from "@/stores/warehouse";

const router = useRouter();
const store = useWarehouseStore();
const weighWeight = ref<Record<string, number>>({});
const scanning = ref(false);
const scanInput = ref("");
const scanInputRef = ref<HTMLInputElement | null>(null);

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
  return weighedOutbounds.value.find(
    (o) => o.order_no?.toUpperCase() === v || o.id?.toUpperCase() === v
  ) || null;
}

async function handleScan() {
  const value = scanInput.value.trim();
  if (!value || scanning.value) return;

  const ob = findOutbound(value);
  if (!ob) {
    ElMessage.warning(`未找到待称重出库单: ${value}`);
    scanInput.value = "";
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
    weighWeight.value[ob.id] = 0;
    scanInput.value = "";
    await store.fetchOutbounds();
  } catch (e: unknown) {
    if (isDuplicateError(e)) {
      ElMessage.info(`已称重（重复操作）: ${ob.order_no || ob.id}`);
      scanInput.value = "";
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
    weighWeight.value[outboundId] = 0;
    await store.fetchOutbounds();
  } catch (e: unknown) {
    if (isDuplicateError(e)) {
      ElMessage.info("已称重（重复操作）");
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
  <div style="padding: 16px; max-width: 480px; margin: 0 auto">
    <div style="display: flex; align-items: center; justify-content: space-between; margin-bottom: 12px">
      <el-page-header @back="router.push('/')" content="称重" />
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
        placeholder="扫描或输入…"
        size="large"
        clearable
        :disabled="scanning"
        @keyup.enter="handleScan"
      >
        <template #append>
          <el-button :loading="scanning" :disabled="!scanInput.trim()" @click="handleScan">
            确认
          </el-button>
        </template>
      </el-input>
    </el-card>

    <el-empty v-if="weighedOutbounds.length === 0" description="暂无待称重出库单" />

    <el-card
      v-for="ob in weighedOutbounds"
      :key="ob.id"
      style="margin-bottom: 8px"
      :body-style="{ padding: '12px 16px' }"
    >
      <div style="font-weight: 500">{{ ob.order_no }}</div>
      <div style="color: #909399; margin: 8px 0; font-size: 13px">{{ ob.id }}</div>
      <div style="display: flex; gap: 8px; align-items: center">
        <el-input-number
          v-model="weighWeight[ob.id]"
          :min="1"
          :step="1"
          placeholder="重量(g)"
          size="small"
        />
        <el-button type="warning" size="small" :loading="scanning" @click="handleWeigh(ob.id)">
          确认称重
        </el-button>
      </div>
    </el-card>
  </div>
</template>
