import React from "react"
import PageResponse from "../../pagination/PageResponse"
import { Select } from "antd"
import { InputStatus } from "antd/es/_util/statusUtils"
import Notification from "../notification/Notification"
import { LoadingOutlined } from "@ant-design/icons"
import debounce from "debounce"
import { I18n, I18nMessage } from "../../i18n/I18n"
import MessageKey from "../../i18n/model/MessageKey.generated"

const PAGE_SIZE = 10

const INITIAL_STATE: PaginatedSelectState<any> = {
    data: [],
    firstLoadComplete: false,
    loading: false,
    hasMorePages: true,
    nextPageNumber: 0,
}

interface SelectOption<T> {
    item?: T
    value: string
    label: React.ReactNode
    disabled?: boolean
}

interface PaginatedSelectState<T> {
    data: T[]
    firstLoadComplete: boolean
    loading: boolean
    hasMorePages: boolean
    nextPageNumber: number
    searchTerms?: string
}

export interface PaginatedSelectProps<T> {
    placeholder?: I18nMessage
    onChange?: (selected?: T) => void
    pageProvider: (pageSize: number, pageNumber: number, searchTerms?: string) => Promise<PageResponse<T>>
    itemKey: (item: T) => string
    itemDescription: (item: T) => React.ReactNode
    disableSearch?: boolean
    disabled?: boolean
    allowEmpty?: boolean
    status?: InputStatus
    value?: T
    autoFocus?: boolean
}

export default class PaginatedSelect<T> extends React.Component<PaginatedSelectProps<T>, PaginatedSelectState<T>> {
    constructor(props: PaginatedSelectProps<T>) {
        super(props)
        this.state = {
            ...INITIAL_STATE,
        }
    }

    reset() {
        this.setState(INITIAL_STATE)
    }

    private loadNextPage() {
        const { loading } = this.state
        if (loading) return

        this.setState({ loading: true }, () => this.fetchNextPage())
    }

    private fetchNextPage() {
        const { pageProvider } = this.props
        const { nextPageNumber, searchTerms } = this.state

        pageProvider(PAGE_SIZE, nextPageNumber, searchTerms)
            .then(page => {
                this.setState(current => ({
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
                    MessageKey.CommonUnableToFetchOptionsTitle,
                    MessageKey.CommonUnableToFetchOptionsDescription,
                )

                this.setState({ loading: false })
            })
    }

    private handleDropdownVisibilityChange(open: boolean) {
        if (!open || this.state.firstLoadComplete) return
        this.loadNextPage()
    }

    private buildOption(item: T): SelectOption<T> {
        const { itemDescription, itemKey } = this.props

        return {
            value: itemKey(item),
            label: itemDescription(item),
            item,
        }
    }

    private buildOptions(): SelectOption<T>[] {
        const { data, loading } = this.state
        const options = data.map(item => this.buildOption(item))

        if (!loading) return options

        const loadingLabel = (
            <>
                <LoadingOutlined style={{ fontSize: 14, marginRight: 8 }} />
                Loading...
            </>
        )
        return [...options, { value: "", label: loadingLabel, disabled: true }]
    }

    private readonly handleSearch = debounce((searchTerms?: string) => {
        this.setState(
            {
                nextPageNumber: 0,
                loading: true,
                hasMorePages: true,
                data: [],
                searchTerms,
            },
            () => this.fetchNextPage(),
        )
    }, 500)

    private handleScrollEvent(event: React.UIEvent<HTMLDivElement>) {
        const { loading, hasMorePages } = this.state
        if (loading || !hasMorePages) return

        const target = event.target as HTMLDivElement
        if (target.scrollTop + target.offsetHeight === target.scrollHeight) {
            this.loadNextPage()
        }
    }

    render() {
        const { allowEmpty, disabled, status, value, onChange, placeholder, disableSearch, autoFocus } = this.props
        return (
            <Select<SelectOption<T>>
                showSearch={disableSearch !== true}
                placeholder={placeholder ? <I18n id={placeholder} /> : undefined}
                disabled={disabled}
                allowClear={allowEmpty}
                options={this.buildOptions()}
                status={status}
                value={value !== undefined && value !== null ? this.buildOption(value) : undefined}
                onChange={(_, option) => onChange?.((option as SelectOption<T>)?.item)}
                onDeselect={() => onChange?.(undefined)}
                onPopupScroll={event => this.handleScrollEvent(event)}
                onOpenChange={open => this.handleDropdownVisibilityChange(open)}
                onSearch={searchTerms => this.handleSearch(searchTerms)}
                filterOption={false}
                autoFocus={autoFocus}
            />
        )
    }
}
