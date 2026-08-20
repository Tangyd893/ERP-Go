import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { toast } from "sonner";
import { Button, Input, isDemo } from "@erp/shared-react";
import { useAuthStore } from "@/stores/auth";

export default function Login() {
  const navigate = useNavigate();
  const login = useAuthStore((s) => s.login);
  const [username, setUsername] = useState("admin");
  const [password, setPassword] = useState("admin123");
  const [loading, setLoading] = useState(false);

  async function handleLogin() {
    if (!username.trim() || !password) {
      toast.warning("请输入工号和密码");
      return;
    }
    setLoading(true);
    try {
      await login(username.trim(), password, "default");
      toast.success("登录成功");
      navigate("/");
    } catch {
      toast.error("登录失败，请检查账号密码或服务是否启动");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="mx-auto max-w-sm px-6 pb-12 pt-12">
      <h1 className="mb-8 text-center text-2xl font-semibold text-pda-text">WMS PDA</h1>
      <form
        className="space-y-4"
        onSubmit={(e) => {
          e.preventDefault();
          void handleLogin();
        }}
      >
        <Input
          value={username}
          placeholder="工号"
          className="min-h-touch text-base"
          autoComplete="username"
          onChange={(e) => setUsername(e.target.value)}
        />
        <Input
          type="password"
          value={password}
          placeholder="密码"
          className="min-h-touch text-base"
          autoComplete="current-password"
          onChange={(e) => setPassword(e.target.value)}
        />
        <Button type="submit" className="min-h-touch w-full text-base" disabled={loading}>
          {loading ? "登录中…" : "登 录"}
        </Button>
      </form>
      <p className="mt-4 text-center text-xs text-pda-text-secondary">
        {isDemo() ? "演示模式：本地 mock，无需后端 · admin / admin123" : "开发默认: admin / admin123"}
      </p>
    </div>
  );
}
