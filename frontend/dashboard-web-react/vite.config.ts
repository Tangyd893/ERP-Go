import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import path from "node:path";

export default defineConfig({
  plugins: [react()],
  define: {
    "import.meta.env.VITE_AUTH_REDIRECT": JSON.stringify("false"),
  },
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "src"),
      "@erp/shared-react": path.resolve(__dirname, "../shared-react/src"),
    },
  },
  server: {
    port: 5275,
    host: true,
    proxy: {
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true,
      },
    },
  },
});
