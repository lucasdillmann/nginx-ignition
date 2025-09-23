import React from "react"
import { StreamBackend } from "../model/StreamRequest"
import ValidationResult from "../../../core/validation/ValidationResult"
import { Button, Flex, Form, FormListFieldData, FormListOperation } from "antd"
import StreamAddressInput from "./StreamAddressInput"
import { DeleteOutlined, PlusOutlined, SettingOutlined } from "@ant-design/icons"
import StreamBackendSettingsModal from "./StreamBackendSettingsModal"
import FormLayout from "../../../core/components/form/FormLayout"
import If from "../../../core/components/flowcontrol/If"
import { StreamBackendDefault } from "../StreamFormDefaults"

interface StreamRouteBackendListState {
    openSettingsModalIndex?: number
}

export interface StreamRouteBackendListProps {
    path: any[]
    routeIndex: number
    backends: StreamBackend[]
    validationResult: ValidationResult
}

export default class StreamRouteBackendList extends React.PureComponent<
    StreamRouteBackendListProps,
    StreamRouteBackendListState
> {
    constructor(props: StreamRouteBackendListProps) {
        super(props)
        this.state = {
            openSettingsModalIndex: undefined,
        }
    }

    private changeOpenModal(index?: number) {
        this.setState({ openSettingsModalIndex: index })
    }

    private updateBackend(backend: StreamBackend, index: number, operations: FormListOperation) {
        operations.add(backend, index)
        operations.remove(index + 1)
    }

    private renderBackend(field: FormListFieldData, operations: FormListOperation) {
        const { validationResult, backends, routeIndex } = this.props
        const { openSettingsModalIndex } = this.state
        const index = field.name as number
        const deleteEnabled = backends.length > 1
        const marginBottom = index === backends.length - 1 ? 0 : 25
        const backend = backends[index]

        return (
            <div key={index}>
                <Flex style={{ flexGrow: 1, marginBottom }}>
                    <Flex style={{ flexGrow: 1, alignContent: "center", flexShrink: 1 }}>
                        <StreamAddressInput
                            basePath={`routes[${routeIndex}].backends[${index}].target`}
                            validationResult={validationResult}
                            address={backend.target}
                            onChange={address => this.updateBackend({ ...backend, target: address }, index, operations)}
                        />
                    </Flex>
                    <Flex style={{ marginLeft: "20px", flexShrink: 1 }}>
                        <SettingOutlined onClick={() => this.changeOpenModal(index)} size={10} />
                        <If condition={deleteEnabled}>
                            <DeleteOutlined
                                onClick={() => operations.remove(index)}
                                style={{ marginLeft: 10 }}
                                size={10}
                            />
                        </If>
                    </Flex>
                </Flex>
                <StreamBackendSettingsModal
                    backend={backend}
                    open={openSettingsModalIndex === index}
                    validationBasePath={`routes[${routeIndex}].backends[${index}]`}
                    validationResult={validationResult}
                    onClose={() => this.changeOpenModal(undefined)}
                    onChange={value => this.updateBackend(value, index, operations)}
                />
            </div>
        )
    }

    private renderBackends(fields: FormListFieldData[], operations: FormListOperation) {
        const backendElements = fields.map(field => this.renderBackend(field, operations))
        const addButton = (
            <Form.Item {...FormLayout.ExpandedUnlabeledItem} style={{ marginTop: 25 }}>
                <Button type="dashed" onClick={() => operations.add(StreamBackendDefault)} icon={<PlusOutlined />}>
                    Add backend
                </Button>
            </Form.Item>
        )

        return [...backendElements, addButton]
    }

    render() {
        const { path } = this.props
        return (
            <Form.List name={[...path, "backends"]}>
                {(fields, operations) => this.renderBackends(fields, operations)}
            </Form.List>
        )
    }
}
