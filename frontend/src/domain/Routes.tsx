import React from "react"
import LoginPage from "./authentication/LoginPage"
import HomePage from "./home/HomePage"
import AppRoute from "../core/components/router/AppRoute"
import OnboardingPage from "./onboarding/OnboardingPage"
import {
    AuditOutlined,
    BlockOutlined,
    FileProtectOutlined,
    FileSearchOutlined,
    HddOutlined,
    MergeCellsOutlined,
    SettingOutlined,
    TeamOutlined,
    DownloadOutlined,
    ApartmentOutlined,
} from "@ant-design/icons"
import HostListPage from "./host/HostListPage"
import HostFormPage from "./host/HostFormPage"
import CertificateListPage from "./certificate/CertificateListPage"
import CertificateIssuePage from "./certificate/CertificateIssuePage"
import CertificateDetailsPage from "./certificate/CertificateDetailsPage"
import LogsPage from "./logs/LogsPage"
import UserListPage from "./user/UserListPage"
import UserFormPage from "./user/UserFormPage"
import SettingsPage from "./settings/SettingsPage"
import NotFoundPage from "./notfound/NotFoundPage"
import AccessListFormPage from "./accesslist/AccessListFormPage"
import AccessListListPage from "./accesslist/AccessListListPage"
import StreamFormPage from "./stream/StreamFormPage"
import StreamListPage from "./stream/StreamListPage"
import ExportPage from "./export/ExportPage"
import IntegrationListPage from "./integration/IntegrationListPage"
import IntegrationFormPage from "./integration/IntegrationFormPage"
import VpnListPage from "./vpn/VpnListPage"
import VpnFormPage from "./vpn/VpnFormPage"

const Routes: AppRoute[] = [
    {
        path: "/login",
        requiresAuthentication: false,
        fullPage: true,
        component: <LoginPage />,
    },
    {
        path: "/onboarding",
        requiresAuthentication: false,
        fullPage: true,
        component: <OnboardingPage />,
    },
    {
        path: "/hosts/:id",
        requiresAuthentication: true,
        fullPage: false,
        component: <HostFormPage />,
        activeMenuItemPath: "/hosts",
    },
    {
        path: "/hosts",
        requiresAuthentication: true,
        fullPage: false,
        component: <HostListPage />,
        menuItem: {
            description: "Hosts",
            icon: <HddOutlined />,
        },
    },
    {
        path: "/streams/:id",
        requiresAuthentication: true,
        fullPage: false,
        component: <StreamFormPage />,
        activeMenuItemPath: "/streams",
    },
    {
        path: "/streams",
        requiresAuthentication: true,
        fullPage: false,
        component: <StreamListPage />,
        menuItem: {
            description: "Streams",
            icon: <MergeCellsOutlined />,
        },
    },
    {
        path: "/certificates/new",
        requiresAuthentication: true,
        fullPage: false,
        component: <CertificateIssuePage />,
        activeMenuItemPath: "/certificates",
    },
    {
        path: "/certificates/:id",
        requiresAuthentication: true,
        fullPage: false,
        component: <CertificateDetailsPage />,
        activeMenuItemPath: "/certificates",
    },
    {
        path: "/certificates",
        requiresAuthentication: true,
        fullPage: false,
        component: <CertificateListPage />,
        menuItem: {
            description: "SSL certificates",
            icon: <AuditOutlined />,
        },
    },
    {
        path: "/logs",
        requiresAuthentication: true,
        fullPage: false,
        component: <LogsPage />,
        menuItem: {
            description: "Logs",
            icon: <FileSearchOutlined />,
        },
    },
    {
        path: "/integrations/:id",
        requiresAuthentication: true,
        fullPage: false,
        component: <IntegrationFormPage />,
        activeMenuItemPath: "/integrations",
    },
    {
        path: "/integrations",
        requiresAuthentication: true,
        fullPage: false,
        component: <IntegrationListPage />,
        menuItem: {
            description: "Integrations",
            icon: <BlockOutlined />,
        },
    },
    {
        path: "/vpns/:id",
        requiresAuthentication: true,
        fullPage: false,
        component: <VpnFormPage />,
        activeMenuItemPath: "/vpns",
    },
    {
        path: "/vpns",
        requiresAuthentication: true,
        fullPage: false,
        component: <VpnListPage />,
        menuItem: {
            description: "VPNs",
            icon: <ApartmentOutlined />,
        },
    },
    {
        path: "/access-lists/:id",
        requiresAuthentication: true,
        fullPage: false,
        component: <AccessListFormPage />,
        activeMenuItemPath: "/access-lists",
    },
    {
        path: "/access-lists",
        requiresAuthentication: true,
        fullPage: false,
        component: <AccessListListPage />,
        menuItem: {
            description: "Access lists",
            icon: <FileProtectOutlined />,
        },
    },
    {
        path: "/export",
        requiresAuthentication: true,
        fullPage: false,
        component: <ExportPage />,
        menuItem: {
            description: "Export and backup",
            icon: <DownloadOutlined />,
        },
    },
    {
        path: "/settings",
        requiresAuthentication: true,
        fullPage: false,
        component: <SettingsPage />,
        menuItem: {
            description: "Settings",
            icon: <SettingOutlined />,
        },
    },
    {
        path: "/users/:id",
        requiresAuthentication: true,
        fullPage: false,
        component: <UserFormPage />,
        activeMenuItemPath: "/users",
    },
    {
        path: "/users",
        requiresAuthentication: true,
        fullPage: false,
        component: <UserListPage />,
        menuItem: {
            description: "Users",
            icon: <TeamOutlined />,
        },
    },
    {
        path: "/",
        requiresAuthentication: true,
        fullPage: false,
        component: <HomePage />,
    },
    {
        path: "*",
        requiresAuthentication: false,
        fullPage: true,
        component: <NotFoundPage />,
    },
]

export default Routes
