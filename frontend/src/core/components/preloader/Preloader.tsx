import React, {PropsWithChildren} from "react";
import {Spin} from "antd";
import {LoadingOutlined} from "@ant-design/icons";

interface PreloaderState {
    loading: boolean
}

export interface PreloaderProps extends PropsWithChildren {
    size?: number
    loading: boolean
}

export default class Preloader extends React.PureComponent<PreloaderProps, PreloaderState> {
    private changeStateTimeoutId?: number

    constructor(props: PreloaderProps) {
        super(props);
        this.state = {
            loading: false,
        }
    }

    private setLoadingStateDelayed(loading: boolean) {
        if (this.changeStateTimeoutId !== undefined)
            window.clearTimeout(this.changeStateTimeoutId)

        if (loading)
            this.changeStateTimeoutId = window.setTimeout(() => this.setState({ loading }), 500)
        else
            this.setState({ loading })
    }

    componentDidUpdate(prevProps: Readonly<PreloaderProps>) {
        const {loading: currentValue} = this.props
        const {loading: previousValue} = prevProps

        if (currentValue !== previousValue)
            this.setLoadingStateDelayed(currentValue)
    }

    render() {
        const {loading} = this.state
        const {children, size} = this.props

        if (!loading)
            return children

        return (
            <Spin indicator={<LoadingOutlined style={{fontSize: size ?? 48}} spin />}>
                {children}
            </Spin>
        )
    }
}
