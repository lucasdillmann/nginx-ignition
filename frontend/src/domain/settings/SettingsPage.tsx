import React from "react"
import AppShellContext from "../../core/components/shell/AppShellContext"
import { LogLevel, TimeUnit } from "./model/SettingsDto"
import SettingsService from "./SettingsService"
import Preloader from "../../core/components/preloader/Preloader"
import Notification from "../../core/components/notification/Notification"
import { Flex, Form, FormInstance, Input, InputNumber, Select, Space, Switch } from "antd"
import { UnexpectedResponseError } from "../../core/apiclient/ApiResponse"
import ValidationResultConverter from "../../core/validation/ValidationResultConverter"
import ReloadNginxAction from "../nginx/actions/ReloadNginxAction"
import ValidationResult from "../../core/validation/ValidationResult"
import ModalPreloader from "../../core/components/preloader/ModalPreloader"
import FormLayout from "../../core/components/form/FormLayout"
import HostBindings from "../host/components/HostBindings"
import "./SettingsPage.css"
import SettingsFormValues from "./model/SettingsFormValues"
import SettingsConverter from "./SettingsConverter"
import CommonNotifications from "../../core/components/notification/CommonNotifications"
import EmptyStates from "../../core/components/emptystate/EmptyStates"
import SettingsDefaults from "./SettingsDefaults"

const INTEGER_MAX = 2147483647

interface SettingsPageState {
    loading: boolean
    validationResult: ValidationResult
    error?: Error
    formValues?: SettingsFormValues
}

export default class SettingsPage extends React.Component<any, SettingsPageState> {
    private readonly service: SettingsService
    private readonly saveModal: ModalPreloader
    private readonly formRef: React.RefObject<FormInstance>

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
        const { nginx, logRotation, certificateAutoRenew } = SettingsDefaults
        const newValues = { nginx, logRotation, certificateAutoRenew }

        this.formRef.current?.setFieldsValue(newValues)
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
        return ReloadNginxAction.execute()
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

    private renderExecutionIntervalFieldset(pathPrefix: string) {
        const { validationResult } = this.state
        return (
            <Form.Item label="Execution interval" required>
                <Space.Compact>
                    <Form.Item
                        name={[pathPrefix, "intervalUnitCount"]}
                        validateStatus={validationResult.getStatus(`${pathPrefix}.intervalUnitCount`)}
                        help={validationResult.getMessage(`${pathPrefix}.intervalUnitCount`)}
                    >
                        <InputNumber min={1} max={INTEGER_MAX} />
                    </Form.Item>
                    <Form.Item
                        name={[pathPrefix, "intervalUnit"]}
                        validateStatus={validationResult.getStatus(`${pathPrefix}.intervalUnit`)}
                        help={validationResult.getMessage(`${pathPrefix}.intervalUnit`)}
                        style={{ minWidth: 100 }}
                    >
                        <Select>
                            <Select.Option value={TimeUnit.DAYS}>days</Select.Option>
                            <Select.Option value={TimeUnit.HOURS}>hours</Select.Option>
                            <Select.Option value={TimeUnit.MINUTES}>minutes</Select.Option>
                        </Select>
                    </Form.Item>
                </Space.Compact>
            </Form.Item>
        )
    }

    private renderScheduledTasksFormPortion() {
        const { validationResult } = this.state
        return (
            <Flex className="settings-form-inner-flex-container-column  settings-form-expanded-label-size">
                <h2 className="settings-form-section-name">Scheduled tasks</h2>
                <p className="settings-form-section-help-text">Definition of the nginx ignition's housekeeping tasks</p>
                <h3 className="settings-form-subsection-name">Log rotation</h3>
                <Form.Item
                    name={["logRotation", "enabled"]}
                    validateStatus={validationResult.getStatus("logRotation.enabled")}
                    help={validationResult.getMessage("logRotation.enabled")}
                    label="Auto rotation enabled"
                    required
                >
                    <Switch />
                </Form.Item>
                <Form.Item
                    name={["logRotation", "maximumLines"]}
                    validateStatus={validationResult.getStatus("logRotation.maximumLines")}
                    help={validationResult.getMessage("logRotation.maximumLines")}
                    label="Lines to keep"
                    required
                >
                    <InputNumber min={0} max={10000} />
                </Form.Item>
                {this.renderExecutionIntervalFieldset("logRotation")}
                <h3 className="settings-form-subsection-name">SSL certificates auto renew</h3>
                <Form.Item
                    name={["certificateAutoRenew", "enabled"]}
                    validateStatus={validationResult.getStatus("certificateAutoRenew.enabled")}
                    help={validationResult.getMessage("certificateAutoRenew.enabled")}
                    label="Auto renew enabled"
                    required
                >
                    <Switch />
                </Form.Item>
                {this.renderExecutionIntervalFieldset("certificateAutoRenew")}
            </Flex>
        )
    }

