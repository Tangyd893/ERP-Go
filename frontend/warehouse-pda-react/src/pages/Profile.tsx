import { useNavigate } from "react-router-dom";
import { Button, Card, CardContent, CardHeader, CardTitle } from "@erp/shared-react";
import { useAuthStore } from "@/stores/auth";

export default function Profile() {
  const navigate = useNavigate();
  const username = useAuthStore((s) => s.username);
  const tenantId = useAuthStore((s) => s.tenantId);
  const logout = useAuthStore((s) => s.logout);

  function handleLogout() {
    logout();
    navigate("/login");
  }

  return (
    <div>
      <Card>
        <CardHeader className="pb-2">
          <CardTitle className="text-base">我的</CardTitle>
        </CardHeader>
        <CardContent className="space-y-3">
          <div className="rounded-md border">
            <div className="flex border-b px-3 py-2.5 text-sm">
              <span className="w-20 shrink-0 text-pda-text-secondary">用户名</span>
              <span>{username || "-"}</span>
            </div>
            <div className="flex px-3 py-2.5 text-sm">
              <span className="w-20 shrink-0 text-pda-text-secondary">租户</span>
              <span>{tenantId || "-"}</span>
            </div>
          </div>
          <div className="pt-2 text-center">
            <Button variant="destructive" className="min-h-touch px-8" onClick={handleLogout}>
              退出登录
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
