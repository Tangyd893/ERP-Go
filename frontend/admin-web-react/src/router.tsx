import { lazy, Suspense, type ReactNode } from "react";
import { Navigate, useLocation } from "react-router-dom";
import MainLayout from "@/layouts/MainLayout";
import Login from "@/pages/Login";
import { Skeleton } from "@erp/shared-react";

const Dashboard = lazy(() => import("@/pages/Dashboard"));
const UserManagement = lazy(() => import("@/pages/system/UserManagement"));
const RoleManagement = lazy(() => import("@/pages/system/RoleManagement"));
const OrgManagement = lazy(() => import("@/pages/system/OrgManagement"));
const AuditLog = lazy(() => import("@/pages/system/AuditLog"));
const ProductList = lazy(() => import("@/pages/product/ProductList"));
const SKUMapping = lazy(() => import("@/pages/product/SKUMapping"));
const StoreManagement = lazy(() => import("@/pages/channel/StoreManagement"));
const OrderImport = lazy(() => import("@/pages/channel/OrderImport"));
const OrderList = lazy(() => import("@/pages/order/OrderList"));
const OrderAudit = lazy(() => import("@/pages/order/OrderAudit"));
const InventoryList = lazy(() => import("@/pages/inventory/InventoryList"));
const WarehouseOutbound = lazy(() => import("@/pages/warehouse/WarehouseOutbound"));
const TransportManagement = lazy(() => import("@/pages/transport/TransportManagement"));
const NotificationCenter = lazy(() => import("@/pages/notification/NotificationCenter"));
const PurchaseManagement = lazy(() => import("@/pages/purchase/PurchaseManagement"));
const FinanceManagement = lazy(() => import("@/pages/finance/FinanceManagement"));

function RequireAuth({ children }: { children: ReactNode }) {
  const location = useLocation();
  const token = localStorage.getItem("access_token")?.trim();
  if (!token) {
    return <Navigate to="/login" state={{ from: location }} replace />;
  }
  return children;
}

function PageFallback() {
  return (
    <div className="space-y-3 p-1">
      <Skeleton className="h-8 w-48" />
      <Skeleton className="h-64 w-full" />
    </div>
  );
}

function LazyPage({ children }: { children: ReactNode }) {
  return <Suspense fallback={<PageFallback />}>{children}</Suspense>;
}

export const appRoutes = [
  { path: "/login", element: <Login /> },
  {
    path: "/",
    element: (
      <RequireAuth>
        <MainLayout />
      </RequireAuth>
    ),
    children: [
      { index: true, element: <Navigate to="/dashboard" replace /> },
      {
        path: "dashboard",
        element: (
          <LazyPage>
            <Dashboard />
          </LazyPage>
        ),
        handle: { title: "首页" },
      },
      {
        path: "system/users",
        element: (
          <LazyPage>
            <UserManagement />
          </LazyPage>
        ),
        handle: { title: "用户管理" },
      },
      {
        path: "system/roles",
        element: (
          <LazyPage>
            <RoleManagement />
          </LazyPage>
        ),
        handle: { title: "角色管理" },
      },
      {
        path: "system/org",
        element: (
          <LazyPage>
            <OrgManagement />
          </LazyPage>
        ),
        handle: { title: "组织管理" },
      },
      {
        path: "system/audit",
        element: (
          <LazyPage>
            <AuditLog />
          </LazyPage>
        ),
        handle: { title: "操作审计" },
      },
      {
        path: "product/list",
        element: (
          <LazyPage>
            <ProductList />
          </LazyPage>
        ),
        handle: { title: "商品列表" },
      },
      {
        path: "product/sku-mapping",
        element: (
          <LazyPage>
            <SKUMapping />
          </LazyPage>
        ),
        handle: { title: "SKU 映射" },
      },
      {
        path: "channel/stores",
        element: (
          <LazyPage>
            <StoreManagement />
          </LazyPage>
        ),
        handle: { title: "店铺授权" },
      },
      {
        path: "channel/import",
        element: (
          <LazyPage>
            <OrderImport />
          </LazyPage>
        ),
        handle: { title: "订单导入" },
      },
      {
        path: "order/list",
        element: (
          <LazyPage>
            <OrderList />
          </LazyPage>
        ),
        handle: { title: "订单列表" },
      },
      {
        path: "order/audit",
        element: (
          <LazyPage>
            <OrderAudit />
          </LazyPage>
        ),
        handle: { title: "订单审核" },
      },
      {
        path: "inventory",
        element: (
          <LazyPage>
            <InventoryList />
          </LazyPage>
        ),
        handle: { title: "库存管理" },
      },
      {
        path: "warehouse/outbound",
        element: (
          <LazyPage>
            <WarehouseOutbound />
          </LazyPage>
        ),
        handle: { title: "出库管理" },
      },
      {
        path: "transport",
        element: (
          <LazyPage>
            <TransportManagement />
          </LazyPage>
        ),
        handle: { title: "物流管理" },
      },
      {
        path: "notification",
        element: (
          <LazyPage>
            <NotificationCenter />
          </LazyPage>
        ),
        handle: { title: "通知中心" },
      },
      {
        path: "purchase",
        element: (
          <LazyPage>
            <PurchaseManagement />
          </LazyPage>
        ),
        handle: { title: "采购管理" },
      },
      {
        path: "finance",
        element: (
          <LazyPage>
            <FinanceManagement />
          </LazyPage>
        ),
        handle: { title: "财务管理" },
      },
    ],
  },
];
