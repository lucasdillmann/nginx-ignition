import LocalStorageRepository from "../../repository/LocalStorageRepository"
import { theme } from "antd"

type ThemeVariant = "light" | "dark"

export type ThemeListener = (darkMode: boolean) => void

class ThemeContext {
    private readonly repository: LocalStorageRepository<ThemeVariant>
    private listeners: ThemeListener[]
    private current: ThemeVariant

    constructor() {
        this.repository = new LocalStorageRepository("nginxIgnition.theme")
        this.listeners = []

        const prefersDarkMode = window.matchMedia && window.matchMedia("(prefers-color-scheme: dark)").matches
        this.current = this.repository.getOrDefault(prefersDarkMode ? "dark" : "light")
    }

    register(listener: ThemeListener) {
        this.listeners.push(listener)
    }

    deregister(listener: ThemeListener) {
        this.listeners = this.listeners.filter(element => element !== listener)
    }

    isDarkMode(): boolean {
        return this.current == "dark"
    }

    toggle() {
        this.current = this.isDarkMode() ? "light" : "dark"
        this.repository.set(this.current)

        this.listeners.forEach(listener => {
            listener(this.isDarkMode())
        })
    }

    algorithm() {
        return this.isDarkMode() ? theme.darkAlgorithm : theme.defaultAlgorithm
    }
}

export default new ThemeContext()
