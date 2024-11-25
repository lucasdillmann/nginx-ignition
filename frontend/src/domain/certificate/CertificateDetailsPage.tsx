import React from "react";
import AppShellContext from "../../core/components/shell/AppShellContext";

export default class CertificateDetailsPage extends React.PureComponent {
    static contextType = AppShellContext
    context!: React.ContextType<typeof AppShellContext>

    componentDidMount() {
        this.context.updateConfig({
            title: "SSL certificate details",
            subtitle: "Details of a uploaded or issued SSL certificate",
            actions: [
                {
                    description: "Delete",
                    color: "danger",
                    onClick: () => {},
                },
                {
                    description: "Renew",
                    onClick: () => {},
                },
            ],
        })
    }

    render() {
        return <>TODO: Implement this (CertificateDetailsPage)</>;
    }
}
