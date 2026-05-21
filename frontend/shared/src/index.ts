import type { App } from "vue";

/**
 * @erp/shared 共享组件库
 * 提供跨项目通用的类型、组件、工具函数和 API 客户端
 */

// 导出所有共享模块
export { default as apiClient } from "./api";
export * from "./types";

declare const _default: {
  install: (app: App) => void;
};

export default _default;
