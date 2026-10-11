import fs from "node:fs"
import os from "node:os"
import path from "node:path"
import { fileURLToPath } from "node:url"

const currentDirectory = path.dirname(fileURLToPath(import.meta.url))

export const repoRoot = path.resolve(currentDirectory, "..", "..", "..")
export const frontendRoot = path.join(repoRoot, "frontend")
export const frontendBuildDirectory = path.join(frontendRoot, "build")
export const migrationsScriptsDirectory = path.join(repoRoot, "internal", "database", "core", "migrations", "scripts")
export const helpVideosDirectory = path.join(frontendRoot, "src", "domain", "help", "videos")
export const docsImagesDirectory = path.join(repoRoot, "docs", "images")
export const readmeHeroImage = path.join(repoRoot, "docs", "readme-screenshots.png")

function numberFromEnvironment(name, fallback) {
    const value = process.env[name]
    if (value === undefined || value === "") return fallback

    const parsed = Number(value)
    if (!Number.isFinite(parsed)) throw new Error(`${name} must be a number, got ${value}`)

    return parsed
}

function booleanFromEnvironment(name, fallback) {
    const value = process.env[name]
    if (value === undefined || value === "") return fallback

    return value === "true" || value === "1"
}

export const port = numberFromEnvironment("CAPTURE_PORT", 8090)
export const baseUrl = `http://localhost:${port}`
export const workspaceDirectory = process.env.CAPTURE_WORKSPACE ?? path.join(os.tmpdir(), "nginx-ignition", "capture")
export const databaseDirectory = path.join(workspaceDirectory, "data")
export const nginxConfigDirectory = path.join(workspaceDirectory, "nginx")
export const binaryPath = path.join(workspaceDirectory, "nginx-ignition")
export const keepStack = booleanFromEnvironment("CAPTURE_KEEP_STACK", false)

export const httpPort = numberFromEnvironment("CAPTURE_HTTP_PORT", 8081)
export const httpsPort = numberFromEnvironment("CAPTURE_HTTPS_PORT", 8444)

function currentVersion() {
    if (process.env.CAPTURE_VERSION) return process.env.CAPTURE_VERSION

    const changelog = fs.readFileSync(path.join(repoRoot, "CHANGELOG.md"), "utf-8")
    const heading = changelog.match(/^##\s+(\S+)/m)

    return heading === null ? "0.0.0" : heading[1]
}

export const version = currentVersion()

export const credentials = {
    name: process.env.CAPTURE_USER_NAME ?? "Example user",
    username: process.env.CAPTURE_USER_USERNAME ?? "user",
    password: process.env.CAPTURE_USER_PASSWORD ?? "example-password-1234",
}

export const video = {
    width: numberFromEnvironment("CAPTURE_VIDEO_WIDTH", 1440),
    height: numberFromEnvironment("CAPTURE_VIDEO_HEIGHT", 900),
    framesPerSecond: numberFromEnvironment("CAPTURE_FPS", 10),
    quality: numberFromEnvironment("CAPTURE_QUALITY", 72),
    effort: numberFromEnvironment("CAPTURE_EFFORT", 6),
    budgetBytes: numberFromEnvironment("CAPTURE_VIDEO_BUDGET_BYTES", 1024 * 1024),
    guardSeconds: numberFromEnvironment("CAPTURE_GUARD_SECONDS", 2.5),
}

export const screenshot = {
    width: numberFromEnvironment("CAPTURE_SCREENSHOT_WIDTH", 1600),
    height: numberFromEnvironment("CAPTURE_SCREENSHOT_HEIGHT", 1000),
    deviceScaleFactor: numberFromEnvironment("CAPTURE_SCREENSHOT_SCALE", 1),
}

export const hero = {
    width: 1600,
    height: 925,
    deviceScaleFactor: 2,
}

export function formatBytes(bytes) {
    if (bytes < 1024) return `${bytes} B`
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`

    return `${(bytes / (1024 * 1024)).toFixed(2)} MB`
}
