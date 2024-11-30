import React from "react"
import { LoadingOutlined } from "@ant-design/icons"

export interface FullPagePreloaderProps {
    title?: string
    message?: string
}

export default class FullPagePreloader extends React.Component<FullPagePreloaderProps> {
    render() {
        const title = this.props.title ?? "nginx ignition"
        const message = this.props.message ?? "Hang on tight, we're loading some stuff and will be ready soon"

        return (
            <div className="preloader-container">
                <div className="preloader-body">
                    <div className="preloader-spinner">
                        <LoadingOutlined />
                    </div>
                    <div className="preloader-text-container">
                        <h2 className="preloader-title">{title}</h2>
                        <p className="preloader-message">{message}</p>
                    </div>
                </div>
            </div>
        )
    }
}
