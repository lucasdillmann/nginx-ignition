import React from "react"
import { Navigate } from "react-router-dom"
import AppShellContext from "../../core/components/shell/AppShellContext"

export default class HomePage extends React.PureComponent {
    static readonly contextType = AppShellContext
    context!: React.ContextType<typeof AppShellContext>

    componentDidMount() {
        this.context.updateConfig({
            title: "Home",
        })
    }

    render() {
        return <Navigate to="/hosts" />
    }
}
