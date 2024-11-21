import React from "react";
import ShellAwareComponent, {ShellConfig} from "../../core/components/shell/ShellAwareComponent";

export default class UserFormPage extends ShellAwareComponent {
    shellConfig(): ShellConfig {
        return {
            title: "User details",
            subtitle: "Full details and configurations of the nginx ignition's user",
            actions: [
                {
                    description: "Save",
                    onClick: () => {},
                },
            ],
        };
    }

    render() {
        return <>TODO: Implement this (UserFormPage)</>;
    }
}
