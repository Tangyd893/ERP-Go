import { test, expect } from "@playwright/test";

/**
 * Admin React 并行端烟测（:5273）
 * 与 Vue admin-smoke 并存；不依赖 Element Plus 选择器。
 */
test.describe("Admin React Smoke", () => {
  test("login page loads", async ({ page }) => {
    await page.goto("/login");
    await expect(page.getByRole("heading", { name: /登录管理后台|ERP-Go/ })).toBeVisible();
    await expect(page.getByLabel("用户名")).toBeVisible();
  });

  test("login navigates to dashboard", async ({ page }) => {
    await page.goto("/login");
    await page.getByLabel("用户名").fill("admin");
    await page.getByLabel("密码").fill("admin123");
    await page.getByRole("button", { name: /登\s*录/ }).click();
    await page.waitForURL("**/dashboard", { timeout: 15000 });
    await expect(page.getByText("今日订单")).toBeVisible({ timeout: 10000 });
  });

  test("command palette opens with Ctrl+K", async ({ page }) => {
    await page.goto("/login");
    await page.getByLabel("用户名").fill("admin");
    await page.getByLabel("密码").fill("admin123");
    await page.getByRole("button", { name: /登\s*录/ }).click();
    await page.waitForURL("**/dashboard", { timeout: 15000 });

    await page.keyboard.press("Control+K");
    await expect(page.getByText("快速导航")).toBeVisible();
  });
});
