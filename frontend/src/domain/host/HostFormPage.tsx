import React from "react"
import AppShellContext, { ShellAction } from "../../core/components/shell/AppShellContext"
import { Empty, Flex, Form, FormInstance, Switch } from "antd"
import FormLayout from "../../core/components/form/FormLayout"
import DomainNamesList from "../certificate/components/DomainNamesList"
import { navigateTo, queryParams, routeParams } from "../../core/components/router/AppRouter"
import HostService from "./HostService"
import ValidationResult from "../../core/validation/ValidationResult"
import { HostBindingType, HostRouteType } from "./model/HostRequest"
import Preloader from "../../core/components/preloader/Preloader"
import DeleteHostAction from "./actions/DeleteHostAction"
import Notification from "../../core/components/notification/Notification"
import { UnexpectedResponseError } from "../../core/apiclient/ApiResponse"
import ValidationResultConverter from "../../core/validation/ValidationResultConverter"
import ModalPreloader from "../../core/components/preloader/ModalPreloader"
import HostRoutes from "./components/HostRoutes"
import HostBindings from "./components/HostBindings"
import "./HostFormPage.css"
import ReloadNginxAction from "../nginx/actions/ReloadNginxAction"
import HostFormValues from "./model/HostFormValues"
import HostConverter from "./HostConverter"
import { IntegrationResponse } from "../integration/model/IntegrationResponse"
import IntegrationService from "../integration/IntegrationService"
import If from "../../core/components/flowcontrol/If"

const DEFAULT_HOST: HostFormValues = {
    enabled: true,
    defaultServer: false,
    useGlobalBindings: true,
    domainNames: [""],
    bindings: [
        {
            ip: "0.0.0.0",
            port: 8080,
            type: HostBindingType.HTTP,
        },
    ],
    routes: [
        {
            priority: 0,
            type: HostRouteType.PROXY,
            sourcePath: "/",
            targetUri: "",
        },
    ],
    featureSet: {
        websocketsSupport: true,
        http2Support: true,
        redirectHttpToHttps: false,
    },
}

interface HostFormPageState {
    formValues: HostFormValues
    validationResult: ValidationResult
    integrations: IntegrationResponse[]
    loading: boolean
    notFound: boolean
}

export default class HostFormPage extends React.Component<any, HostFormPageState> {
    private readonly hostService: HostService
    private readonly integrationService: IntegrationService
    private readonly hostId?: string
    private readonly saveModal: ModalPreloader
    private readonly formRef: React.RefObject<FormInstance>

    constructor(props: any) {
        super(props)

        const hostId = routeParams().id
        this.hostId = hostId === "new" ? undefined : hostId
        this.hostService = new HostService()
        this.integrationService = new IntegrationService()
        this.saveModal = new ModalPreloader()
        this.formRef = React.createRef()
        this.state = {
            validationResult: new ValidationResult(),
            loading: true,
            notFound: false,
            formValues: DEFAULT_HOST,
            integrations: [],
        }
    }

    private async delete() {
        if (this.hostId === undefined) return

        return DeleteHostAction.execute(this.hostId).then(() => navigateTo("/hosts"))
    }

    private submit() {
        const payload = HostConverter.formValuesToRequest(this.state.formValues)
        this.saveModal.show("Hang on tight", "We're saving the host")

        const action =
            this.hostId === undefined
                ? this.hostService.create(payload)
                : this.hostService.updateById(this.hostId, payload)

        action
            .then(() => this.handleSuccess())
            .catch(error => this.handleError(error))
            .then(() => this.saveModal.close())
    }

    private handleSuccess() {
        Notification.success("Host saved", "The host was saved successfully")
        ReloadNginxAction.execute()
        navigateTo("/hosts")
    }

    private handleError(error: Error) {
        if (error instanceof UnexpectedResponseError) {
            const validationResult = ValidationResultConverter.parse(error.response)
            if (validationResult != null) this.setState({ validationResult })
        }

        Notification.error("That didn't work", "Please check the form to see if everything seems correct")
    }

    private updateShellConfig(enableActions: boolean) {
        const actions: ShellAction[] = [
            {
                description: "Save",
                disabled: !enableActions,
                onClick: () => this.submit(),
            },
        ]

        if (this.hostId !== undefined)
            actions.unshift({
                description: "Delete",
                disabled: !enableActions,
                color: "danger",
                onClick: () => this.delete(),
            })

        AppShellContext.get().updateConfig({
            title: "Host details",
            subtitle: "Full details and configurations of the nginx's virtual host",
            actions,
        })
    }

    private handleChange(host: HostFormValues) {
        const { bindings, routes, useGlobalBindings } = host

        const injectBindingNeeded = !useGlobalBindings && bindings.length === 0
        const sortedRoutes = routes.sort((left, right) => (left.priority > right.priority ? 1 : -1))
        const orderedData: HostFormValues = {
            ...host,
            routes: sortedRoutes,
            bindings: injectBindingNeeded ? DEFAULT_HOST.bindings : bindings,
        }

        this.setState({ formValues: orderedData }, () => {
            if (injectBindingNeeded) this.formRef.current?.resetFields()
        })
    }

