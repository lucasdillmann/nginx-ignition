import React from "react";
import {Flex} from "antd";
import Preloader from "./Preloader";
import styles from "./FullPagePreloader.styles"

export interface FullPagePreloaderProps {
    title?: string
    message?: string
}

export default class FullPagePreloader extends React.Component<FullPagePreloaderProps> {
    render() {
        const title = this.props.title ?? "nginx ignition"
        const message = this.props.message ?? "Hang on tight, we're loading some stuff"

        return (
            <Flex align="center" justify="center" style={styles.mainContainer}>
                <Flex align="center">
                    <Preloader />
                    <Flex style={styles.textContainer} vertical>
                        <h2 style={styles.title}>{title}</h2>
                        <p style={styles.message}>{message}</p>
                    </Flex>
                </Flex>
            </Flex>
        )
    }
}
