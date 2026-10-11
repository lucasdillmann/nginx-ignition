import { addTag, clickButton, delay, dismissReloadPrompt, fillField, scrollTo } from "./shared.mjs"

export default async function caches(page) {
    await clickButton(page, "New cache configuration")
    await delay(1800)

    await fillField(page, "Name", "Static assets cache")
    await delay(1000)

    await fillField(page, "Maximum size", "128")
    await delay(1000)

    await addTag(page, "File extensions", "woff2")
    await addTag(page, "File extensions", "webp")

    await scrollTo(page, "Stale contents and revalidation")
    await delay(1600)

    await clickButton(page, "Save")
    await delay(2200)
    await dismissReloadPrompt(page)
    await delay(1200)
}
