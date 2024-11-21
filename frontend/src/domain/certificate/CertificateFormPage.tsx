import React from "react";
import ShellAwareComponent, {ShellConfig} from "../../core/components/shell/ShellAwareComponent";

export default class CertificateFormPage extends ShellAwareComponent {
    shellConfig(): ShellConfig {
        return {
            title: "New SSL certificate",
            subtitle: "Issue or upload a SSL certificate for use with the nginx's virtual hosts",
            actions: [
                {
                    description: "Issue and save",
                    onClick: () => {},
                },
            ],
        }
    }

    render() {
        return <>TODO: Implement this (CertificateFormPage)</>;
    }
}
