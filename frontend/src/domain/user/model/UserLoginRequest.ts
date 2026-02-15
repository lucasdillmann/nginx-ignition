export default interface UserLoginRequest {
    username: string
    password: string
    totp?: string
}
