import LoginPage from "./authentication/LoginPage";
import HomePage from "./home/HomePage";
import AppRoute from "../core/components/router/AppRoute";
import OnboardingPage from "./onboarding/OnboardingPage";

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
        path: "/",
        requiresAuthentication: true,
        fullPage: false,
        component: <HomePage />,
    }
]

export default Routes
