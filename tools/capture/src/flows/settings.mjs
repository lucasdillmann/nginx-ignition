import { clickButton, delay, dismissReloadPrompt, toggleSwitch } from "./shared.mjs"

export default async function settings(page) {
    await toggleSwitch(page, "Server tokens")

    await page.mouse.wheel(0, 320)
    await delay(1600)

    await clickButton(page, "Save")
    await delay(2200)
    await dismissReloadPrompt(page)
    await delay(1200)
}
