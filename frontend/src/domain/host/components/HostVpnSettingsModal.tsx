import { Form, FormItemProps, Input, Modal } from "antd"
import React from "react"
import ValidationResult from "../../../core/validation/ValidationResult"
import { I18n } from "../../../core/i18n/I18n"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"

const ItemProps: FormItemProps = {
    labelCol: {
        sm: { span: 8 },
    },
    wrapperCol: {
        sm: { span: 16 },
    },
}

export interface HostVpnSettingsProps {
    open: boolean
    index: number
    fieldPath: string | string[]
    onClose: () => void
    validationResult: ValidationResult
}

export default class HostVpnSettingsModal extends React.Component<HostVpnSettingsProps> {
    render() {
        const { open, onClose, fieldPath, validationResult, index } = this.props
        return (
            <Modal
                title={<I18n id={MessageKey.FrontendHostComponentsVpnBindingSettingsTitle} />}
                open={open}
                afterClose={onClose}
                onCancel={onClose}
                footer={null}
                width={750}
            >
                <div style={{ height: 40 }} />
                <Form.Item
                    {...ItemProps}
                    name={[fieldPath, "host"]}
                    validateStatus={validationResult.getStatus(`vpns[${index}].host`)}
                    help={
                        validationResult.getMessage(`vpns[${index}].host`) ?? (
                            <I18n id={MessageKey.FrontendHostComponentsHostvpnsTargetHostHelp} />
                        )
                    }
                    label={<I18n id={MessageKey.FrontendHostComponentsHostvpnsTargetHost} />}
                >
                    <Input />
                </Form.Item>
            </Modal>
        )
    }
}
