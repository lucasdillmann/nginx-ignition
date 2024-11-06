import React from 'react';
import ApiClientEventDispatcher from "../core/apiclient/event/ApiClientEventDispatcher";
import AuthenticationApiClientEventListener from "../core/authentication/AuthenticationApiClientEventListener";

export default class NginxIgnite extends React.PureComponent {
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
                <p>Hello there</p>
            </React.StrictMode>
        );
    }
}
