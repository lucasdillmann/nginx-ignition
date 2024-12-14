import React from "react"
import { Button, Flex, Layout, Menu, Tooltip } from "antd"
import { MenuItemType } from "antd/es/menu/interface"
import { Link } from "react-router-dom"
import AppRoute from "../router/AppRoute"
import "./AppShell.css"
import If from "../flowcontrol/If"
import AppShellContext, { ShellAction, ShellConfig } from "./AppShellContext"
import { GithubFilled, LinkedinFilled } from "@ant-design/icons"
const { Sider, Content } = Layout

const EmptyConfig: ShellConfig = {
    title: "",
}

export interface AppShellMenuItem {
    icon: React.ReactNode
    description: string
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
            config: EmptyConfig,
        }
    }

    private buildMenuItemsAdapters(): MenuItemType[] {
        const { menuItems } = this.props
        return menuItems.map(item => ({
            key: item.path,
            icon: item.icon,
            label: <Link to={item.path}>{item.description}</Link>,
            className: "shell-sider-menu-item",
        }))
    }

    private renderActionButton(action: ShellAction): React.ReactNode {
        const { description, type, color, onClick, disabled, disabledReason } = action
        if (typeof onClick === "string") {
            return (
                <Tooltip title={disabledReason}>
                    <Link to={onClick} key={action.description}>
                        <Button
                            className="shell-content-actions-action-item"
                            variant={type ?? "solid"}
                            color={color ?? "primary"}
                            disabled={disabled}
                        >
                            {description}
                        </Button>
                    </Link>
                </Tooltip>
            )
        } else {
            return (
                <Tooltip title={disabledReason}>
                    <Button
                        className="shell-content-actions-action-item"
                        key={action.description}
                        variant={type ?? "solid"}
                        color={color ?? "primary"}
                        onClick={onClick}
                        disabled={disabled}
                    >
                        {description}
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
        window.open("https://linkedin.com/in/lucasdillmann", "_blank")
    }

    private handleGithubClick() {
        window.open("https://github.com/lucasdillmann/nginx-ignition", "_blank")
    }

    shouldComponentUpdate(nextProps: Readonly<AppShellProps>): boolean {
        const { children: previous } = nextProps
        const { children: current } = this.props

        if (previous !== current) {
            this.setState({ config: EmptyConfig })
        }

        return true
    }

    render() {
        AppShellContext.replace({
            updateConfig: config => this.setState({ config }),
        })

        const { activeRoute, children, userMenu, serverControl } = this.props
        const { config } = this.state
        const activeMenuItemPath = activeRoute.activeMenuItemPath ?? activeRoute.path
        const { title, subtitle } = config

        return (
            <Layout className="shell-container">
                <Sider trigger={null} width={250} className="shell-sider-container">
                    <div className="shell-sider-logo">
                        <Link to="/" className="shell-sider-logo-link">
                            nginx ignition
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
                            Made by Lucas Dillmann
                            <LinkedinFilled onClick={() => this.handleLinkedInClick()} />
                            <GithubFilled onClick={() => this.handleGithubClick()} />
                        </div>
                        <div className="shell-sider-bottom-menu">{userMenu}</div>
                    </div>
                </Sider>
                <Layout className="shell-content-container">
                    <Flex className="shell-content-header-container">
                        <Flex className="shell-content-header" vertical>
                            <h1 className="shell-content-header-title">{title}</h1>
                            <If condition={subtitle !== undefined}>
                                <h2 className="shell-content-header-subtitle">{subtitle}</h2>
                            </If>
                        </Flex>
                        <Flex className="shell-content-header-actions-container">{this.renderActions()}</Flex>
                    </Flex>
                    <Content className="shell-content-main">{children}</Content>
                </Layout>
            </Layout>
        )
    }
}
