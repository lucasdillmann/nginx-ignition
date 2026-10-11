import { clickButton, delay, dismissReloadPrompt, fillAddressPair, fillField } from "./shared.mjs"

export default async function streams(page) {
    await clickButton(page, "New stream")
    await delay(1800)

    await fillField(page, "Name", "PostgreSQL")
    await delay(1200)

    await fillAddressPair(page, 0, "postgres.example.com", "5432")
    await fillAddressPair(page, 1, "0.0.0.0", "15432")

    await clickButton(page, "Save")
    await delay(2200)
    await dismissReloadPrompt(page)
    await delay(1200)
}
