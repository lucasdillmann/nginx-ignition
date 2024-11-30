import { UserRole } from "./UserRole"

export default interface UserResponse {
    id: string
    enabled: boolean
    name: string
    username: string
    role: UserRole
}