    private renderErrorLogFieldset(enabledField: string, levelField: string) {
        const { validationResult } = this.state
        return (
            <>
                <Form.Item
                    name={["nginx", "logs", enabledField]}
                    validateStatus={validationResult.getStatus(`nginx.logs.${enabledField}`)}
                    help={validationResult.getMessage(`nginx.logs.${enabledField}`)}
                    label="Error logs enabled"
                    required
                >
                    <Switch />
                </Form.Item>
                <Form.Item
                    name={["nginx", "logs", levelField]}
                    validateStatus={validationResult.getStatus(`nginx.logs.${levelField}`)}
                    help={validationResult.getMessage(`nginx.logs.${levelField}`)}
                    label="Level"
                    required
                >
                    <Select>
                        <Select.Option value={LogLevel.WARN}>warn</Select.Option>
                        <Select.Option value={LogLevel.ALERT}>alert</Select.Option>
                        <Select.Option value={LogLevel.ERROR}>error</Select.Option>
                        <Select.Option value={LogLevel.CRIT}>crit</Select.Option>
                        <Select.Option value={LogLevel.EMERG}>emerg</Select.Option>
                    </Select>
                </Form.Item>
            </>
        )
    }

    private renderNginxLogsFormPortion() {
        const { validationResult } = this.state
        return (
            <Flex className="settings-form-inner-flex-container-column  settings-form-expanded-label-size">
                <h2 className="settings-form-section-name">Logs</h2>
                <p className="settings-form-section-help-text">
                    Logging settings for the nginx server and virtual hosts
                </p>
                <h3 className="settings-form-subsection-name">Server</h3>
                {this.renderErrorLogFieldset("serverLogsEnabled", "serverLogsLevel")}

                <h3 className="settings-form-subsection-name">Virtual hosts</h3>
                <Form.Item
                    name={["nginx", "logs", "accessLogsEnabled"]}
                    validateStatus={validationResult.getStatus("nginx.logs.accessLogsEnabled")}
                    help={validationResult.getMessage("nginx.logs.accessLogsEnabled")}
                    label="Access logs enabled"
                    required
                >
                    <Switch />
                </Form.Item>
                {this.renderErrorLogFieldset("errorLogsEnabled", "errorLogsLevel")}
            </Flex>
        )
    }

    private renderGeneralSecondColumn() {
        const { validationResult } = this.state
        return (
            <Flex className="settings-form-inner-flex-container-column settings-form-expanded-label-size">
                <Form.Item
                    name={["nginx", "maximumBodySizeMb"]}
                    validateStatus={validationResult.getStatus("nginx.maximumBodySizeMb")}
                    help={validationResult.getMessage("nginx.maximumBodySizeMb")}
                    label="Maximum body size (MB)"
                    required
                >
                    <InputNumber min={1} max={INTEGER_MAX} />
                </Form.Item>
                <Form.Item
                    name={["nginx", "timeouts", "connect"]}
                    validateStatus={validationResult.getStatus("nginx.timeouts.connect")}
                    help={validationResult.getMessage("nginx.timeouts.connect")}
                    label="Connect timeout (seconds)"
                    required
                >
                    <InputNumber min={1} max={INTEGER_MAX} />
                </Form.Item>
                <Form.Item
                    name={["nginx", "timeouts", "read"]}
                    validateStatus={validationResult.getStatus("nginx.timeouts.read")}
                    help={validationResult.getMessage("nginx.timeouts.read")}
                    label="Read timeout (seconds)"
                    required
                >
                    <InputNumber min={1} max={INTEGER_MAX} />
                </Form.Item>
                <Form.Item
                    name={["nginx", "timeouts", "send"]}
                    validateStatus={validationResult.getStatus("nginx.timeouts.send")}
                    help={validationResult.getMessage("nginx.timeouts.send")}
                    label="Send timeout (seconds)"
                    required
                >
                    <InputNumber min={1} max={INTEGER_MAX} />
                </Form.Item>
                <Form.Item
                    name={["nginx", "timeouts", "keepalive"]}
                    validateStatus={validationResult.getStatus("nginx.timeouts.keepalive")}
                    help={validationResult.getMessage("nginx.timeouts.keepalive")}
                    label="Keepalive timeout (seconds)"
                    required
                >
                    <InputNumber min={1} max={INTEGER_MAX} />
                </Form.Item>
            </Flex>
        )
    }

