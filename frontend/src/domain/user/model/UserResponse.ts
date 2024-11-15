export default interface UserResponse {
    id: string
    enabled: boolean
    name: string
    username: string
    role: UserRole
}

export enum UserRole {
    ADMIN = "ADMIN",
    REGULAR_USER = "REGULAR_USER",
}
