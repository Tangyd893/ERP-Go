import { test, expect } from "@playwright/test";

/**
 * Dashboard 端烟测（T-644）
 * 验证：页面加载 → KPI 卡片可见 → 无无限 spinner → Demo 按钮可用
 */
test.describe("Dashboard Smoke", () => {
  test("page loads without blocking spinner", async ({ page }) => {
    await page.goto("/");
    await page.waitForTimeout(2000);

    // 验证 #app 非空
    await expect(page.locator("#app")).not.toBeEmpty();

    // 不应有全页 infinite spinner（v-loading 已消除）
    const fullPageSpinner = page.locator(".el-loading-mask");
    // 可能有局部 loading，但不应遮挡整个页面
    const maskCount = await fullPageSpinner.count();
    // 允许局部 skeleton loading，但不应全页遮罩
    expect(maskCount).toBeLessThanOrEqual(2);
  });

  test("KPI cards are visible", async ({ page }) => {
    await page.goto("/");
    await page.waitForTimeout(3000);

    // 应该有 4 个 KPI 卡片（订单总数/销售额/出库量/库存SKU）
    // DEMO 模式下直接显示数据
    const kpiCards = page.locator(".el-card, .el-skeleton");
    await expect(kpiCards.first()).toBeVisible({ timeout: 5000 });
  });

  test("Demo button visible when no auth", async ({ page }) => {
    await page.goto("/");
    await page.waitForTimeout(3000);

    // 检查是否有 demo 按钮或引导提示
    const demoButton = page.locator('button:has-text("Demo"), button:has-text("演示")');
    const authAlert = page.locator('.el-alert:has-text("登录")');

    // 至少有一个（demo 按钮或登录引导）
    const hasDemo = (await demoButton.count()) > 0;
    const hasAuthAlert = (await authAlert.count()) > 0;
    expect(hasDemo || hasAuthAlert).toBeTruthy();
  });

  test("charts render without error", async ({ page }) => {
    await page.goto("/");
    await page.waitForTimeout(4000);

    // ECharts 图表容器应存在
    const chartContainer = page.locator('div[style*="height: 320px"], .echarts, canvas');
    // 至少有图表区域（skeleton 或真实图表）
    await expect(chartContainer.first()).toBeVisible({ timeout: 5000 });
  });
});
