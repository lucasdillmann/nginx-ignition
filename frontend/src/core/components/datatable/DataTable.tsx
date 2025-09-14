import PageResponse from "../../pagination/PageResponse"
import React from "react"
import { ColumnType, TablePaginationConfig } from "antd/es/table"
import { AlignType } from "rc-table/lib/interface"
import Preloader from "../preloader/Preloader"
import { Pagination, Table } from "antd"
import "./DataTable.css"
import DataTableSearchBar from "./DataTableSearchBar"
import CommonNotifications from "../notification/CommonNotifications"
import EmptyStates from "../emptystate/EmptyStates"

const DEFAULT_PAGE_SIZE = 10
const PAGE_SIZES = [10, 25, 50, 100, 250, 500]
const DEFAULT_DATA: PageResponse<any> = {
    pageSize: DEFAULT_PAGE_SIZE,
    pageNumber: 0,
    totalItems: 0,
    contents: [],
}

export interface DataTableColumn<T> {
    id: string
    description: string
    renderer: (row: T, index: number) => React.ReactNode
    width?: number
    minWidth?: number
    align?: AlignType
}

export interface DataTableProps<T> {
    columns: DataTableColumn<T>[]
    dataProvider: (pageSize: number, pageNumber: number, searchTerms?: string) => Promise<PageResponse<T>>
    rowKey: (row: T) => React.Key
    disableSearch?: boolean
}

interface DataTableState<T> {
    loading: boolean
    data: PageResponse<T>
    error?: Error
    searchTerms?: string
}

export default class DataTable<T> extends React.Component<DataTableProps<T>, DataTableState<T>> {
    constructor(props: DataTableProps<T>) {
        super(props)
        this.state = {
            loading: true,
            data: DEFAULT_DATA,
        }
    }

    private buildColumnAdapters(): ColumnType<T>[] {
        const { columns } = this.props
        return columns.map(column => ({
            key: column.id,
            title: column.description,
            render: (_, row, index) => column.renderer(row, index),
            width: column.width,
            minWidth: column.minWidth,
            align: column.align,
        }))
    }

    private changePage(pageSize: number, pageNumber: number) {
        this.setState({ loading: true }, () => this.fetchData(pageSize, pageNumber))
    }

    private fetchData(pageSize: number, pageNumber: number): Promise<void> {
        const { dataProvider } = this.props
        const { searchTerms } = this.state
        return dataProvider(pageSize, pageNumber, searchTerms)
            .then(data => this.setState({ loading: false, data }))
            .catch(error => {
                CommonNotifications.failedToFetch()
                this.setState({ loading: false, error })
            })
    }

    private buildPaginationProps(): TablePaginationConfig {
        const { data } = this.state

        return {
            align: "end",
            className: "pagination-container",
            defaultCurrent: data.pageNumber + 1,
            total: data.totalItems,
            pageSize: data.pageSize,
            pageSizeOptions: PAGE_SIZES,
            onChange: (pageNumber, pageSize) => this.changePage(pageSize, pageNumber - 1),
            showTotal: (total, [start, end]) => `Showing items ${start} to ${end} from a total of ${total}`,
            showSizeChanger: true,
            showQuickJumper: true,
            responsive: true,
        }
    }

    private handleSearchTerms(searchTerms?: string) {
        this.setState(
            {
                loading: true,
                data: DEFAULT_DATA,
                searchTerms,
            },
            () => this.refresh(),
        )
    }

    refresh(): Promise<void> {
        const { data } = this.state
        return this.fetchData(data.pageSize, data.pageNumber)
    }

    componentDidMount() {
        this.fetchData(DEFAULT_PAGE_SIZE, 0)
    }

    render() {
        const { loading, data, error } = this.state
        const { rowKey } = this.props

        if (error !== undefined) return EmptyStates.FailedToFetch

        return (
            <Preloader loading={loading}>
                <DataTableSearchBar onSearch={searchTerms => this.handleSearchTerms(searchTerms)} />
                <Table
                    className="data-table"
                    columns={this.buildColumnAdapters()}
                    dataSource={data.contents}
                    rowKey={row => rowKey(row)}
                    pagination={false}
                    tableLayout="fixed"
                    scroll={{
                        scrollToFirstRowOnChange: true,
                    }}
                    virtual
                    bordered
                />
                <Pagination {...this.buildPaginationProps()} />
            </Preloader>
        )
    }
}
