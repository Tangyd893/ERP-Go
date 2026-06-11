import { createRouter, createWebHashHistory } from "vue-router";

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: "/",
      name: "Dashboard",
      component: () => import("@/views/DashboardView.vue"),
    },
  ],
});

export default router;
