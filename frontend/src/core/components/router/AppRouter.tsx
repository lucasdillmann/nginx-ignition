import AppRoute from "./AppRoute";
import React from "react";
import AppContext from "../context/AppContext";
import {RouteObject, createBrowserRouter, RouterProvider, Navigate} from "react-router-dom";
import AppShell, {AppShellMenuItem} from "../shell/AppShell";

export interface AppRouterProps {
    routes: AppRoute[]
}

export default class AppRouter extends React.Component<AppRouterProps> {
    static contextType = AppContext
    context!: React.ContextType<typeof AppContext>

    private isRouteVisible(route: AppRoute): boolean {
        if (!Array.isArray(route.visibleRoles))
            return true

        const {user} = this.context
        if (user === undefined)
            return false

        return route.visibleRoles.includes(user.role)
    }

    private buildMenuItemsAdapter(): AppShellMenuItem[] {
        const {routes} = this.props
        return routes
            .filter(route => route.menuItem !== undefined)
            .filter(route => this.isRouteVisible(route))
            .map(route => ({
                icon: route.menuItem!!.icon,
                description: route.menuItem!!.description,
                path: route.path,
            }))
    }

    private buildRouteComponent(route: AppRoute): any {
        const {component, requiresAuthentication, fullPage} = route
        const {user} = this.context

        if (requiresAuthentication && user?.id == null) {
            return <Navigate to="/login" replace />
        }

        if (fullPage) {
            return component
        }

        return (
            <AppShell
                menuItems={this.buildMenuItemsAdapter()}
                activeRoute={route}
            >
                {component}
            </AppShell>
        )
    }

    private buildRouteAdapters(): RouteObject[] {
        const {routes} = this.props
        return routes
            .filter(route => this.isRouteVisible(route))
            .map(route => ({
                path: route.path,
                element: this.buildRouteComponent(route),
            }))
    }

    render() {
        const routes = this.buildRouteAdapters()
        const router = createBrowserRouter(routes)
        return <RouterProvider router={router}/>
    }
}
