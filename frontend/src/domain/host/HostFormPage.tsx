import React from "react";
import AppShellContext from "../../core/components/shell/AppShellContext";

export default class HostFormPage extends React.PureComponent {
    static contextType = AppShellContext
    context!: React.ContextType<typeof AppShellContext>

    componentDidMount() {
        this.context.updateConfig({
            title: "Host details",
            subtitle: "Full details and configurations of the nginx's virtual host",
            actions: [
                {
                    description: "Delete",
                    color: "danger",
                    onClick: () => {},
                },
                {
                    description: "Save",
                    onClick: () => {},
                },
            ],
        })
    }

    render() {
        return <>TODO: Implement this (HostFormPage)</>;
    }
}
