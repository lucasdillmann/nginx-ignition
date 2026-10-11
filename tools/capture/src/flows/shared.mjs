export function delay(milliseconds = 900) {
    return new Promise(resolve => setTimeout(resolve, milliseconds))
}

export async function clickButton(page, name) {
    await page.getByRole("button", { name, exact: true }).first().click()
}

export async function editRow(page, rowText) {
    const row = page.locator("tr", { hasText: rowText }).first()
    await row.locator("td").last().locator("button, a[href]").first().click()
}

export async function scrollTo(page, text) {
    await page.getByText(text, { exact: true }).first().scrollIntoViewIfNeeded()
}

function escapeForRegExp(value) {
    return value.replace(/[.*+?^${}()|[\]\\]/g, "\\$&")
}

function formItem(page, label) {
    const pattern = new RegExp(`^${escapeForRegExp(label)}(\\s*\\(optional\\))?$`)

    return page
        .locator("label")
        .filter({ hasText: pattern })
        .first()
        .locator("xpath=ancestor::div[contains(concat(' ', normalize-space(@class), ' '), ' ant-form-item ')][1]")
}

function controlOf(page, label) {
    return formItem(page, label).locator("input, textarea").first()
}

export async function fillField(page, label, value) {
    const field = controlOf(page, label)
    await field.scrollIntoViewIfNeeded()
    await field.click()
    await field.fill(value)
}

export async function selectOption(page, label, option) {
    const dropdown = formItem(page, label).locator(".ant-select").first()
    await dropdown.scrollIntoViewIfNeeded()
    await dropdown.click()
    await page.waitForTimeout(600)
    await page.locator(".ant-select-dropdown:visible .ant-select-item-option", { hasText: option }).first().click()
    await page.waitForTimeout(600)
}

export async function addTag(page, label, value) {
    const field = formItem(page, label).locator(".ant-select").first()
    await field.scrollIntoViewIfNeeded()
    await field.click()
    await page.waitForTimeout(500)
    await page.keyboard.type(value, { delay: 55 })
    await page.waitForTimeout(500)
    await page.keyboard.press("Enter")
    await page.waitForTimeout(500)
}

export async function toggleSwitch(page, label) {
    const control = formItem(page, label).getByRole("switch").first()
    await control.scrollIntoViewIfNeeded()
    await control.click()
    await delay(700)
}

export async function dismissReloadPrompt(page) {
    const skip = page.getByRole("button", { name: /don.t reload/i }).first()

    if ((await skip.count()) === 0) return

    await skip.click()
    await delay(1200)
}

export async function fillAddressPair(page, index, address, port) {
    const addressField = page.getByPlaceholder("Address").nth(index)
    await addressField.scrollIntoViewIfNeeded()
    await addressField.fill(address)
    await page.getByPlaceholder("Port").nth(index).fill(port)
    await delay(700)
}
