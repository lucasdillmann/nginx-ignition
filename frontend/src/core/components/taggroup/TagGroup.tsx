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
        const tags = values.slice(0, limit).map(name => <Tag>{name}</Tag>)
        const tooltipContents = (<>{values.slice(limit).map(tag => <>{tag}<br /></>)}</>)
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
