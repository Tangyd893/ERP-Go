import { defineConfig, devices } from "@playwright/test";

export default defineConfig({
  testDir: "./testing/e2e/specs",
  timeout: 60000,
  retries: 1,
  workers: 1,
  reporter: [["html", { outputFolder: "testing/e2e/playwright-report" }], ["list"]],

  use: {
    baseURL: "http://localhost:5173",
    trace: "on-first-retry",
    screenshot: "only-on-failure",
  },

  projects: [
    {
      name: "admin",
      testMatch: /admin.*\.spec\.ts/,
      use: {
        ...devices["Desktop Chrome"],
        baseURL: "http://localhost:5173",
      },
    },
    {
      name: "pda",
      testMatch: /pda.*\.spec\.ts/,
      use: {
        ...devices["Pixel 7"],
        baseURL: "http://localhost:5174",
      },
    },
    {
      name: "dashboard",
      testMatch: /dashboard.*\.spec\.ts/,
      use: {
        ...devices["Desktop Chrome"],
        baseURL: "http://localhost:5175",
      },
    },
  ],

  webServer: [
    {
      command: "npm run dev:admin",
      port: 5173,
      reuseExistingServer: true,
      timeout: 30000,
    },
    {
      command: "npm run dev:pda",
      port: 5174,
      reuseExistingServer: true,
      timeout: 30000,
    },
    {
      command: "npm run dev:dashboard",
      port: 5175,
      reuseExistingServer: true,
      timeout: 30000,
    },
  ],
});
