import { UserAccessLevel } from "../../../domain/user/model/UserAccessLevel"
import UserPermissions from "../../../domain/user/model/UserPermissions"
import AppContext from "../context/AppContext"

export function isAccessGranted(
    requiredAccessLevel: UserAccessLevel,
    permissionResolver: (permissions: UserPermissions) => UserAccessLevel,
): boolean {
    const { user } = AppContext.get()
    if (!user) {
        return false
    }

    const currentAccessLevel = permissionResolver(user.permissions)
    switch (requiredAccessLevel) {
        case UserAccessLevel.READ_WRITE:
            return currentAccessLevel === UserAccessLevel.READ_WRITE
        case UserAccessLevel.READ_ONLY:
            return currentAccessLevel === UserAccessLevel.READ_WRITE || currentAccessLevel === UserAccessLevel.READ_ONLY
        default:
            return false
    }
}
