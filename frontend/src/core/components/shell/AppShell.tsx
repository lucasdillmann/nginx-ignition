import React, {PropsWithChildren} from "react";
import {Layout, Menu} from "antd";
import {MenuItemType} from "antd/es/menu/interface";
import {NavLink} from "react-router-dom";
import NginxControl from "../nginx/NginxControl";
import AppRoute from "../router/AppRoute";
import "./AppShell.css"
const {Sider, Content} = Layout;

export interface AppShellMenuItem {
    icon: React.ReactNode
    description: string
    path: string
}

export interface AppShellProps extends PropsWithChildren {
    menuItems: AppShellMenuItem[]
    activeRoute: AppRoute
}

export default class AppShell extends React.Component<AppShellProps> {
    private buildMenuItemsAdapters(): MenuItemType[] {
        const {menuItems} = this.props
        return menuItems.map(item => ({
            key: item.path,
            icon: item.icon,
            label: <NavLink to={item.path}>{item.description}</NavLink>,
            className: "shell-menu-item",
        }))
    }

    render() {
        const {children, activeRoute} = this.props
        const activeMenuItemPath = activeRoute.activeMenuItemPath ?? activeRoute.path

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
                        theme="dark"
                        mode="inline"
                        defaultSelectedKeys={activeMenuItemPath ? [activeMenuItemPath] : undefined}
                        items={this.buildMenuItemsAdapters()}
                    />
                </Sider>
                <Layout>
                    <h1 className="shell-title">{activeRoute.title}</h1>
                    <Content className="shell-content">
                        {children}
                    </Content>
                </Layout>
            </Layout>
        );
    }
}
