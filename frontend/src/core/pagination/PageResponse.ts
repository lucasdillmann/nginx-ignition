export function emptyPageResponse<T>(): PageResponse<T> {
    return {
        pageSize: 0,
        pageNumber: 0,
        totalItems: 0,
        contents: [],
    }
}

export default interface PageResponse<T> {
    pageSize: number
    pageNumber: number
    totalItems: number
    contents: T[]
}
