import { createRouter, createWebHistory } from "vue-router";

const routes = [
  {
    path: "/login",
    name: "Login",
    component: () => import("@/views/Login.vue"),
    meta: { title: "登录" },
  },
  {
    path: "/",
    name: "Layout",
    component: () => import("@/layouts/MainLayout.vue"),
    redirect: "/dashboard",
    children: [
      {
        path: "dashboard",
        name: "Dashboard",
        component: () => import("@/views/Dashboard.vue"),
        meta: { title: "首页", icon: "DataAnalysis" },
      },
      {
        path: "system/users",
        name: "UserManagement",
        component: () => import("@/views/system/UserManagement.vue"),
        meta: { title: "用户管理", icon: "User" },
      },
      {
        path: "system/roles",
        name: "RoleManagement",
        component: () => import("@/views/system/RoleManagement.vue"),
        meta: { title: "角色管理", icon: "Avatar" },
      },
      {
        path: "system/org",
        name: "OrgManagement",
        component: () => import("@/views/system/OrgManagement.vue"),
        meta: { title: "组织管理", icon: "OfficeBuilding" },
      },
      {
        path: "system/audit",
        name: "AuditLog",
        component: () => import("@/views/system/AuditLog.vue"),
        meta: { title: "操作审计", icon: "DocumentChecked" },
      },
      {
        path: "product/list",
        name: "ProductList",
        component: () => import("@/views/product/ProductList.vue"),
        meta: { title: "商品列表", icon: "List" },
      },
      {
        path: "product/sku-mapping",
        name: "SKUMapping",
        component: () => import("@/views/product/SKUMapping.vue"),
        meta: { title: "SKU 映射", icon: "Connection" },
      },
      {
        path: "channel/stores",
        name: "StoreManagement",
        component: () => import("@/views/channel/StoreManagement.vue"),
        meta: { title: "店铺授权", icon: "Shop" },
      },
      {
        path: "channel/import",
        name: "OrderImport",
        component: () => import("@/views/channel/OrderImport.vue"),
        meta: { title: "订单导入", icon: "Upload" },
      },
      {
        path: "order/list",
        name: "OrderList",
        component: () => import("@/views/order/OrderList.vue"),
        meta: { title: "订单列表", icon: "Tickets" },
      },
      {
        path: "inventory",
        name: "InventoryList",
        component: () => import("@/views/inventory/InventoryList.vue"),
        meta: { title: "库存管理", icon: "Box" },
      },
    ],
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
