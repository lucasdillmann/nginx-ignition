import React from "react";
import ShellAwareComponent, {ShellConfig} from "../../core/components/shell/ShellAwareComponent";

export default class LogsPage extends ShellAwareComponent {
    shellConfig(): ShellConfig {
        return {
            title: "Logs",
            subtitle: "Query the nginx's produced logs for the main process or each virtual host",
        }
    }

    render() {
        return <>TODO: Implement this (LogsPage)</>;
    }
}
