import React from "react"
import { SunOutlined, MoonOutlined } from "@ant-design/icons"
import "./ThemeToggle.css"
import ThemeContext from "../context/ThemeContext"

interface ThemeToggleState {
    darkMode: boolean
}

export default class ThemeToggle extends React.Component<unknown, ThemeToggleState> {
    constructor(props: unknown) {
        super(props)

        this.state = {
            darkMode: ThemeContext.isDarkMode(),
        }
    }

    private handleThemeChange(darkMode: boolean) {
        this.setState({ darkMode })
    }

    componentDidMount() {
        ThemeContext.register(this.handleThemeChange.bind(this))
    }

    componentWillUnmount() {
        ThemeContext.deregister(this.handleThemeChange.bind(this))
    }

    render() {
        const darkMode = ThemeContext.isDarkMode()
        const ToggleIcon = darkMode ? SunOutlined : MoonOutlined

        return <ToggleIcon onClick={() => ThemeContext.toggle()} />
    }
}
