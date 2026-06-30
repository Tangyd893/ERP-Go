<script setup lang="ts">
import { ref, watch, onMounted, nextTick } from "vue";

const props = withDefaults(
  defineProps<{
    placeholder?: string;
    hint?: string;
    disabled?: boolean;
    modelValue?: string;
    /** 扫描后自动清空 */
    clearAfterScan?: boolean;
  }>(),
  {
    placeholder: "扫描或输入…",
    hint: "",
    disabled: false,
    modelValue: "",
    clearAfterScan: true,
  }
);

const emit = defineEmits<{
  (e: "update:modelValue", value: string): void;
  (e: "scan", value: string): void;
}>();

const inputRef = ref<HTMLInputElement | null>(null);
const localValue = ref(props.modelValue);

watch(
  () => props.modelValue,
  (v) => {
    localValue.value = v;
  }
);

watch(localValue, (v) => {
  emit("update:modelValue", v);
});

/** 执行扫码操作 */
function triggerScan() {
  const value = localValue.value.trim();
  if (!value || props.disabled) return;
  emit("scan", value);
  if (props.clearAfterScan) {
    localValue.value = "";
  }
}

/** 外部调用：聚焦输入框 */
function focus() {
  nextTick(() => inputRef.value?.focus());
}

defineExpose({ focus });
</script>

<template>
  <div class="scan-input-card">
    <div v-if="hint" class="scan-input-card__hint">{{ hint }}</div>
    <el-input
      ref="inputRef"
      v-model="localValue"
      :placeholder="placeholder"
      size="large"
      clearable
      :disabled="disabled"
      class="scan-input-card__input"
      @keyup.enter="triggerScan"
    >
      <template #append>
        <el-button
          :disabled="!localValue.trim() || disabled"
          class="scan-input-card__btn"
          @click="triggerScan"
        >
          确认
        </el-button>
      </template>
    </el-input>
  </div>
</template>

<style scoped>
.scan-input-card {
  margin-bottom: 12px;
  padding: 16px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
}

.scan-input-card__hint {
  font-size: 13px;
  color: var(--pda-text-secondary, #909399);
  margin-bottom: 10px;
}

.scan-input-card__input :deep(.el-input__inner) {
  font-size: 16px;
  letter-spacing: 0.5px;
}

.scan-input-card__btn {
  min-height: var(--pda-touch-min, 44px);
}
</style>
