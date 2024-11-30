import React, { PropsWithChildren } from "react"
import FullPageError from "../error/FullPageError"

interface ErrorBoundaryState {
    error?: Error
}

export default class ErrorBoundary extends React.Component<PropsWithChildren, ErrorBoundaryState> {
    constructor(props: PropsWithChildren) {
        super(props)
        this.state = {}
    }

    componentDidCatch(error: Error) {
        this.setState(current => ({
            ...current,
            error,
        }))
    }

    render() {
        const { error } = this.state
        const { children } = this.props
        return error === undefined ? children : <FullPageError error={error} />
    }
}
