import React from "react"

export interface IfProps {
    condition: boolean
    children: React.ReactNode | (() => React.ReactNode)
}

export default class If extends React.PureComponent<IfProps> {
    render() {
        const { condition, children } = this.props
        if (!condition) return undefined

        return typeof children === "function" ? children() : children
    }
}
