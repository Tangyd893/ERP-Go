import { test, expect } from "@playwright/test";

test.describe("Dashboard React Smoke", () => {
  test("page renders KPI region", async ({ page }) => {
    await page.goto("/");
    await expect(page.getByText(/订单|销售|出库|SKU|经营|看板/)).toBeVisible({
      timeout: 15000,
    });
  });
});
