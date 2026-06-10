<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import { useWarehouseStore, isDuplicateError, isNetworkError, getErrorMessage } from "@/stores/warehouse";

const router = useRouter();
const store = useWarehouseStore();
const packingWeight = ref<Record<string, number>>({});
const scanning = ref(false);
const scanInput = ref("");
const scanInputRef = ref<HTMLInputElement | null>(null);

onMounted(() => {
  store.fetchOutbounds();
  scanInputRef.value?.focus();
});

const packedOutbounds = computed(() =>
  store.outbounds.filter((o) => o.status === "checked")
);

const recentPacks = computed(() =>
  store.scanHistory.filter((r) => r.type === "check").slice(0, 10)
);

function findOutbound(value: string) {
  const v = value.trim().toUpperCase();
  if (!v) return null;
  return packedOutbounds.value.find(
    (o) => o.order_no?.toUpperCase() === v || o.id?.toUpperCase() === v
  ) || null;
}

async function handleScan() {
  const value = scanInput.value.trim();
  if (!value || scanning.value) return;

  const ob = findOutbound(value);
  if (!ob) {
    ElMessage.warning(`未找到待打包出库单: ${value}`);
    scanInput.value = "";
    scanInputRef.value?.focus();
    return;
  }

  scanning.value = true;
  const weight = packingWeight.value[ob.id] || 0;
  try {
    await store.pack(ob.id, weight);
    ElMessage.success(`✓ ${ob.order_no || ob.id} 打包完成`);
    packingWeight.value[ob.id] = 0;
    scanInput.value = "";
    await store.fetchOutbounds();
  } catch (e: unknown) {
    if (isDuplicateError(e)) {
      ElMessage.info(`已打包（重复操作）: ${ob.order_no || ob.id}`);
      scanInput.value = "";
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
    packingWeight.value[outboundId] = 0;
    await store.fetchOutbounds();
  } catch (e: unknown) {
    if (isDuplicateError(e)) {
      ElMessage.info("已打包（重复操作）");
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
  <div style="padding: 16px; max-width: 480px; margin: 0 auto">
    <div style="display: flex; align-items: center; justify-content: space-between; margin-bottom: 12px">
      <el-page-header @back="router.push('/')" content="打包" />
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

    <el-empty v-if="packedOutbounds.length === 0" description="暂无待打包出库单" />

    <el-card
      v-for="ob in packedOutbounds"
      :key="ob.id"
      style="margin-bottom: 8px"
      :body-style="{ padding: '12px 16px' }"
    >
      <div style="font-weight: 500">{{ ob.order_no }}</div>
      <div style="color: #909399; margin: 8px 0; font-size: 13px">{{ ob.id }}</div>
      <div style="display: flex; gap: 8px; align-items: center">
        <el-input-number
          v-model="packingWeight[ob.id]"
          :min="0"
          :step="1"
          placeholder="重量(g)"
          size="small"
        />
        <el-button type="primary" size="small" :loading="scanning" @click="handlePack(ob.id)">
          确认打包
        </el-button>
      </div>
    </el-card>

    <!-- Recent History -->
    <div v-if="recentPacks.length > 0" style="margin-top: 16px">
      <div style="font-size: 13px; color: #909399; margin-bottom: 8px">最近打包记录</div>
      <div
        v-for="r in recentPacks"
        :key="r.id"
        style="font-size: 12px; padding: 4px 0; border-bottom: 1px solid #f0f0f0; display: flex; justify-content: space-between"
      >
        <span>
          <span :style="{ color: r.success ? '#67c23a' : '#f56c6c' }">{{ r.success ? "✓" : "✗" }}</span>
          {{ r.targetLabel }} {{ r.quantity }}g
        </span>
        <span style="color: #c0c4cc">{{ new Date(r.time).toLocaleTimeString() }}</span>
      </div>
    </div>
  </div>
</template>
