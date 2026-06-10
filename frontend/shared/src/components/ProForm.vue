<script setup lang="ts">
import { reactive, ref, watch } from "vue";

interface FormField {
  label: string;
  prop: string;
  type: "input" | "textarea" | "number" | "select" | "date" | "switch" | "radio";
  placeholder?: string;
  required?: boolean;
  rules?: Record<string, unknown>[];
  options?: { label: string; value: string | number }[];
  disabled?: boolean;
}

interface Props {
  modelValue: boolean;
  title: string;
  formData: Record<string, unknown>;
  fields: FormField[];
  submitText?: string;
  width?: string;
}

const props = withDefaults(defineProps<Props>(), {
  submitText: "确定",
  width: "520px",
});

const emit = defineEmits<{
  (e: "update:modelValue", value: boolean): void;
  (e: "submit", data: Record<string, unknown>): void;
}>();

const formRef = ref();
const visible = ref(props.modelValue);
const localForm = reactive({ ...props.formData });

watch(
  () => props.modelValue,
  (val) => {
    visible.value = val;
    if (val) {
      Object.assign(localForm, props.formData);
    }
  }
);

watch(
  () => props.formData,
  (val) => {
    Object.assign(localForm, val);
  }
);

function handleClose() {
  emit("update:modelValue", false);
}

function handleSubmit() {
  formRef.value?.validate((valid: boolean) => {
    if (valid) {
      emit("submit", { ...localForm });
    }
  });
}

</script>

<template>
  <el-dialog
    :model-value="visible"
    :title="title"
    :width="width"
    :close-on-click-modal="false"
    @update:model-value="(v: boolean) => emit('update:modelValue', v)"
  >
    <el-form ref="formRef" :model="localForm" label-width="100px">
      <el-form-item
        v-for="field in fields"
        :key="field.prop"
        :label="field.label"
        :prop="field.prop"
        :required="field.required"
        :rules="field.rules"
      >
        <el-input
          v-if="field.type === 'input'"
          v-model="localForm[field.prop]"
          :placeholder="field.placeholder"
          :disabled="field.disabled"
        />
        <el-input
          v-else-if="field.type === 'textarea'"
          v-model="localForm[field.prop]"
          type="textarea"
          :placeholder="field.placeholder"
          :disabled="field.disabled"
        />
        <el-input-number
          v-else-if="field.type === 'number'"
          v-model="localForm[field.prop]"
          :placeholder="field.placeholder"
          :disabled="field.disabled"
        />
        <el-select
          v-else-if="field.type === 'select'"
          v-model="localForm[field.prop]"
          :placeholder="field.placeholder"
          :disabled="field.disabled"
        >
          <el-option
            v-for="opt in field.options"
            :key="opt.value"
            :label="opt.label"
            :value="opt.value"
          />
        </el-select>
        <el-date-picker
          v-else-if="field.type === 'date'"
          v-model="localForm[field.prop]"
          type="date"
          :placeholder="field.placeholder"
          :disabled="field.disabled"
        />
        <el-switch
          v-else-if="field.type === 'switch'"
          v-model="localForm[field.prop]"
          :disabled="field.disabled"
        />
        <el-radio-group
          v-else-if="field.type === 'radio'"
          v-model="localForm[field.prop]"
          :disabled="field.disabled"
        >
          <el-radio
            v-for="opt in field.options"
            :key="opt.value"
            :value="opt.value"
          >
            {{ opt.label }}
          </el-radio>
        </el-radio-group>
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" @click="handleSubmit">{{ submitText }}</el-button>
    </template>
  </el-dialog>
</template>
