import React from "react"
import { Tag, Tooltip } from "antd"
import { I18n, I18nMessage, raw } from "../../i18n/I18n"
import MessageKey from "../../i18n/model/MessageKey.generated"

const DEFAULT_MAXIMUM_SIZE = 3

export interface TagGroupItem {
    name: string
    url: string
}

export interface TagGroupProps {
    values: string[] | TagGroupItem[]
    maximumSize?: number
}

export default class TagGroup extends React.Component<TagGroupProps> {
    private isTagGroupItem(value: string | TagGroupItem): value is TagGroupItem {
        return typeof value === "object" && "name" in value && "url" in value
    }

    private getKey(value: string | TagGroupItem): string {
        return this.isTagGroupItem(value) ? value.name : value
    }

    private renderTag(value: string | TagGroupItem) {
        if (this.isTagGroupItem(value)) {
            return (
                <Tag key={value.name}>
                    <a href={value.url} target="_blank" rel="noopener noreferrer">
                        {value.name}
                    </a>
                </Tag>
            )
        }

        return <Tag key={value}>{value}</Tag>
    }

    private renderTooltipItem(value: string | TagGroupItem) {
        const key = this.getKey(value)
        if (this.isTagGroupItem(value)) {
            return (
                <span key={key}>
                    <a href={value.url} target="_blank" rel="noopener noreferrer" style={{ color: "inherit" }}>
                        {value.name}
                    </a>
                    <br />
                </span>
            )
        }

        return (
            <span key={key}>
                {value}
                <br />
            </span>
        )
    }

    render() {
        const { values, maximumSize } = this.props
        const limit = maximumSize ?? DEFAULT_MAXIMUM_SIZE
        const tags = values.slice(0, limit).map(value => this.renderTag(value))
        const tooltipContents = <>{values.slice(limit).map(value => this.renderTooltipItem(value))}</>
        const remaining = values.length - limit
        const additionalMessage: I18nMessage =
            remaining > 0 ? { id: MessageKey.FrontendComponentsTaggroupRemaining, params: { remaining } } : raw("")

        return (
            <>
                {tags}
                <Tooltip title={tooltipContents}>
                    <I18n id={additionalMessage} />
                </Tooltip>
            </>
        )
    }
}
