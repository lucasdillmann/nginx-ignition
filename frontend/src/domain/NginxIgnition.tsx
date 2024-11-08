import React from 'react';
import ApiClientEventDispatcher from "../core/apiclient/event/ApiClientEventDispatcher";
import AuthenticationApiClientEventListener from "../core/authentication/AuthenticationApiClientEventListener";
import { App, ConfigProvider } from "antd";

export default class NginxIgnition extends React.PureComponent {
    constructor(props: any) {
        super(props);
        ApiClientEventDispatcher.register(new AuthenticationApiClientEventListener())
    }

    componentDidMount() {
        const preloader = document.getElementById('preloader') as HTMLElement
        preloader.remove()
    }

    render() {
        return (
            <React.StrictMode>
                <ConfigProvider>
                    <App>
                        <p>Hello there</p>
                    </App>
                </ConfigProvider>
            </React.StrictMode>
        );
    }
}
