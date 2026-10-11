import { delay, editRow, scrollTo } from "./shared.mjs"

export default async function sslCertificates(page) {
    await editRow(page, "example.com")
    await delay(1800)

    await scrollTo(page, "Validity")
    await delay(2000)
}
