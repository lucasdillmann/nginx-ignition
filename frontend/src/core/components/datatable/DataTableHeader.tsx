import React from "react"
import { Flex, Input } from "antd"
import debounce from "debounce"
import { SearchOutlined, SettingOutlined } from "@ant-design/icons"
import { i18n } from "../../i18n/I18n"
import MessageKey from "../../i18n/model/MessageKey.generated"
import DataTableOptions from "./DataTableOptions"

const ICON_STYLE = {
    fontSize: 18,
    margin: 10,
    marginLeft: 15,
} satisfies React.CSSProperties

export interface DataTableHeaderProps {
    id: string
    initialSearchTerms?: string
    onSearch: (searchTerms?: string) => void
}

interface DataTableHeaderState {
    optionsOpen: boolean
}

export default class DataTableHeader extends React.Component<DataTableHeaderProps, DataTableHeaderState> {
    constructor(props: DataTableHeaderProps) {
        super(props)
        this.state = {
            optionsOpen: false,
        }
    }

    private setOptionsVisibility(optionsOpen: boolean) {
        this.setState({ optionsOpen })
    }

    render() {
        const { id, initialSearchTerms } = this.props
        const { optionsOpen } = this.state
        const handleChange = debounce((searchTerms?: string) => this.props.onSearch(searchTerms), 500)

        return (
            <>
                <DataTableOptions id={id} open={optionsOpen} onClose={() => this.setOptionsVisibility(false)} />
                <Flex className="data-table-search-bar-container">
                    <Input
                        // @ts-expect-error target is generic, but in this scenario is safe to use the value attribute
                        onInput={event => handleChange(event.nativeEvent.target!!.value)}
                        onClear={() => handleChange()}
                        defaultValue={initialSearchTerms}
                        placeholder={i18n(MessageKey.FrontendComponentsDatatableSearchPlaceholder)}
                        className="data-table-search-bar"
                        autoFocus
                        allowClear
                    />
                    <SearchOutlined style={ICON_STYLE} />
                    <SettingOutlined style={ICON_STYLE} onClick={() => this.setOptionsVisibility(true)} />
                </Flex>
            </>
        )
    }
}
