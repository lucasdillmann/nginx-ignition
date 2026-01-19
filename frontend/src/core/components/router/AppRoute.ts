import React from "react"
import { I18nMessage } from "../../i18n/I18n"

export interface AppRouteMenuItem {
    icon: React.ReactNode
    description: I18nMessage
}

export default interface AppRoute {
    path: string
    requiresAuthentication: boolean
    fullPage: boolean
    component: React.ReactElement
    menuItem?: AppRouteMenuItem
    activeMenuItemPath?: string
}
