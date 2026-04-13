// @ts-expect-error Types not available
import path from "path"
// @ts-expect-error Types are available even though the IDE says otherwise
import viteReact from "@vitejs/plugin-react"
// @ts-expect-error Types are available even though the IDE says otherwise
import { defineConfig } from "vite"

// https://vitejs.dev/config/
export default defineConfig({
    base: "/",
    plugins: [viteReact()],
    server: {
        host: "0.0.0.0",
        port: 8080,
        proxy: {
            "/api": {
                target: "http://localhost:8090",
            },
        },
    },
    build: {
        outDir: "build",
        emptyOutDir: true,
        chunkSizeWarningLimit: 1750,
        rolldownOptions: {
            output: {
                entryFileNames: "assets/[name].[hash].js",
                chunkFileNames: "assets/[name].[hash].js",
                assetFileNames: "assets/[name].[hash].[ext]",
                codeSplitting: {
                    groups: [
                        {
                            name: "vendor.react",
                            test: /node_modules[\\/](react|react-dom|react-router-dom)(\/|$)/,
                            priority: 30,
                        },
                        {
                            name: "vendor.monaco",
                            test: /node_modules[\\/](monaco-editor|@monaco-editor\/react)(\/|$)/,
                            priority: 29,
                        },
                        {
                            name: "vendor.antd-pro",
                            test: /node_modules[\\/]@ant-design\/(pro-components|pro-provider|pro-utils|pro-card|pro-descriptions|pro-field|pro-layout|pro-list|pro-skeleton|pro-table|pro-form)\//,
                            priority: 28,
                        },
                        {
                            name: "vendor.antd-charts",
                            test: /node_modules[\\/]@ant-design\/charts\//,
                            priority: 27,
                        },
                        {
                            name: "vendor.antd-core",
                            test: /node_modules[\\/](antd|@ant-design\/(v5-patch-for-react-19|icons-svg|icons|colors|fast-color|cssinjs-utils|cssinjs|react-slick))\//,
                            priority: 26,
                        },
                        {
                            name: "vendor.antv",
                            test: /node_modules[\\/]@antv\//,
                            priority: 25,
                        },
                        {
                            name: "vendor.utils",
                            test: /node_modules[\\/](qs|debounce|typescript)\//,
                            priority: 21,
                        },
                        {
                            name: "vendor.common",
                            test: /node_modules/,
                            priority: 10,
                        },
                        {
                            name: "common",
                            minShareCount: 2,
                            minSize: 10000,
                            priority: 5,
                        },
                    ],
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
