import React from "react";
import ShellAwareComponent, {ShellConfig} from "../../core/components/shell/ShellAwareComponent";

export default class HostFormPage extends ShellAwareComponent {
    shellConfig(): ShellConfig {
        return {
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
        }
    }

    render() {
        return <>TODO: Implement this (HostFormPage)</>;
    }
}
