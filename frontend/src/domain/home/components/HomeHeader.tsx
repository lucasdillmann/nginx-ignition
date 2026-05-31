import React from "react"
import { Button, Flex } from "antd"
import AppContext from "../../../core/components/context/AppContext"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"
import { I18n } from "../../../core/i18n/I18n"
import NginxMetadata from "../../nginx/model/NginxMetadata"
import "./HomeHeader.css"
import If from "../../../core/components/flowcontrol/If"

export interface HomeHeaderProps {
    metadata?: NginxMetadata
    onRefresh: () => void
}

export default class HomeHeader extends React.PureComponent<HomeHeaderProps> {
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

    private renderAppVersion() {
        const { current } = AppContext.get().configuration.version

        return (
            <span className="home-header-meta-line">
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
            <span className="home-header-meta-line">
                <I18n id={{ id: MessageKey.FrontendHomeNginxVersion, params: { version: metadata.version } }} />
            </span>
        )
    }

    private renderActions(onRefresh: () => void) {
        const { current, latest } = AppContext.get().configuration.version
        const updateAvailable = Boolean(current) && Boolean(latest) && current !== latest

        return (
            <Flex className="home-header-actions" align="center" gap={8}>
                <If condition={updateAvailable}>
                    <Button
                        className="home-header-update-button"
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
            <div className="home-header">
                <Flex className="home-header-columns" align="center" justify="space-between" wrap="wrap">
                    <Flex className="home-header-greeting" vertical justify="center">
                        <h1 className="home-header-title">
                            <I18n
                                id={{
                                    id: MessageKey.FrontendHomeGreetingTitle,
                                    params: { userName: this.firstName() },
                                }}
                            />
                        </h1>
                        <p className="home-header-subtitle">
                            <I18n id={MessageKey.FrontendHomeGreetingSubtitle} />
                        </p>
                    </Flex>
                    <Flex className="home-header-sidebar" vertical align="flex-end">
                        <Flex className="home-header-meta-group" vertical align="flex-end">
                            {this.renderAppVersion()}
                            {this.renderNginxVersion()}
                        </Flex>
                        {this.renderActions(onRefresh)}
                    </Flex>
                </Flex>
            </div>
        )
    }
}
