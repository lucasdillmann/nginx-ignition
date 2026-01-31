import { DataTablePersistentStateMode } from "./DataTablePersistentStateMode"
import { DataTablePageSize } from "./DataTablePageSize"

export default interface DataTablePersistentStateConfig {
    mode: DataTablePersistentStateMode
    pageSize: DataTablePageSize
    persistPageNumber: boolean
    persistSearchTerms: boolean
}
