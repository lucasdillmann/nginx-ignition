import { StreamAddress, StreamProtocol } from "../model/StreamRequest"
import { Flex, Form, Input, InputNumber, Select, Space } from "antd"
import If from "../../../core/components/flowcontrol/If"
import React from "react"
import ValidationResult from "../../../core/validation/ValidationResult"
import CompatibleStreamProtocolResolver from "../utils/CompatibleStreamProtocolResolver"
import { I18n, i18n } from "../../../core/i18n/I18n"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"

const INPUT_STYLE = {
    width: 150,
    maxWidth: 150,
    minWidth: 150,
}

export interface StreamAddressInputProps {
    basePath: string
    validationResult: ValidationResult
    address: StreamAddress
    onChange: (address: StreamAddress) => void
    parentProtocol?: StreamProtocol
}

export default class StreamAddressInput extends React.Component<StreamAddressInputProps> {
    private handleChange(attribute: string, value: any) {
        const { address, onChange } = this.props

        onChange({
            ...address,
            [attribute]: value,
        })
    }

    private protocolLabel(protocol: StreamProtocol) {
        switch (protocol) {
            case StreamProtocol.TCP:
                return <I18n id={MessageKey.CommonProtocolTcp} />
            case StreamProtocol.UDP:
                return <I18n id={MessageKey.FrontendStreamComponentsAddressinputProtocolUdp} />
            case StreamProtocol.SOCKET:
                return <I18n id={MessageKey.FrontendStreamComponentsAddressinputProtocolSocket} />
        }
    }

    private buildOptions(): any[] {
        const { parentProtocol } = this.props
        const possibleProtocols = parentProtocol
            ? CompatibleStreamProtocolResolver.resolve(parentProtocol)
            : Object.values(StreamProtocol)

        return possibleProtocols.map(protocol => ({
            value: protocol,
            label: this.protocolLabel(protocol),
        }))
    }

    render() {
        const { basePath, address, validationResult } = this.props
        const backendSocket = address.protocol === StreamProtocol.SOCKET

        const validationStatus =
            validationResult.getStatus(`${basePath}.protocol`) ??
            validationResult.getStatus(`${basePath}.address`) ??
            validationResult.getStatus(`${basePath}.port`)
        const validationMessage =
            validationResult.getMessage(`${basePath}.protocol`) ??
            validationResult.getMessage(`${basePath}.address`) ??
            validationResult.getMessage(`${basePath}.port`)

        return (
            <Form.Item validateStatus={validationStatus} noStyle>
                <Flex style={{ flexDirection: "column", flexGrow: 1, width: "100%" }}>
                    <Space direction="vertical" style={{ flexGrow: 1 }}>
                        <Space.Compact block>
                            <Select
                                value={address.protocol}
                                onChange={value => this.handleChange("protocol", value)}
                                style={INPUT_STYLE}
                                options={this.buildOptions()}
                            />
                            <Input
                                value={address.address}
                                onChange={event => this.handleChange("address", event.target.value)}
                                placeholder={
                                    backendSocket
                                        ? i18n(MessageKey.FrontendStreamComponentsAddressinputSocketPlaceholder)
                                        : i18n(MessageKey.FrontendStreamComponentsAddressinputAddressPlaceholder)
                                }
                            />
                            <If condition={!backendSocket}>
                                <Space.Addon>:</Space.Addon>
                                <InputNumber
                                    placeholder={i18n(MessageKey.CommonPort)}
                                    style={INPUT_STYLE}
                                    value={address.port}
                                    onChange={value => this.handleChange("port", value)}
                                    min={1}
                                    max={65535}
                                />
                            </If>
                        </Space.Compact>
                    </Space>
                    <If condition={validationMessage !== undefined}>
                        <div className="ant-form-item-additional">
                            <div
                                id="description_help"
                                className="ant-form-item-explain ant-form-item-explain-connected"
                                style={{ color: "var(--nginxIgnition-colorError)" }}
                            >
                                <div className="ant-form-item-explain-error">{validationMessage}</div>
                            </div>
                        </div>
                    </If>
                </Flex>
            </Form.Item>
        )
    }
}
