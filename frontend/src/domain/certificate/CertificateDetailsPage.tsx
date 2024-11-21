import React from "react";
import ShellAwareComponent, {ShellConfig} from "../../core/components/shell/ShellAwareComponent";

export default class CertificateDetailsPage extends ShellAwareComponent {
    shellConfig(): ShellConfig {
        return {
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
        }
    }

    render() {
        return <>TODO: Implement this (CertificateDetailsPage)</>;
    }
}
