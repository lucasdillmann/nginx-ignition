import { StreamRoute } from "../model/StreamRequest"
import ValidationResult from "../../../core/validation/ValidationResult"
import React from "react"
import { Card, Flex } from "antd"
import DomainNamesList from "../../../core/components/domainnames/DomainNamesList"
import StreamRouteBackendList from "./StreamRouteBackendList"
import { DeleteOutlined } from "@ant-design/icons"

export interface StreamRouteFormProps {
    route: StreamRoute
    index: number
    validationResult: ValidationResult
    onChange: (routes: StreamRoute) => void
    onRemove?: () => void
}

export default class StreamRouteForm extends React.Component<StreamRouteFormProps> {
    private handleChange(attribute: string, value: any) {
        const { route, onChange } = this.props
        onChange({
            ...route,
            [attribute]: value,
        })
    }

    render() {
        const { validationResult, index, route, onRemove } = this.props
        const removeButton =
            onRemove !== undefined
                ? [
                      <div key={0} style={{ textAlign: "right", width: "100%", paddingRight: 20 }}>
                          <span onClick={onRemove}>
                              <DeleteOutlined /> Remove route
                          </span>
                      </div>,
                  ]
                : undefined

        return (
            <Card actions={removeButton}>
                <Flex style={{ flexGrow: 1, flexShrink: 1 }}>
                    <Flex style={{ width: "35%", paddingRight: 10 }} vertical>
                        <DomainNamesList
                            pathPrefix={{
                                merged: `routes[${index}]`,
                                name: ["routes", index],
                            }}
                            validationResult={validationResult}
                        />
                    </Flex>
                    <Flex style={{ width: "65%" }} vertical>
                        <StreamRouteBackendList
                            routeIndex={index}
                            backends={route.backends}
                            validationResult={validationResult}
                            onChange={backends => this.handleChange("backends", backends)}
                        />
                    </Flex>
                </Flex>
            </Card>
        )
    }
}
