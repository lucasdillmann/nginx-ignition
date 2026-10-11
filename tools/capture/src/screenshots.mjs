import fs from "node:fs"
import path from "node:path"
import { authenticate, createContext, launchBrowser, openPage, waitForApp } from "./browser.mjs"
import { docsImagesDirectory, formatBytes, screenshot } from "./config.mjs"

const log = message => console.log(`[screenshots] ${message}`)

const sleep = milliseconds => new Promise(resolve => setTimeout(resolve, milliseconds))

const gallery = [
    { name: "dashboard", route: "/" },
    { name: "hosts", route: "/hosts", companion: route => `/hosts/${route.hosts.blog}` },
    { name: "streams", route: "/streams", companion: route => `/streams/${route.streams.postgres}` },
    {
        name: "ssl-certificates",
        route: "/certificates",
        companion: route => `/certificates/${route.certificate.id}`,
    },
    {
        name: "access-lists",
        route: "/access-lists",
        companion: route => `/access-lists/${route.accessLists.restricted}`,
    },
    { name: "caches", route: "/caches", companion: route => `/caches/${route.cache.id}` },
    { name: "integrations", route: "/integrations", companion: route => `/integrations/${route.integration.id}` },
    { name: "vpns", route: "/vpns", companion: route => `/vpns/${route.vpn.id}` },
    { name: "logs", route: "/logs", prepare: showServerLogs },
    { name: "export", route: "/export" },
    { name: "settings", route: "/settings" },
    { name: "users", route: "/users", companion: route => `/users/${route.user.id}` },
    { name: "help", route: "/help" },
]

async function showServerLogs(page) {
    await page.getByText("Server logs", { exact: true }).click()
    await sleep(1200)
}

async function open(context, routePath, prepare) {
    const page = await openPage(context, routePath)
    await waitForApp(page)
    if (prepare) await prepare(page)
    await sleep(1600)

    return page
}

async function captureGallery(browser, token, seeded) {
    const context = await createContext(browser, {
        theme: "light",
        viewport: { width: screenshot.width, height: screenshot.height },
        deviceScaleFactor: screenshot.deviceScaleFactor,
    })
    await authenticate(context, token, "light")

    try {
        for (const item of gallery) {
            const page = await open(context, item.route, item.prepare)

            const destination = path.join(docsImagesDirectory, `${item.name}.png`)
            await page.screenshot({ path: destination })
            log(`  ${item.name.padEnd(18)} ${formatBytes(fs.statSync(destination).size)}`)
            await page.close()

            if (item.companion === undefined) continue

            const form = await open(context, item.companion(seeded))
            const formDestination = path.join(docsImagesDirectory, `${item.name}-form.png`)
            await form.screenshot({ path: formDestination })
            log(`  ${`${item.name}-form`.padEnd(18)} ${formatBytes(fs.statSync(formDestination).size)}`)
            await form.close()
        }
    } finally {
        await context.close()
    }
}

async function captureDarkGallery(browser, token) {
    const context = await createContext(browser, {
        theme: "dark",
        viewport: { width: screenshot.width, height: screenshot.height },
        deviceScaleFactor: screenshot.deviceScaleFactor,
    })
    await authenticate(context, token, "dark")

    try {
        for (const item of [gallery[0], gallery[1]]) {
            const page = await open(context, item.route)
            const destination = path.join(docsImagesDirectory, `${item.name}-dark.png`)
            await page.screenshot({ path: destination })
            log(`  ${`${item.name}-dark`.padEnd(18)} ${formatBytes(fs.statSync(destination).size)}`)
            await page.close()
        }
    } finally {
        await context.close()
    }
}

async function captureLogin(browser) {
    const context = await createContext(browser, {
        theme: "light",
        viewport: { width: screenshot.width, height: screenshot.height },
        deviceScaleFactor: screenshot.deviceScaleFactor,
    })

    try {
        const page = await open(context, "/login")
        const destination = path.join(docsImagesDirectory, "login.png")
        await page.screenshot({ path: destination })
        log(`  ${"login".padEnd(18)} ${formatBytes(fs.statSync(destination).size)}`)
    } finally {
        await context.close()
    }
}

export async function captureScreenshots(api, seeded) {
    const browser = await launchBrowser()

    try {
        log("capturing the documentation gallery")
        await captureGallery(browser, api.token, seeded)
        await captureDarkGallery(browser, api.token)
        await captureLogin(browser)
    } finally {
        await browser.close()
    }

    const { captureHero } = await import("./hero.mjs")
    await captureHero(api)
}
