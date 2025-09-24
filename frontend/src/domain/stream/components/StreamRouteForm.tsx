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
    onRemove?: () => void
}

export default class StreamRouteForm extends React.Component<StreamRouteFormProps> {
    private buildCardHeader(): { title?: string; extra?: React.ReactNode } {
        const { onRemove, index } = this.props

        return {
            title: `Group ${index + 1}`,
            extra: (
                <Button onClick={onRemove} disabled={onRemove === undefined}>
                    <DeleteOutlined />
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
                        <p>Domains</p>
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
                        <p>Backends</p>
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
