/**
 * Demo 模式检测
 *
 * 当 VITE_DEMO_MODE=true 时，前端不发起真实 API 请求。
 * 各页面/组件应据此渲染占位 UI，而非静默失败或显示假数据。
 *
 * 用法:
 *   import { isDemo } from "@erp/shared";
 *   if (isDemo()) { return; }
 */

export function isDemo(): boolean {
  return import.meta.env.VITE_DEMO_MODE === "true";
}
