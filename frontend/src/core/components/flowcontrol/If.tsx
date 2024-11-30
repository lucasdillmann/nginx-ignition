import React, { PropsWithChildren } from "react"

export interface IfProps extends PropsWithChildren {
    condition: boolean
}

export default class If extends React.PureComponent<IfProps> {
    render() {
        const { condition, children } = this.props
        return condition ? children : undefined
    }
}
