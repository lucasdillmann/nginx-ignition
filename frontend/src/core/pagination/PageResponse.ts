export default interface PageResponse<T> {
    pageSize: number
    pageNumber: number
    totalItems: number
    contents: T[]
}
