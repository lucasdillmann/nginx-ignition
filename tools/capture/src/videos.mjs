import fs from "node:fs"
import os from "node:os"
import path from "node:path"
import { authenticate, createContext, launchBrowser, navigateToMenu, openPage, waitForApp } from "./browser.mjs"
import { formatBytes, helpVideosDirectory, video } from "./config.mjs"
import { encodeVideo } from "./encode.mjs"

const log = message => console.log(`[videos] ${message}`)

const sleep = milliseconds => new Promise(resolve => setTimeout(resolve, milliseconds))

const flows = {
    "access-lists": {
        route: "/access-lists",
        menuItem: "Access list",
        module: () => import("./flows/access-lists.mjs"),
    },
    caches: { route: "/caches", menuItem: "Cache configurations", module: () => import("./flows/caches.mjs") },
    hosts: { route: "/hosts", menuItem: "Hosts", module: () => import("./flows/hosts.mjs") },
    integrations: {
        route: "/integrations",
        menuItem: "Integrations",
        module: () => import("./flows/integrations.mjs"),
    },
    logs: { route: "/logs", menuItem: "Logs", module: () => import("./flows/logs.mjs") },
    settings: { route: "/settings", menuItem: "Settings", module: () => import("./flows/settings.mjs") },
    "ssl-certificates": {
        route: "/certificates",
        menuItem: "SSL certificates",
        module: () => import("./flows/ssl-certificates.mjs"),
    },
    streams: { route: "/streams", menuItem: "Streams", module: () => import("./flows/streams.mjs") },
    vpns: { route: "/vpns", menuItem: "VPN", module: () => import("./flows/vpns.mjs") },
}

const startScreenHoldMilliseconds = (video.guardSeconds + 1.5) * 1000

if (video.guardSeconds * 1000 >= startScreenHoldMilliseconds) {
    throw new Error(
        `CAPTURE_GUARD_SECONDS (${video.guardSeconds}) must be shorter than the ${startScreenHoldMilliseconds}ms ` +
            "the recorder waits on the starting screen, otherwise the clips would start after the first action.",
    )
}

async function recordFlow(browser, token, name) {
    const flow = flows[name]
    const module = await flow.module()
    const temporaryDirectory = fs.mkdtempSync(path.join(os.tmpdir(), "nginx-ignition-video-"))

    try {
        const context = await createContext(browser, {
            theme: "light",
            viewport: { width: video.width, height: video.height },
            deviceScaleFactor: 1,
            video: { dir: temporaryDirectory, size: { width: video.width, height: video.height } },
        })
        await authenticate(context, token, "light")

        const startedAt = Date.now()
        const page = await openPage(context, "/")
        await waitForApp(page)

        await navigateToMenu(page, flow.menuItem)

        if (!page.url().endsWith(flow.route)) {
            throw new Error(`"${name}" landed on ${page.url()} instead of ${flow.route}`)
        }

        await sleep(startScreenHoldMilliseconds)

        const startOffset = (Date.now() - startedAt) / 1000 - video.guardSeconds

        try {
            await module.default(page)
        } finally {
            await context.close()
        }

        const recorded = fs
            .readdirSync(temporaryDirectory)
            .map(file => path.join(temporaryDirectory, file))
            .find(file => file.endsWith(".webm"))

        if (recorded === undefined) throw new Error(`No video was recorded for "${name}"`)

        return { recorded, temporaryDirectory, startOffset }
    } catch (error) {
        fs.rmSync(temporaryDirectory, { recursive: true, force: true })
        throw error
    }
}

export async function recordVideos(api) {
    fs.mkdirSync(helpVideosDirectory, { recursive: true })

    const browser = await launchBrowser()
    const results = []

    try {
        for (const name of Object.keys(flows)) {
            const destination = path.join(helpVideosDirectory, `${name}.webp`)
            const { recorded, temporaryDirectory, startOffset } = await recordFlow(browser, api.token, name)

            try {
                const result = await encodeVideo(recorded, destination, { startOffset })
                results.push({ name, ...result })
            } finally {
                fs.rmSync(temporaryDirectory, { recursive: true, force: true })
            }
        }
    } finally {
        await browser.close()
    }

    report(results)

    return results
}

function report(results) {
    const total = results.reduce((accumulator, result) => accumulator + result.bytes, 0)

    for (const result of results) log(`  ${result.name.padEnd(18)} ${formatBytes(result.bytes)}`)
    log(`  ${"total".padEnd(18)} ${formatBytes(total)}`)

    const oversized = results.filter(result => result.bytes > video.budgetBytes)
    if (oversized.length > 0) {
        log(`warning: ${oversized.length} clip(s) above the ${formatBytes(video.budgetBytes)} budget`)
    }
}

export { flows }
