import { ButtonColorType, ButtonVariantType } from "antd/es/button"
import ContextHolder from "../../context/ContextHolder"
import { I18nMessage } from "../../i18n/I18n"

export interface ShellAction {
    description: I18nMessage
    onClick: string | (() => Promise<void>) | (() => void)
    disabled?: boolean
    disabledReason?: string
    type?: ButtonVariantType
    color?: ButtonColorType
}

export interface ShellConfig {
    title?: I18nMessage
    subtitle?: I18nMessage
    actions?: ShellAction[]
    noContainerPadding?: boolean
}

export interface ShellOperations {
    updateConfig(config: ShellConfig): void
}

export default new ContextHolder<ShellOperations>({
    updateConfig: (_: ShellConfig) => {},
})
