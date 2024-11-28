import AppRoute from "./AppRoute";
import React from "react";
import AppContext, {AppContextData} from "../context/AppContext";
import {RouteObject, createBrowserRouter, RouterProvider, Navigate, Params} from "react-router-dom";
import AppShell, {AppShellMenuItem} from "../shell/AppShell";
import ErrorBoundary from "../errorboundary/ErrorBoundary";
import {Router} from "@remix-run/router/dist/router";
import qs, {ParsedQs} from "qs";

let currentInstance: AppRouter

function router() {
    return currentInstance.state.router
}

export function navigateTo(destination: string) {
    return router()?.navigate(destination)
}

export function routeParams(): Params {
    return router()?.state.matches[0]?.params ?? {}
}

export function queryParams(): ParsedQs {
    let queryString = router()?.state.location.search ?? ""
    if (queryString.startsWith("?"))
        queryString = queryString.substring(1)

    return qs.parse(queryString)
}

interface AppRouterState {
    router?: Router
}

export interface AppRouterProps {
    routes: AppRoute[]
}

export default class AppRouter extends React.Component<AppRouterProps, AppRouterState> {
    static contextType = AppContext
    context!: React.ContextType<typeof AppContext>

    constructor(props: AppRouterProps, context: AppContextData) {
        super(props, context);
        this.state = {}
        currentInstance = this
    }

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
            <ErrorBoundary>
                <AppShell
                    menuItems={this.buildMenuItemsAdapter()}
                    activeRoute={route}
                >
                    {component}
                </AppShell>
            </ErrorBoundary>
        )
    }

    private buildRouteAdapters(): RouteObject[] {
        const {routes} = this.props
        return routes
            .filter(route => this.isRouteVisible(route))
            .map(route => ({
                path: route.path,
                element: this.buildRouteComponent(route),
                hasErrorBoundary: true,
            }))
    }

    componentDidMount() {
        if (this.state.router !== undefined)
            return

        const routes = this.buildRouteAdapters()
        const router = createBrowserRouter(routes, { window })
        this.setState({router})
    }

    render() {
        const {router} = this.state
        if (router === undefined)
            return <></>

        return <RouterProvider router={router}/>
    }
}
