import globals from "globals"
import pluginJs from "@eslint/js"
import tseslint from "typescript-eslint"
import pluginReact from "eslint-plugin-react"

/** @type {import("eslint").Linter.Config[]} */
export default [
    {
        ignores: ["**/*.generated.ts"],
    },
    {
        files: ["src/**/*.{js,mjs,cjs,ts,jsx,tsx}"],
    },
    {
        languageOptions: {
            globals: {
                ...globals.browser,
                ...globals.node,
            },
            parserOptions: {
                ecmaFeatures: {
                    jsx: true,
                },
            },
        },
    },
    {
        ...pluginJs.configs.recommended,
    },
    ...tseslint.configs.recommended,
    {
        ...pluginReact.configs.flat.recommended,
        settings: {
            react: {
                version: "detect",
            },
        },
    },
    {
        files: ["src/**/*.{ts,tsx}"],
        rules: {
            "@typescript-eslint/no-explicit-any": "off",
            "@typescript-eslint/no-extra-non-null-assertion": "off",
            "@typescript-eslint/no-this-alias": "off",
            "react/no-unescaped-entities": "off",
            "@typescript-eslint/no-unused-vars": [
                "error",
                {
                    argsIgnorePattern: "^_",
                    varsIgnorePattern: "^_",
                    caughtErrorsIgnorePattern: "^_",
                },
            ],
        },
    },
]
