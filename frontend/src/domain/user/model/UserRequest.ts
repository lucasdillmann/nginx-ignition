import UserPermissions from "./UserPermissions"

export default interface UserRequest {
    enabled: boolean
    name: string
    username: string
    password?: string
    permissions: UserPermissions
}
