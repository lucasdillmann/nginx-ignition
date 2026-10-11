import { spawn, spawnSync } from "node:child_process"
import fs from "node:fs"
import os from "node:os"
import path from "node:path"
import sharp from "sharp"
import { video } from "./config.mjs"

const log = message => console.log(`[encode] ${message}`)

function run(command, arguments_) {
    return new Promise((resolve, reject) => {
        const child = spawn(command, arguments_, { stdio: ["ignore", "ignore", "pipe"] })

        let error = ""
        child.stderr.on("data", chunk => (error += chunk))
        child.once("error", reject)
        child.once("exit", code => {
            if (code === 0) resolve()
            else reject(new Error(`${command} failed with ${code}:\n${error}`))
        })
    })
}

function probe(source) {
    const result = spawnSync(
        "ffprobe",
        ["-v", "error", "-select_streams", "v:0", "-show_entries", "stream=width,height", "-of", "csv=p=0", source],
        { encoding: "utf-8" },
    )

    if (result.status !== 0) throw new Error(`ffprobe failed on ${source}:\n${result.stderr}`)

    const [width, height] = result.stdout.trim().split(",").map(Number)

    return { width, height }
}

async function extractFrames(source, framesDirectory, width, framesPerSecond, startOffset) {
    const filters = []

    if (startOffset > 0) filters.push(`trim=start=${startOffset.toFixed(2)}`, "setpts=PTS-STARTPTS")

    filters.push(`fps=${framesPerSecond}`)

    if (probe(source).width !== width) filters.push(`scale=${width}:-2:flags=lanczos`)

    await run("ffmpeg", [
        "-y",
        "-i",
        source,
        "-vf",
        filters.join(","),
        "-compression_level",
        "3",
        path.join(framesDirectory, "%05d.png"),
    ])
}

function listFrames(framesDirectory) {
    return fs
        .readdirSync(framesDirectory)
        .filter(name => name.endsWith(".png"))
        .sort()
        .map(name => path.join(framesDirectory, name))
}

export async function encodeAnimatedWebp(source, destination, options = {}) {
    const {
        width = video.width,
        framesPerSecond = video.framesPerSecond,
        quality = video.quality,
        effort = video.effort,
        startOffset = 0,
    } = options

    const framesDirectory = fs.mkdtempSync(path.join(os.tmpdir(), "nginx-ignition-frames-"))

    try {
        await extractFrames(source, framesDirectory, width, framesPerSecond, startOffset)

        const frames = listFrames(framesDirectory)
        if (frames.length === 0) throw new Error(`No frames were extracted from ${source}`)

        const delay = new Array(frames.length).fill(Math.round(1000 / framesPerSecond))

        await sharp(frames, { join: { animated: true } })
            .webp({ quality, effort, loop: 0, delay })
            .toFile(destination)

        const { size } = fs.statSync(destination)

        return { frames: frames.length, bytes: size }
    } finally {
        fs.rmSync(framesDirectory, { recursive: true, force: true })
    }
}

export async function encodeVideo(source, destination, options) {
    log(`encoding ${path.basename(destination)}`)
    const result = await encodeAnimatedWebp(source, destination, options)

    log(`  ${result.frames} frames, ${(result.bytes / 1024).toFixed(0)} KB`)

    return result
}
