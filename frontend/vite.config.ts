import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import { sentryVitePlugin } from "@sentry/vite-plugin";

export default defineConfig({
  build: {
    sourcemap: true,
  },
  plugins: [
    react(),
    sentryVitePlugin({
      authToken: process.env.SENTRY_TOKEN,
      org: "kalisto",
      project: "javascript-react",
    }),
  ],
});
