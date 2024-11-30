import React from "react"
import { Button, Flex, Layout, Menu } from "antd"
import { MenuItemType } from "antd/es/menu/interface"
import { Link } from "react-router-dom"
import AppRoute from "../router/AppRoute"
import "./AppShell.css"
import If from "../flowcontrol/If"
import AppShellContext, { ShellAction, ShellConfig } from "./AppShellContext"
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
            className: "shell-menu-item",
        }))
    }

    private renderActionButton(action: ShellAction): React.ReactNode {
        const { description, type, color, onClick, disabled } = action
        if (typeof onClick === "string") {
            return (
                <Link to={onClick} key={action.description}>
                    <Button variant={type ?? "solid"} color={color ?? "primary"} disabled={disabled}>
                        {description}
                    </Button>
                </Link>
            )
        } else {
            return (
                <Button
                    key={action.description}
                    variant={type ?? "solid"}
                    color={color ?? "primary"}
                    onClick={onClick}
                    disabled={disabled}
                >
                    {description}
                </Button>
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

    shouldComponentUpdate(nextProps: Readonly<AppShellProps>): boolean {
        const { children: previous } = nextProps
        const { children: current } = this.props

        if (previous !== current) {
            this.setState({ config: EmptyConfig })
        }

        return true
    }

    render() {
        const { activeRoute, children, userMenu, serverControl } = this.props
        const { config } = this.state
        const activeMenuItemPath = activeRoute.activeMenuItemPath ?? activeRoute.path
        const { title, subtitle } = config

        return (
            <Layout className="shell-container">
                <Sider trigger={null} width={250} className="shell-sider-container">
                    <div className="shell-logo">
                        <Link to="/" className="shell-logo-link">
                            nginx ignition
                        </Link>
                    </div>
                    <div className="shell-server-control">{serverControl}</div>
                    <Menu
                        className="shell-menu-container"
                        theme="dark"
                        mode="inline"
                        defaultSelectedKeys={activeMenuItemPath ? [activeMenuItemPath] : undefined}
                        items={this.buildMenuItemsAdapters()}
                    />
                    <div className="shell-user-menu">{userMenu}</div>
                </Sider>
                <Layout>
                    <Flex>
                        <Flex vertical>
                            <h1 className="shell-title">{title}</h1>
                            <If condition={subtitle !== undefined}>
                                <h2 className="shell-subtitle">{subtitle}</h2>
                            </If>
                        </Flex>
                        <Flex className="shell-actions-container">{this.renderActions()}</Flex>
                    </Flex>
                    <Content className="shell-content">
                        <AppShellContext.Provider
                            value={{
                                updateConfig: (config: ShellConfig) => this.setState({ config }),
                            }}
                        >
                            {children}
                        </AppShellContext.Provider>
                    </Content>
                </Layout>
            </Layout>
        )
    }
}
