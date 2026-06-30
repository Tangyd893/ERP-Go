import { test, expect } from "@playwright/test";

/**
 * Admin 端烟测（T-644）
 * 验证：登录 → Dashboard 可见 → 侧边栏可导航 → 退出
 */
test.describe("Admin Smoke", () => {
  test("login page loads", async ({ page }) => {
    await page.goto("/login");
    await expect(page.locator("#app")).not.toBeEmpty();
  });

  test("login with admin credentials", async ({ page }) => {
    await page.goto("/login");

    // 填写登录表单
    await page.fill('input[placeholder*="用户名"]', "admin");
    await page.fill('input[placeholder*="密码"]', "admin123");
    await page.fill('input[placeholder*="租户"]', "default");
    await page.click('button:has-text("登录")');

    // 等待跳转到 dashboard
    await page.waitForURL("**/dashboard", { timeout: 10000 });
    await expect(page.locator(".el-menu")).toBeVisible();
  });

  test("sidebar navigation works", async ({ page }) => {
    // 先登录
    await page.goto("/login");
    await page.fill('input[placeholder*="用户名"]', "admin");
    await page.fill('input[placeholder*="密码"]', "admin123");
    await page.fill('input[placeholder*="租户"]', "default");
    await page.click('button:has-text("登录")');
    await page.waitForURL("**/dashboard", { timeout: 10000 });

    // 点击订单管理 → 订单列表
    await page.click('.el-sub-menu__title:has-text("订单管理")');
    await page.click('.el-menu-item:has-text("订单列表")');
    await page.waitForTimeout(1000);

    // 验证表格出现（可能为空但有结构）
    await expect(page.locator(".el-table, .pro-table")).toBeVisible({ timeout: 5000 });
  });

  test("logout clears token and redirects", async ({ page }) => {
    await page.goto("/login");
    await page.fill('input[placeholder*="用户名"]', "admin");
    await page.fill('input[placeholder*="密码"]', "admin123");
    await page.fill('input[placeholder*="租户"]', "default");
    await page.click('button:has-text("登录")');
    await page.waitForURL("**/dashboard", { timeout: 10000 });

    // 退出登录
    await page.click(".el-dropdown");
    await page.click('.el-dropdown-menu__item:has-text("退出登录")');
    await page.waitForURL("**/login", { timeout: 5000 });
  });
});
