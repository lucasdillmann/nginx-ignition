import LoginPage from "./authentication/LoginPage";
import HomePage from "./home/HomePage";
import AppRoute from "../core/components/router/AppRoute";
import OnboardingPage from "./onboarding/OnboardingPage";
import {HddOutlined, FileProtectOutlined, FileSearchOutlined, TeamOutlined} from "@ant-design/icons"

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
        path: "/hosts",
        requiresAuthentication: true,
        fullPage: false,
        component: <HomePage />,
        menuItem: {
            description: "Hosts",
            icon: <HddOutlined />,
        },
    },
    {
        path: "/certificates",
        requiresAuthentication: true,
        fullPage: false,
        component: <HomePage />,
        menuItem: {
            description: "SSL certificates",
            icon: <FileProtectOutlined />,
        },
    },
    {
        path: "/logs",
        requiresAuthentication: true,
        fullPage: false,
        component: <HomePage />,
        menuItem: {
            description: "Logs",
            icon: <FileSearchOutlined />,
        },
    },
    {
        path: "/users",
        requiresAuthentication: true,
        fullPage: false,
        component: <HomePage />,
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
]

export default Routes
