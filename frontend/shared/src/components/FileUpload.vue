<script setup lang="ts">
import { ref } from "vue";
import { UploadFilled } from "@element-plus/icons-vue";
import type { UploadInstance, UploadProps, UploadFile } from "element-plus";

interface Props {
  accept?: string;
  limit?: number;
  action?: string;
  disabled?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  accept: "*",
  limit: 1,
  action: "",
  disabled: false,
});

const emit = defineEmits<{
  (e: "success", file: UploadFile, fileList: UploadFile[]): void;
  (e: "error", error: Error, file: UploadFile): void;
  (e: "exceed", files: File[]): void;
  (e: "remove", file: UploadFile): void;
}>();

const uploadRef = ref<UploadInstance>();

const handleExceed: UploadProps["onExceed"] = (files) => {
  emit("exceed", files);
};

const handleSuccess: UploadProps["onSuccess"] = (response, file, fileList) => {
  emit("success", file, fileList);
};

const handleError: UploadProps["onError"] = (error, file) => {
  emit("error", error, file);
};
</script>

<template>
  <el-upload
    ref="uploadRef"
    drag
    :accept="accept"
    :limit="limit"
    :action="action"
    :disabled="disabled"
    :on-exceed="handleExceed"
    :on-success="handleSuccess"
    :on-error="handleError"
    :auto-upload="false"
  >
    <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
    <div class="el-upload__text">
      将文件拖到此处，或<em>点击上传</em>
    </div>
    <template #tip>
      <div class="el-upload__tip" v-if="$slots.tip">
        <slot name="tip" />
      </div>
    </template>
  </el-upload>
</template>
