// @ts-expect-error Types not available
import path from "path"
import viteReact from "@vitejs/plugin-react"
import { defineConfig } from "vite"

// https://vitejs.dev/config/
export default defineConfig({
    base: "/",
    plugins: [viteReact()],
    build: {
        outDir: "build",
        emptyOutDir: true,
        chunkSizeWarningLimit: 1750,
        rollupOptions: {
            output: {
                entryFileNames: "assets/[name].[hash].js",
                chunkFileNames: "assets/[name].[hash].js",
                assetFileNames: "assets/[name].[hash].[ext]",
                manualChunks: {
                    "vendor.react": ["react", "react-dom", "react-router-dom"],
                    "vendor.codeium": ["@codeium/react-code-editor"],
                    "vendor.antd-core": [
                        "antd",
                        "@ant-design/v5-patch-for-react-19",
                        "@ant-design/icons",
                        "@ant-design/icons-svg",
                        "@ant-design/colors",
                        "@ant-design/fast-color",
                        "@ant-design/cssinjs",
                        "@ant-design/cssinjs-utils",
                        "@ant-design/cssinjs-utils",
                        "@ant-design/react-slick",
                    ],
                    "vendor.antd-pro": [
                        "@ant-design/pro-components",
                        "@ant-design/pro-provider",
                        "@ant-design/pro-utils",
                        "@ant-design/pro-card",
                        "@ant-design/pro-descriptions",
                        "@ant-design/pro-field",
                        "@ant-design/pro-layout",
                        "@ant-design/pro-list",
                        "@ant-design/pro-skeleton",
                        "@ant-design/pro-table",
                        "@ant-design/pro-form",
                    ],
                    "vendor.utils": ["qs", "debounce", "typescript"],
                },
            },
        },
    },
    resolve: {
        alias: {
            // @ts-expect-error Types not available
            "@": path.resolve(__dirname, "src"),
        },
    },
})
