import React from "react"
import AppShellContext, { ShellAction } from "../../core/components/shell/AppShellContext"
import { Flex, Form, FormInstance, Switch } from "antd"
import FormLayout from "../../core/components/form/FormLayout"
import DomainNamesList from "../certificate/components/DomainNamesList"
import { navigateTo, queryParams, routeParams } from "../../core/components/router/AppRouter"
import HostService from "./HostService"
import ValidationResult from "../../core/validation/ValidationResult"
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
import CommonNotifications from "../../core/components/notification/CommonNotifications"
import EmptyStates from "../../core/components/emptystate/EmptyStates"
import HostFormValuesDefaults from "./model/HostFormValuesDefaults"
import PaginatedSelect from "../../core/components/select/PaginatedSelect"
import AccessListResponse from "../accesslist/model/AccessListResponse"
import PageResponse from "../../core/pagination/PageResponse"
import AccessListService from "../accesslist/AccessListService"

interface HostFormPageState {
    formValues: HostFormValues
    validationResult: ValidationResult
    integrations: IntegrationResponse[]
    loading: boolean
    notFound: boolean
    error?: Error
}

export default class HostFormPage extends React.Component<any, HostFormPageState> {
    private readonly hostService: HostService
    private readonly integrationService: IntegrationService
    private readonly accessListService: AccessListService
    private readonly saveModal: ModalPreloader
    private readonly formRef: React.RefObject<FormInstance | null>
    private hostId?: string

    constructor(props: any) {
        super(props)

        const hostId = routeParams().id
        this.hostId = hostId === "new" ? undefined : hostId
        this.hostService = new HostService()
        this.integrationService = new IntegrationService()
        this.accessListService = new AccessListService()
        this.saveModal = new ModalPreloader()
        this.formRef = React.createRef()
        this.state = {
            validationResult: new ValidationResult(),
            loading: true,
            notFound: false,
            formValues: HostFormValuesDefaults,
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
        this.setState({ validationResult: new ValidationResult() })

        const action =
            this.hostId === undefined
                ? this.hostService.create(payload).then(response => this.updateId(response.id))
                : this.hostService.updateById(this.hostId, payload)

        action
            .then(() => this.handleSuccess())
            .catch(error => this.handleError(error))
            .then(() => this.saveModal.close())
    }

    private updateId(id: string) {
        this.hostId = id
        navigateTo(`/hosts/${id}`, true)
        this.updateShellConfig(true)
    }

    private handleSuccess() {
        Notification.success("Host saved", "The host was saved successfully")
        ReloadNginxAction.execute()
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
        const { bindings, routes, useGlobalBindings, defaultServer, domainNames } = host

        const injectBindingNeeded = !useGlobalBindings && bindings.length === 0
        const sortedRoutes = routes.sort((left, right) => (left.priority > right.priority ? 1 : -1))

        let injectDomainNeeded = false
        let newDomainNames = domainNames
        if (defaultServer && domainNames.length > 0) {
            injectDomainNeeded = true
            newDomainNames = []
        } else if (!defaultServer && domainNames.length === 0) {
            injectDomainNeeded = true
            newDomainNames = HostFormValuesDefaults.domainNames
        }

        const updatedData: HostFormValues = {
            ...host,
            routes: sortedRoutes,
            bindings: injectBindingNeeded ? HostFormValuesDefaults.bindings : bindings,
            domainNames: newDomainNames,
        }

        this.setState({ formValues: updatedData }, () => {
            if (injectBindingNeeded || injectDomainNeeded) this.formRef.current?.resetFields()
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

    private fetchAccessLists(
        pageSize: number,
        pageNumber: number,
        searchTerms?: string,
    ): Promise<PageResponse<AccessListResponse>> {
        return this.accessListService.list(pageSize, pageNumber, searchTerms)
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
                        <Form.Item
                            name="accessList"
                            validateStatus={validationResult.getStatus("accessListId")}
                            help={validationResult.getMessage("accessListId")}
                            label="Access list"
                        >
                            <PaginatedSelect<AccessListResponse>
                                itemDescription={item => item?.name}
                                itemKey={item => item?.id}
                                pageProvider={(pageSize, pageNumber, searchTerms) =>
                                    this.fetchAccessLists(pageSize, pageNumber, searchTerms)
                                }
                                allowEmpty
                            />
                        </Form.Item>
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
                    onChange={() => this.formRef.current?.resetFields()}
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
            integrations
                .then(response => {
                    this.setState({ loading: false, integrations: response })
                    this.updateShellConfig(true)
                })
                .catch(() => {
                    this.setState({ loading: false })
                    CommonNotifications.failedToFetch()
                })
            return
        }

        const formValues = this.hostService
            .getById((this.hostId ?? copyFrom)!!)
            .then(response => (response === undefined ? undefined : HostConverter.responseToFormValues(response)))

        Promise.all([formValues, integrations])
            .then(([formValues, integrations]) => {
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
            .catch(error => {
                CommonNotifications.failedToFetch()
                this.setState({ loading: false, error })
            })

        this.updateShellConfig(false)
    }

    render() {
        const { loading, notFound, error } = this.state

        if (error !== undefined) return EmptyStates.FailedToFetch
        if (notFound) return EmptyStates.NotFound
        if (loading) return <Preloader loading />

        return this.renderForm()
    }
}
