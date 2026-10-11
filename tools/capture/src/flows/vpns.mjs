import { clickButton, delay, fillField, scrollTo } from "./shared.mjs"

export default async function vpns(page) {
    await clickButton(page, "New connection")
    await delay(1800)

    await fillField(page, "Name", "Corporate network")
    await delay(1200)

    await fillField(page, "NetBird setup key", "nbp_example_setup_key")
    await delay(1200)

    await scrollTo(page, "Important instructions")
    await delay(1800)
}
