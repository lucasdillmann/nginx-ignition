import AppRoute from "./AppRoute"
import React from "react"
import AppContext from "../context/AppContext"
import { RouteObject, createBrowserRouter, RouterProvider, Navigate, Params } from "react-router-dom"
import AppShell, { AppShellMenuItem } from "../shell/AppShell"
import ErrorBoundary from "../errorboundary/ErrorBoundary"
import qs, { ParsedQs } from "qs"
import { buildLoginUrl } from "../../authentication/buildLoginUrl"

let currentInstance: AppRouter

function router() {
    return currentInstance.state.router
}

export function navigateTo(destination: string, replace?: boolean) {
    return router()?.navigate(destination, { replace })
}

export function routeParams(): Params {
    return router()?.state.matches[0]?.params ?? {}
}

export function queryParams(): ParsedQs {
    let queryString = router()?.state.location.search ?? ""
    if (queryString.startsWith("?")) queryString = queryString.substring(1)

    return qs.parse(queryString)
}

interface AppRouterState {
    router?: any
}

export interface AppRouterProps {
    routes: AppRoute[]
    userMenu: React.ReactNode
    serverControl: React.ReactNode
}

export default class AppRouter extends React.Component<AppRouterProps, AppRouterState> {
    constructor(props: AppRouterProps) {
        super(props)
        this.state = {}
        currentInstance = this
    }

    private buildMenuItemsAdapter(): AppShellMenuItem[] {
        const { routes } = this.props
        return routes
            .filter(route => route.menuItem !== undefined)
            .map(route => ({
                icon: route.menuItem!!.icon,
                description: route.menuItem!!.description,
                path: route.path,
            }))
    }

    private buildRouteComponent(route: AppRoute): any {
        const { userMenu, serverControl } = this.props
        const { component, requiresAuthentication, fullPage } = route
        const { user } = AppContext.get()

        if (requiresAuthentication && user?.id == null) {
            return <Navigate to={buildLoginUrl()} replace />
        }

        if (fullPage) {
            return component
        }

        return (
            <ErrorBoundary>
                <AppShell
                    menuItems={this.buildMenuItemsAdapter()}
                    activeRoute={route}
                    userMenu={userMenu}
                    serverControl={serverControl}
                >
                    {component}
                </AppShell>
            </ErrorBoundary>
        )
    }

    private buildRouteAdapters(): RouteObject[] {
        const { routes } = this.props
        return routes.map(route => ({
            path: route.path,
            element: this.buildRouteComponent(route),
            hasErrorBoundary: true,
        }))
    }

    componentDidMount() {
        if (this.state.router !== undefined) return

        const routes = this.buildRouteAdapters()
        const router = createBrowserRouter(routes, { window })
        this.setState({ router })
    }

    render() {
        const { router } = this.state
        if (router === undefined) return <></>

        return <RouterProvider router={router} />
    }
}
