import React from "react"
import { navigateTo, routeParams } from "../../core/components/router/AppRouter"
import StreamService from "./StreamService"
import { Flex, Form, FormInstance, Input, Segmented, Switch } from "antd"
import Preloader from "../../core/components/preloader/Preloader"
import ValidationResult from "../../core/validation/ValidationResult"
import ModalPreloader from "../../core/components/preloader/ModalPreloader"
import Notification from "../../core/components/notification/Notification"
import { UnexpectedResponseError } from "../../core/apiclient/ApiResponse"
import ValidationResultConverter from "../../core/validation/ValidationResultConverter"
import AppShellContext, { ShellAction } from "../../core/components/shell/AppShellContext"
import CommonNotifications from "../../core/components/notification/CommonNotifications"
import EmptyStates from "../../core/components/emptystate/EmptyStates"
import DeleteStreamAction from "./actions/DeleteStreamAction"
import ReloadNginxAction from "../nginx/actions/ReloadNginxAction"
import StreamRequest, { StreamAddress, StreamProtocol, StreamType } from "./model/StreamRequest"
import "./StreamFormPage.css"
import StreamAddressInput from "./components/StreamAddressInput"
import CompatibleStreamProtocolResolver from "./utils/CompatibleStreamProtocolResolver"
import { UserAccessLevel } from "../user/model/UserAccessLevel"
import { isAccessGranted } from "../../core/components/accesscontrol/IsAccessGranted"
import AccessDeniedPage from "../../core/components/accesscontrol/AccessDeniedPage"
import { ArrowRightOutlined, QuestionCircleFilled, SettingOutlined, SwapOutlined } from "@ant-design/icons"
import StreamTypeDescription from "./utils/StreamTypeDescription"
import If from "../../core/components/flowcontrol/If"
import StreamBackendSettingsModal from "./components/StreamBackendSettingsModal"
import FormLayout from "../../core/components/form/FormLayout"
import StreamRoutesForm from "./components/StreamRoutesForm"
import { streamFormDefaults, streamRouteDefaults } from "./StreamFormDefaults"
import StreamSupportWarning from "./components/StreamSupportWarning"
import MessageKey from "../../core/i18n/model/MessageKey.generated"
import { I18n } from "../../core/i18n/I18n"

interface StreamFormPageState {
    formValues: StreamRequest
    validationResult: ValidationResult
    loading: boolean
    notFound: boolean
    defaultBackendSettingsOpen: boolean
    error?: Error
}

export default class StreamFormPage extends React.Component<unknown, StreamFormPageState> {
    private readonly service: StreamService
    private readonly saveModal: ModalPreloader
    private readonly formRef: React.RefObject<FormInstance | null>
    private streamId?: string

    constructor(props: any) {
        super(props)
        const streamId = routeParams().id
        this.formRef = React.createRef()
        this.streamId = streamId === "new" ? undefined : streamId
        this.service = new StreamService()
        this.saveModal = new ModalPreloader()
        this.state = {
            formValues: streamFormDefaults(),
            validationResult: new ValidationResult(),
            loading: true,
            notFound: false,
            defaultBackendSettingsOpen: false,
        }
    }

    private submit() {
        const { formValues } = this.state
        this.saveModal.show(MessageKey.CommonHangOnTight, {
            id: MessageKey.CommonSavingType,
            params: { type: MessageKey.CommonStream },
        })
        this.setState({ validationResult: new ValidationResult() })

        const action =
            this.streamId === undefined
                ? this.service.create(formValues).then(response => this.updateId(response.id))
                : this.service.updateById(this.streamId, formValues)

        action.then(() => this.handleSuccess()).catch(error => this.handleError(error))
    }

    private updateId(id: string) {
        this.streamId = id
        navigateTo(`/streams/${id}`, true)
        this.updateShellConfig(true)
    }

    private handleSuccess() {
        this.saveModal.close()
        Notification.success(
            { id: MessageKey.CommonTypeSaved, params: { type: MessageKey.CommonStream } },
            MessageKey.CommonSuccessMessage,
        )
        ReloadNginxAction.execute()
    }

    private handleError(error: Error) {
        if (error instanceof UnexpectedResponseError) {
            const validationResult = ValidationResultConverter.parse(error.response)
            if (validationResult != null) this.setState({ validationResult })
        }

        this.saveModal.close()
        Notification.error(MessageKey.CommonThatDidntWork, MessageKey.CommonFormCheckMessage)
    }

