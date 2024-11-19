import React from "react";
import {Flex} from "antd";
import Preloader from "./Preloader";
import "./FullPagePreloader.css"

export interface FullPagePreloaderProps {
    title?: string
    message?: string
}

export default class FullPagePreloader extends React.Component<FullPagePreloaderProps> {
    render() {
        const title = this.props.title ?? "nginx ignition"
        const message = this.props.message ?? "Hang on tight, we're loading some stuff"

        return (
            <Flex align="center" justify="center" className="preloader-container">
                <Flex align="center">
                    <Preloader loading />
                    <Flex className="preloader-text-container" vertical>
                        <h2 className="preloader-title">{title}</h2>
                        <p className="preloader-message">{message}</p>
                    </Flex>
                </Flex>
            </Flex>
        )
    }
}
