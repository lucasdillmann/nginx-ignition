import React from "react"
import AppShellContext from "../../core/components/shell/AppShellContext"
import "./HomePage.css"
import { BlockOutlined, FileProtectOutlined, FileSearchOutlined, HddOutlined, SettingOutlined } from "@ant-design/icons"
import { Flex } from "antd"
import Videos from "./videos/Videos"

export default class HomePage extends React.PureComponent {
    componentDidMount() {
        AppShellContext.get().updateConfig({})
    }

    render() {
        return (
            <div className="home-guide-container">
                <h1>Hello, and welcome to nginx ignition 👋</h1>
                <p className="home-guide-subtitle">
                    Here are some quick start info to help you make the most of the app. We hope you enjoy it.
                </p>

                <Flex className="home-guide-section">
                    <Flex className="home-guide-section-content" vertical>
                        <h2>
                            <HddOutlined /> Hosts
                        </h2>
                        <p>
                            A virtual host, host in short, is a website that the nginx server will make available to be
                            opened in a browser. You can find all hosts managed at the left, in the main menu.
                        </p>
                        <p>
                            nginx ignition provides an intuitive way to configure such websites. For example, if you
                            have a NAS and have some services running in it (like Jellyfin, Vaultwarden and more),
                            ignition enables a easy way to access them from a domain like jellyfin.myhome.com, way more
                            easy to remember and use that an IP and port.
                        </p>
                        <p>
                            Each host will have a set of routes, which are rules that defines which requests patterns
                            should be forwarded and to where, and bindings, which are the definitions in what ports the
                            nginx should listen for such requests.
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
                            If you want to or need to protect you domains with HTTPS encryption, the app comes with a
                            easy way to manage such SSL certificates too.
                        </p>
                        <p>
                            Either if you need a valid certificate backed by Let's Encrypt, a self-signed or bring your
                            custom one for a third-party provider, ignition will allow it with ease. Even when the
                            certificate is about to expire, the app will automatically renew it for you.
                        </p>
                        <p>
                            Once a SSL certificated is created or imported, you can use them on the hosts by simple
                            selecting a option in the form.
                        </p>
                    </Flex>
                </Flex>

                <Flex className="home-guide-section">
                    <Flex className="home-guide-section-content" vertical>
                        <h2>
                            <FileSearchOutlined /> Logs
                        </h2>
                        <p>TODO: Write this</p>
                    </Flex>
                    <Flex className="home-guide-right-side-video">
                        <video src={Videos.Certificates} autoPlay loop controls />
                    </Flex>
                </Flex>

                <Flex className="home-guide-section">
                    <Flex className="home-guide-left-side-video">
                        <video src={Videos.Certificates} autoPlay loop controls />
                    </Flex>
                    <Flex className="home-guide-section-content" vertical>
                        <h2>
                            <BlockOutlined /> Integrations
                        </h2>
                        <p>TODO: Write this</p>
                    </Flex>
                </Flex>

                <Flex className="home-guide-section">
                    <Flex className="home-guide-section-content" vertical>
                        <h2>
                            <SettingOutlined /> Settings
                        </h2>
                        <p>TODO: Write this</p>
                    </Flex>
                    <Flex className="home-guide-right-side-video">
                        <video src={Videos.Certificates} autoPlay loop controls />
                    </Flex>
                </Flex>
            </div>
        )
    }
}
