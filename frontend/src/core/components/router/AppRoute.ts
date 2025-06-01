import React from "react"

export interface AppRouteMenuItem {
    icon: React.ReactNode
    description: string
}

export default interface AppRoute {
    path: string
    requiresAuthentication: boolean
    fullPage: boolean
    component: React.ReactElement
    menuItem?: AppRouteMenuItem
    activeMenuItemPath?: string
}
