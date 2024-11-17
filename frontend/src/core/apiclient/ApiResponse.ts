import Header from "./Header";

export default interface ApiResponse<T> {
    statusCode: number
    headers: Header[]
    body?: T
}

export class UnexpectedResponseError<T> extends Error {
    readonly response: ApiResponse<T>

    constructor(response: ApiResponse<T>) {
        super();
        this.response = response;
    }
}

export function requireSuccessResponse<T>(response: ApiResponse<T>): T | undefined {
    if (response.statusCode < 200 || response.statusCode > 299)
        throw new UnexpectedResponseError(response)

    return response.body
}

export function requireSuccessPayload<T>(response: ApiResponse<T>): T {
    const payload = requireSuccessResponse(response)
    if (payload == null)
        throw new UnexpectedResponseError(response)

    return payload
}