    private handleChange(attribute: string, value: any) {
        const newFormValues: any = { ...this.state.formValues }

        if (attribute.includes(".")) {
            const [path, subPath] = attribute.split(".")
            newFormValues[path][subPath] = value
        } else {
            newFormValues[attribute] = value
        }

        this.setState(
            current => ({
                ...current,
                formValues: newFormValues,
            }),
            () => this.updateFormValues(),
        )
    }

    private updateFormValues() {
        const { formValues } = this.state
        const { defaultBackend, binding, featureSet } = formValues

        this.syncAddressPort(defaultBackend.target)
        this.syncAddressPort(binding)
        this.syncAddressProtocol(defaultBackend.target.protocol, binding)

        if (binding.protocol !== StreamProtocol.TCP) {
            featureSet.tcpDeferred = false
            featureSet.tcpNoDelay = false
            featureSet.tcpKeepAlive = false
        }

        this.setState(
            current => ({ ...current, formValues }),
            () => this.formRef.current?.setFieldsValue(formValues),
        )
    }

    private syncAddressPort(address: StreamAddress) {
        if (address.protocol == StreamProtocol.SOCKET) {
            address.port = undefined
        }
    }

    private syncAddressProtocol(parentProtocol: StreamProtocol, address: StreamAddress) {
        const candidates = CompatibleStreamProtocolResolver.resolve(parentProtocol)

        if (!candidates.includes(address.protocol)) {
            address.protocol = candidates[0]
        }
    }

    private handleUpdate(newValues: StreamRequest) {
        const { formValues } = this.state

        this.setState({
            formValues: {
                ...formValues,
                ...newValues,
            },
        })
    }

    private buildTypeTooltipContents() {
        return (
            <>
                <p>
                    <I18n id={MessageKey.FrontendStreamFormTypeTooltipIntro} />
                </p>
                <p>
                    <I18n id={MessageKey.FrontendStreamFormTypeTooltipSimple} />
                </p>
                <p>
                    <I18n id={MessageKey.FrontendStreamFormTypeTooltipSni} />
                </p>
            </>
        )
    }

    private changeDefaultBackendModalState(open: boolean) {
        this.setState({ defaultBackendSettingsOpen: open })
    }

    private renderDefaultBackendForm(): React.ReactNode {
        const { formValues, validationResult, defaultBackendSettingsOpen } = this.state
        const { type } = formValues

        const { title, subtitle } =
            type == StreamType.SIMPLE
                ? {
                      title: MessageKey.FrontendStreamFormBackendTitle,
                      subtitle: MessageKey.FrontendStreamFormBackendSubtitle,
                  }
                : {
                      title: MessageKey.FrontendStreamFormDefaultBackendTitle,
                      subtitle: MessageKey.FrontendStreamFormDefaultBackendSubtitle,
                  }

        return (
            <>
                <h2 className="streams-form-section-name">
                    <I18n id={title} />
                </h2>
                <p className="streams-form-section-help-text" style={{ height: 50 }}>
                    <I18n id={subtitle} />
                </p>
                <Flex style={{ flexGrow: 1 }}>
                    <Flex style={{ flexGrow: 1, alignContent: "center", flexShrink: 1 }}>
                        <StreamAddressInput
                            basePath="defaultBackend.target"
                            validationResult={validationResult}
                            address={formValues.defaultBackend.target}
                            onChange={value => this.handleChange("defaultBackend.target", value)}
                        />
                    </Flex>
                    <Flex style={{ marginLeft: "20px", flexShrink: 1 }}>
                        <SettingOutlined onClick={() => this.changeDefaultBackendModalState(true)} size={10} />
                    </Flex>
                </Flex>
                <StreamBackendSettingsModal
                    backend={formValues.defaultBackend}
                    open={defaultBackendSettingsOpen}
                    validationBasePath="defaultBackend"
                    validationResult={validationResult}
                    onClose={() => this.changeDefaultBackendModalState(false)}
                    onChange={value => this.handleChange("defaultBackend", value)}
                    hideWeight
                />
            </>
        )
    }

