import { DataTablePersistentStateMode } from "./DataTablePersistentStateMode"
import { DataTablePageSize } from "./DataTablePageSize"

export default interface DataTablePersistentStateConfig {
    paginationMode: DataTablePersistentStateMode
    defaultPageSize: DataTablePageSize
    rememberPageNumber: boolean
    rememberSearchTerms: boolean
}
