import { spawn, spawnSync } from "node:child_process"
import fs from "node:fs"
import net from "node:net"
import path from "node:path"
import {
    baseUrl,
    binaryPath,
    databaseDirectory,
    frontendBuildDirectory,
    helpVideosDirectory,
    keepStack,
    migrationsScriptsDirectory,
    nginxConfigDirectory,
    port,
    repoRoot,
    version,
    workspaceDirectory,
} from "./config.mjs"

const expectedVideos = [
    "access-lists",
    "caches",
    "hosts",
    "integrations",
    "logs",
    "settings",
    "ssl-certificates",
    "streams",
    "vpns",
]

const serverPidFile = path.join(workspaceDirectory, "server.pid")

const log = message => console.log(`[stack] ${message}`)

const sleep = milliseconds => new Promise(resolve => setTimeout(resolve, milliseconds))

function ensureDirectory(directory) {
    fs.mkdirSync(directory, { recursive: true })
}

function ensurePortIsFree() {
    return new Promise((resolve, reject) => {
        const probe = net.connect({ port, host: "localhost" })
        probe.once("connect", () => {
            probe.destroy()
            reject(
                new Error(
                    `Port ${port} is already in use. Stop the process using it or run the capture with CAPTURE_PORT=<other port>.`,
                ),
            )
        })
        probe.once("error", () => resolve())
    })
}

async function waitForServer(timeoutMilliseconds = 90_000) {
    const deadline = Date.now() + timeoutMilliseconds

    while (Date.now() < deadline) {
        try {
            const response = await fetch(`${baseUrl}/api/i18n`)
            if (response.ok) return
        } catch {}

        await sleep(250)
    }

    throw new Error(`The server did not become ready at ${baseUrl} within ${timeoutMilliseconds}ms`)
}

function buildBinary() {
    log("building the server binary")
    const result = spawnSync(
        "go",
        [
            "build",
            "-ldflags",
            `-X 'github.com/lucasdillmann/nginx-ignition/internal/business/core/version.Number=${version}'`,
            "-o",
            binaryPath,
            "./cmd",
        ],
        {
            cwd: repoRoot,
            stdio: ["ignore", "ignore", "pipe"],
            encoding: "utf-8",
        },
    )

    if (result.status !== 0) throw new Error(`go build failed:\n${result.stderr}`)
}

function readPid(pidFile) {
    if (!fs.existsSync(pidFile)) return undefined

    const pid = Number.parseInt(fs.readFileSync(pidFile, "utf-8").trim(), 10)

    return Number.isInteger(pid) && pid > 0 ? pid : undefined
}

function belongsToCapture(pid, marker) {
    const result = spawnSync("ps", ["-p", String(pid), "-o", "command="], { encoding: "utf-8" })

    return result.status === 0 && result.stdout.includes(marker)
}

function stopProcess(pidFile, marker, signal, description) {
    const pid = readPid(pidFile)
    if (pid === undefined || pid === process.pid || !belongsToCapture(pid, marker)) return

    try {
        process.kill(pid, signal)
        log(`stopped the ${description} (pid ${pid})`)
    } catch {}
}

function stopManagedNginx() {
    stopProcess(path.join(nginxConfigDirectory, "nginx.pid"), nginxConfigDirectory, "SIGQUIT", "capture nginx")
}

function stopLeftoverStack() {
    stopProcess(serverPidFile, binaryPath, "SIGTERM", "server from a previous run")
    stopProcess(path.join(nginxConfigDirectory, "nginx.pid"), nginxConfigDirectory, "SIGTERM", "nginx from a previous run")
}

export async function startStack() {
    if (!fs.existsSync(path.join(frontendBuildDirectory, "index.html"))) {
        throw new Error(
            `The frontend build is missing at ${frontendBuildDirectory}. Run "make .frontend-build" before capturing.`,
        )
    }

    const missingVideos = expectedVideos.filter(name => !fs.existsSync(path.join(helpVideosDirectory, `${name}.webp`)))
    if (missingVideos.length > 0) {
        throw new Error(
            `The help page tours are missing (${missingVideos.join(", ")}). They are imported by the frontend, ` +
                "so the app cannot be built until they exist again.",
        )
    }

    stopLeftoverStack()
    await ensurePortIsFree()

    if (!keepStack) {
        log("resetting the capture workspace")
        fs.rmSync(workspaceDirectory, { recursive: true, force: true })
    }

    ensureDirectory(databaseDirectory)
    ensureDirectory(nginxConfigDirectory)

    buildBinary()

    log(`starting the server on ${baseUrl}`)
    const server = spawn(binaryPath, [], {
        cwd: repoRoot,
        stdio: ["ignore", "pipe", "pipe"],
        env: {
            ...process.env,
            NGINX_IGNITION_SERVER_PORT: String(port),
            NGINX_IGNITION_SERVER_FRONTEND_PATH: frontendBuildDirectory,
            NGINX_IGNITION_DATABASE_MIGRATIONS_PATH: migrationsScriptsDirectory,
            NGINX_IGNITION_DATABASE_DRIVER: "sqlite",
            NGINX_IGNITION_DATABASE_DATA_PATH: databaseDirectory,
            NGINX_IGNITION_NGINX_CONFIG_PATH: nginxConfigDirectory,
        },
    })

    fs.writeFileSync(serverPidFile, String(server.pid))

    let output = ""
    server.stdout.on("data", chunk => (output += chunk))
    server.stderr.on("data", chunk => (output += chunk))

    let exited = false
    const onExit = new Promise(resolve =>
        server.once("exit", () => {
            exited = true
            resolve()
        }),
    )

    try {
        await waitForServer()
    } catch (error) {
        server.kill("SIGKILL")
        throw new Error(`${error.message}\n\nServer output:\n${output}`)
    }

    if (exited) throw new Error(`The server exited during startup.\n\nServer output:\n${output}`)

    log("server is ready")

    return {
        output: () => output,
        async stop() {
            log("stopping the server")
            stopManagedNginx()

            if (exited) return

            server.kill("SIGTERM")
            await Promise.race([onExit, sleep(10_000)])

            if (!exited) server.kill("SIGKILL")

            fs.rmSync(serverPidFile, { force: true })

            if (!keepStack) fs.rmSync(workspaceDirectory, { recursive: true, force: true })
        },
    }
}
