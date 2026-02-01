import LocalStorageRepository from "../../repository/LocalStorageRepository"
import SessionStorageRepository from "../../repository/SessionStorageRepository"
import { DataTableInitialState } from "./model/DataTableInitialState"
import { DataTablePersistentStateMode } from "./model/DataTablePersistentStateMode"
import DataTablePersistentStateConfig from "./model/DataTablePersistentStateConfig"
import { DataTablePageSize } from "./model/DataTablePageSize"

const PERSISTENT_STORAGE_KEY = "nginxIgnition.datatable.preferences"
const DEFAULT_PAGE_SIZE = 10

const DEFAULT_SHORT_TERM_DATA = {
    pageNumberByTable: {},
    searchTermsByTable: {},
} satisfies ShortTermData

const DEFAULT_LONG_TERM_DATA = {
    paginationMode: DataTablePersistentStateMode.GLOBAL,
    defaultPageSize: DEFAULT_PAGE_SIZE,
    rememberPageNumber: true,
    rememberSearchTerms: true,
    pageSizeByTable: {},
} satisfies LongTermData

interface ShortTermData {
    pageNumberByTable: {
        [tableId: string]: number
    }
    searchTermsByTable: {
        [tableId: string]: string | undefined
    }
}

interface LongTermData {
    paginationMode: DataTablePersistentStateMode
    defaultPageSize: DataTablePageSize
    rememberPageNumber: boolean
    rememberSearchTerms: boolean
    pageSizeByTable: {
        [tableId: string]: DataTablePageSize
    }
}

export default class DataTableService {
    private readonly longTermRepository: LocalStorageRepository<LongTermData>
    private readonly shortTermRepository: SessionStorageRepository<ShortTermData>

    constructor() {
        this.longTermRepository = new LocalStorageRepository(PERSISTENT_STORAGE_KEY)
        this.shortTermRepository = new SessionStorageRepository(PERSISTENT_STORAGE_KEY)
    }

    public currentConfig(): DataTablePersistentStateConfig {
        const { paginationMode, defaultPageSize, rememberPageNumber, rememberSearchTerms } =
            this.longTermRepository.getOrDefault(DEFAULT_LONG_TERM_DATA)

        return { paginationMode, defaultPageSize, rememberPageNumber, rememberSearchTerms }
    }

    public updateConfig(config: DataTablePersistentStateConfig) {
        const longTerm = this.longTermRepository.getOrDefault(DEFAULT_LONG_TERM_DATA)
        const { paginationMode, defaultPageSize, rememberPageNumber, rememberSearchTerms } = config

        this.longTermRepository.set({
            ...longTerm,
            paginationMode,
            defaultPageSize,
            rememberPageNumber,
            rememberSearchTerms,
        })
    }

    public paginationChanged(tableId: string, pageSize: DataTablePageSize, pageNumber: number) {
        const longTerm = this.longTermRepository.getOrDefault(DEFAULT_LONG_TERM_DATA)
        const { rememberPageNumber, paginationMode } = longTerm

        if (rememberPageNumber) {
            const shortTerm = this.shortTermRepository.getOrDefault(DEFAULT_SHORT_TERM_DATA)
            this.shortTermRepository.set({
                ...shortTerm,
                pageNumberByTable: {
                    ...shortTerm.pageNumberByTable,
                    [tableId]: pageNumber,
                },
            })
        }

        switch (paginationMode) {
            case DataTablePersistentStateMode.GLOBAL:
                this.longTermRepository.set({
                    ...longTerm,
                    defaultPageSize: pageSize,
                })
                break

            case DataTablePersistentStateMode.BY_TABLE:
                this.longTermRepository.set({
                    ...longTerm,
                    pageSizeByTable: {
                        ...longTerm.pageSizeByTable,
                        [tableId]: pageSize,
                    },
                })
                break
        }
    }

    public searchTermsChanged(tableId: string, searchTerms?: string) {
        const longTerm = this.longTermRepository.getOrDefault(DEFAULT_LONG_TERM_DATA)
        if (!longTerm.rememberSearchTerms) return

        const shortTerm = this.shortTermRepository.getOrDefault(DEFAULT_SHORT_TERM_DATA)
        this.shortTermRepository.set({
            ...shortTerm,
            searchTermsByTable: {
                ...shortTerm.searchTermsByTable,
                [tableId]: searchTerms,
            },
        })
    }

    public getInitialState(tableId: string): DataTableInitialState {
        const longTermData = this.longTermRepository.getOrDefault(DEFAULT_LONG_TERM_DATA)
        const shortTermData = this.shortTermRepository.getOrDefault(DEFAULT_SHORT_TERM_DATA)

        let pageSize, pageNumber: number
        let searchTerms: string | undefined

        if (longTermData.paginationMode !== DataTablePersistentStateMode.BY_TABLE)
            pageSize = longTermData.defaultPageSize
        else pageSize = longTermData.pageSizeByTable[tableId] ?? longTermData.defaultPageSize

        if (!longTermData.rememberPageNumber) pageNumber = 0
        else pageNumber = shortTermData.pageNumberByTable[tableId] ?? 0

        if (longTermData.rememberSearchTerms) searchTerms = shortTermData.searchTermsByTable[tableId]

        return { searchTerms, pageSize, pageNumber }
    }
}
