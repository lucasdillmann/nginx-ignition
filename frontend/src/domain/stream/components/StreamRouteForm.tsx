import { StreamRoute } from "../model/StreamRequest"
import ValidationResult from "../../../core/validation/ValidationResult"
import React from "react"
import { Button, Card, Flex } from "antd"
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

    private buildCardHeader(): { title?: string; extra?: React.ReactNode } {
        const { onRemove, index } = this.props
        if (onRemove === undefined) return {}

        return {
            title: `Routing group ${index + 1}`,
            extra: (
                <Button onClick={onRemove}>
                    <DeleteOutlined /> Remove group
                </Button>
            ),
        }
    }

    render() {
        const { validationResult, index, route } = this.props
        const { title, extra } = this.buildCardHeader()

        return (
            <Card key={`route-${index}`} title={title} extra={extra}>
                <Flex style={{ flexGrow: 1, flexShrink: 1 }}>
                    <Flex style={{ width: "45%" }} vertical>
                        <DomainNamesList
                            pathPrefix={{
                                merged: `routes[${index}]`,
                                name: ["routes", index],
                            }}
                            validationResult={validationResult}
                            expandedLabelSize
                            disableTitle
                        />
                    </Flex>
                    <Flex style={{ width: "55%" }} vertical>
                        <StreamRouteBackendList
                            routeIndex={index}
                            backends={route.backends}
                            validationResult={validationResult}
                            path={["routes", index]}
                        />
                    </Flex>
                </Flex>
            </Card>
        )
    }
}
