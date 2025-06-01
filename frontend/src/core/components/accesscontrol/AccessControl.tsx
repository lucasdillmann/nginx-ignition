import { UserAccessLevel } from "../../../domain/user/model/UserAccessLevel"
import UserPermissions from "../../../domain/user/model/UserPermissions"
import React from "react"
import { isAccessGranted } from "./IsAccessGranted"
import AccessDeniedPage from "./AccessDeniedPage"

export interface AccessControlProps {
    children: React.ReactNode
    requiredAccessLevel: UserAccessLevel
    permissionResolver: (permissions: UserPermissions) => UserAccessLevel
}

export default class AccessControl extends React.Component<AccessControlProps> {
    render() {
        const { children, requiredAccessLevel, permissionResolver } = this.props
        if (!isAccessGranted(requiredAccessLevel, permissionResolver)) {
            return <AccessDeniedPage />
        }

        return children
    }
}
