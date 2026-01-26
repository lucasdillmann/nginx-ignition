import React from "react"
import AppShellContext from "../../core/components/shell/AppShellContext"
import "./HomePage.css"
import {
    AuditOutlined,
    BlockOutlined,
    FileProtectOutlined,
    FileSearchOutlined,
    HddOutlined,
    MergeCellsOutlined,
    SettingOutlined,
    ApartmentOutlined,
    RocketOutlined,
} from "@ant-design/icons"
import { Flex } from "antd"
import Videos from "./videos/Videos"
import { Link } from "react-router-dom"
import MessageKey from "../../core/i18n/model/MessageKey.generated"
import { I18n } from "../../core/i18n/I18n"

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
                    <h1>
                        <I18n id={MessageKey.FrontendHomeWelcomeTitle} />
                    </h1>
                    <p className="home-guide-subtitle">
                        <I18n id={MessageKey.FrontendHomeWelcomeSubtitle} />
                    </p>
                </div>

                <Flex className="home-guide-section">
                    <Flex className="home-guide-section-content" vertical>
                        <h2>
                            <HddOutlined /> <I18n id={MessageKey.CommonHosts} />
                        </h2>
                        <p>
                            <I18n id={MessageKey.FrontendHomeHostsDescription1} />
                        </p>
                        <p>
                            <I18n id={MessageKey.FrontendHomeHostsDescription2} />
                        </p>
                        <p>
                            <I18n id={MessageKey.FrontendHomeHostsDescription3} />
                        </p>
                    </Flex>
                    <Flex className="home-guide-right-side-video">
                        <div className="home-guide-video-mask">
                            <video src={Videos.Hosts} autoPlay loop controls />
                        </div>
                    </Flex>
                </Flex>

                <Flex className="home-guide-section">
                    <Flex className="home-guide-left-side-video">
                        <div className="home-guide-video-mask">
                            <video src={Videos.Streams} autoPlay loop controls />
                        </div>
                    </Flex>
                    <Flex className="home-guide-section-content" vertical>
                        <h2>
                            <MergeCellsOutlined /> <I18n id={MessageKey.CommonStreams} />
                        </h2>
                        <p>
                            <I18n id={MessageKey.FrontendHomeStreamsDescription1} />
                        </p>
                        <p>
                            <I18n id={MessageKey.FrontendHomeStreamsDescription2} />
                        </p>
                    </Flex>
                </Flex>

                <Flex className="home-guide-section">
                    <Flex className="home-guide-section-content" vertical>
                        <h2>
                            <AuditOutlined /> <I18n id={MessageKey.CommonSslCertificates} />
                        </h2>
                        <p>
                            <I18n id={MessageKey.FrontendHomeSslDescription1} />
                        </p>
                        <p>
                            <I18n id={MessageKey.FrontendHomeSslDescription2} />
                        </p>
                        <p>
                            <I18n id={MessageKey.FrontendHomeSslDescription3} />
                        </p>
                    </Flex>
                    <Flex className="home-guide-right-side-video">
                        <div className="home-guide-video-mask">
                            <video src={Videos.SslCertificates} autoPlay loop controls />
                        </div>
                    </Flex>
                </Flex>

                <Flex className="home-guide-section">
                    <Flex className="home-guide-left-side-video">
                        <div className="home-guide-video-mask">
                            <video src={Videos.Logs} autoPlay loop controls />
                        </div>
                    </Flex>
                    <Flex className="home-guide-section-content" vertical>
                        <h2>
                            <FileSearchOutlined /> <I18n id={MessageKey.CommonLogs} />
                        </h2>
                        <p>
                            <I18n id={MessageKey.FrontendHomeLogsDescription1} />
                        </p>
                        <p>
                            <I18n id={MessageKey.FrontendHomeLogsDescription2} />
                        </p>
                    </Flex>
                </Flex>

                <Flex className="home-guide-section">
                    <Flex className="home-guide-section-content" vertical>
                        <h2>
                            <BlockOutlined /> <I18n id={MessageKey.CommonIntegrations} />
                        </h2>
                        <p>
                            <I18n id={MessageKey.FrontendHomeIntegrationsDescription} />
                        </p>
                    </Flex>
                    <Flex className="home-guide-right-side-video">
                        <div className="home-guide-video-mask">
                            <video src={Videos.Integrations} autoPlay loop controls />
                        </div>
                    </Flex>
                </Flex>

                <Flex className="home-guide-section">
                    <Flex className="home-guide-left-side-video">
                        <div className="home-guide-video-mask">
                            <video src={Videos.VPNs} autoPlay loop controls />
                        </div>
                    </Flex>
                    <Flex className="home-guide-section-content" vertical>
                        <h2>
                            <ApartmentOutlined /> <I18n id={MessageKey.CommonVpns} />
                        </h2>
                        <p>
                            <I18n id={MessageKey.FrontendHomeVpnsDescription} />
                        </p>
                    </Flex>
                </Flex>

                <Flex className="home-guide-section">
                    <Flex className="home-guide-section-content" vertical>
                        <h2>
                            <FileProtectOutlined /> <I18n id={MessageKey.CommonAccessLists} />
                        </h2>
                        <p>
                            <I18n id={MessageKey.FrontendHomeAccessListsDescription} />
                        </p>
                    </Flex>
                    <Flex className="home-guide-right-side-video">
                        <div className="home-guide-video-mask">
                            <video src={Videos.AccessLists} autoPlay loop controls />
                        </div>
                    </Flex>
                </Flex>

                <Flex className="home-guide-section">
                    <Flex className="home-guide-left-side-video">
                        <div className="home-guide-video-mask">
                            <video src={Videos.Caches} autoPlay loop controls />
                        </div>
                    </Flex>
                    <Flex className="home-guide-section-content" vertical>
                        <h2>
                            <RocketOutlined /> <I18n id={MessageKey.CommonCacheConfiguration} />
                        </h2>
                        <p>
                            <I18n id={MessageKey.FrontendHomeCacheDescription1} />
                        </p>
                        <p>
                            <I18n id={MessageKey.FrontendHomeCacheDescription2} />
                        </p>
                    </Flex>
                </Flex>

                <Flex className="home-guide-section">
                    <Flex className="home-guide-section-content" vertical>
                        <h2>
                            <SettingOutlined /> <I18n id={MessageKey.CommonSettings} />
                        </h2>
                        <p>
                            <I18n id={MessageKey.FrontendHomeSettingsDescription1} />
                        </p>
                        <p>
                            <I18n id={MessageKey.FrontendHomeSettingsDescription2} />
                        </p>
                        <p>
                            <I18n id={MessageKey.FrontendHomeSettingsDescription3} />
                        </p>
                    </Flex>
                    <Flex className="home-guide-right-side-video">
                        <div className="home-guide-video-mask">
                            <video src={Videos.Settings} autoPlay loop controls />
                        </div>
                    </Flex>
                </Flex>

                <div className="home-guide-footer-container">
                    <h1>
                        <I18n id={MessageKey.FrontendHomeFooterTitle} />
                    </h1>
                    <p className="home-guide-subtitle">
                        <Link to="https://github.com/lucasdillmann/nginx-ignition" target="_blank">
                            <I18n id={MessageKey.FrontendHomeFooterLink} />
                        </Link>
                        . <I18n id={MessageKey.FrontendHomeFooterSubtitle} />
                    </p>
                </div>
            </div>
        )
    }
}
