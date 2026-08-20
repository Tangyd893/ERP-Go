import { lazy, Suspense, type ReactNode } from "react";
import { Navigate, useLocation } from "react-router-dom";
import PdaLayout from "@/layouts/PdaLayout";
import Login from "@/pages/Login";
import { Skeleton } from "@erp/shared-react";
import { useAuthStore } from "@/stores/auth";

const Home = lazy(() => import("@/pages/Home"));
const PickList = lazy(() => import("@/pages/PickList"));
const PickScan = lazy(() => import("@/pages/PickScan"));
const CheckScan = lazy(() => import("@/pages/CheckScan"));
const ShipConfirm = lazy(() => import("@/pages/ShipConfirm"));
const PackScan = lazy(() => import("@/pages/PackScan"));
const WeighScan = lazy(() => import("@/pages/WeighScan"));
const Profile = lazy(() => import("@/pages/Profile"));

function RequireAuth({ children }: { children: ReactNode }) {
  const location = useLocation();
  const token = useAuthStore((s) => s.token);
  if (!token) {
    return <Navigate to="/login" state={{ from: location }} replace />;
  }
  return children;
}

function PublicOnly({ children }: { children: ReactNode }) {
  const token = useAuthStore((s) => s.token);
  if (token) {
    return <Navigate to="/" replace />;
  }
  return children;
}

function PageFallback() {
  return (
    <div className="space-y-3">
      <Skeleton className="h-8 w-40" />
      <Skeleton className="h-24 w-full" />
      <Skeleton className="h-24 w-full" />
    </div>
  );
}

function LazyPage({ children }: { children: ReactNode }) {
  return <Suspense fallback={<PageFallback />}>{children}</Suspense>;
}

export const appRoutes = [
  {
    path: "/login",
    element: (
      <PublicOnly>
        <Login />
      </PublicOnly>
    ),
  },
  {
    path: "/",
    element: (
      <RequireAuth>
        <PdaLayout />
      </RequireAuth>
    ),
    children: [
      {
        index: true,
        element: (
          <LazyPage>
            <Home />
          </LazyPage>
        ),
      },
      {
        path: "pick",
        element: (
          <LazyPage>
            <PickList />
          </LazyPage>
        ),
      },
      {
        path: "pick/scan",
        element: (
          <LazyPage>
            <PickScan />
          </LazyPage>
        ),
      },
      {
        path: "check",
        element: (
          <LazyPage>
            <CheckScan />
          </LazyPage>
        ),
      },
      {
        path: "ship",
        element: (
          <LazyPage>
            <ShipConfirm />
          </LazyPage>
        ),
      },
      {
        path: "pack",
        element: (
          <LazyPage>
            <PackScan />
          </LazyPage>
        ),
      },
      {
        path: "weigh",
        element: (
          <LazyPage>
            <WeighScan />
          </LazyPage>
        ),
      },
      {
        path: "profile",
        element: (
          <LazyPage>
            <Profile />
          </LazyPage>
        ),
      },
    ],
  },
];
