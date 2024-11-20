import LoginPage from "./authentication/LoginPage";
import HomePage from "./home/HomePage";
import AppRoute from "../core/components/router/AppRoute";
import OnboardingPage from "./onboarding/OnboardingPage";
import {FileProtectOutlined, FileSearchOutlined, HddOutlined, TeamOutlined} from "@ant-design/icons"
import HostListPage from "./host/HostListPage";
import HostFormPage from "./host/HostFormPage";
import CertificateListPage from "./certificate/CertificateListPage";
import CertificateFormPage from "./certificate/CertificateFormPage";
import CertificateDetailsPage from "./certificate/CertificateDetailsPage";
import LogsPage from "./logs/LogsPage";
import UserListPage from "./user/UserListPage";
import UserFormPage from "./user/UserFormPage";
import {UserRole} from "./user/model/UserRole";

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
        title: "Host details",
        requiresAuthentication: true,
        fullPage: false,
        component: <HostFormPage />,
        activeMenuItemPath: "/hosts",
    },
    {
        path: "/hosts",
        title: "Hosts",
        requiresAuthentication: true,
        fullPage: false,
        component: <HostListPage />,
        menuItem: {
            description: "Hosts",
            icon: <HddOutlined />,
        },
    },
    {
        path: "/certificates/new",
        title: "New SSL certificate",
        requiresAuthentication: true,
        fullPage: false,
        component: <CertificateFormPage />,
        activeMenuItemPath: "/certificates",
    },
    {
        path: "/certificates/:id",
        title: "SSL certificate details",
        requiresAuthentication: true,
        fullPage: false,
        component: <CertificateDetailsPage />,
        activeMenuItemPath: "/certificates",
    },
    {
        path: "/certificates",
        title: "SSL certificates",
        requiresAuthentication: true,
        fullPage: false,
        component: <CertificateListPage />,
        menuItem: {
            description: "SSL certificates",
            icon: <FileProtectOutlined />,
        },
    },
    {
        path: "/logs",
        title: "Logs",
        requiresAuthentication: true,
        fullPage: false,
        component: <LogsPage />,
        menuItem: {
            description: "Logs",
            icon: <FileSearchOutlined />,
        },
    },
    {
        path: "/users/:id",
        title: "User details",
        requiresAuthentication: true,
        fullPage: false,
        component: <UserFormPage />,
        visibleRoles: [UserRole.ADMIN],
        activeMenuItemPath: "/users",
    },
    {
        path: "/users",
        title: "Users",
        requiresAuthentication: true,
        fullPage: false,
        component: <UserListPage />,
        visibleRoles: [UserRole.ADMIN],
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
        activeMenuItemPath: "/hosts",
    },
]

export default Routes
