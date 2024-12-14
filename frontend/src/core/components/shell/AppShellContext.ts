import { ButtonColorType, ButtonVariantType } from "antd/es/button"
import ContextHolder from "../../context/ContextHolder"

export interface ShellAction {
    description: string
    onClick: string | (() => Promise<void>) | (() => void)
    disabled?: boolean
    disabledReason?: string
    type?: ButtonVariantType
    color?: ButtonColorType
}

export interface ShellConfig {
    title?: string
    subtitle?: string
    actions?: ShellAction[]
}

export interface ShellOperations {
    updateConfig(config: ShellConfig): void
}

// eslint-disable-next-line import/no-anonymous-default-export
export default new ContextHolder<ShellOperations>({
    updateConfig: (_: ShellConfig) => {},
})
