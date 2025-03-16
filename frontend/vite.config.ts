// @ts-expect-error Types not available
import path from "path"
import react from "@vitejs/plugin-react"
import { defineConfig } from "vite"

// https://vitejs.dev/config/
export default defineConfig({
    base: "/",
    plugins: [react()],
    build: {
        outDir: "build",
        emptyOutDir: true,
        chunkSizeWarningLimit: 1500,
        rollupOptions: {
            output: {
                entryFileNames: "assets/[name].[hash].js",
                chunkFileNames: "assets/[name].[hash].js",
                assetFileNames: "assets/[name].[hash].[ext]",
                manualChunks: {
                    "vendor.react": ["react", "react-dom", "react-router-dom"],
                    "vendor.antd-core": ["antd", "@ant-design/icons"],
                    "vendor.antd-pro": ["@ant-design/pro-form", "@ant-design/pro-components"],
                    "vendor.codeium": ["@codeium/react-code-editor"],
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
