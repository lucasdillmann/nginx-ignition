import AppRoute from "./AppRoute";
import React from "react";
import AppContext from "../context/AppContext";
import {RouteObject, createBrowserRouter, RouterProvider, Navigate} from "react-router-dom";

export interface AppRouterProps {
    routes: AppRoute[]
}

export default class AppRouter extends React.Component<AppRouterProps> {
    static contextType = AppContext
    context!: React.ContextType<typeof AppContext>

    private buildRouteComponent(route: AppRoute): any {
        const {component, requiresAuthentication} = route
        const {user} = this.context

        if (requiresAuthentication && user?.id == null) {
            return <Navigate to="/login" replace />
        }

        return component
    }

    private buildRouteAdapters(): RouteObject[] {
        const {routes} = this.props
        return routes.map(route => ({
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
