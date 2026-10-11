import fs from "node:fs"
import path from "node:path"
import sharp from "sharp"
import { authenticate, createContext, launchBrowser, openPage, waitForApp } from "./browser.mjs"
import { formatBytes, hero, readmeHeroImage } from "./config.mjs"

const log = message => console.log(`[hero] ${message}`)

const sleep = milliseconds => new Promise(resolve => setTimeout(resolve, milliseconds))

const canvas = { width: 2560, height: 1440 }
const cornerRadius = 14
const shadowBlur = 18
const shadowOffset = 8
const shadowOpacity = 0.28
const padding = shadowBlur * 2

const panels = [
    { source: "host-details", left: 89, top: 46, width: 2087 },
    { source: "integrations", left: 1365, top: 208, width: 936 },
    { source: "logs", left: 1507, top: 594, width: 905 },
    { source: "certificate-details", left: 1277, top: 862, width: 905 },
]

const captures = [
    { source: "host-details", route: "/hosts", prepare: openFirstHostDetails },
    { source: "integrations", route: "/integrations" },
    {
        source: "logs",
        route: "/logs",
        prepare: async page => {
            await page.getByText("Server logs", { exact: true }).click()
        },
    },
    { source: "certificate-details", route: "/certificates", prepare: openFirstCertificateDetails },
]

async function openFirstHostDetails(page) {
    await page
        .locator("tr", { hasText: "blog.example.com" })
        .first()
        .locator("td")
        .last()
        .locator("a[href^='/hosts/']")
        .first()
        .click()
    await waitForApp(page)
}

async function openFirstCertificateDetails(page) {
    await page.locator("tr", { hasText: "example.com" }).first().locator("td").last().locator("a[href]").first().click()
    await waitForApp(page)
}

async function withRoundedCorners(buffer) {
    const { width, height } = await sharp(buffer).metadata()

    const mask = Buffer.from(
        `<svg width="${width}" height="${height}"><rect width="${width}" height="${height}" ` +
            `rx="${cornerRadius}" ry="${cornerRadius}" fill="#ffffff"/></svg>`,
    )

    return sharp(buffer)
        .composite([{ input: mask, blend: "dest-in" }])
        .png()
        .toBuffer()
}

async function buildPanel(buffer, width) {
    const resized = await sharp(buffer).resize({ width }).png().toBuffer()
    const { width: contentWidth, height: contentHeight } = await sharp(resized).metadata()

    const content = await withRoundedCorners(resized)

    const shape = Buffer.from(
        `<svg width="${contentWidth + padding * 2}" height="${contentHeight + padding * 2}">` +
            `<rect x="${padding}" y="${padding + shadowOffset}" width="${contentWidth}" height="${contentHeight}" ` +
            `rx="${cornerRadius}" ry="${cornerRadius}" fill="#0b1b2b" fill-opacity="${shadowOpacity}"/>` +
            `</svg>`,
    )

    const shadow = await sharp(shape).blur(shadowBlur).png().toBuffer()

    return { shadow, content, offset: -padding }
}

async function capturePanels(browser, token) {
    const context = await createContext(browser, {
        theme: "light",
        viewport: { width: hero.width, height: hero.height },
        deviceScaleFactor: hero.deviceScaleFactor,
    })
    await authenticate(context, token, "light")

    const screenshots = {}

    try {
        for (const capture of captures) {
            const page = await openPage(context, capture.route)
            await waitForApp(page)
            if (capture.prepare) await capture.prepare(page)
            await sleep(1800)
            screenshots[capture.source] = await page.screenshot()
            await page.close()
        }
    } finally {
        await context.close()
    }

    return screenshots
}

export async function captureHero(api) {
    const browser = await launchBrowser()

    try {
        log("capturing the README header panels")
        const screenshots = await capturePanels(browser, api.token)

        const composites = []
        const built = await Promise.all(panels.map(panel => buildPanel(screenshots[panel.source], panel.width)))

        built.forEach((panel, index) => {
            const placement = panels[index]
            composites.push(
                {
                    input: panel.shadow,
                    left: Math.max(0, placement.left + panel.offset),
                    top: Math.max(0, placement.top + panel.offset),
                },
                { input: panel.content, left: placement.left, top: placement.top },
            )
        })

        await sharp({
            create: {
                width: canvas.width,
                height: canvas.height,
                channels: 4,
                background: { r: 0, g: 0, b: 0, alpha: 0 },
            },
        })
            .composite(composites)
            .png({ compressionLevel: 9 })
            .toFile(readmeHeroImage)

        log(`wrote ${path.basename(readmeHeroImage)} (${formatBytes(fs.statSync(readmeHeroImage).size)})`)
    } finally {
        await browser.close()
    }
}
