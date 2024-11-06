import Header from "./Header";

export default interface ApiResponse<T> {
    statusCode: number
    headers: Header[]
    body?: T | null
}
