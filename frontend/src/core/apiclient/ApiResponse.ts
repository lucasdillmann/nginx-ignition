import Header from "./Header";

export default interface ApiResponse<T> {
    statusCode: number
    headers: Header[]
    body?: T
}

export function requireSuccessResponse<T>(response: ApiResponse<T>): T {
    if (response.statusCode < 200 || response.statusCode > 299)
        throw Error(`Unexpected status code: ${response.statusCode}`)

    if (response.body == null)
        throw Error("Null or empty response body")

    return response.body
}
