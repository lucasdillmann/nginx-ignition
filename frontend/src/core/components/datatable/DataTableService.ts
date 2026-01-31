import LocalStorageRepository from "../../repository/LocalStorageRepository"
import SessionStorageRepository from "../../repository/SessionStorageRepository"
import { DataTableInitialState } from "./model/DataTableInitialState"
import { DataTablePersistentStateMode } from "./model/DataTablePersistentStateMode"
import DataTablePersistentStateConfig from "./model/DataTablePersistentStateConfig"
import { DataTablePageSize } from "./model/DataTablePageSize"

const DEFAULT_PAGE_SIZE = 10

const DEFAULT_SHORT_TERM_DATA = {
    pageNumber: {
        enabled: true,
        tables: {},
    },
    searchTerms: {
        enabled: true,
        tables: {},
    },
} satisfies ShortTermData

const DEFAULT_LONG_TERM_DATA = {
    mode: DataTablePersistentStateMode.GLOBAL,
    global: {
        pageSize: DEFAULT_PAGE_SIZE,
    },
    pageSize: {},
} satisfies LongTermData

interface ShortTermData {
    pageNumber: {
        enabled: boolean
        tables: {
            [tableId: string]: number
        }
    }
    searchTerms: {
        enabled: boolean
        tables: {
            [tableId: string]: string | undefined
        }
    }
}

interface LongTermData {
    mode: DataTablePersistentStateMode
    global: {
        pageSize: DataTablePageSize
    }
    pageSize: {
        [tableId: string]: DataTablePageSize
    }
}

export default class DataTableService {
    private readonly longTermRepository: LocalStorageRepository<LongTermData>
    private readonly shortTermRepository: SessionStorageRepository<ShortTermData>

    constructor() {
        this.longTermRepository = new LocalStorageRepository("nginxIgnition.datatable.longTerm")
        this.shortTermRepository = new SessionStorageRepository("nginxIgnition.datatable.shortTerm")
    }

    public currentConfig(): DataTablePersistentStateConfig {
        const { mode, global } = this.longTermRepository.getOrDefault(DEFAULT_LONG_TERM_DATA)
        const { pageNumber, searchTerms } = this.shortTermRepository.getOrDefault(DEFAULT_SHORT_TERM_DATA)

        return {
            mode,
            pageSize: global.pageSize,
            persistPageNumber: pageNumber.enabled,
            persistSearchTerms: searchTerms.enabled,
        }
    }

    public updateConfig(config: DataTablePersistentStateConfig) {
        const longTerm = this.longTermRepository.getOrDefault(DEFAULT_LONG_TERM_DATA)
        const shortTerm = this.shortTermRepository.getOrDefault(DEFAULT_SHORT_TERM_DATA)

        this.longTermRepository.set({
            mode: config.mode,
            global: {
                pageSize: config.pageSize,
            },
            pageSize: longTerm.pageSize,
        })

        this.shortTermRepository.set({
            pageNumber: {
                enabled: config.persistPageNumber,
                tables: shortTerm.pageNumber.tables,
            },
            searchTerms: {
                enabled: config.persistSearchTerms,
                tables: shortTerm.searchTerms.tables,
            },
        })
    }

    public paginationChanged(tableId: string, pageSize: DataTablePageSize, pageNumber: number) {
        const shortTerm = this.shortTermRepository.getOrDefault(DEFAULT_SHORT_TERM_DATA)
        if (shortTerm.pageNumber.enabled) {
            this.shortTermRepository.set({
                ...shortTerm,
                pageNumber: {
                    ...shortTerm.pageNumber,
                    tables: {
                        ...shortTerm.pageNumber.tables,
                        [tableId]: pageNumber,
                    },
                },
            })
        }

        const longTerm = this.longTermRepository.getOrDefault(DEFAULT_LONG_TERM_DATA)
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
        const currentData = this.shortTermRepository.getOrDefault(DEFAULT_SHORT_TERM_DATA)
        if (!currentData.searchTerms.enabled) return

        this.shortTermRepository.set({
            ...currentData,
            searchTerms: {
                ...currentData.searchTerms,
                tables: {
                    ...currentData.searchTerms.tables,
                    [tableId]: searchTerms,
                },
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

        if (!shortTermData.pageNumber.enabled) pageNumber = 0
        else pageNumber = shortTermData.pageNumber.tables[tableId] ?? 0

        if (shortTermData.searchTerms.enabled) searchTerms = shortTermData.searchTerms.tables[tableId]

        return { searchTerms, pageSize, pageNumber }
    }
}
