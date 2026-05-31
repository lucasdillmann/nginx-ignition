import React from "react"
import { Button, Flex } from "antd"
import { GithubFilled, LinkedinFilled } from "@ant-design/icons"
import AppContext from "../../../core/components/context/AppContext"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"
import { I18n } from "../../../core/i18n/I18n"
import NginxMetadata from "../../nginx/model/NginxMetadata"
import "./HomeHeaderCard.css"
import If from "../../../core/components/flowcontrol/If"

export interface HomeHeaderCardProps {
    metadata?: NginxMetadata
    onRefresh: () => void
}

export default class HomeHeaderCard extends React.PureComponent<HomeHeaderCardProps> {
    private firstName(): string {
        const user = AppContext.get().user
        if (user === undefined) return ""

        const trimmedName = user.name.trim()
        if (trimmedName.length > 0) return trimmedName.split(/\s+/)[0]

        return user.username
    }

    private releaseUrl(latest: string): string {
        return `https://github.com/lucasdillmann/nginx-ignition/releases/${latest}`
    }

    private handleGithubClick() {
        window.open("https://github.com/lucasdillmann/nginx-ignition", "_blank", "noopener")
    }

    private handleLinkedInClick() {
        window.open("https://linkedin.com/in/lucasdillmann", "_blank", "noopener")
    }

    private renderAppVersion() {
        const { current } = AppContext.get().configuration.version

        return (
            <span className="home-header-card-meta-line">
                <I18n
                    id={
                        current
                            ? MessageKey.FrontendComponentsShellVersionFormat
                            : MessageKey.FrontendComponentsShellVersionDev
                    }
                    params={{ version: current }}
                />
            </span>
        )
    }

    private renderNginxVersion() {
        const { metadata } = this.props
        if (metadata === undefined) return null

        return (
            <span className="home-header-card-meta-line">
                <I18n id={{ id: MessageKey.FrontendHomeNginxVersion, params: { version: metadata.version } }} />
            </span>
        )
    }

    private renderSocialIcons() {
        return (
            <Flex className="home-header-card-social-row" align="center" gap={10}>
                <LinkedinFilled className="home-header-card-social-icon" onClick={() => this.handleLinkedInClick()} />
                <GithubFilled className="home-header-card-social-icon" onClick={() => this.handleGithubClick()} />
            </Flex>
        )
    }

    private renderActions(onRefresh: () => void) {
        const { current, latest } = AppContext.get().configuration.version
        const updateAvailable = current !== undefined && latest !== undefined && current !== latest

        return (
            <Flex className="home-header-card-actions" align="center" gap={8}>
                <If condition={updateAvailable}>
                    <Button
                        className="home-header-card-update-button"
                        onClick={() => window.open(this.releaseUrl(latest!), "_blank", "noopener")}
                    >
                        <I18n id={MessageKey.FrontendHomeUpdateAvailable} />
                    </Button>
                </If>
                <Button type="primary" onClick={() => onRefresh()}>
                    <I18n id={MessageKey.CommonRefresh} />
                </Button>
            </Flex>
        )
    }

    render() {
        const { onRefresh } = this.props

        return (
            <div className="home-header-card">
                <Flex className="home-header-card-columns" align="center" justify="space-between" wrap="wrap">
                    <Flex className="home-header-card-greeting" vertical justify="center">
                        <h1 className="home-header-card-title">
                            <I18n
                                id={{
                                    id: MessageKey.FrontendHomeGreetingTitle,
                                    params: { userName: this.firstName() },
                                }}
                            />
                        </h1>
                        <p className="home-header-card-subtitle">
                            <I18n id={MessageKey.FrontendHomeGreetingSubtitle} />
                        </p>
                    </Flex>
                    <Flex className="home-header-card-sidebar" vertical align="flex-end">
                        <Flex className="home-header-card-meta-group" vertical align="flex-end">
                            {this.renderAppVersion()}
                            {this.renderNginxVersion()}
                            {this.renderSocialIcons()}
                        </Flex>
                        {this.renderActions(onRefresh)}
                    </Flex>
                </Flex>
            </div>
        )
    }
}
