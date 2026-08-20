import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Lock, User } from "lucide-react";
import { toast } from "sonner";
import { Button, Card, CardContent, CardHeader, CardTitle, Input, Label } from "@erp/shared-react";
import { useAuthStore } from "@/stores/auth";

export default function Login() {
  const navigate = useNavigate();
  const login = useAuthStore((s) => s.login);
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  async function handleLogin(e: React.FormEvent) {
    e.preventDefault();
    setError("");
    if (!username || !password) {
      setError("请输入用户名和密码");
      return;
    }
    setLoading(true);
    try {
      await login(username, password, "default");
      toast.success("登录成功");
      navigate("/dashboard");
    } catch {
      setError("登录失败，请检查用户名和密码");
      toast.error("登录失败，请检查用户名和密码");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="grid min-h-dvh lg:grid-cols-[1.1fr_0.9fr]">
      <aside className="relative hidden overflow-hidden bg-sidebar px-10 py-12 text-sidebar-foreground lg:flex lg:flex-col lg:justify-between">
        <div className="relative space-y-3">
          <p className="text-sm font-medium text-sidebar-accent">ERP-Go</p>
          <h1 className="max-w-md text-3xl font-semibold leading-tight text-balance">
            跨境电商运营后台
          </h1>
          <p className="max-w-sm text-sm text-sidebar-muted text-pretty">
            订单、库存、仓储与财务在同一工作台完成审核与履约。
          </p>
        </div>
        <p className="relative text-xs text-sidebar-muted">React 并行版 · 端口 5273</p>
      </aside>

      <main className="flex items-center justify-center bg-background px-4 py-10">
        <Card className="w-full max-w-md border shadow-sm">
          <CardHeader className="space-y-1">
            <CardTitle className="text-xl font-semibold tracking-tight text-balance">
              登录管理后台
            </CardTitle>
            <p className="text-sm text-muted-foreground text-pretty">使用 IAM 账号进入系统</p>
          </CardHeader>
          <CardContent>
            <form className="space-y-4" onSubmit={handleLogin}>
              <div className="space-y-2">
                <Label htmlFor="username">用户名</Label>
                <div className="relative">
                  <User className="pointer-events-none absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
                  <Input
                    id="username"
                    className="pl-9"
                    placeholder="用户名"
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                    autoComplete="username"
                  />
                </div>
              </div>
              <div className="space-y-2">
                <Label htmlFor="password">密码</Label>
                <div className="relative">
                  <Lock className="pointer-events-none absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
                  <Input
                    id="password"
                    className="pl-9"
                    type="password"
                    placeholder="密码"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    autoComplete="current-password"
                  />
                </div>
              </div>
              <div className="space-y-2">
                <Label htmlFor="tenant">租户</Label>
                <Input id="tenant" value="default" disabled />
              </div>
              {error ? <p className="text-sm text-destructive">{error}</p> : null}
              <Button type="submit" className="w-full" disabled={loading}>
                {loading ? "登录中…" : "登 录"}
              </Button>
            </form>
          </CardContent>
        </Card>
      </main>
    </div>
  );
}
