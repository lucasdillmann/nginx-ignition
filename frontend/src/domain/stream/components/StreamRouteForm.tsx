import { StreamRoute } from "../model/StreamRequest"
import ValidationResult from "../../../core/validation/ValidationResult"
import React from "react"
import { Button, Card, Flex } from "antd"
import DomainNamesList from "../../../core/components/domainnames/DomainNamesList"
import StreamRouteBackendList from "./StreamRouteBackendList"
import { DeleteOutlined } from "@ant-design/icons"
import { I18n } from "../../../core/i18n/I18n"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"

export interface StreamRouteFormProps {
    route: StreamRoute
    index: number
    validationResult: ValidationResult
    onRemove?: () => void
}

export default class StreamRouteForm extends React.Component<StreamRouteFormProps> {
    private buildCardHeader(): { title?: React.ReactNode; extra?: React.ReactNode } {
        const { onRemove, index } = this.props

        return {
            title: <I18n id={MessageKey.FrontendStreamComponentsRouteformGroupTitle} params={{ index: index + 1 }} />,
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
                        <p>
                            <I18n id={MessageKey.FrontendStreamComponentsRouteformDomains} />
                        </p>
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
                        <p>
                            <I18n id={MessageKey.FrontendStreamComponentsRouteformBackends} />
                        </p>
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
