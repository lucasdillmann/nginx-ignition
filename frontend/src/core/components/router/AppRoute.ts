import React from "react";

export default interface AppRoute {
    path: string
    requiresAuthentication: boolean
    fullPage: boolean
    component: React.ReactElement
}
