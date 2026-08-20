import { test, expect } from "@playwright/test";

test.describe("PDA React Smoke", () => {
  test("login page loads", async ({ page }) => {
    await page.goto("/login");
    await expect(page.getByLabel(/用户名|账号/)).toBeVisible({ timeout: 10000 });
  });

  test("login shows bottom tabs", async ({ page }) => {
    await page.goto("/login");
    await page.getByLabel(/用户名|账号/).fill("admin");
    await page.getByLabel(/密码/).fill("admin123");
    await page.getByRole("button", { name: /登\s*录|登录/ }).click();
    await page.waitForURL((url) => !url.pathname.includes("/login"), { timeout: 15000 });
    await expect(page.getByRole("link", { name: "首页" }).or(page.getByText("首页"))).toBeVisible({
      timeout: 10000,
    });
  });
});
