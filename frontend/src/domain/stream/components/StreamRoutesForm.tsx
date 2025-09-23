import { StreamRoute } from "../model/StreamRequest"
import ValidationResult from "../../../core/validation/ValidationResult"
import React from "react"
import { StreamRouteDefaults } from "../StreamFormDefaults"
import StreamRouteForm from "./StreamRouteForm"
import { PlusOutlined } from "@ant-design/icons"
import { Button } from "antd"
import If from "../../../core/components/flowcontrol/If"

export interface StreamRoutesFormProps {
    routes: StreamRoute[]
    validationResult: ValidationResult
    onChange: (routes: StreamRoute[]) => void
}

export default class StreamRoutesForm extends React.Component<StreamRoutesFormProps> {
    private removeRoute(index: number) {
        const { routes, onChange } = this.props
        const newRoutes = [...routes]
        newRoutes.splice(index, 1)
        onChange(newRoutes)
    }

    private addRoute() {
        const { routes, onChange } = this.props
        onChange([...routes, { ...StreamRouteDefaults }])
    }

    private renderRoute(route: StreamRoute, index: number): React.ReactNode {
        const { routes, validationResult } = this.props
        const disableDelete = routes.length <= 1
        const renderSeparator = routes.length > 1 && index < routes.length - 1

        return (
            <>
                <StreamRouteForm
                    key={index}
                    route={route}
                    index={index}
                    validationResult={validationResult}
                    onRemove={disableDelete ? undefined : () => this.removeRoute(index)}
                />

                <If condition={renderSeparator}>
                    <div style={{ marginBottom: 10 }} />
                </If>
            </>
        )
    }

    private renderRoutes() {
        const { routes } = this.props

        return routes.map((route, index) => this.renderRoute(route, index))
    }

    render() {
        return (
            <>
                <h2 className="streams-form-section-name">Routes</h2>
                <p className="streams-form-section-help-text">
                    Routes to be configured in the stream. The nginx will lookup for the SNI (Server Name Indication) in
                    TLS connections and route the request to the first backend that matches the domain name. When either
                    no SNI value is available or no match is found, the request will be forwarded to the default
                    backend.
                </p>
                {this.renderRoutes()}
                <p>
                    <Button type="dashed" onClick={() => this.addRoute()} icon={<PlusOutlined />}>
                        Add routing group
                    </Button>
                </p>
            </>
        )
    }
}
