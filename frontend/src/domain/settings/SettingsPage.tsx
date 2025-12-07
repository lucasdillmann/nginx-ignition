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
import { settingsDefaults } from "./SettingsDefaults"
import NginxSettingsTab from "./tabs/NginxSettingsTab"
import NginxIgnitionSettingsTab from "./tabs/NginxIgnitionSettingsTab"
import AdvancedNginxSettingsTab from "./tabs/AdvancedNginxSettingsTab"

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

        Notification.success(
            "Values reset",
            "Settings were changed back to the default values (except the global bindings), but not yet saved",
        )
    }

    private async submit() {
        const { formValues } = this.state
        this.saveModal.show("Hang on tight", "We're saving the settings")
        this.setState({ validationResult: new ValidationResult() })

        const settings = SettingsConverter.formValuesToSettings(formValues!!)
        return this.service
            .save(settings)
            .then(() => this.handleSuccess())
            .catch(error => this.handleError(error))
    }

    private async handleSuccess() {
        Notification.success("Settings saved", "Global settings were updated successfully")
        this.saveModal.close()
        ReloadNginxAction.execute()
    }

    private handleError(error: Error) {
        if (error instanceof UnexpectedResponseError) {
            const validationResult = ValidationResultConverter.parse(error.response)
            if (validationResult != null) this.setState({ validationResult })
        }

        Notification.error("That didn't work", "Please check the form to see if everything seems correct")
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
                label: "nginx",
                forceRender: true,
                children: <NginxSettingsTab formValues={formValues!!} validationResult={validationResult} />,
            },
            {
                key: "advanced-settings",
                label: "nginx (advanced)",
                forceRender: true,
                children: <AdvancedNginxSettingsTab formValues={formValues!!} validationResult={validationResult} />,
            },
            {
                key: "nginx-ignition-settings",
                label: "nginx ignition",
                forceRender: true,
                children: <NginxIgnitionSettingsTab formValues={formValues!!} validationResult={validationResult} />,
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
                <Tabs items={this.buildTabs()} destroyOnHidden={false} animated />
            </Form>
        )
    }

    private updateShellConfig(enableActions: boolean) {
        if (!isAccessGranted(UserAccessLevel.READ_WRITE, permissions => permissions.settings)) {
            enableActions = false
        }

        AppShellContext.get().updateConfig({
            title: "Settings",
            subtitle: "Globals settings for the nginx server and nginx ignition",
            actions: [
                {
                    description: "Reset to defaults",
                    disabled: !enableActions,
                    onClick: () => this.resetToDefaultValues(),
                    color: "default",
                    type: "outlined",
                },
                {
                    description: "Save",
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
