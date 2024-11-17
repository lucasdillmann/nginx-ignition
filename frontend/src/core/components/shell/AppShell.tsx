import React, {CSSProperties, PropsWithChildren} from "react";
import {MenuFoldOutlined, MenuUnfoldOutlined} from "@ant-design/icons";
import {Button, Layout, Menu} from "antd";
import LocalStorageRepository from "../../repository/LocalStorageRepository";
import {MenuItemType} from "antd/es/menu/interface";
import {NavLink} from "react-router-dom";
import NginxService from "../../../domain/nginx/NginxService";
import NginxStatus from "./NginxStatus";
const {Header, Sider, Content} = Layout;

const styles: { [key: string]: CSSProperties } = {
    container: {
        position: "absolute",
        width: "100%",
        height: "100%",
        top: 0,
        left: 0,
    },
    header: {
        padding: 0,
        background: "#FFF",
    },
    content: {
        margin: "24px 16px",
        padding: 24,
        minHeight: 280,
        background: "#FFF",
        borderRadius: 4,
    },
    toggleButton: {
        fontSize: "16px",
        width: 64,
        height: 64,
    },
    logo: {
        color: "#FFF",
        fontSize: 18,
        padding: "30px 20px 20px",
    },
    logoLink: {
        color: "#FFF",
    },
}

interface AppShellState {
    collapsed: boolean,
}

export interface AppShellMenuItem {
    icon: React.ReactNode
    description: string
    path: string
}

export interface AppShellProps extends PropsWithChildren {
    menuItems: AppShellMenuItem[]
    activeMenuItemPath?: string
}

export default class AppShell extends React.Component<AppShellProps, AppShellState> {
    private repository: LocalStorageRepository<boolean>
    private nginxService: NginxService

    constructor(props: AppShellProps) {
        super(props);
        this.repository = new LocalStorageRepository("nginxIgnition.shell.collapsed")
        this.nginxService = new NginxService()
        this.state = {
            collapsed: this.repository.getOrDefault(false),
        }
    }

    private toggleCollapsed() {
        const {collapsed} = this.state
        this.repository.set(!collapsed)
        this.setState({collapsed: !collapsed})
    }

    private buildMenuItemsAdapters(): MenuItemType[] {
        const {menuItems} = this.props
        return menuItems.map(item => ({
            key: item.path,
            icon: item.icon,
            label: <NavLink to={item.path}>{item.description}</NavLink>,
        }))
    }

    render() {
        const {collapsed} = this.state
        const {children, activeMenuItemPath} = this.props

        return (
            <Layout style={styles.container}>
                <Sider trigger={null} width={250} collapsible collapsed={collapsed}>
                    <div style={styles.logo}>
                        <NavLink to="/" style={styles.logoLink}>
                            nginx ignition
                        </NavLink>
                    </div>
                    <div>
                        <NginxStatus />
                    </div>
                    <Menu
                        theme="dark"
                        mode="inline"
                        defaultSelectedKeys={activeMenuItemPath ? [activeMenuItemPath] : undefined}
                        items={this.buildMenuItemsAdapters()}
                    />
                </Sider>
                <Layout>
                    <Header style={styles.header}>
                        <Button
                            type="text"
                            icon={collapsed ? <MenuUnfoldOutlined/> : <MenuFoldOutlined/>}
                            onClick={() => this.toggleCollapsed()}
                            style={styles.toggleButton}
                        />
                    </Header>
                    <Content style={styles.content}>
                        {children}
                    </Content>
                </Layout>
            </Layout>
        );
    }
}
