import React from "react"
import AppShellContext from "../../core/components/shell/AppShellContext"
import "./HomePage.css"

export default class HomePage extends React.PureComponent {
    componentDidMount() {
        AppShellContext.get().updateConfig({})
    }

    render() {
        return (
            <>
                <h1 className="home-title">Hello, and welcome to nginx ignition 👋</h1>
                <p className="home-subtitle">
                    Here are some quick start tips and information to help you make the most of the app. If you're
                    already know your ways, feel free to just go and do what you need. 😃
                </p>
                <p>TODO: Implement this</p>
            </>
        )
    }
}
