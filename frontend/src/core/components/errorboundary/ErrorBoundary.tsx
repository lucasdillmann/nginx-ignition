import React, {PropsWithChildren} from "react";

interface ErrorBoundaryState {
    error?: Error
}

export default class ErrorBoundary extends React.Component<PropsWithChildren, ErrorBoundaryState> {
    constructor(props: PropsWithChildren) {
        super(props);
        this.state = {}
    }

    componentDidCatch(error: Error) {
        this.setState((current) => ({
            ...current,
            error,
        }))
    }

    render() {
        const { error } = this.state
        const { children } = this.props

        if (error === undefined) {
            return children
        }

        // TODO: Improve this
        return (
            <p>Fatal error: ${error.message}</p>
        )
    }
}
