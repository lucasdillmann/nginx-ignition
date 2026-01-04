import React from "react"
import { HostRouteType } from "../model/HostRequest"
import { HostFormRoute } from "../model/HostFormValues"
import HideableFormInput from "../../../core/components/form/HideableFormInput"

const NOT_AVAILABLE_REASON = "Not available for this route type"

export interface HostRouteConditionalConfigProps {
    route: HostFormRoute
    types: HostRouteType[]
    children: any
}

export default class HostRouteConditionalConfig extends React.Component<HostRouteConditionalConfigProps> {
    render() {
        const { route, types, children } = this.props
        return (
            <HideableFormInput hidden={!types.includes(route.type)} reason={NOT_AVAILABLE_REASON}>
                {children}
            </HideableFormInput>
        )
    }
}
