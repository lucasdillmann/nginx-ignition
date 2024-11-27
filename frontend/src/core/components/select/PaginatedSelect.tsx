import React from "react";
import PageResponse from "../../pagination/PageResponse";
import {Select} from "antd";
import {InputStatus} from "antd/es/_util/statusUtils";
import Notification from "../notification/Notification";

const PAGE_SIZE = 10

interface SelectOption<T> {
    item: T,
    value: string,
    label: string,
    disabled?: boolean
}

interface PaginatedSelectState<T> {
    data: T[]
    firstLoadComplete: boolean
    loading: boolean
    hasMorePages: boolean
    nextPageNumber: number
}

export interface PaginatedSelectProps<T> {
    placeholder?: string
    onChange?: (selected?: T) => void
    pageProvider: (pageSize: number, pageNumber: number) => Promise<PageResponse<T>>
    itemKey: (item: T) => string
    itemDescription: (item: T) => string
    disabled?: boolean
    allowEmpty?: boolean
    status?: InputStatus
    value?: T
}

export default class PaginatedSelect<T> extends React.Component<PaginatedSelectProps<T>, PaginatedSelectState<T>> {
    constructor(props: PaginatedSelectProps<T>) {
        super(props);
        this.state = {
            data: [],
            firstLoadComplete: false,
            loading: false,
            hasMorePages: true,
            nextPageNumber: 0,
        }
    }

    private loadNextPage() {
        const {loading} = this.state
        if (loading) return

        this.setState(
            {loading: true},
            () => this.fetchNextPage(),
        )
    }

    private fetchNextPage() {
        const {pageProvider} = this.props
        const {nextPageNumber} = this.state

        pageProvider(PAGE_SIZE, nextPageNumber)
            .then(page => {
                this.setState((current) => ({
                    ...current,
                    data: [...current.data, ...page.contents],
                    nextPageNumber: page.pageNumber + 1,
                    hasMorePages: current.data.length + page.contents.length < page.totalItems,
                    loading: false,
                    firstLoadComplete: true,
                }))
            })
            .catch(() => {
                Notification.error(
                    "Unable to fetch the next options",
                    "We're unable to fetch the next available options to select at this time. Please try again later.",
                )

                this.setState({ loading: false })
            })
    }

    private handleDropdownVisibilityChange(open: boolean) {
        if (!open || this.state.firstLoadComplete) return
        this.loadNextPage()
    }

    private buildOption(item: T): SelectOption<T> {
        const {itemDescription, itemKey} = this.props

        return {
            value: itemKey(item),
            label: itemDescription(item),
            item,
        }
    }

    private buildOptions() {
        const {data, loading} = this.state
        const options = data.map(item => this.buildOption(item))

        if (!loading)
            return options
        else
            return [...options, { value: "", label: "Loading...", disabled: true }]
    }

    private handleScrollEvent(event: React.UIEvent<HTMLDivElement>) {
        const {loading, hasMorePages} = this.state
        if (loading || !hasMorePages) return

        const target = event.target as HTMLDivElement
        if (target.scrollTop + target.offsetHeight === target.scrollHeight) {
            this.loadNextPage()
        }
    }

    render() {
        const {allowEmpty, disabled, status, value, onChange, placeholder} = this.props
        const {loading} = this.state

        return (
            <Select
                placeholder={placeholder}
                disabled={disabled}
                allowClear={allowEmpty}
                options={this.buildOptions()}
                status={status}
                value={value !== undefined && value !== null ? this.buildOption(value) : undefined}
                onChange={(_, option) => onChange?.((option as SelectOption<T>).item)}
                onDeselect={() => onChange?.(undefined)}
                onPopupScroll={event => this.handleScrollEvent(event)}
                onDropdownVisibleChange={open => this.handleDropdownVisibilityChange(open)}
                loading={loading}
            />
        )
    }
}
