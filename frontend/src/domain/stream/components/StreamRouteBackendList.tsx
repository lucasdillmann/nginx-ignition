import React from "react"
import { StreamBackend } from "../model/StreamRequest"
import ValidationResult from "../../../core/validation/ValidationResult"
import { Flex, Form } from "antd"
import StreamAddressInput from "./StreamAddressInput"
import { SettingOutlined } from "@ant-design/icons"
import StreamBackendSettingsModal from "./StreamBackendSettingsModal"

export interface StreamRouteBackendListProps {
    routeIndex: number
    backends: StreamBackend[]
    validationResult: ValidationResult
    onChange: (backends: StreamBackend[]) => void
    openSettingsModalIndex?: number
}

export default class StreamRouteBackendList extends React.PureComponent<StreamRouteBackendListProps> {
    private changeOpenModal(index?: number) {
        this.setState({ openSettingsModalIndex: index })
    }

    private renderBackend(backend: StreamBackend, index: number) {
        const { validationResult, openSettingsModalIndex } = this.props

        return (
            <>
                <Flex style={{ flexGrow: 1 }}>
                    <Flex style={{ flexGrow: 1, alignContent: "center", flexShrink: 1 }}>
                        <StreamAddressInput
                            basePath="defaultBackend.target"
                            validationResult={validationResult}
                            address={backend.target}
                            onChange={() => {}}
                        />
                    </Flex>
                    <Flex style={{ marginLeft: "20px", flexShrink: 1 }}>
                        <SettingOutlined onClick={() => this.changeOpenModal(index)} size={10} />
                    </Flex>
                </Flex>
                <StreamBackendSettingsModal
                    backend={backend}
                    open={openSettingsModalIndex === index}
                    validationBasePath="defaultBackend"
                    validationResult={validationResult}
                    onClose={() => this.changeOpenModal(undefined)}
                    onChange={() => {}}
                    hideWeight
                />
            </>
        )
    }

    render() {
        const { backends, routeIndex } = this.props

        return (
            <Form.Item label="Backend servers" key={routeIndex} required>
                {backends.map((backend, index) => this.renderBackend(backend, index))}
            </Form.Item>
        )
    }
}
