import { fileURLToPath, URL } from "node:url";

import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import VueI18nPlugin from "@intlify/unplugin-vue-i18n/vite";
import { compression } from "vite-plugin-compression2";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue({
      template: {
        compilerOptions: {
          compatConfig: {
            MODE: 2,
          },
        },
      },
    }),
    VueI18nPlugin(),
    compression({ include: /\.js$/i, deleteOriginalAssets: true }),
  ],
  resolve: {
    alias: {
      vue: "@vue/compat",
      "@/": fileURLToPath(new URL("./src/", import.meta.url)),
    },
  },
  base: "",
  experimental: {
    renderBuiltUrl(filename, { hostType }) {
      if (hostType === "js") {
        return { runtime: `window.__prependStaticUrl("${filename}")` };
      } else if (hostType === "html") {
        return `[{[ .StaticURL ]}]/${filename}`;
      } else {
        return { relative: true };
      }
    },
  },
});
