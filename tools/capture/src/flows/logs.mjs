import { delay } from "./shared.mjs"

export default async function logs(page) {
    await page.getByText("Server logs", { exact: true }).click()
    await delay(3200)
}
