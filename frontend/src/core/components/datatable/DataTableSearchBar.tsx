import React from "react"
import { Flex, Input } from "antd"
import debounce from "debounce"
import { SearchOutlined } from "@ant-design/icons"

export interface SearchBarProps {
    onSearch: (searchTerms?: string) => void
}

export default class DataTableSearchBar extends React.Component<SearchBarProps> {
    render() {
        const handleChange = debounce((searchTerms?: string) => this.props.onSearch(searchTerms), 500)

        return (
            <Flex className="data-table-search-bar-container">
                <Input
                    // @ts-ignore
                    onInput={event => handleChange(event.nativeEvent.target!!.value)}
                    onClear={() => handleChange()}
                    placeholder="Search terms"
                    className="data-table-search-bar"
                    autoFocus
                    allowClear
                />
                <SearchOutlined style={{ fontSize: 18, margin: 10, marginLeft: 15 }} />
            </Flex>
        )
    }
}
