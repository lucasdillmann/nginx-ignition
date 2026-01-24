import React from "react"
import { Button, Flex, Layout, Menu, Tooltip } from "antd"
import { MenuItemType } from "antd/es/menu/interface"
import { Link } from "react-router-dom"
import AppRoute from "../router/AppRoute"
import "./AppShell.css"
import If from "../flowcontrol/If"
import AppShellContext, { ShellAction, ShellConfig } from "./AppShellContext"
import { GithubFilled, LinkedinFilled } from "@ant-design/icons"
import AppContext from "../context/AppContext"
import MessageKey from "../../i18n/model/MessageKey.generated"
import { I18n, I18nMessage } from "../../i18n/I18n"

const { Sider, Content } = Layout

export interface AppShellMenuItem {
    icon: React.ReactNode
    description: I18nMessage
    path: string
}

export interface AppShellProps {
    menuItems: AppShellMenuItem[]
    activeRoute: AppRoute
    children: React.ReactElement
    userMenu: React.ReactNode
    serverControl: React.ReactNode
}

interface AppShellState {
    config: ShellConfig
}

export default class AppShell extends React.Component<AppShellProps, AppShellState> {
    constructor(props: AppShellProps) {
        super(props)
        this.state = {
            config: {},
        }
    }

    private buildMenuItemsAdapters(): MenuItemType[] {
        const { menuItems } = this.props
        return menuItems.map(item => ({
            key: item.path,
            icon: item.icon,
            label: (
                <Link to={item.path}>
                    <I18n id={item.description} />
                </Link>
            ),
            className: "shell-sider-menu-item",
        }))
    }

    private renderActionButton(action: ShellAction): React.ReactNode {
        const { description, type, color, onClick, disabled, disabledReason } = action
        if (typeof onClick === "string") {
            return (
                <Tooltip title={disabledReason}>
                    <Link to={onClick} key={action.description as string}>
                        <Button
                            className="shell-content-actions-action-item"
                            variant={type ?? "solid"}
                            color={color ?? "primary"}
                            disabled={disabled}
                        >
                            <I18n id={description} />
                        </Button>
                    </Link>
                </Tooltip>
            )
        } else {
            return (
                <Tooltip title={disabledReason}>
                    <Button
                        className="shell-content-actions-action-item"
                        key={action.description as string}
                        variant={type ?? "solid"}
                        color={color ?? "primary"}
                        onClick={onClick}
                        disabled={disabled}
                    >
                        <I18n id={description} />
                    </Button>
                </Tooltip>
            )
        }
    }

    private renderActions() {
        const {
            config: { actions },
        } = this.state
        if (!actions) return null

        return <>{actions.map(action => this.renderActionButton(action))}</>
    }

    private handleLinkedInClick() {
        window.open("https://linkedin.com/in/lucasdillmann", "_blank", "noopener")
    }

    private handleGithubClick() {
        window.open("https://github.com/lucasdillmann/nginx-ignition", "_blank", "noopener")
    }

    shouldComponentUpdate(nextProps: Readonly<AppShellProps>): boolean {
        const { children: previous } = nextProps
        const { children: current } = this.props

        if (previous !== current) {
            this.setState({ config: {} })
        }

        return true
    }

    render() {
        AppShellContext.replace({
            updateConfig: config => this.setState({ config }),
        })

        const { version } = AppContext.get().configuration
        const { activeRoute, children, userMenu, serverControl } = this.props
        const { config } = this.state
        const activeMenuItemPath = activeRoute.activeMenuItemPath ?? activeRoute.path
        const { title, subtitle, noContainerPadding } = config
        const mainContentClassNames = !noContainerPadding
            ? "shell-content-main"
            : "shell-content-main shell-content-main-no-padding"

        return (
            <Layout className="shell-container">
                <Sider trigger={null} width={250} className="shell-sider-container">
                    <div className="shell-sider-logo">
                        <Link to="/" className="shell-sider-logo-link">
                            <I18n id={MessageKey.CommonAppName} />
                        </Link>
                    </div>
                    <div className="shell-sider-server-control">{serverControl}</div>
                    <Menu
                        className="shell-sider-menu-container"
                        theme="dark"
                        mode="inline"
                        selectedKeys={activeMenuItemPath ? [activeMenuItemPath] : undefined}
                        items={this.buildMenuItemsAdapters()}
                    />
                    <div className="shell-sider-bottom">
                        <div className="shell-sider-bottom-credits">
                            <I18n
                                id={
                                    version.current
                                        ? MessageKey.FrontendComponentsShellVersionFormat
                                        : MessageKey.FrontendComponentsShellVersionDev
                                }
                                params={{ version: version.current }}
                            />
                            <br />
                            <I18n id={MessageKey.FrontendComponentsShellCreditsMadeBy} />
                            <LinkedinFilled onClick={() => this.handleLinkedInClick()} />
                            <GithubFilled onClick={() => this.handleGithubClick()} />
                        </div>
                        <div className="shell-sider-bottom-menu">{userMenu}</div>
                    </div>
                </Sider>
                <Layout className="shell-content-container">
                    <If condition={title !== undefined}>
                        <Flex className="shell-content-header-container">
                            <Flex className="shell-content-header" vertical>
                                <If condition={title !== undefined}>
                                    <h1 className="shell-content-header-title">
                                        <I18n id={title!} />
                                    </h1>
                                </If>
                                <If condition={subtitle !== undefined}>
                                    <h2 className="shell-content-header-subtitle">
                                        <I18n id={subtitle!} />
                                    </h2>
                                </If>
                            </Flex>
                            <Flex className="shell-content-header-actions-container">{this.renderActions()}</Flex>
                        </Flex>
                    </If>
                    <Content className={mainContentClassNames}>{children}</Content>
                </Layout>
            </Layout>
        )
    }
}
