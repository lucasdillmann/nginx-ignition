import React from "react"
import { Flex, Input } from "antd"
import debounce from "debounce"
import { SearchOutlined } from "@ant-design/icons"
import { i18n } from "../../i18n/I18n"
import MessageKey from "../../i18n/model/MessageKey.generated"

export interface SearchBarProps {
    onSearch: (searchTerms?: string) => void
}

export default class DataTableSearchBar extends React.Component<SearchBarProps> {
    render() {
        const handleChange = debounce((searchTerms?: string) => this.props.onSearch(searchTerms), 500)

        return (
            <Flex className="data-table-search-bar-container">
                <Input
                    // @ts-expect-error target is generic, but in this scenario is safe to use the value attribute
                    onInput={event => handleChange(event.nativeEvent.target!!.value)}
                    onClear={() => handleChange()}
                    placeholder={i18n(MessageKey.FrontendComponentsDatatableSearchPlaceholder)}
                    className="data-table-search-bar"
                    autoFocus
                    allowClear
                />
                <SearchOutlined style={{ fontSize: 18, margin: 10, marginLeft: 15 }} />
            </Flex>
        )
    }
}
