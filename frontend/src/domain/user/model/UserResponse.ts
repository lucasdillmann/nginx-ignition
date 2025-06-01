import UserPermissions from "./UserPermissions"

export default interface UserResponse {
    id: string
    enabled: boolean
    name: string
    username: string
    permissions: UserPermissions
}
