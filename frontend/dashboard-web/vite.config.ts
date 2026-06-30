import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import path from "node:path";

export default defineConfig({
  plugins: [vue()],
  define: {
    "import.meta.env.VITE_AUTH_REDIRECT": JSON.stringify("false"),
  },
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "src"),
      "@erp/shared": path.resolve(__dirname, "../shared/src"),
    },
  },
  server: {
    port: 5175,
    host: true,
    proxy: {
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true,
      },
    },
  },
});
