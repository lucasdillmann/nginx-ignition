import PageResponse from "../../pagination/PageResponse"
import React from "react"
import { ColumnType, TablePaginationConfig } from "antd/es/table"
import { AlignType } from "rc-table/lib/interface"
import Preloader from "../preloader/Preloader"
import { Pagination, Table } from "antd"
import "./DataTable.css"
import DataTableHeader from "./DataTableHeader"
import CommonNotifications from "../notification/CommonNotifications"
import EmptyStates from "../emptystate/EmptyStates"
import { I18n, i18n, I18nMessage } from "../../i18n/I18n"
import MessageKey from "../../i18n/model/MessageKey.generated"
import DataTableService from "./DataTableService"
import { DataTableInitialState } from "./model/DataTableInitialState"
import { DATA_TABLE_PAGE_SIZES, DataTablePageSize, DEFAULT_PAGE_SIZE } from "./model/DataTablePageSize"

const DEFAULT_DATA: PageResponse<any> = {
    pageSize: DEFAULT_PAGE_SIZE,
    pageNumber: 0,
    totalItems: 0,
    contents: [],
}

export interface DataTableColumn<T> {
    id: string
    description: I18nMessage
    renderer: (row: T, index: number) => React.ReactNode
    width?: number
    minWidth?: number
    align?: AlignType
}

export interface DataTableProps<T> {
    id: string
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
    initialValues: DataTableInitialState
}

export default class DataTable<T> extends React.Component<DataTableProps<T>, DataTableState<T>> {
    private readonly service: DataTableService

    constructor(props: DataTableProps<T>) {
        super(props)

        this.service = new DataTableService()
        const initialValues = this.service.getInitialState(props.id)

        this.state = {
            loading: true,
            data: DEFAULT_DATA,
            searchTerms: initialValues.searchTerms,
            initialValues,
        }
    }

    private buildColumnAdapters(): ColumnType<T>[] {
        const { columns } = this.props
        return columns.map(column => ({
            key: column.id,
            title: i18n(column.description),
            render: (_, row, index) => column.renderer(row, index),
            width: column.width,
            minWidth: column.minWidth,
            align: column.align,
        }))
    }

    private changePage(pageSize: DataTablePageSize, pageNumber: number) {
        const { id } = this.props

        this.service.paginationChanged(id, pageSize, pageNumber)
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
            pageSizeOptions: DATA_TABLE_PAGE_SIZES,
            onChange: (pageNumber, pageSize) => this.changePage(pageSize as DataTablePageSize, pageNumber - 1),
            showTotal: (total, [start, end]) => (
                <I18n id={MessageKey.FrontendComponentsDatatablePaginationSummary} params={{ start, end, total }} />
            ),
            showSizeChanger: true,
            showQuickJumper: true,
            responsive: true,
        }
    }

    private handleSearchTerms(searchTerms?: string) {
        const { id } = this.props

        this.service.searchTermsChanged(id, searchTerms)
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
        const { initialValues } = this.state
        this.fetchData(initialValues.pageSize, initialValues.pageNumber)
    }

    render() {
        const { loading, data, error, initialValues } = this.state
        const { id, rowKey } = this.props

        if (error !== undefined) return EmptyStates.FailedToFetch

        return (
            <Preloader loading={loading}>
                <DataTableHeader
                    id={id}
                    initialSearchTerms={initialValues.searchTerms}
                    onSearch={searchTerms => this.handleSearchTerms(searchTerms)}
                />
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
