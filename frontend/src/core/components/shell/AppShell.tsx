import React, {PropsWithChildren} from "react";
import {Layout, Menu} from "antd";
import {MenuItemType} from "antd/es/menu/interface";
import {NavLink} from "react-router-dom";
import NginxControl from "../nginxcontrol/NginxControl";
import styles from "./AppShell.styles"
import AppRoute from "../router/AppRoute";
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
            style: styles.menuItem,
        }))
    }

    render() {
        const {children, activeRoute} = this.props
        const activeMenuItemPath = activeRoute.activeMenuItemPath ?? activeRoute.path

        return (
            <Layout style={styles.container}>
                <Sider trigger={null} width={250}>
                    <div style={styles.logo}>
                        <NavLink to="/" style={styles.logoLink}>
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
                    <h1 style={styles.title}>{activeRoute.title}</h1>
                    <Content style={styles.content}>
                        {children}
                    </Content>
                </Layout>
            </Layout>
        );
    }
}
