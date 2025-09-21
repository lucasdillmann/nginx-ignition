import React from "react"
import { navigateTo, routeParams } from "../../core/components/router/AppRouter"
import StreamService from "./StreamService"
import { Flex, Form, FormInstance, Input, Segmented, Switch } from "antd"
import Preloader from "../../core/components/preloader/Preloader"
import FormLayout from "../../core/components/form/FormLayout"
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
import StreamFormDefaults from "./StreamFormDefaults"
import StreamRequest, { StreamAddress, StreamProtocol, StreamType } from "./model/StreamRequest"
import "./StreamFormPage.css"
import StreamAddressInput from "./components/StreamAddressInput"
import CompatibleStreamProtocolResolver from "./utils/CompatibleStreamProtocolResolver"
import { UserAccessLevel } from "../user/model/UserAccessLevel"
import { isAccessGranted } from "../../core/components/accesscontrol/IsAccessGranted"
import AccessDeniedPage from "../../core/components/accesscontrol/AccessDeniedPage"
import { ArrowRightOutlined, QuestionCircleFilled, SwapOutlined } from "@ant-design/icons"
import StreamTypeDescription from "./utils/StreamTypeDescription"

interface StreamFormPageState {
    formValues: StreamRequest
    validationResult: ValidationResult
    loading: boolean
    notFound: boolean
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
            formValues: StreamFormDefaults,
            validationResult: new ValidationResult(),
            loading: true,
            notFound: false,
        }
    }

    private submit() {
        const { formValues } = this.state
        this.saveModal.show("Hang on tight", "We're saving the stream")
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
        Notification.success("Stream saved", "The stream was saved successfully")
        ReloadNginxAction.execute()
    }

    private handleError(error: Error) {
        if (error instanceof UnexpectedResponseError) {
            const validationResult = ValidationResultConverter.parse(error.response)
            if (validationResult != null) this.setState({ validationResult })
        }

        this.saveModal.close()
        Notification.error("That didn't work", "Please check the form to see if everything seems correct")
    }

    private handleChange(attribute: string, value: any) {
        this.setState(
            current => ({
                ...current,
                formValues: {
                    ...current.formValues,
                    [attribute]: value,
                },
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
                <p>The type defines how the requests should be handled</p>
                <p>
                    <b>Simple:</b> Proxies requests to the backend server as-is, without any modifications or
                    evaluations.
                </p>
                <p>
                    <b>Domain-based router:</b> Uses the SNI (Server Name Indication) from the TLS protocol to detect
                    the domain name requested by the client and routes the request to the corresponding backend server.
                    Please note that this type of routing only works with TLS (like HTTPS) connections, all remaining
                    connections will be forwarded to the default backend server.
                </p>
            </>
        )
    }

    private renderForm() {
        const { formValues, validationResult } = this.state
        const bindingTcp = formValues.binding.protocol === StreamProtocol.TCP

        return (
            <Form<StreamRequest>
                {...FormLayout.FormDefaults}
                ref={this.formRef}
                initialValues={formValues}
                onValuesChange={(_, values) => this.handleUpdate(values)}
            >
                <Flex className="streams-form-inner-flex-container">
                    <Flex className="streams-form-inner-flex-container-column" style={{ maxWidth: "60%" }}>
                        <h2 className="streams-form-section-name">General</h2>
                        <p className="streams-form-section-help-text">Main definitions and properties of the stream.</p>
                        <Form.Item
                            name="enabled"
                            validateStatus={validationResult.getStatus("enabled")}
                            help={validationResult.getMessage("enabled")}
                            label="Enabled"
                            required
                        >
                            <Switch />
                        </Form.Item>
                        <Form.Item
                            name="type"
                            validateStatus={validationResult.getStatus("type")}
                            help={validationResult.getMessage("type")}
                            label="Type"
                            tooltip={{
                                title: this.buildTypeTooltipContents(),
                                icon: <QuestionCircleFilled />,
                            }}
                            required
                        >
                            <Segmented
                                options={[
                                    {
                                        label: StreamTypeDescription[StreamType.SIMPLE],
                                        value: StreamType.SIMPLE,
                                        icon: <ArrowRightOutlined />,
                                    },
                                    {
                                        label: StreamTypeDescription[StreamType.SNI_ROUTER],
                                        value: StreamType.SNI_ROUTER,
                                        icon: <SwapOutlined />,
                                    },
                                ]}
                                value={formValues.type}
                            />
                        </Form.Item>
                        <Form.Item
                            name="name"
                            validateStatus={validationResult.getStatus("name")}
                            help={validationResult.getMessage("name")}
                            label="Name"
                            required
                        >
                            <Input />
                        </Form.Item>

                        <h2 className="streams-form-section-name">Backend</h2>
                        <p className="streams-form-section-help-text">
                            Address or socket file of the backing service that will reply to the requests
                        </p>
                        <StreamAddressInput
                            basePath="backend"
                            validationResult={validationResult}
                            address={formValues.defaultBackend.target}
                            onChange={value => this.handleChange("backend", value)}
                        />

                        <h2 className="streams-form-section-name">Binding</h2>
                        <p className="streams-form-section-help-text">
                            Address or socket file where the nginx's stream will listen for requests
                        </p>
                        <StreamAddressInput
                            basePath="binding"
                            validationResult={validationResult}
                            address={formValues.binding}
                            onChange={value => this.handleChange("binding", value)}
                            parentProtocol={formValues.defaultBackend.target.protocol}
                        />
                    </Flex>
                    <Flex className="streams-form-inner-flex-container-column" style={{ maxWidth: "40%" }}>
                        <h2 className="streams-form-section-name">Features</h2>
                        <p className="streams-form-section-help-text">
                            Personalization of the behaviours of the stream.
                        </p>
                        <Form.Item
                            name={["featureSet", "useProxyProtocol"]}
                            validateStatus={validationResult.getStatus("featureSet.useProxyProtocol")}
                            help={validationResult.getMessage("featureSet.useProxyProtocol")}
                            className="streams-form-expanded-label-size"
                            label="Use the PROXY protocol"
                            required
                        >
                            <Switch />
                        </Form.Item>
                        <Form.Item
                            name={["featureSet", "socketKeepAlive"]}
                            validateStatus={validationResult.getStatus("featureSet.socketKeepAlive")}
                            help={validationResult.getMessage("featureSet.socketKeepAlive")}
                            className="streams-form-expanded-label-size"
                            label="Socket keep alive"
                            required
                        >
                            <Switch />
                        </Form.Item>
                        <Form.Item
                            name={["featureSet", "tcpKeepAlive"]}
                            validateStatus={validationResult.getStatus("featureSet.tcpKeepAlive")}
                            help={validationResult.getMessage("featureSet.tcpKeepAlive")}
                            className="streams-form-expanded-label-size"
                            label="TCP keep alive"
                            required
                        >
                            <Switch disabled={!bindingTcp} />
                        </Form.Item>
                        <Form.Item
                            name={["featureSet", "tcpNoDelay"]}
                            validateStatus={validationResult.getStatus("featureSet.tcpNoDelay")}
                            help={validationResult.getMessage("featureSet.tcpNoDelay")}
                            className="streams-form-expanded-label-size"
                            label="TCP no delay"
                            required
                        >
                            <Switch disabled={!bindingTcp} />
                        </Form.Item>
                        <Form.Item
                            name={["featureSet", "tcpDeferred"]}
                            validateStatus={validationResult.getStatus("featureSet.tcpDeferred")}
                            help={validationResult.getMessage("featureSet.tcpDeferred")}
                            className="streams-form-expanded-label-size"
                            label="TCP deferred"
                            required
                        >
                            <Switch disabled={!bindingTcp} />
                        </Form.Item>
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
                description: "Save",
                disabled: !enableActions,
                onClick: () => this.submit(),
            },
        ]

        if (this.streamId !== undefined)
            actions.unshift({
                description: "Delete",
                disabled: !enableActions,
                color: "danger",
                onClick: () => this.delete(),
            })

        AppShellContext.get().updateConfig({
            title: "Stream details",
            subtitle: "Full details and configurations of the nginx stream",
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
