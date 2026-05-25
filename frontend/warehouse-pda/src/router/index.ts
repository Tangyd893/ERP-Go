import { createRouter, createWebHistory } from "vue-router";

const routes = [
  {
    path: "/login",
    name: "Login",
    component: () => import("@/views/Login.vue"),
    meta: { title: "PDA 登录" },
  },
  {
    path: "/",
    name: "Home",
    component: () => import("@/views/Home.vue"),
    meta: { title: "PDA 首页" },
  },
  {
    path: "/pick",
    name: "PickList",
    component: () => import("@/views/PickList.vue"),
    meta: { title: "拣货任务" },
  },
  {
    path: "/pick/scan",
    name: "PickScan",
    component: () => import("@/views/PickScan.vue"),
    meta: { title: "拣货扫码" },
  },
  {
    path: "/check",
    name: "CheckScan",
    component: () => import("@/views/CheckScan.vue"),
    meta: { title: "复核" },
  },
  {
    path: "/ship",
    name: "ShipConfirm",
    component: () => import("@/views/ShipConfirm.vue"),
    meta: { title: "出库确认" },
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
