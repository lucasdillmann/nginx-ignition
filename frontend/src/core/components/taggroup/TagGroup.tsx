import React from "react";
import {Tag, Tooltip} from "antd";

const DEFAULT_MAXIMUM_SIZE = 3

export interface TagGroupProps {
    values: string[]
    maximumSize?: number
}

export default class TagGroup extends React.Component<TagGroupProps> {
    render() {
        const {values, maximumSize} = this.props
        const limit = maximumSize ?? DEFAULT_MAXIMUM_SIZE
        const tags = values.slice(0, limit).map(name => <Tag key={name}>{name}</Tag>)
        const tooltipContents = (<>{values.slice(limit).map(tag => <span key={tag}>{tag}<br /></span>)}</>)
        const additionalMessage = values.length > limit ? `and ${values.length - limit} more` : ""

        return (
            <>
                {tags}
                <Tooltip title={tooltipContents}>
                    {additionalMessage}
                </Tooltip>
            </>
        )
    }
}
