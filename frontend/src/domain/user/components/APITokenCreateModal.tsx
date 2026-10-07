import React, { createRef } from "react"
import { Alert, Button, DatePicker, Flex, Form, FormInstance, Input, Modal } from "antd"
import { CopyOutlined } from "@ant-design/icons"
import dayjs, { Dayjs } from "dayjs"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"
import { I18n, i18n } from "../../../core/i18n/I18n"
import ValidationResult from "../../../core/validation/ValidationResult"
import ValidationResultConverter from "../../../core/validation/ValidationResultConverter"
import { UnexpectedResponseError } from "../../../core/apiclient/ApiResponse"
import Notification from "../../../core/components/notification/Notification"
import FormLayout from "../../../core/components/form/FormLayout"
import UserService from "../UserService"
import APITokenRequest from "../model/APITokenRequest"
import "./APITokenCreateModal.css"

const DEFAULT_EXPIRATION_DAYS = 7

interface APITokenFormValues {
    name: string
    expiration?: Dayjs
}

export interface APITokenCreateModalProps {
    open: boolean
    onCancel: () => void
    onCreated: () => void
}

export interface APITokenCreateModalState {
    loading: boolean
    step: number
    token: string
    validationResult: ValidationResult
}

export default class APITokenCreateModal extends React.Component<APITokenCreateModalProps, APITokenCreateModalState> {
    private readonly formRef: React.RefObject<FormInstance<APITokenFormValues> | null>
    private readonly service: UserService

    constructor(props: APITokenCreateModalProps) {
        super(props)
        this.service = new UserService()
        this.formRef = createRef()
        this.state = {
            loading: false,
            step: 0,
            token: "",
            validationResult: new ValidationResult(),
        }
    }

    componentDidUpdate(prevProps: Readonly<APITokenCreateModalProps>) {
        if (this.props.open && !prevProps.open) {
            this.reset()
        }
    }

    private reset() {
        this.formRef.current?.resetFields()
        this.setState({
            loading: false,
            step: 0,
            token: "",
            validationResult: new ValidationResult(),
        })
    }

    private async createToken(values: APITokenFormValues) {
        this.setState({ validationResult: new ValidationResult(), loading: true })

        const request: APITokenRequest = {
            name: values.name,
            expiration: values.expiration?.toISOString(),
        }

        return this.service
            .createToken(request)
            .then(response => {
                this.setState({ step: 1, token: response.token })
                this.props.onCreated()
            })
            .catch(error => this.handleErrorResponse(error))
            .then(() => this.setState({ loading: false }))
    }

    private handleErrorResponse(error: Error) {
        if (error instanceof UnexpectedResponseError) {
            const validationResult = ValidationResultConverter.parse(error.response)
            if (validationResult != null) this.setState({ validationResult })
        }

        Notification.error(MessageKey.CommonThatDidntWork, MessageKey.CommonFormCheckMessage)
    }

    private async copyToken() {
        try {
            await navigator.clipboard.writeText(this.state.token)
            Notification.success(MessageKey.CommonCopy, MessageKey.FrontendUserTokensCopiedDescription)
        } catch {
            Notification.error(MessageKey.CommonThatDidntWork, MessageKey.CommonUnexpectedErrorTryAgain)
        }
    }

    private renderForm() {
        const { loading, validationResult } = this.state

        return (
            <Form<APITokenFormValues>
                {...FormLayout.FormDefaults}
                labelCol={FormLayout.ExpandedLabeledItem.labelCol}
                wrapperCol={FormLayout.ExpandedLabeledItem.wrapperCol}
                ref={this.formRef}
                layout="vertical"
                initialValues={{ expiration: dayjs().add(DEFAULT_EXPIRATION_DAYS, "day") }}
                className="api-token-create-form"
            >
                <Form.Item
                    name="name"
                    validateStatus={validationResult.getStatus("name")}
                    help={validationResult.getMessage("name")}
                    label={<I18n id={MessageKey.CommonName} />}
                    required
                >
                    <Input />
                </Form.Item>

                <Form.Item
                    name="expiration"
                    validateStatus={validationResult.getStatus("expiration")}
                    help={
                        validationResult.getMessage("expiration") ?? i18n(MessageKey.FrontendUserTokensExpirationHelp)
                    }
                    label={<I18n id={MessageKey.FrontendUserTokensExpirationLabel} />}
                >
                    <DatePicker showTime allowClear className="api-token-create-expiration" />
                </Form.Item>

                <Flex justify="end" style={{ marginTop: 24 }}>
                    <Button type="primary" loading={loading} onClick={() => this.formRef.current?.submit()}>
                        <I18n id={MessageKey.FrontendUserTokensCreateButton} />
                    </Button>
                </Flex>
            </Form>
        )
    }

    private renderGeneratedToken() {
        return (
            <Flex vertical className="api-token-generated-container">
                <Alert
                    type="warning"
                    showIcon
                    message={<I18n id={MessageKey.FrontendUserTokensGeneratedTitle} />}
                    description={<I18n id={MessageKey.FrontendUserTokensGeneratedDescription} />}
                />

                <Flex className="api-token-generated-value-container" onClick={() => this.copyToken()}>
                    <code className="api-token-generated-value">{this.state.token}</code>
                    <CopyOutlined className="api-token-generated-copy-icon" />
                </Flex>

                <Button type="primary" icon={<CopyOutlined />} onClick={() => this.copyToken()}>
                    <I18n id={MessageKey.CommonCopy} />
                </Button>
            </Flex>
        )
    }

    render() {
        const { open, onCancel } = this.props
        const { step } = this.state

        return (
            <Modal
                title={<I18n id={MessageKey.FrontendUserTokensCreateModalTitle} />}
                onCancel={onCancel}
                footer={null}
                open={open}
                destroyOnHidden
            >
                {step === 0 ? this.renderForm() : this.renderGeneratedToken()}
            </Modal>
        )
    }
}
