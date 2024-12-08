import React from "react"
import AppShellContext from "../../core/components/shell/AppShellContext"

export default class SettingsPage extends React.PureComponent {
    componentDidMount() {
        AppShellContext.get().updateConfig({
            title: "Settings",
            subtitle: "Globals settings for the nginx server",
        })
    }

    render() {
        return "TODO: Implement this"
    }
}
