import React from "react";
import {Button, Flex, Layout, Menu} from "antd";
import {MenuItemType} from "antd/es/menu/interface";
import {Link, NavLink} from "react-router-dom";
import NginxControl from "../nginx/NginxControl";
import AppRoute from "../router/AppRoute";
import "./AppShell.css"
import ShellAwareComponent, {ShellAction, ShellConfig} from "./ShellAwareComponent";
import If from "../flowcontrol/If";
const {Sider, Content} = Layout;

export interface AppShellMenuItem {
    icon: React.ReactNode
    description: string
    path: string
}

export interface AppShellProps {
    menuItems: AppShellMenuItem[]
    activeRoute: AppRoute
    children: React.ReactElement
}

interface AppShellState {
    config: ShellConfig
}

export default class AppShell extends React.Component<AppShellProps, AppShellState> {
    private childRef: React.RefObject<ShellAwareComponent>

    constructor(props: AppShellProps) {
        super(props);
        this.childRef = React.createRef()
        this.state = {
            config: {
                title: "",
            }
        }
    }

    private buildMenuItemsAdapters(): MenuItemType[] {
        const {menuItems} = this.props
        return menuItems.map(item => ({
            key: item.path,
            icon: item.icon,
            label: <NavLink to={item.path}>{item.description}</NavLink>,
            className: "shell-menu-item",
        }))
    }

    private renderActionButton(action: ShellAction): React.ReactNode {
        const {description, type, color, onClick} = action
        if (typeof onClick === "string") {
            return (
                <Link to={onClick}>
                    <Button
                        variant={type ?? "solid"}
                        color={color ?? "primary"}>
                        {description}
                    </Button>
                </Link>
            )
        } else {
            return (
                <Button
                    variant={type ?? "solid"}
                    color={color ?? "primary"}
                    onClick={() => onClick()}>
                    {description}
                </Button>
            )
        }
    }

    private renderActions() {
        const {config: {actions}} = this.state
        if (!actions)
            return null

        return (
            <>
                {actions.map(action => this.renderActionButton(action))}
            </>
        )
    }

    private applyShellConfig() {
        this.setState({
            config: this.childRef.current!!.shellConfig(),
        })
    }

    componentDidMount() {
       this.applyShellConfig()
    }

    componentDidUpdate(prevProps: Readonly<AppShellProps>) {
        const {children: previous} = prevProps
        const {children: current} = this.props

        if (previous !== current) {
            this.applyShellConfig()
        }
    }

    render() {
        const {activeRoute, children} = this.props
        const {config} = this.state
        const activeMenuItemPath = activeRoute.activeMenuItemPath ?? activeRoute.path
        const {title, subtitle} = config
        const enhancedChild = React.cloneElement(children, { ref: this.childRef })

        return (
            <Layout className="shell-container">
                <Sider trigger={null} width={250}>
                    <div className="shell-logo">
                        <NavLink to="/" className="shell-logo-link">
                            nginx ignition
                        </NavLink>
                    </div>
                    <div>
                        <NginxControl />
                    </div>
                    <Menu
                        className="shell-menu-container"
                        theme="dark"
                        mode="inline"
                        defaultSelectedKeys={activeMenuItemPath ? [activeMenuItemPath] : undefined}
                        items={this.buildMenuItemsAdapters()}
                    />
                </Sider>
                <Layout>
                    <Flex>
                        <Flex vertical>
                            <h1 className="shell-title">{title}</h1>
                            <If condition={subtitle !== undefined}>
                                <h2 className="shell-subtitle">{subtitle}</h2>
                            </If>
                        </Flex>
                        <Flex className="shell-actions-container">
                            {this.renderActions()}
                        </Flex>
                    </Flex>
                    <Content className="shell-content">
                        {enhancedChild}
                    </Content>
                </Layout>
            </Layout>
        );
    }
}