    private removeRoute(index: number) {
        const { formValues } = this.state
        const { routes } = formValues

        let priority = 0
        const updatedValues = routes
            .filter((_, itemIndex) => itemIndex !== index)
            .map(route => ({
                ...route,
                priority: priority++,
            }))

        this.formRef.current?.setFieldValue("routes", updatedValues)
        this.setState({
            formValues: {
                ...formValues,
                routes: updatedValues,
            },
        })
    }

    private renderForm() {
        const { validationResult, formValues, integrations } = this.state

        return (
            <Form<HostFormValues>
                {...FormLayout.FormDefaults}
                ref={this.formRef}
                onValuesChange={(_, formValues) => this.handleChange(formValues)}
                initialValues={formValues}
            >
                <h2 className="hosts-form-section-name">General</h2>
                <p className="hosts-form-section-help-text">
                    General configurations properties of the nginx's virtual host
                </p>
                <Flex className="hosts-form-inner-flex-container">
                    <Flex
                        className="hosts-form-inner-flex-container-column hosts-form-expanded-label-size"
                        style={{ maxWidth: "30%" }}
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
                        <Form.Item
                            name="defaultServer"
                            validateStatus={validationResult.getStatus("defaultServer")}
                            help={validationResult.getMessage("defaultServer")}
                            label="Default server"
                            required
                        >
                            <Switch />
                        </Form.Item>
                        <Form.Item
                            name={["featureSet", "websocketsSupport"]}
                            validateStatus={validationResult.getStatus("featureSet.websocketsSupport")}
                            help={validationResult.getMessage("featureSet.websocketsSupport")}
                            label="Websockets support"
                            required
                        >
                            <Switch />
                        </Form.Item>
                        <Form.Item
                            name={["featureSet", "http2Support"]}
                            validateStatus={validationResult.getStatus("featureSet.http2Support")}
                            help={validationResult.getMessage("featureSet.http2Support")}
                            label="HTTP2 support"
                            required
                        >
                            <Switch />
                        </Form.Item>
                        <Form.Item
                            name={["featureSet", "redirectHttpToHttps"]}
                            validateStatus={validationResult.getStatus("featureSet.redirectHttpToHttps")}
                            help={validationResult.getMessage("featureSet.redirectHttpToHttps")}
                            label="Redirect HTTP to HTTPS"
                            required
                        >
                            <Switch />
                        </Form.Item>
                    </Flex>
                    <Flex className="hosts-form-inner-flex-container-column">
                        <If condition={formValues.defaultServer}>
                            <Form.Item label="Domain names" required>
                                <Flex>Not available for the default server</Flex>
                            </Form.Item>
                        </If>
                        <DomainNamesList
                            validationResult={validationResult}
                            className={formValues.defaultServer ? "hosts-form-invisible-input" : undefined}
                        />
                    </Flex>
                </Flex>

                <h2 className="hosts-form-section-name">Routing</h2>
                <p className="hosts-form-section-help-text">
                    Routes to be configured in the host. The nginx will evaluate them from top to bottom, executing the
                    first one that matches the source path.
                </p>
                <HostRoutes
                    routes={formValues.routes}
                    validationResult={validationResult}
                    integrations={integrations}
                    onRouteRemove={index => this.removeRoute(index)}
                />

                <h2 className="hosts-form-section-name">Bindings</h2>
                <p className="hosts-form-section-help-text">
                    Relation of IPs and ports where the host will listen for requests
                </p>
                <Form.Item
                    name="useGlobalBindings"
                    validateStatus={validationResult.getStatus("useGlobalBindings")}
                    help={validationResult.getMessage("useGlobalBindings")}
                    label="Use global bindings"
                    wrapperCol={{ style: { flexGrow: 0, minWidth: 65 } }}
                    labelCol={{ style: { order: 2, flexGrow: 1 } }}
                    required
                >
                    <Switch />
                </Form.Item>
                <HostBindings
                    className={formValues.useGlobalBindings ? "hosts-form-invisible-input" : undefined}
                    pathPrefix="bindings"
                    bindings={formValues.bindings}
                    validationResult={validationResult}
                />
            </Form>
        )
    }

    componentDidMount() {
        const copyFrom = queryParams().copyFrom as string | undefined
        const integrations = this.integrationService.getAll(true)

        if (this.hostId === undefined && copyFrom === undefined) {
            integrations.then(response => {
                this.setState({ loading: false, integrations: response })
                this.updateShellConfig(true)
            })
            return
        }

        const formValues = this.hostService
            .getById((this.hostId ?? copyFrom)!!)
            .then(response => (response === undefined ? undefined : HostConverter.responseToFormValues(response)))

        Promise.all([formValues, integrations]).then(([formValues, integrations]) => {
            if (formValues === undefined) this.setState({ loading: false, notFound: true })
            else {
                if (copyFrom !== undefined)
                    Notification.success(
                        "Host values copied",
                        "The values from the selected host where successfully copied as a new host",
                    )

                this.setState({ loading: false, formValues, integrations })
                this.updateShellConfig(true)
            }
        })

        this.updateShellConfig(false)
    }

    render() {
        const { loading, notFound } = this.state

        if (notFound) return <Empty description="Not found" />

        if (loading) return <Preloader loading />

        return this.renderForm()
    }
}
