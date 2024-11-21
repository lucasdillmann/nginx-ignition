import React from "react";
import {Navigate} from "react-router-dom";
import ShellAwareComponent, {ShellConfig} from "../../core/components/shell/ShellAwareComponent";

export default class HomePage extends ShellAwareComponent {
    shellConfig(): ShellConfig {
        return {
            title: "Home",
        };
    }

    render() {
        return (
            <Navigate to="/hosts" />
        )
    }
}
