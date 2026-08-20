import { useEffect, useMemo, useState } from "react";
import { Link, Outlet, useLocation, useNavigate } from "react-router-dom";
import {
  Bell,
  Box,
  CheckSquare,
  ChevronDown,
  ChevronLeft,
  ChevronRight,
  ClipboardList,
  DollarSign,
  FileText,
  Home,
  Link2,
  List,
  LogOut,
  Package,
  Search,
  Settings,
  ShoppingCart,
  Store,
  Truck,
  Upload,
  User,
  Users,
  Warehouse,
} from "lucide-react";
import {
  Button,
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  Input,
  cn,
  isDemo,
} from "@erp/shared-react";
import { useAuthStore } from "@/stores/auth";

type MenuLeaf = { path: string; title: string; icon: typeof Home };
type MenuGroup = { title: string; icon: typeof Home; children: MenuLeaf[] };
type MenuItem = MenuLeaf | MenuGroup;

const menuItems: MenuItem[] = [
  { path: "/dashboard", title: "首页", icon: Home },
  {
    title: "系统管理",
    icon: Settings,
    children: [
      { path: "/system/users", title: "用户管理", icon: User },
      { path: "/system/roles", title: "角色管理", icon: Users },
      { path: "/system/org", title: "组织管理", icon: Users },
      { path: "/system/audit", title: "操作审计", icon: FileText },
    ],
  },
  {
    title: "商品管理",
    icon: Package,
    children: [
      { path: "/product/list", title: "商品列表", icon: List },
      { path: "/product/sku-mapping", title: "SKU 映射", icon: Link2 },
    ],
  },
  {
    title: "渠道管理",
    icon: Link2,
    children: [
      { path: "/channel/stores", title: "店铺授权", icon: Store },
      { path: "/channel/import", title: "订单导入", icon: Upload },
    ],
  },
  {
    title: "订单管理",
    icon: ClipboardList,
    children: [
      { path: "/order/list", title: "订单列表", icon: FileText },
      { path: "/order/audit", title: "订单审核", icon: CheckSquare },
    ],
  },
  { path: "/inventory", title: "库存管理", icon: Box },
  { path: "/warehouse/outbound", title: "出库管理", icon: Warehouse },
  { path: "/transport", title: "物流管理", icon: Truck },
  { path: "/notification", title: "通知中心", icon: Bell },
  { path: "/purchase", title: "采购管理", icon: ShoppingCart },
  { path: "/finance", title: "财务管理", icon: DollarSign },
];

function flattenMenu(items: MenuItem[]): MenuLeaf[] {
  return items.flatMap((item) => ("children" in item ? item.children : [item]));
}

