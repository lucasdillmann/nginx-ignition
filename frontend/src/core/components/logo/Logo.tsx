import React from "react"
import LogoImage from "./logo.svg?react"
import "./Logo.css"

export interface LogoProps {
    size?: number
}

export default class Logo extends React.Component<LogoProps> {
    render() {
        const { size = 40 } = this.props

        return <LogoImage aria-hidden="true" className="app-logo" style={{ width: size, height: size }} />
    }
}
