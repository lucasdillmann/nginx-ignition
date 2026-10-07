import React from "react"
import { Button, Flex, Form, FormInstance, Input } from "antd"
import UserService from "../../../../UserService"
import Notification from "../../../../../../core/components/notification/Notification"
import ValidationResult from "../../../../../../core/validation/ValidationResult"
import UserUpdateProfileRequest from "../../../../model/UserUpdateProfileRequest"
import FormLayout from "../../../../../../core/components/form/FormLayout"
import { UnexpectedResponseError } from "../../../../../../core/apiclient/ApiResponse"
import ValidationResultConverter from "../../../../../../core/validation/ValidationResultConverter"
import MessageKey from "../../../../../../core/i18n/model/MessageKey.generated"
import { I18n } from "../../../../../../core/i18n/I18n"
import AppContext from "../../../../../../core/components/context/AppContext"
import "./UserProfileTab.css"

const DEFAULT_FORM_VALUES: UserUpdateProfileRequest = {
    name: "",
    username: "",
}

export interface UserProfileTabProps {
    onClose: () => void
}

interface UserProfileTabState {
    loading: boolean
    validationResult: ValidationResult
    formValues: UserUpdateProfileRequest
}

export default class UserProfileTab extends React.Component<UserProfileTabProps, UserProfileTabState> {
    private readonly formRef: React.RefObject<FormInstance | null>
    private readonly service: UserService

    constructor(props: UserProfileTabProps) {
        super(props)
        this.service = new UserService()
        this.formRef = React.createRef()
        this.state = {
            loading: false,
            validationResult: new ValidationResult(),
            formValues: DEFAULT_FORM_VALUES,
        }
    }

    componentDidMount() {
        this.loadFormValues()
    }

    private loadFormValues() {
        const user = AppContext.get().user
        const formValues = {
            name: user?.name ?? "",
            username: user?.username ?? "",
        }

        this.setState({ formValues, validationResult: new ValidationResult() })
        this.formRef.current?.setFieldsValue(formValues)
    }

    private async executeUpdate() {
        const { formValues } = this.state
        this.setState({ validationResult: new ValidationResult(), loading: true })

        return this.service
            .updateProfile(formValues)
            .then(() =>
                Notification.success(
                    { id: MessageKey.CommonTypeSaved, params: { type: MessageKey.CommonUser } },
                    MessageKey.CommonSuccessMessage,
                ),
            )
            .then(() => AppContext.get().container!!.reload())
            .then(() => this.props.onClose())
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

    render() {
        const { loading, validationResult, formValues } = this.state

        return (
            <Form<UserUpdateProfileRequest>
                {...FormLayout.FormDefaults}
                labelCol={FormLayout.ExpandedLabeledItem.labelCol}
                wrapperCol={FormLayout.ExpandedLabeledItem.wrapperCol}
                ref={this.formRef}
                layout="vertical"
                onValuesChange={(_, formValues) => this.setState({ formValues })}
                initialValues={formValues}
                className="user-profile-tab-form"
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
                    name="username"
                    validateStatus={validationResult.getStatus("username")}
                    help={validationResult.getMessage("username")}
                    label={<I18n id={MessageKey.CommonUsername} />}
                    required
                >
                    <Input />
                </Form.Item>

                <Flex justify="end" style={{ marginTop: 24 }}>
                    <Button type="primary" loading={loading} onClick={() => this.executeUpdate()}>
                        <I18n id={MessageKey.CommonSave} />
                    </Button>
                </Flex>
            </Form>
        )
    }
}
