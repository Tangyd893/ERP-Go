import { createRouter, createWebHistory } from "vue-router";
import PdaLayout from "@/layouts/PdaLayout.vue";

const routes = [
  {
    path: "/login",
    name: "Login",
    component: () => import("@/views/Login.vue"),
    meta: { title: "PDA 登录", public: true },
  },
  {
    path: "/",
    component: PdaLayout,
    children: [
      {
        path: "",
        name: "Home",
        component: () => import("@/views/Home.vue"),
        meta: { title: "PDA 首页" },
      },
      {
        path: "pick",
        name: "PickList",
        component: () => import("@/views/PickList.vue"),
        meta: { title: "拣货任务" },
      },
      {
        path: "pick/scan",
        name: "PickScan",
        component: () => import("@/views/PickScan.vue"),
        meta: { title: "拣货扫码" },
      },
      {
        path: "check",
        name: "CheckScan",
        component: () => import("@/views/CheckScan.vue"),
        meta: { title: "复核" },
      },
      {
        path: "ship",
        name: "ShipConfirm",
        component: () => import("@/views/ShipConfirm.vue"),
        meta: { title: "出库确认" },
      },
      {
        path: "pack",
        name: "PackScan",
        component: () => import("@/views/PackScan.vue"),
        meta: { title: "打包" },
      },
      {
        path: "weigh",
        name: "WeighScan",
        component: () => import("@/views/WeighScan.vue"),
        meta: { title: "称重" },
      },
      {
        path: "profile",
        name: "Profile",
        component: () => import("@/views/Profile.vue"),
        meta: { title: "我的" },
      },
    ],
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach((to, _from, next) => {
  const token = localStorage.getItem("access_token");
  if (!to.meta.public && !token) {
    next({ path: "/login", query: { redirect: to.fullPath } });
    return;
  }
  if (to.path === "/login" && token) {
    next("/");
    return;
  }
  next();
});

export default router;