    private renderSniRouterForm(): React.ReactNode {
        const { formValues, validationResult } = this.state
        if (!Array.isArray(formValues.routes)) {
            this.handleChange("routes", [streamRouteDefaults()])
            return <></>
        }

        return (
            <StreamRoutesForm
                routes={formValues.routes}
                validationResult={validationResult}
                onChange={routes => this.handleChange("routes", routes)}
            />
        )
    }

    private renderGeneralSettingsForm(): React.ReactNode {
        const { formValues, validationResult } = this.state

        return (
            <>
                <h2 className="streams-form-section-name">
                    <I18n id={MessageKey.CommonGeneral} />
                </h2>
                <p className="streams-form-section-help-text">
                    <I18n id={MessageKey.FrontendStreamFormSectionGeneralHelp} />
                </p>
                <Form.Item
                    name="enabled"
                    validateStatus={validationResult.getStatus("enabled")}
                    help={validationResult.getMessage("enabled")}
                    label={<I18n id={MessageKey.CommonEnabled} />}
                    required
                >
                    <Switch />
                </Form.Item>
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
                    name="type"
                    validateStatus={validationResult.getStatus("type")}
                    help={validationResult.getMessage("type")}
                    label={<I18n id={MessageKey.CommonType} />}
                    tooltip={{
                        title: this.buildTypeTooltipContents(),
                        icon: <QuestionCircleFilled />,
                    }}
                    required
                >
                    <Segmented
                        options={[
                            {
                                label: <I18n id={StreamTypeDescription[StreamType.SIMPLE]} />,
                                value: StreamType.SIMPLE,
                                icon: <ArrowRightOutlined />,
                            },
                            {
                                label: <I18n id={StreamTypeDescription[StreamType.SNI_ROUTER]} />,
                                value: StreamType.SNI_ROUTER,
                                icon: <SwapOutlined />,
                            },
                        ]}
                        value={formValues.type}
                    />
                </Form.Item>
            </>
        )
    }

    private renderBindingForm(): React.ReactNode {
        const { formValues, validationResult } = this.state

        return (
            <>
                <h2 className="streams-form-section-name">
                    <I18n id={MessageKey.CommonBinding} />
                </h2>
                <p className="streams-form-section-help-text" style={{ height: 50 }}>
                    <I18n id={MessageKey.FrontendStreamFormSectionBindingHelp} />
                </p>
                <StreamAddressInput
                    basePath="binding"
                    validationResult={validationResult}
                    address={formValues.binding}
                    onChange={value => this.handleChange("binding", value)}
                    parentProtocol={formValues.defaultBackend.target.protocol}
                />
            </>
        )
    }

    private renderFeatureSetForm(): React.ReactNode {
        const { formValues, validationResult } = this.state
        const bindingTcp = formValues.binding.protocol === StreamProtocol.TCP

        return (
            <>
                <h2 className="streams-form-section-name">
                    <I18n id={MessageKey.FrontendStreamFormSectionFeatures} />
                </h2>
                <p className="streams-form-section-help-text">
                    <I18n id={MessageKey.FrontendStreamFormSectionFeaturesHelp} />
                </p>
                <Form.Item
                    name={["featureSet", "useProxyProtocol"]}
                    validateStatus={validationResult.getStatus("featureSet.useProxyProtocol")}
                    help={validationResult.getMessage("featureSet.useProxyProtocol")}
                    className="streams-form-expanded-label-size"
                    label={<I18n id={MessageKey.FrontendStreamFormUseProxyProtocol} />}
                    required
                >
                    <Switch />
                </Form.Item>
                <Form.Item
                    name={["featureSet", "socketKeepAlive"]}
                    validateStatus={validationResult.getStatus("featureSet.socketKeepAlive")}
                    help={validationResult.getMessage("featureSet.socketKeepAlive")}
                    className="streams-form-expanded-label-size"
                    label={<I18n id={MessageKey.FrontendStreamFormSocketKeepAlive} />}
                    required
                >
                    <Switch />
                </Form.Item>
                <Form.Item
                    name={["featureSet", "tcpKeepAlive"]}
                    validateStatus={validationResult.getStatus("featureSet.tcpKeepAlive")}
                    help={validationResult.getMessage("featureSet.tcpKeepAlive")}
                    className="streams-form-expanded-label-size"
                    label={<I18n id={MessageKey.FrontendStreamFormTcpKeepAlive} />}
                    required
                >
                    <Switch disabled={!bindingTcp} />
                </Form.Item>
                <Form.Item
                    name={["featureSet", "tcpNoDelay"]}
                    validateStatus={validationResult.getStatus("featureSet.tcpNoDelay")}
                    help={validationResult.getMessage("featureSet.tcpNoDelay")}
                    className="streams-form-expanded-label-size"
                    label={<I18n id={MessageKey.FrontendStreamFormTcpNoDelay} />}
                    required
                >
                    <Switch disabled={!bindingTcp} />
                </Form.Item>
                <Form.Item
                    name={["featureSet", "tcpDeferred"]}
                    validateStatus={validationResult.getStatus("featureSet.tcpDeferred")}
                    help={validationResult.getMessage("featureSet.tcpDeferred")}
                    className="streams-form-expanded-label-size"
                    label={<I18n id={MessageKey.FrontendStreamFormTcpDeferred} />}
                    required
                >
                    <Switch disabled={!bindingTcp} />
                </Form.Item>
            </>
        )
    }

