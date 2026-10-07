import React from "react"
import { Button, Flex, Form, FormInstance } from "antd"
import Password from "antd/es/input/Password"
import UserService from "../../../../UserService"
import Notification from "../../../../../../core/components/notification/Notification"
import ValidationResult from "../../../../../../core/validation/ValidationResult"
import UserUpdatePasswordRequest from "../../../../model/UserUpdatePasswordRequest"
import FormLayout from "../../../../../../core/components/form/FormLayout"
import { UnexpectedResponseError } from "../../../../../../core/apiclient/ApiResponse"
import ValidationResultConverter from "../../../../../../core/validation/ValidationResultConverter"
import MessageKey from "../../../../../../core/i18n/model/MessageKey.generated"
import { I18n } from "../../../../../../core/i18n/I18n"
import "./UserPasswordTab.css"

const DEFAULT_FORM_VALUES: UserUpdatePasswordRequest = {
    currentPassword: "",
    newPassword: "",
}

export interface UserPasswordTabProps {
    onClose: () => void
}

interface UserPasswordTabState {
    loading: boolean
    validationResult: ValidationResult
    formValues: UserUpdatePasswordRequest
}

export default class UserPasswordTab extends React.Component<UserPasswordTabProps, UserPasswordTabState> {
    private readonly formRef: React.RefObject<FormInstance | null>
    private readonly service: UserService

    constructor(props: UserPasswordTabProps) {
        super(props)
        this.service = new UserService()
        this.formRef = React.createRef()
        this.state = {
            loading: false,
            validationResult: new ValidationResult(),
            formValues: DEFAULT_FORM_VALUES,
        }
    }

    private async executeChange() {
        const { formValues } = this.state
        this.setState({ validationResult: new ValidationResult(), loading: true })

        return this.service
            .changePassword(formValues)
            .then(() => Notification.success(MessageKey.CommonPasswordChanged, MessageKey.CommonSuccessMessage))
            .then(() => this.resetForm())
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

    private resetForm() {
        this.setState({ validationResult: new ValidationResult(), formValues: DEFAULT_FORM_VALUES })
        this.formRef.current?.resetFields()
        this.props.onClose()
    }

    render() {
        const { loading, validationResult, formValues } = this.state

        return (
            <Form<UserUpdatePasswordRequest>
                {...FormLayout.FormDefaults}
                labelCol={FormLayout.ExpandedLabeledItem.labelCol}
                wrapperCol={FormLayout.ExpandedLabeledItem.wrapperCol}
                ref={this.formRef}
                layout="vertical"
                onValuesChange={(_, formValues) => this.setState({ formValues })}
                initialValues={formValues}
                className="user-password-tab-form"
            >
                <Form.Item
                    name="currentPassword"
                    validateStatus={validationResult.getStatus("currentPassword")}
                    help={validationResult.getMessage("currentPassword")}
                    label={<I18n id={MessageKey.FrontendUserMenuCurrentPassword} />}
                    required
                >
                    <Password />
                </Form.Item>
                <Form.Item
                    name="newPassword"
                    validateStatus={validationResult.getStatus("newPassword")}
                    help={validationResult.getMessage("newPassword")}
                    label={<I18n id={MessageKey.FrontendUserMenuNewPassword} />}
                    required
                >
                    <Password />
                </Form.Item>

                <Flex justify="end" style={{ marginTop: 24 }}>
                    <Button type="primary" loading={loading} onClick={() => this.executeChange()}>
                        <I18n id={MessageKey.CommonSave} />
                    </Button>
                </Flex>
            </Form>
        )
    }
}
