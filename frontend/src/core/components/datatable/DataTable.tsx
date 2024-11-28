import PageResponse from "../../pagination/PageResponse";
import React from "react";
import {ColumnType, TablePaginationConfig} from "antd/es/table";
import {AlignType} from "rc-table/lib/interface";
import Preloader from "../preloader/Preloader";
import {Empty, Pagination, Table} from "antd";
import "./DataTable.css"
import {ExclamationCircleOutlined} from "@ant-design/icons";
import Notification from "../notification/Notification";

const DEFAULT_PAGE_SIZE = 10
const PAGE_SIZES = [10, 25, 50, 100, 250, 500]

export interface DataTableColumn<T> {
    id: string
    description: string
    renderer: (row: T, index: number) => React.ReactNode
    width?: number
    minWidth?: number
    align?: AlignType,
}

export interface DataTableProps<T> {
    columns: DataTableColumn<T>[]
    dataProvider: (pageSize: number, pageNumber: number) => Promise<PageResponse<T>>
    rowKey: (row: T) => React.Key
}

interface DataTableState<T> {
    loading: boolean
    data: PageResponse<T>
    error?: Error
}

export default class DataTable<T> extends React.Component<DataTableProps<T>, DataTableState<T>> {
    constructor(props: DataTableProps<T>) {
        super(props);
        this.state = {
            loading: true,
            data: {
                pageSize: DEFAULT_PAGE_SIZE,
                pageNumber: 0,
                totalItems: 0,
                contents: [],
            }
        }
    }

    private buildColumnAdapters(): ColumnType<T>[] {
        const {columns} = this.props
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
        this.setState(
            { loading: true,},
            () => this.fetchData(pageSize, pageNumber),
        )
    }

    private fetchData(pageSize: number, pageNumber: number): Promise<void> {
        const {dataProvider} = this.props
        return dataProvider(pageSize, pageNumber)
            .then(data => this.setState({ loading: false, data }))
            .catch(error =>{
                Notification.error(
                    "Unable to fetch the data",
                    "We're unable to fetch the data at this time. Please try again later.",
                )

                this.setState({ loading: false, error })
            })
    }

    private buildPaginationProps(): TablePaginationConfig {
        const {data} = this.state

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

    refresh(): Promise<void> {
        const {data} = this.state
        return this.fetchData(data.pageSize, data.pageNumber)
    }

    componentDidMount() {
        this.fetchData(DEFAULT_PAGE_SIZE, 0)
    }

    render() {
        const {loading, data, error} = this.state
        const {rowKey} = this.props

        if (error !== undefined)
            return (
                <Empty
                    image={<ExclamationCircleOutlined style={{ fontSize: 70, color: "#b8b8b8" }} />}
                    description="Unable to fetch the data. Please try again later."
                />
            )

        return (
            <Preloader loading={loading}>
                <Table
                    columns={this.buildColumnAdapters()}
                    dataSource={data.contents}
                    rowKey={row => rowKey(row)}
                    pagination={false}
                    tableLayout="fixed"
                    virtual
                    bordered
                />
                <Pagination
                    {...this.buildPaginationProps()}
                />
            </Preloader>
        )
    }
}
