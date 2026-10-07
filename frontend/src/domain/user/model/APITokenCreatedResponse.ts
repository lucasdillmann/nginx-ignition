import APITokenResponse from "./APITokenResponse"

export default interface APITokenCreatedResponse extends APITokenResponse {
    token: string
}
