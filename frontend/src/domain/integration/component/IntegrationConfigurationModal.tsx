import React from "react"
import { IntegrationConfigurationRequest } from "../model/IntegrationConfigurationRequest"
import { IntegrationConfigurationResponse } from "../model/IntegrationConfigurationResponse"
import { Form, Modal, Switch } from "antd"
import Preloader from "../../../core/components/preloader/Preloader"
import IntegrationService from "../IntegrationService"
import ValidationResult from "../../../core/validation/ValidationResult"
import Notification from "../../../core/components/notification/Notification"
import { UnexpectedResponseError } from "../../../core/apiclient/ApiResponse"
import ValidationResultConverter from "../../../core/validation/ValidationResultConverter"
import FormLayout from "../../../core/components/form/FormLayout"
import DynamicInput from "../../../core/components/dynamicfield/DynamicInput"
import CommonNotifications from "../../../core/components/notification/CommonNotifications"
import { isAccessGranted } from "../../../core/components/accesscontrol/IsAccessGranted"
import { UserAccessLevel } from "../../user/model/UserAccessLevel"
import AccessDeniedModal from "../../../core/components/accesscontrol/AccessDeniedModal"

export interface IntegrationConfigurationModalProps {
    integrationId: string
    onClose: (updated: boolean) => void
}

interface IntegrationConfigurationModalState {
    integration?: IntegrationConfigurationResponse
    formValues: IntegrationConfigurationRequest
    validationResult: ValidationResult
}

export default class IntegrationConfigurationModal extends React.Component<
    IntegrationConfigurationModalProps,
    IntegrationConfigurationModalState
> {
    private readonly service: IntegrationService

    constructor(props: IntegrationConfigurationModalProps) {
        super(props)
        this.service = new IntegrationService()
        this.state = {
            validationResult: new ValidationResult(),
            formValues: {
                enabled: false,
                parameters: {},
            },
        }
    }

    private saveConfiguration() {
        if (!isAccessGranted(UserAccessLevel.READ_WRITE, permissions => permissions.integrations)) {
            return AccessDeniedModal.show()
        }

        const { integrationId } = this.props
        const { formValues } = this.state

        this.setState({ validationResult: new ValidationResult() })
        this.service
            .setConfiguration(integrationId, formValues)
            .then(() => this.handleSuccess())
            .catch(error => this.handleError(error))
    }

    private handleSuccess() {
        const { onClose } = this.props

        Notification.success("Configuration saved", "The integration settings were updated successfully")
        onClose(true)
    }

    private handleError(error: Error) {
        if (error instanceof UnexpectedResponseError) {
            const validationResult = ValidationResultConverter.parse(error.response)
            if (validationResult != null) this.setState({ validationResult })
        }

        Notification.error("That didn't work", "Please check the form to see if everything seems correct")
    }

    private renderDynamicFields() {
        const { integration, formValues, validationResult } = this.state
        return integration?.configurationFields
            .sort((left, right) => (left.priority > right.priority ? 1 : -1))
            .map(field => (
                <DynamicInput
                    key={field.id}
                    formValues={formValues}
                    validationResult={validationResult}
                    field={field}
                />
            ))
    }

    private renderForm() {
        const { formValues, integration, validationResult } = this.state
        if (integration === undefined) return <Preloader loading />

        return (
            <Form<IntegrationConfigurationRequest>
                {...FormLayout.FormDefaults}
                onValuesChange={(_, formValues) => this.setState({ formValues })}
                initialValues={formValues}
            >
                <Form.Item
                    name="enabled"
                    validateStatus={validationResult.getStatus("enabled")}
                    help={validationResult.getMessage("enabled")}
                    label="Enabled"
                    required
                >
                    <Switch />
                </Form.Item>
                {this.renderDynamicFields()}
            </Form>
        )
    }

    private buildParametersInitialValue(integration: IntegrationConfigurationResponse) {
        const { parameters, configurationFields } = integration
        const output = { ...parameters }

        configurationFields.forEach(field => {
            const currentValue = output[field.id]
            if (currentValue === undefined) output[field.id] = field.defaultValue
        })

        return output
    }

    componentDidMount() {
        const { integrationId } = this.props
        this.service
            .getConfiguration(integrationId)
            .then(integration => {
                const formValues = {
                    enabled: integration.enabled,
                    parameters: this.buildParametersInitialValue(integration),
                }

                this.setState({
                    integration,
                    formValues,
                })
            })
            .catch(() => CommonNotifications.failedToFetch())
    }

    render() {
        const { onClose } = this.props
        const { integration } = this.state
        return (
            <Modal
                afterClose={() => onClose(false)}
                onCancel={() => onClose(false)}
                onOk={() => this.saveConfiguration()}
                okButtonProps={{
                    disabled: !isAccessGranted(UserAccessLevel.READ_WRITE, permissions => permissions.integrations),
                }}
                title={integration?.name}
                width={800}
                open
            >
                {this.renderForm()}
            </Modal>
        )
    }
}
