<script setup lang="ts">
import { computed } from "vue";
import { useWarehouseStore } from "@/stores/warehouse";

const props = defineProps<{
  /** 过滤 scanHistory 的 type 字段 */
  type?: string;
  /** 最多显示条数 */
  limit?: number;
}>();

const store = useWarehouseStore();

const records = computed(() => {
  let list = store.scanHistory;
  if (props.type) {
    list = list.filter((r) => r.type === props.type);
  }
  return list.slice(0, props.limit ?? 50);
});

const emit = defineEmits<{
  (e: "close"): void;
}>();

/** 触发设备振动（如果支持） */
function vibrate(pattern: number | number[] = 50) {
  try {
    if (typeof navigator !== "undefined" && "vibrate" in navigator) {
      (navigator as Navigator).vibrate(pattern);
    }
  } catch {
    // 不支持振动的环境静默忽略
  }
}

/** 扫码结果反馈（成功/失败/重复） */
function scanFeedback(success: boolean, duplicate = false) {
  if (duplicate) {
    vibrate([30, 50, 30]);
  } else if (success) {
    vibrate(80);
  } else {
    vibrate([80, 100, 80, 100, 200]);
  }
}

defineExpose({ vibrate, scanFeedback });
</script>

<template>
  <el-drawer
    :model-value="true"
    direction="btt"
    size="60%"
    title="扫码历史"
    @close="emit('close')"
  >
    <div v-if="records.length === 0" style="text-align: center; padding: 24px; color: #909399">
      暂无扫码记录
    </div>
    <div
      v-for="r in records"
      :key="r.id"
      class="scan-record"
    >
      <span class="scan-record__icon" :class="r.success ? 'scan-record__icon--ok' : 'scan-record__icon--fail'">
        {{ r.success ? "✓" : "✗" }}
      </span>
      <span class="scan-record__label">{{ r.targetLabel }}</span>
      <span class="scan-record__time">{{ new Date(r.time).toLocaleTimeString() }}</span>
    </div>
  </el-drawer>
</template>

<style scoped>
.scan-record {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 0;
  border-bottom: 1px solid #f0f0f0;
  font-size: 13px;
}

.scan-record__icon {
  font-weight: bold;
  width: 24px;
  text-align: center;
  flex-shrink: 0;
}

.scan-record__icon--ok {
  color: var(--pda-success, #67c23a);
}

.scan-record__icon--fail {
  color: var(--pda-danger, #f56c6c);
}

.scan-record__label {
  flex: 1;
  margin: 0 8px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.scan-record__time {
  color: var(--pda-text-secondary, #c0c4cc);
  font-size: 11px;
  flex-shrink: 0;
}
</style>
