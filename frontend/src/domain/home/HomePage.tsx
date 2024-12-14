import React from "react"
import AppShellContext from "../../core/components/shell/AppShellContext"
import "./HomePage.css"
import { BlockOutlined, FileProtectOutlined, FileSearchOutlined, HddOutlined, SettingOutlined } from "@ant-design/icons"
import { Flex } from "antd"
import Videos from "./videos/Videos"
import { Link } from "react-router-dom"

export default class HomePage extends React.PureComponent {
    componentDidMount() {
        AppShellContext.get().updateConfig({
            noContainerPadding: true,
        })
    }

    render() {
        return (
            <div className="home-guide-container">
                <div className="home-guide-header-container">
                    <h1>Hello, and welcome to nginx ignition 👋</h1>
                    <p className="home-guide-subtitle">
                        Here are some quick start info to help you make the most of the app. We hope you enjoy it.
                    </p>
                </div>

                <Flex className="home-guide-section">
                    <Flex className="home-guide-section-content" vertical>
                        <h2>
                            <HddOutlined /> Hosts
                        </h2>
                        <p>
                            A virtual host, host in short, is a website that the nginx server will make available to
                            be opened in a browser. You can find all hosts managed at the left, on the main menu.
                        </p>
                        <p>
                            nginx ignition provides an intuitive way to configure such websites. For example, if you
                            have a NAS and have some services running on it (like Jellyfin, Vaultwarden and more),
                            ignition enables an easy way to access it from a domain like jellyfin.myhome.com, way
                            easier to remember and use than an IP and port.
                        </p>
                        <p>
                            Each host will have a set of routes, which are rules that define which requests patterns
                            should be forwarded and to where, and bindings, which are the definitions at what ports
                            the nginx should listen for such requests.
                        </p>
                    </Flex>
                    <Flex className="home-guide-right-side-video">
                        <video src={Videos.Hosts} autoPlay loop controls />
                    </Flex>
                </Flex>

                <Flex className="home-guide-section">
                    <Flex className="home-guide-left-side-video">
                        <video src={Videos.Certificates} autoPlay loop controls />
                    </Flex>
                    <Flex className="home-guide-section-content" vertical>
                        <h2>
                            <FileProtectOutlined /> SSL certificates
                        </h2>
                        <p>
                            If you want to or need to protect your domains with HTTPS encryption, the app comes with
                            an easy way to manage such SSL certificates too.
                        </p>
                        <p>
                            Either if you need a valid certificate backed by Let's Encrypt, a self-signed, or bring
                            your custom one for a third-party provider, ignition will allow it with ease. Even when the
                            certificate is about to expire, the app will automatically renew it for you.

                        </p>
                        <p>
                            Once an SSL certificate is created or imported, you can use it on the hosts by simply
                            selecting an option in the form.
                        </p>
                    </Flex>
                </Flex>

                <Flex className="home-guide-section">
                    <Flex className="home-guide-section-content" vertical>
                        <h2>
                            <FileSearchOutlined /> Logs
                        </h2>
                        <p>
                            Need to know what has been requested or why something isn't working the way it was expected?
                            ignition provides a simple way for you to check the nginx logs too. Just select the host
                            (or the nginx server itself) that you want to check what's going on, and we will get the
                            logs for you.
                        </p>
                        <p>
                            And you don't need to worry about your disk getting full of logs, ignition will rotate them
                            automatically for you (and you can control how much to keep, or even disable the rotation
                            if it suits you best).
                        </p>
                    </Flex>
                    <Flex className="home-guide-right-side-video">
                        <video src={Videos.Logs} autoPlay loop controls />
                    </Flex>
                </Flex>

                <Flex className="home-guide-section">
                    <Flex className="home-guide-left-side-video">
                        <video src={Videos.Integrations} autoPlay loop controls />
                    </Flex>
                    <Flex className="home-guide-section-content" vertical>
                        <h2>
                            <BlockOutlined /> Integrations
                        </h2>
                        <p>
                            Is your app running in a Docker container or in a TrueNAS? You can enable the native
                            integration nginx ignition offers and easily pick a container or app as the destination
                            app that the host should forward the requests to.
                        </p>
                    </Flex>
                </Flex>

                <Flex className="home-guide-section">
                    <Flex className="home-guide-section-content" vertical>
                        <h2>
                            <SettingOutlined /> Settings
                        </h2>
                        <p>
                            ignition abstracts away the complexity of the nginx's configuration files, but that doesn't
                            mean that you lose the ability to apply some fine adjustments.
                        </p>
                        <p>
                            With the settings page, you can define some important definitions such as the timeout
                            values, the maximum upload/body size of the requests, the default ports that the nginx will
                            listen for requests and more.
                        </p>
                        <p>
                            Beyond that, you can also configure some of the nginx ignition's features there, like the
                            automatic renewal of SSL certificates and log rotation.
                        </p>
                    </Flex>
                    <Flex className="home-guide-right-side-video">
                        <video src={Videos.Settings} autoPlay loop controls />
                    </Flex>
                </Flex>

                <div className="home-guide-footer-container">
                    <h1>Still have questions or miss something?</h1>
                    <p className="home-guide-subtitle">
                        <Link to="https://github.com/lucasdillmann/nginx-ignition" target="_blank">
                            Reach us out at our GitHub page
                        </Link>.
                        We'd love some constructive feedback from you.
                    </p>
                </div>
            </div>
        )
    }
}