    private renderGeneralFirstColumn() {
        const { validationResult } = this.state
        return (
            <Flex className="settings-form-inner-flex-container-column settings-form-expanded-label-size">
                <Form.Item
                    name={["nginx", "serverTokensEnabled"]}
                    validateStatus={validationResult.getStatus("nginx.serverTokensEnabled")}
                    help={validationResult.getMessage("nginx.serverTokensEnabled")}
                    label="Server tokens"
                    required
                >
                    <Switch />
                </Form.Item>
                <Form.Item
                    name={["nginx", "gzipEnabled"]}
                    validateStatus={validationResult.getStatus("nginx.gzipEnabled")}
                    help={validationResult.getMessage("nginx.gzipEnabled")}
                    label="GZIP enabled"
                    required
                >
                    <Switch />
                </Form.Item>
                <Form.Item
                    name={["nginx", "sendfileEnabled"]}
                    validateStatus={validationResult.getStatus("nginx.sendfileEnabled")}
                    help={validationResult.getMessage("nginx.sendfileEnabled")}
                    label="Sendfile enabled"
                    required
                >
                    <Switch />
                </Form.Item>
                <Form.Item
                    name={["nginx", "defaultContentType"]}
                    validateStatus={validationResult.getStatus("nginx.defaultContentType")}
                    help={validationResult.getMessage("nginx.defaultContentType")}
                    label="Default content type"
                    required
                >
                    <Input maxLength={128} minLength={1} />
                </Form.Item>
                <Form.Item
                    name={["nginx", "workerProcesses"]}
                    validateStatus={validationResult.getStatus("nginx.workerProcesses")}
                    help={validationResult.getMessage("nginx.workerProcesses")}
                    label="Worker processes"
                    required
                >
                    <InputNumber min={1} max={100} />
                </Form.Item>
                <Form.Item
                    name={["nginx", "workerConnections"]}
                    validateStatus={validationResult.getStatus("nginx.workerConnections")}
                    help={validationResult.getMessage("nginx.workerConnections")}
                    label="Connections per worker"
                    required
                >
                    <InputNumber min={32} max={4096} />
                </Form.Item>
            </Flex>
        )
    }

    private renderGeneralFormPortion() {
        return (
            <>
                <h2 className="settings-form-section-name" style={{ marginTop: 0 }}>
                    General
                </h2>
                <p className="settings-form-section-help-text">General configurations properties of the nginx server</p>
                <Flex className="settings-form-inner-flex-container">
                    {this.renderGeneralFirstColumn()}
                    {this.renderGeneralSecondColumn()}
                </Flex>
            </>
        )
    }

    private renderBindingsFormPortion() {
        const { validationResult } = this.state
        const formValues = this.state.formValues!!

        return (
            <>
                <h2 className="settings-form-section-name">Global bindings</h2>
                <p className="settings-form-section-help-text">
                    Relation of IPs and ports where the virtual hosts will listen for requests by default (can be
                    overwritten on every host if needed)
                </p>
                <HostBindings
                    pathPrefix="globalBindings"
                    bindings={formValues.globalBindings}
                    validationResult={validationResult}
                />
            </>
        )
    }

    private renderLogsAndTaskFormPortion() {
        return (
            <Flex className="settings-form-inner-flex-container">
                {this.renderNginxLogsFormPortion()}
                {this.renderScheduledTasksFormPortion()}
            </Flex>
        )
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
                {this.renderGeneralFormPortion()}
                {this.renderLogsAndTaskFormPortion()}
                {this.renderBindingsFormPortion()}
            </Form>
        )
    }

    private updateShellConfig(enableActions: boolean) {
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
        const { loading, error } = this.state
        if (loading) return <Preloader loading />
        if (error !== undefined) return EmptyStates.FailedToFetch

        return this.renderForm()
    }
}
