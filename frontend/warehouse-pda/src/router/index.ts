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
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
