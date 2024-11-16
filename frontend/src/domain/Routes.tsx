import LoginPage from "./authentication/LoginPage";
import HomePage from "./home/HomePage";
import AppRoute from "../core/components/router/AppRoute";

const Routes: AppRoute[] = [
    {
        path: "/login",
        requiresAuthentication: false,
        fullPage: true,
        component: <LoginPage />,
    },
    {
        path: "/",
        requiresAuthentication: true,
        fullPage: false,
        component: <HomePage />,
    }
]

export default Routes
