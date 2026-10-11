import React from "react"
import AppShellContext from "../../core/components/shell/AppShellContext"
import "./HelpPage.css"
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

export default class HelpPage extends React.PureComponent {
    componentDidMount() {
        AppShellContext.get().updateConfig({
            noContainerPadding: true,
        })
    }

    render() {
        return (
            <div className="help-guide-container">
                <div className="help-guide-header-container">
                    <h1>
                        <I18n id={MessageKey.FrontendHelpWelcomeTitle} />
                    </h1>
                    <p className="help-guide-subtitle">
                        <I18n id={MessageKey.FrontendHelpWelcomeSubtitle} />
                    </p>
                </div>

                <Flex className="help-guide-section">
                    <Flex className="help-guide-section-content" vertical>
                        <h2>
                            <HddOutlined /> <I18n id={MessageKey.CommonHosts} />
                        </h2>
                        <p>
                            <I18n id={MessageKey.FrontendHelpHostsDescription1} />
                        </p>
                        <p>
                            <I18n id={MessageKey.FrontendHelpHostsDescription2} />
                        </p>
                        <p>
                            <I18n id={MessageKey.FrontendHelpHostsDescription3} />
                        </p>
                    </Flex>
                    <Flex className="help-guide-right-side-video">
                        <div className="help-guide-video-mask">
                            <img src={Videos.Hosts} alt="" />
                        </div>
                    </Flex>
                </Flex>

                <Flex className="help-guide-section">
                    <Flex className="help-guide-left-side-video">
                        <div className="help-guide-video-mask">
                            <img src={Videos.Streams} alt="" />
                        </div>
                    </Flex>
                    <Flex className="help-guide-section-content" vertical>
                        <h2>
                            <MergeCellsOutlined /> <I18n id={MessageKey.CommonStreams} />
                        </h2>
                        <p>
                            <I18n id={MessageKey.FrontendHelpStreamsDescription1} />
                        </p>
                        <p>
                            <I18n id={MessageKey.FrontendHelpStreamsDescription2} />
                        </p>
                    </Flex>
                </Flex>

                <Flex className="help-guide-section">
                    <Flex className="help-guide-section-content" vertical>
                        <h2>
                            <AuditOutlined /> <I18n id={MessageKey.CommonSslCertificates} />
                        </h2>
                        <p>
                            <I18n id={MessageKey.FrontendHelpSslDescription1} />
                        </p>
                        <p>
                            <I18n id={MessageKey.FrontendHelpSslDescription2} />
                        </p>
                        <p>
                            <I18n id={MessageKey.FrontendHelpSslDescription3} />
                        </p>
                    </Flex>
                    <Flex className="help-guide-right-side-video">
                        <div className="help-guide-video-mask">
                            <img src={Videos.SslCertificates} alt="" />
                        </div>
                    </Flex>
                </Flex>

                <Flex className="help-guide-section">
                    <Flex className="help-guide-left-side-video">
                        <div className="help-guide-video-mask">
                            <img src={Videos.Logs} alt="" />
                        </div>
                    </Flex>
                    <Flex className="help-guide-section-content" vertical>
                        <h2>
                            <FileSearchOutlined /> <I18n id={MessageKey.CommonLogs} />
                        </h2>
                        <p>
                            <I18n id={MessageKey.FrontendHelpLogsDescription1} />
                        </p>
                        <p>
                            <I18n id={MessageKey.FrontendHelpLogsDescription2} />
                        </p>
                    </Flex>
                </Flex>

                <Flex className="help-guide-section">
                    <Flex className="help-guide-section-content" vertical>
                        <h2>
                            <BlockOutlined /> <I18n id={MessageKey.CommonIntegrations} />
                        </h2>
                        <p>
                            <I18n id={MessageKey.FrontendHelpIntegrationsDescription} />
                        </p>
                    </Flex>
                    <Flex className="help-guide-right-side-video">
                        <div className="help-guide-video-mask">
                            <img src={Videos.Integrations} alt="" />
                        </div>
                    </Flex>
                </Flex>

                <Flex className="help-guide-section">
                    <Flex className="help-guide-left-side-video">
                        <div className="help-guide-video-mask">
                            <img src={Videos.VPNs} alt="" />
                        </div>
                    </Flex>
                    <Flex className="help-guide-section-content" vertical>
                        <h2>
                            <ApartmentOutlined /> <I18n id={MessageKey.CommonVpns} />
                        </h2>
                        <p>
                            <I18n id={MessageKey.FrontendHelpVpnsDescription} />
                        </p>
                    </Flex>
                </Flex>

                <Flex className="help-guide-section">
                    <Flex className="help-guide-section-content" vertical>
                        <h2>
                            <FileProtectOutlined /> <I18n id={MessageKey.CommonAccessLists} />
                        </h2>
                        <p>
                            <I18n id={MessageKey.FrontendHelpAccessListsDescription} />
                        </p>
                    </Flex>
                    <Flex className="help-guide-right-side-video">
                        <div className="help-guide-video-mask">
                            <img src={Videos.AccessLists} alt="" />
                        </div>
                    </Flex>
                </Flex>

                <Flex className="help-guide-section">
                    <Flex className="help-guide-left-side-video">
                        <div className="help-guide-video-mask">
                            <img src={Videos.Caches} alt="" />
                        </div>
                    </Flex>
                    <Flex className="help-guide-section-content" vertical>
                        <h2>
                            <RocketOutlined /> <I18n id={MessageKey.CommonCacheConfiguration} />
                        </h2>
                        <p>
                            <I18n id={MessageKey.FrontendHelpCacheDescription1} />
                        </p>
                        <p>
                            <I18n id={MessageKey.FrontendHelpCacheDescription2} />
                        </p>
                    </Flex>
                </Flex>

                <Flex className="help-guide-section">
                    <Flex className="help-guide-section-content" vertical>
                        <h2>
                            <SettingOutlined /> <I18n id={MessageKey.CommonSettings} />
                        </h2>
                        <p>
                            <I18n id={MessageKey.FrontendHelpSettingsDescription1} />
                        </p>
                        <p>
                            <I18n id={MessageKey.FrontendHelpSettingsDescription2} />
                        </p>
                        <p>
                            <I18n id={MessageKey.FrontendHelpSettingsDescription3} />
                        </p>
                    </Flex>
                    <Flex className="help-guide-right-side-video">
                        <div className="help-guide-video-mask">
                            <img src={Videos.Settings} alt="" />
                        </div>
                    </Flex>
                </Flex>

                <div className="help-guide-footer-container">
                    <h1>
                        <I18n id={MessageKey.FrontendHelpFooterTitle} />
                    </h1>
                    <p className="help-guide-subtitle">
                        <Link to="https://github.com/lucasdillmann/nginx-ignition" target="_blank">
                            <I18n id={MessageKey.FrontendHelpFooterLink} />
                        </Link>
                        . <I18n id={MessageKey.FrontendHelpFooterSubtitle} />
                    </p>
                </div>
            </div>
        )
    }
}
