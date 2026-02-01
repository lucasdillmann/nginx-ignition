import LocalStorageRepository from "../../repository/LocalStorageRepository"
import SessionStorageRepository from "../../repository/SessionStorageRepository"
import { DataTableInitialState } from "./model/DataTableInitialState"
import { DataTablePersistentStateMode } from "./model/DataTablePersistentStateMode"
import DataTablePersistentStateConfig from "./model/DataTablePersistentStateConfig"
import { DataTablePageSize } from "./model/DataTablePageSize"

const PERSISTENT_STORAGE_KEY = "nginxIgnition.datatable.preferences"
const DEFAULT_PAGE_SIZE = 10

const DEFAULT_SHORT_TERM_DATA = {
    pageNumber: {},
    searchTerms: {},
} satisfies ShortTermData

const DEFAULT_LONG_TERM_DATA = {
    mode: DataTablePersistentStateMode.GLOBAL,
    global: {
        pageSize: DEFAULT_PAGE_SIZE,
        rememberPageNumber: true,
        rememberSearchTerms: true,
    },
    pageSize: {},
} satisfies LongTermData

interface ShortTermData {
    pageNumber: {
        [tableId: string]: number
    }
    searchTerms: {
        [tableId: string]: string | undefined
    }
}

interface LongTermData {
    mode: DataTablePersistentStateMode
    global: {
        pageSize: DataTablePageSize
        rememberPageNumber: boolean
        rememberSearchTerms: boolean
    }
    pageSize: {
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
        const { mode, global } = this.longTermRepository.getOrDefault(DEFAULT_LONG_TERM_DATA)
        const { pageSize, rememberPageNumber, rememberSearchTerms } = global

        return { mode, pageSize, rememberPageNumber, rememberSearchTerms }
    }

    public updateConfig(config: DataTablePersistentStateConfig) {
        const longTerm = this.longTermRepository.getOrDefault(DEFAULT_LONG_TERM_DATA)
        const { mode, pageSize, rememberPageNumber, rememberSearchTerms } = config

        this.longTermRepository.set({
            mode,
            global: { pageSize, rememberPageNumber, rememberSearchTerms },
            pageSize: longTerm.pageSize,
        })
    }

    public paginationChanged(tableId: string, pageSize: DataTablePageSize, pageNumber: number) {
        const longTerm = this.longTermRepository.getOrDefault(DEFAULT_LONG_TERM_DATA)

        if (longTerm.global.rememberPageNumber) {
            const shortTerm = this.shortTermRepository.getOrDefault(DEFAULT_SHORT_TERM_DATA)
            this.shortTermRepository.set({
                ...shortTerm,
                pageNumber: {
                    ...shortTerm.pageNumber,
                    [tableId]: pageNumber,
                },
            })
        }

        switch (longTerm.mode) {
            case DataTablePersistentStateMode.GLOBAL:
                this.longTermRepository.set({
                    ...longTerm,
                    global: {
                        ...longTerm.global,
                        pageSize,
                    },
                })
                break

            case DataTablePersistentStateMode.BY_TABLE:
                this.longTermRepository.set({
                    ...longTerm,
                    pageSize: {
                        ...longTerm.pageSize,
                        [tableId]: pageSize,
                    },
                })
                break
        }
    }

    public searchTermsChanged(tableId: string, searchTerms?: string) {
        const longTerm = this.longTermRepository.getOrDefault(DEFAULT_LONG_TERM_DATA)
        if (!longTerm.global.rememberSearchTerms) return

        const shortTerm = this.shortTermRepository.getOrDefault(DEFAULT_SHORT_TERM_DATA)
        this.shortTermRepository.set({
            ...shortTerm,
            searchTerms: {
                ...shortTerm.searchTerms,
                [tableId]: searchTerms,
            },
        })
    }

    public getInitialState(tableId: string): DataTableInitialState {
        const longTermData = this.longTermRepository.getOrDefault(DEFAULT_LONG_TERM_DATA)
        const shortTermData = this.shortTermRepository.getOrDefault(DEFAULT_SHORT_TERM_DATA)

        let pageSize, pageNumber: number
        let searchTerms: string | undefined

        if (longTermData.mode !== DataTablePersistentStateMode.BY_TABLE) pageSize = longTermData.global.pageSize
        else pageSize = longTermData.pageSize[tableId] ?? longTermData.global.pageSize

        if (!longTermData.global.rememberPageNumber) pageNumber = 0
        else pageNumber = shortTermData.pageNumber[tableId] ?? 0

        if (longTermData.global.rememberSearchTerms) searchTerms = shortTermData.searchTerms[tableId]

        return { searchTerms, pageSize, pageNumber }
    }
}
