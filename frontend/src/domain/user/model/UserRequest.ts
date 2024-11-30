import { UserRole } from "./UserRole"

export default interface UserRequest {
    enabled: boolean
    name: string
    username: string
    password?: string
    role: UserRole
}
