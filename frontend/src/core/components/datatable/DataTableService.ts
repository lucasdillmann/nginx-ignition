import LocalStorageRepository from "../../repository/LocalStorageRepository"
import SessionStorageRepository from "../../repository/SessionStorageRepository"
import { DataTableInitialState } from "./model/DataTableInitialState"

enum Mode {
    GLOBAL = "GLOBAL",
    BY_TABLE = "BY_TABLE",
    DISABLED = "DISABLED",
}

const DEFAULT_PAGE_SIZE = 10

const DEFAULT_SHORT_TERM_DATA: ShortTermData = {
    pageNumber: {
        enabled: true,
        tables: {},
    },
    searchTerms: {
        enabled: true,
        tables: {},
    },
}

const DEFAULT_LONG_TERM_DATA: LongTermData = {
    mode: Mode.GLOBAL,
    global: {
        pageSize: DEFAULT_PAGE_SIZE,
    },
    pageSize: {},
}

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
    mode: Mode
    global: {
        pageSize: number
    }
    pageSize: {
        [tableId: string]: number
    }
}

export default class DataTableService {
    private readonly longTermRepository: LocalStorageRepository<LongTermData>
    private readonly shortTermRepository: SessionStorageRepository<ShortTermData>

    constructor() {
        this.longTermRepository = new LocalStorageRepository("nginxIgnition.datatable.longTerm")
        this.shortTermRepository = new SessionStorageRepository("nginxIgnition.datatable.shortTerm")
    }

    public paginationChanged(tableId: string, pageSize: number, pageNumber: number) {
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
            case Mode.GLOBAL:
                this.longTermRepository.set({
                    ...longTerm,
                    global: {
                        ...longTerm.global,
                        pageSize,
                    },
                })
                break

            case Mode.BY_TABLE:
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

        switch (longTermData.mode) {
            case Mode.DISABLED:
                pageSize = DEFAULT_PAGE_SIZE
                break
            case Mode.GLOBAL:
                pageSize = longTermData.global.pageSize
                break
            case Mode.BY_TABLE:
                pageSize = longTermData.pageSize[tableId] ?? longTermData.global.pageSize
                break
        }

        if (shortTermData.pageNumber.enabled) pageNumber = shortTermData.pageNumber.tables[tableId] ?? 0
        else pageNumber = 0

        if (shortTermData.searchTerms.enabled) searchTerms = shortTermData.searchTerms.tables[tableId]

        return { searchTerms, pageSize, pageNumber }
    }
}
