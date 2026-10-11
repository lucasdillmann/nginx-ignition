import { chromium } from "playwright"
import { baseUrl } from "./config.mjs"

export async function launchBrowser() {
    return chromium.launch({ args: ["--force-color-profile=srgb", "--font-render-hinting=none"] })
}

export async function createContext(browser, options = {}) {
    const {
        theme = "light",
        viewport = { width: 1280, height: 800 },
        deviceScaleFactor = 1,
        video = undefined,
        storageState = {},
    } = options

    return browser.newContext({
        viewport,
        deviceScaleFactor,
        locale: "en-US",
        timezoneId: "America/Sao_Paulo",
        colorScheme: theme === "dark" ? "dark" : "light",
        caret: "hide",
        ...(video === undefined ? {} : { recordVideo: video }),
        storageState: {
            cookies: [],
            origins: storageState.origins ?? [],
        },
    })
}

export async function authenticate(context, token, theme = "light") {
    await context.addInitScript(
        ([value, themeValue]) => {
            window.localStorage.setItem("nginxIgnition.authentication.token", JSON.stringify(value))
            window.localStorage.setItem("nginxIgnition.theme", JSON.stringify(themeValue))
            window.localStorage.setItem("nginxIgnition.i18n.language", JSON.stringify("en"))
        },
        [token, theme],
    )
}

async function pinReportedVersion(page) {
    await page.route("**/api/frontend/configuration", async route => {
        const response = await route.fetch()
        if (!response.ok()) return route.fulfill({ response })

        const configuration = await response.json()
        if (configuration?.version?.current == null) return route.fulfill({ response })

        await route.fulfill({
            response,
            json: {
                ...configuration,
                version: { ...configuration.version, latest: configuration.version.current },
            },
        })
    })
}

async function overrideNginxCapabilities(page) {
    await page.route("**/api/nginx/metadata", async route => {
        const response = await route.fetch()
        if (!response.ok()) return route.fulfill({ response })

        const metadata = await response.json()
        if (metadata?.availableSupport == null) return route.fulfill({ response })

        await route.fulfill({
            response,
            json: {
                ...metadata,
                availableSupport: {
                    ...metadata.availableSupport,
                    runCode: "DYNAMIC",
                },
            },
        })
    })
}

async function preparePage(page) {
    page.setDefaultTimeout(20_000)
    await pinReportedVersion(page)
    await overrideNginxCapabilities(page)
}

export async function openPage(context, routePath = "/") {
    const page = await context.newPage()
    await preparePage(page)

    await page.goto(`${baseUrl}${routePath}`, { waitUntil: "commit" })

    return page
}

export async function waitForApp(page) {
    await page.waitForFunction(() => document.getElementById("preloader") === null, undefined, { timeout: 30_000 })
    await page.waitForLoadState("networkidle").catch(() => undefined)
    await page.evaluate(() => document.fonts.ready)
}

export async function navigateToMenu(page, menuItem) {
    await page.locator(".ant-menu-item").filter({ hasText: menuItem }).first().click()
    await page.waitForLoadState("networkidle").catch(() => undefined)
    await page.evaluate(() => document.fonts.ready)
}
