import React from "react";
import {UserRole} from "../../../domain/user/model/UserRole";

export interface AppRouteMenuItem {
    icon: React.ReactNode
    description: string
}

export default interface AppRoute {
    path: string
    title?: string
    requiresAuthentication: boolean
    fullPage: boolean
    component: React.ReactElement
    menuItem?: AppRouteMenuItem
    visibleRoles?: UserRole[]
    activeMenuItemPath?: string
}