export default function MainLayout() {
  const location = useLocation();
  const navigate = useNavigate();
  const logout = useAuthStore((s) => s.logout);
  const [collapsed, setCollapsed] = useState(false);
  const [openGroups, setOpenGroups] = useState<Record<string, boolean>>({});
  const [paletteOpen, setPaletteOpen] = useState(false);
  const [paletteQuery, setPaletteQuery] = useState("");
  const [userMenuOpen, setUserMenuOpen] = useState(false);

  const flatItems = useMemo(() => flattenMenu(menuItems), []);
  const pageTitle =
    (location.pathname &&
      flatItems.find((item) => item.path === location.pathname)?.title) ||
    "ERP-Go";

  const paletteResults = useMemo(() => {
    const q = paletteQuery.trim().toLowerCase();
    if (!q) return flatItems;
    return flatItems.filter(
      (item) => item.title.toLowerCase().includes(q) || item.path.includes(q)
    );
  }, [flatItems, paletteQuery]);

  useEffect(() => {
    const onKeyDown = (e: KeyboardEvent) => {
      if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === "k") {
        e.preventDefault();
        setPaletteOpen(true);
      }
      if (e.key === "Escape") {
        setPaletteOpen(false);
        setUserMenuOpen(false);
      }
    };
    window.addEventListener("keydown", onKeyDown);
    return () => window.removeEventListener("keydown", onKeyDown);
  }, []);

  useEffect(() => {
    menuItems.forEach((item) => {
      if ("children" in item) {
        const active = item.children.some((child) => child.path === location.pathname);
        if (active) {
          setOpenGroups((prev) => ({ ...prev, [item.title]: true }));
        }
      }
    });
  }, [location.pathname]);

  const handleLogout = () => {
    logout();
    navigate("/login");
  };

  return (
    <div className="flex h-dvh bg-background">
      <aside
        className={cn(
          "flex shrink-0 flex-col border-r border-sidebar/20 bg-sidebar text-sidebar-foreground",
          collapsed ? "w-[68px]" : "w-[240px]"
        )}
      >
        <div className="flex h-14 items-center justify-between border-b border-white/10 px-4">
          {!collapsed ? (
            <span className="text-sm font-semibold tracking-tight">ERP-Go</span>
          ) : (
            <span className="mx-auto text-xs font-bold text-sidebar-accent">E</span>
          )}
          <Button
            variant="ghost"
            size="sm"
            className="size-8 p-0 text-sidebar-muted hover:bg-white/10 hover:text-sidebar-foreground"
            onClick={() => setCollapsed((v) => !v)}
            aria-label={collapsed ? "展开侧边栏" : "收起侧边栏"}
          >
            {collapsed ? <ChevronRight className="size-4" /> : <ChevronLeft className="size-4" />}
          </Button>
        </div>

        <nav className="flex-1 space-y-1 overflow-y-auto p-2">
          {menuItems.map((item) => {
            if ("children" in item) {
              const Icon = item.icon;
              const open = openGroups[item.title];
              return (
                <div key={item.title}>
                  <button
                    type="button"
                    className={cn(
                      "flex w-full items-center gap-2 rounded-md px-2 py-2 text-sm text-sidebar-muted transition hover:bg-white/5 hover:text-sidebar-foreground",
                      collapsed && "justify-center px-0"
                    )}
                    onClick={() =>
                      setOpenGroups((prev) => ({ ...prev, [item.title]: !prev[item.title] }))
                    }
                    title={item.title}
                  >
                    <Icon className="size-4 shrink-0" />
                    {!collapsed ? (
                      <>
                        <span className="flex-1 text-left">{item.title}</span>
                        <ChevronDown
                          className={cn("size-4 transition", open ? "rotate-180" : "")}
                        />
                      </>
                    ) : null}
                  </button>
                  {!collapsed && open ? (
                    <div className="ml-2 space-y-0.5 border-l border-white/10 pl-2">
                      {item.children.map((child) => {
                        const ChildIcon = child.icon;
                        const active = location.pathname === child.path;
                        return (
                          <Link
                            key={child.path}
                            to={child.path}
                            className={cn(
                              "flex items-center gap-2 rounded-md px-2 py-1.5 text-sm transition",
                              active
                                ? "bg-sidebar-accent/20 text-white"
                                : "text-sidebar-muted hover:bg-white/5 hover:text-sidebar-foreground"
                            )}
                          >
                            <ChildIcon className="size-4 shrink-0" />
                            {child.title}
                          </Link>
                        );
                      })}
                    </div>
                  ) : null}
                </div>
              );
            }

            const Icon = item.icon;
            const active = location.pathname === item.path;
            return (
              <Link
                key={item.path}
                to={item.path}
                title={item.title}
                className={cn(
                  "flex items-center gap-2 rounded-md px-2 py-2 text-sm transition",
                  active
                    ? "bg-sidebar-accent/20 text-white"
                    : "text-sidebar-muted hover:bg-white/5 hover:text-sidebar-foreground",
                  collapsed && "justify-center px-0"
                )}
              >
                <Icon className="size-4 shrink-0" />
                {!collapsed ? item.title : null}
              </Link>
            );
          })}
        </nav>
      </aside>

      <div className="flex min-w-0 flex-1 flex-col">
        <header className="flex h-14 items-center justify-between border-b bg-card px-4">
          <div className="flex min-w-0 items-center gap-2 text-sm text-muted-foreground">
            <Link to="/dashboard" className="hover:text-foreground">
              首页
            </Link>
            {pageTitle !== "首页" ? (
              <>
                <span>/</span>
                <span className="truncate font-medium text-foreground">{pageTitle}</span>
              </>
            ) : null}
          </div>

          <div className="flex items-center gap-2">
            {isDemo() ? (
              <span className="hidden rounded-md border border-amber-200 bg-amber-50 px-2 py-0.5 text-xs font-medium text-amber-900 sm:inline">
                DEMO
              </span>
            ) : null}
            <Button
              variant="outline"
              size="sm"
              className="hidden gap-2 sm:inline-flex"
              onClick={() => setPaletteOpen(true)}
            >
              <Search className="size-4" />
              <span className="text-muted-foreground">搜索</span>
              <kbd className="rounded border bg-muted px-1.5 py-0.5 text-[10px] font-medium">
                ⌘K
              </kbd>
            </Button>

            <div className="relative">
              <Button
                variant="ghost"
                size="sm"
                className="gap-1"
                onClick={() => setUserMenuOpen((v) => !v)}
              >
                管理员
                <ChevronDown className="size-4" />
              </Button>
              {userMenuOpen ? (
                <div className="absolute right-0 z-[var(--z-dropdown)] mt-1 w-40 rounded-md border bg-popover p-1 shadow-md">
                  <button
                    type="button"
                    className="w-full rounded px-2 py-1.5 text-left text-sm text-muted-foreground hover:bg-accent"
                    onClick={() => setUserMenuOpen(false)}
                  >
                    个人设置
                  </button>
                  <button
                    type="button"
                    className="flex w-full items-center gap-2 rounded px-2 py-1.5 text-left text-sm text-destructive hover:bg-destructive/10"
                    onClick={handleLogout}
                  >
                    <LogOut className="size-4" />
                    退出登录
                  </button>
                </div>
              ) : null}
            </div>
          </div>
        </header>

        <main className="flex-1 overflow-y-auto bg-muted/30 p-4 md:p-6">
          <div className="page-shell">
            <Outlet />
          </div>
        </main>
      </div>

      <Dialog open={paletteOpen} onOpenChange={setPaletteOpen}>
        <DialogContent className="max-w-md gap-0 overflow-hidden p-0">
          <DialogHeader className="border-b px-4 py-3">
            <DialogTitle className="text-base">快速导航</DialogTitle>
          </DialogHeader>
          <div className="border-b px-4 py-2">
            <Input
              autoFocus
              placeholder="搜索页面…"
              value={paletteQuery}
              onChange={(e) => setPaletteQuery(e.target.value)}
            />
          </div>
          <div className="max-h-72 overflow-y-auto p-2">
            {paletteResults.length === 0 ? (
              <p className="px-2 py-6 text-center text-sm text-muted-foreground">无匹配页面</p>
            ) : (
              paletteResults.map((item) => {
                const Icon = item.icon;
                return (
                  <button
                    key={item.path}
                    type="button"
                    className="flex w-full items-center gap-2 rounded-md px-2 py-2 text-left text-sm hover:bg-accent"
                    onClick={() => {
                      navigate(item.path);
                      setPaletteOpen(false);
                      setPaletteQuery("");
                    }}
                  >
                    <Icon className="size-4 text-muted-foreground" />
                    <span>{item.title}</span>
                    <span className="ml-auto text-xs text-muted-foreground">{item.path}</span>
                  </button>
                );
              })
            )}
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
}
