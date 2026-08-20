import { Outlet, useLocation, useNavigate } from "react-router-dom";
import { Box, CheckSquare, Home, List, User } from "lucide-react";
import { useAuthStore } from "@/stores/auth";
import { cn } from "@erp/shared-react";

const tabs = [
  { path: "/", label: "首页", icon: Home },
  { path: "/pick", label: "拣货", icon: List },
  { path: "/check", label: "复核", icon: CheckSquare },
  { path: "/pack", label: "打包", icon: Box },
  { path: "/profile", label: "我的", icon: User },
] as const;

function resolveActiveTab(pathname: string): string {
  for (const tab of tabs) {
    if (tab.path === "/") {
      if (pathname === "/") return "/";
    } else if (pathname.startsWith(tab.path)) {
      return tab.path;
    }
  }
  return "/";
}

export default function PdaLayout() {
  const navigate = useNavigate();
  const location = useLocation();
  const username = useAuthStore((s) => s.username);
  const tenantId = useAuthStore((s) => s.tenantId);
  const logout = useAuthStore((s) => s.logout);

  const activeTab = resolveActiveTab(location.pathname);

  function handleLogout() {
    logout();
    navigate("/login");
  }

  return (
    <div className="flex h-[100dvh] flex-col bg-pda-bg">
      <header className="pda-safe-top z-[var(--z-sticky,40)] flex h-12 shrink-0 items-center justify-between border-b border-border bg-white px-3">
        <span className="text-base font-semibold text-pda-text">WMS PDA</span>
        <div className="flex items-center gap-1">
          <span className="max-w-[120px] truncate text-xs text-pda-text-secondary">
            {username || tenantId || "未登录"}
          </span>
          <button
            type="button"
            className="min-h-touch px-2 text-xs text-primary"
            onClick={handleLogout}
          >
            退出
          </button>
        </div>
      </header>

      <main className="flex-1 overflow-y-auto overscroll-contain p-3 [-webkit-overflow-scrolling:touch]">
        <Outlet />
      </main>

      <footer className="pda-safe-bottom z-[var(--z-sticky,40)] flex h-14 shrink-0 border-t border-border bg-white">
        {tabs.map((tab) => {
          const Icon = tab.icon;
          const active = activeTab === tab.path;
          return (
            <button
              key={tab.path}
              type="button"
              className={cn(
                "flex min-h-touch flex-1 flex-col items-center justify-center gap-0.5 transition-colors",
                active ? "text-pda-primary" : "text-pda-text-secondary"
              )}
              onClick={() => {
                if (location.pathname !== tab.path) {
                  navigate(tab.path);
                }
              }}
            >
              <Icon className="h-5 w-5" strokeWidth={active ? 2.25 : 2} />
              <span className="text-[10px] leading-none">{tab.label}</span>
            </button>
          );
        })}
      </footer>
    </div>
  );
}
