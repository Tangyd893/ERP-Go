import ProTable from "./components/ProTable.vue";
import ProForm from "./components/ProForm.vue";
import FileUpload from "./components/FileUpload.vue";
import ErrorState from "./components/ErrorState.vue";

/**
 * @erp/shared 共享组件库
 * 提供跨项目通用的类型、组件、工具函数和 API 客户端
 */

export { ProTable, ProForm, FileUpload, ErrorState };

export { default as apiClient } from "./api";
export { isDemo } from "./demo";
export { useApiState, useApiPage } from "./composables/useApiState";
export * from "./types";
