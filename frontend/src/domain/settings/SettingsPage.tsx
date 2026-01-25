import React from "react"
import AppShellContext from "../../core/components/shell/AppShellContext"
import SettingsService from "./SettingsService"
import Preloader from "../../core/components/preloader/Preloader"
import Notification from "../../core/components/notification/Notification"
import { Form, FormInstance, Tabs } from "antd"
import { UnexpectedResponseError } from "../../core/apiclient/ApiResponse"
import ValidationResultConverter from "../../core/validation/ValidationResultConverter"
import ReloadNginxAction from "../nginx/actions/ReloadNginxAction"
import ValidationResult from "../../core/validation/ValidationResult"
import ModalPreloader from "../../core/components/preloader/ModalPreloader"
import FormLayout from "../../core/components/form/FormLayout"
import "./SettingsPage.css"
import SettingsFormValues from "./model/SettingsFormValues"
import SettingsConverter from "./SettingsConverter"
import CommonNotifications from "../../core/components/notification/CommonNotifications"
import EmptyStates from "../../core/components/emptystate/EmptyStates"
import AccessDeniedPage from "../../core/components/accesscontrol/AccessDeniedPage"
import { UserAccessLevel } from "../user/model/UserAccessLevel"
import { isAccessGranted } from "../../core/components/accesscontrol/IsAccessGranted"
import { I18n } from "../../core/i18n/I18n"
import { settingsDefaults } from "./SettingsDefaults"
import NginxSettingsTab from "./tabs/NginxSettingsTab"
import NginxIgnitionSettingsTab from "./tabs/NginxIgnitionSettingsTab"
import AdvancedNginxSettingsTab from "./tabs/AdvancedNginxSettingsTab"
import MessageKey from "../../core/i18n/model/MessageKey.generated"

interface SettingsPageState {
    loading: boolean
    validationResult: ValidationResult
    error?: Error
    formValues?: SettingsFormValues
}

export default class SettingsPage extends React.Component<any, SettingsPageState> {
    private readonly service: SettingsService
    private readonly saveModal: ModalPreloader
    private readonly formRef: React.RefObject<FormInstance | null>

    constructor(props: any) {
        super(props)
        this.service = new SettingsService()
        this.saveModal = new ModalPreloader()
        this.formRef = React.createRef()
        this.state = {
            loading: true,
            validationResult: new ValidationResult(),
        }
    }

    private resetToDefaultValues() {
        const { nginx, logRotation, certificateAutoRenew } = settingsDefaults()
        const newValues = { nginx, logRotation, certificateAutoRenew }

        this.formRef.current?.setFieldsValue(newValues)
        this.setState(current => ({
            formValues: {
                ...current.formValues!!,
                ...newValues,
            },
        }))

        Notification.success(MessageKey.FrontendSettingsValuesReset, MessageKey.FrontendSettingsResetDescription)
    }

    private async submit() {
        const { formValues } = this.state
        this.saveModal.show(MessageKey.CommonHangOnTight, MessageKey.FrontendSettingsSaving)
        this.setState({ validationResult: new ValidationResult() })

        const settings = SettingsConverter.formValuesToSettings(formValues!!)
        return this.service
            .save(settings)
            .then(() => this.handleSuccess())
            .catch(error => this.handleError(error))
    }

    private async handleSuccess() {
        Notification.success(MessageKey.FrontendSettingsSaved, MessageKey.FrontendSettingsSavedDescription)
        this.saveModal.close()
        ReloadNginxAction.execute()
    }

    private handleError(error: Error) {
        if (error instanceof UnexpectedResponseError) {
            const validationResult = ValidationResultConverter.parse(error.response)
            if (validationResult != null) this.setState({ validationResult })
        }

        Notification.error(MessageKey.CommonThatDidntWork, MessageKey.CommonFormCheckMessage)
        this.saveModal.close()
    }

    private handleChange(formValues: SettingsFormValues) {
        this.setState({ formValues })
    }

    private buildTabs(): any[] {
        const { validationResult, formValues } = this.state
        return [
            {
                key: "nginx-settings",
                label: <I18n id={MessageKey.CommonNginx} />,
                forceRender: true,
                children: <NginxSettingsTab formValues={formValues!!} validationResult={validationResult} />,
            },
            {
                key: "advanced-settings",
                label: <I18n id={MessageKey.FrontendSettingsPageTabAdvancedNginx} />,
                forceRender: true,
                children: <AdvancedNginxSettingsTab formValues={formValues!!} validationResult={validationResult} />,
            },
            {
                key: "nginx-ignition-settings",
                label: <I18n id={MessageKey.CommonAppName} />,
                forceRender: true,
                children: <NginxIgnitionSettingsTab validationResult={validationResult} />,
            },
        ]
    }

    private renderForm() {
        const formValues = this.state.formValues!!

        return (
            <Form<SettingsFormValues>
                ref={this.formRef}
                {...FormLayout.FormDefaults}
                onValuesChange={(_, formValues) => this.handleChange(formValues)}
                initialValues={formValues}
            >
                <Tabs items={this.buildTabs()} destroyOnHidden={false} />
            </Form>
        )
    }

    private updateShellConfig(enableActions: boolean) {
        if (!isAccessGranted(UserAccessLevel.READ_WRITE, permissions => permissions.settings)) {
            enableActions = false
        }

        AppShellContext.get().updateConfig({
            title: MessageKey.CommonSettings,
            subtitle: MessageKey.FrontendSettingsSubtitle,
            actions: [
                {
                    description: MessageKey.CommonResetToDefaults,
                    disabled: !enableActions,
                    onClick: () => this.resetToDefaultValues(),
                    color: "default",
                    type: "outlined",
                },
                {
                    description: MessageKey.CommonSave,
                    disabled: !enableActions,
                    onClick: () => this.submit(),
                },
            ],
        })
    }

    componentDidMount() {
        this.updateShellConfig(false)

        this.service
            .get()
            .then(settings => SettingsConverter.settingsToFormValues(settings))
            .then(formValues => {
                this.updateShellConfig(true)
                this.setState({
                    formValues,
                    loading: false,
                })
            })
            .catch(error => {
                CommonNotifications.failedToFetch()
                this.setState({ error, loading: false })
            })
    }

    render() {
        if (!isAccessGranted(UserAccessLevel.READ_ONLY, permissions => permissions.settings)) {
            return <AccessDeniedPage />
        }

        const { loading, error } = this.state
        if (loading) return <Preloader loading />
        if (error !== undefined) return EmptyStates.FailedToFetch

        return this.renderForm()
    }
}
