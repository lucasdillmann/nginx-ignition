import { clickButton, delay, fillField, selectOption, toggleSwitch } from "./shared.mjs"

export default async function integrations(page) {
    await clickButton(page, "New integration")
    await delay(1800)

    await fillField(page, "Name", "Docker socket")
    await delay(1200)

    await selectOption(page, "Connection mode", "Socket")
    await fillField(page, "Socket path", "/var/run/docker.sock")
    await delay(1200)

    await toggleSwitch(page, "Use container name as ID")

    await page.mouse.wheel(0, 260)
    await delay(1600)
}
