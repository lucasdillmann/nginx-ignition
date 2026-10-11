import { authenticate, seed } from "./api.mjs"
import { startStack } from "./stack.mjs"

const log = message => console.log(`[capture] ${message}`)

const targets = process.argv.slice(2)
const wantsScreenshots = targets.length === 0 || targets.includes("screenshots")
const wantsVideos = targets.length === 0 || targets.includes("videos")

const unknown = targets.filter(target => !["screenshots", "videos"].includes(target))
if (unknown.length > 0) throw new Error(`Unknown target(s): ${unknown.join(", ")}`)

log("starting the capture stack")
const stack = await startStack()

try {
    const api = await authenticate()
    const seeded = await seed(api)

    if (wantsScreenshots) {
        const { captureScreenshots } = await import("./screenshots.mjs")
        await captureScreenshots(api, seeded)
    }

    if (wantsVideos) {
        const { recordVideos } = await import("./videos.mjs")
        await recordVideos(api)
    }
} finally {
    await stack.stop()
}

log("done")
