import { clickButton, delay, dismissReloadPrompt, fillField, selectOption } from "./shared.mjs"

export default async function accessLists(page) {
    await clickButton(page, "New access list")
    await delay(1800)

    await fillField(page, "Name", "Internal tools only")
    await fillField(page, "Realm name", "Restricted area")
    await delay(1200)

    await selectOption(page, "Default outcome", "Deny access")

    await page.locator("button", { hasText: "Add credential" }).first().click()
    await delay(1000)

    await fillField(page, "Username", "demo")
    await fillField(page, "Password", "demo-password")
    await delay(1200)

    await page.locator("button", { hasText: "Add IP address list" }).first().click()
    await delay(1000)

    await selectOption(page, "Outcome", "Allow")
    await fillField(page, "IP addresses or ranges", "192.168.0.0/16\n10.0.0.0/8")

    await clickButton(page, "Save")
    await delay(2200)
    await dismissReloadPrompt(page)
    await delay(1200)
}
