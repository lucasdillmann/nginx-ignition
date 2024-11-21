import React from "react";
import {ButtonColorType, ButtonVariantType} from "antd/es/button";

export default abstract class ShellAwareComponent<P = any, S = any, SS = any> extends React.Component<P, S, SS>{
    abstract shellConfig(): ShellConfig
}

export interface ShellAction {
    description: string
    onClick: string | (() => void)
    type?: ButtonVariantType
    color?: ButtonColorType
}

export interface ShellConfig {
    title: string
    subtitle?: string
    actions?: ShellAction[]
}
