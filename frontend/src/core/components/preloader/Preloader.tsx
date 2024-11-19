import React, {PropsWithChildren} from "react";
import {Spin} from "antd";
import {LoadingOutlined} from "@ant-design/icons";

export interface PreloaderProps extends PropsWithChildren {
    size?: number
    loading: boolean
}

export default class Preloader extends React.PureComponent<PreloaderProps> {
    render() {
        const {children, size, loading} = this.props

        if (!loading)
            return children

        return (
            <Spin indicator={<LoadingOutlined style={{fontSize: size ?? 48}} spin />}>
                {children}
            </Spin>
        )
    }
}