    private renderForm() {
        const { formValues } = this.state

        return (
            <Form<StreamRequest>
                {...FormLayout.FormDefaults}
                ref={this.formRef}
                initialValues={formValues}
                onValuesChange={(_, values) => this.handleUpdate(values)}
            >
                <StreamSupportWarning />

                <Flex className="streams-form-inner-flex-container">
                    <Flex className="streams-form-inner-flex-container-column" style={{ width: "50%" }}>
                        {this.renderGeneralSettingsForm()}
                    </Flex>
                    <Flex className="streams-form-inner-flex-container-column" style={{ width: "50%" }}>
                        {this.renderFeatureSetForm()}
                    </Flex>
                </Flex>
                <Flex className="streams-form-inner-flex-container" style={{ marginTop: 20 }}>
                    <Flex className="streams-form-inner-flex-container-column" style={{ width: "100%" }}>
                        <If condition={formValues.type === StreamType.SNI_ROUTER}>
                            {() => this.renderSniRouterForm()}
                        </If>
                    </Flex>
                </Flex>
                <Flex className="streams-form-inner-flex-container" style={{ marginTop: 20 }}>
                    <Flex className="streams-form-inner-flex-container-column" style={{ width: "50%" }}>
                        {this.renderDefaultBackendForm()}
                    </Flex>
                    <Flex className="streams-form-inner-flex-container-column" style={{ width: "50%" }}>
                        {this.renderBindingForm()}
                    </Flex>
                </Flex>
            </Form>
        )
    }

    private async delete() {
        if (this.streamId === undefined) return

        return DeleteStreamAction.execute(this.streamId).then(() => navigateTo("/streams"))
    }

    private updateShellConfig(enableActions: boolean) {
        if (!isAccessGranted(UserAccessLevel.READ_WRITE, permissions => permissions.streams)) {
            enableActions = false
        }

        const actions: ShellAction[] = [
            {
                description: MessageKey.CommonSave,
                disabled: !enableActions,
                onClick: () => this.submit(),
            },
        ]

        if (this.streamId !== undefined)
            actions.unshift({
                description: MessageKey.CommonDelete,
                disabled: !enableActions,
                color: "danger",
                onClick: () => this.delete(),
            })

        AppShellContext.get().updateConfig({
            title: MessageKey.FrontendStreamFormTitle,
            subtitle: MessageKey.FrontendStreamFormSubtitle,
            actions,
        })
    }

    componentDidMount() {
        if (this.streamId === undefined) {
            this.setState({ loading: false })
            this.updateShellConfig(true)
            return
        }

        this.service
            .getById(this.streamId!!)
            .then(streamDetails => {
                if (streamDetails === undefined) this.setState({ loading: false, notFound: true })
                else {
                    this.setState({ loading: false, formValues: streamDetails })
                    this.updateShellConfig(true)
                }
            })
            .catch(error => {
                CommonNotifications.failedToFetch()
                this.setState({ loading: false, error })
            })

        this.updateShellConfig(false)
    }

    render() {
        if (!isAccessGranted(UserAccessLevel.READ_ONLY, permissions => permissions.streams)) {
            return <AccessDeniedPage />
        }

        const { loading, notFound, error } = this.state

        if (error !== undefined) return EmptyStates.FailedToFetch
        if (notFound) return EmptyStates.NotFound
        if (loading) return <Preloader loading />

        return this.renderForm()
    }
}
