import React, {CSSProperties} from "react";
import {Flex} from "antd";
import Preloader from "./Preloader";

export interface FullPagePreloaderProps {
    title?: string
    message?: string
}

const styles: {[key: string]: CSSProperties} = {
    mainContainer: {
        width: "100%",
        height: "100%",
        position: "absolute",
        top: 0,
        left: 0,
    },
    textContainer: {
        paddingLeft: 40,
    },
    title: {
        marginBottom: 0,
    },
    message: {
        marginTop: 0,
    }
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
