import { ButtonColorType, ButtonVariantType } from "antd/es/button"
import React from "react"

export interface ShellAction {
    description: string
    onClick: string | (() => Promise<void>) | (() => void)
    disabled?: boolean
    type?: ButtonVariantType
    color?: ButtonColorType
}

export interface ShellConfig {
    title: string
    subtitle?: string
    actions?: ShellAction[]
}

export interface ShellOperations {
    updateConfig(config: ShellConfig): void
}

const AppShellContext = React.createContext<ShellOperations>({
    updateConfig: (_: ShellConfig) => {},
})
export default AppShellContext
