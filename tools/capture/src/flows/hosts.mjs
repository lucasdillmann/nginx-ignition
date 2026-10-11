import { clickButton, delay, dismissReloadPrompt, fillField, scrollTo, toggleSwitch } from "./shared.mjs"

export default async function hosts(page) {
    await clickButton(page, "New host")
    await delay(1800)

    await fillField(page, "Domain names", "shop.example.com")
    await delay(1200)

    await fillField(page, "Destination URL", "http://127.0.0.1:3000")
    await delay(1200)

    await toggleSwitch(page, "Redirect HTTP to HTTPS")

    await scrollTo(page, "Standard bindings")
    await delay(1600)

    await clickButton(page, "Save")
    await delay(2200)
    await dismissReloadPrompt(page)
    await delay(1200)
}
