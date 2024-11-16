import React, {PropsWithChildren} from "react";
import {Spin} from "antd";
import {LoadingOutlined} from "@ant-design/icons";

export interface PreloaderProps extends PropsWithChildren {
    size?: number
}

export default class Preloader extends React.Component<PreloaderProps> {
    render() {
        const {children, size} = this.props
        return (
            <Spin indicator={<LoadingOutlined style={{fontSize: size ?? 48}} spin/>}>
                {children}
            </Spin>
        )
    }
}
