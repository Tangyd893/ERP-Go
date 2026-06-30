import { test, expect } from "@playwright/test";

/**
 * PDA 端烟测（T-644）
 * 验证：登录 → 首页任务计数 → 底部 Tab 导航 → 退出
 */
test.describe("PDA Smoke", () => {
  test("login page loads", async ({ page }) => {
    await page.goto("/login");
    await expect(page.locator("#app")).not.toBeEmpty();
    // PDA 有独立登录页样式
    await expect(page.locator(".pda-login__title, h2, .el-form")).toBeVisible();
  });

  test("login and see home dashboard", async ({ page }) => {
    await page.goto("/login");

    // 填写 PDA 登录
    await page.fill('input[placeholder*="工号"]', "admin");
    await page.fill('input[placeholder*="密码"]', "admin123");
    await page.click('button:has-text("登录")');

    // 等待跳转到首页
    await page.waitForURL("**/", { timeout: 10000 });
    await page.waitForTimeout(1000);

    // 验证标题 + 底部导航
    await expect(page.locator(".pda-header__title, h3")).toBeVisible();
    await expect(page.locator(".pda-tabbar")).toBeVisible();
  });

  test("bottom tab navigation works", async ({ page }) => {
    // 先登录
    await page.goto("/login");
    await page.fill('input[placeholder*="工号"]', "admin");
    await page.fill('input[placeholder*="密码"]', "admin123");
    await page.click('button:has-text("登录")');
    await page.waitForURL("**/", { timeout: 10000 });

    // 点击拣货 Tab
    await page.click('.pda-tabbar__item:has-text("拣货")');
    await page.waitForTimeout(500);
    expect(page.url()).toContain("/pick");

    // 点击复核 Tab
    await page.click('.pda-tabbar__item:has-text("复核")');
    await page.waitForTimeout(500);
    expect(page.url()).toContain("/check");

    // 回到首页
    await page.click('.pda-tabbar__item:has-text("首页")');
    await page.waitForTimeout(500);
    expect(page.url()).toBe(page.url().replace(/\/pick|\/check/, ""));
  });

  test("logout returns to login", async ({ page }) => {
    await page.goto("/login");
    await page.fill('input[placeholder*="工号"]', "admin");
    await page.fill('input[placeholder*="密码"]', "admin123");
    await page.click('button:has-text("登录")');
    await page.waitForURL("**/", { timeout: 10000 });

    // 退出
    await page.click('.pda-header__right button:has-text("退出")');
    await page.waitForURL("**/login", { timeout: 5000 });
  });
});
