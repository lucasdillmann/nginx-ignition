import React from "react";
import AppShellContext, {ShellAction, ShellOperations} from "../../core/components/shell/AppShellContext";
import {Empty, Flex, Form, Switch} from "antd";
import FormLayout from "../../core/components/form/FormLayout";
import DomainNamesList from "../certificate/components/DomainNamesList";
import {navigateTo, routeParams} from "../../core/components/router/AppRouter";
import HostService from "./HostService";
import ValidationResult from "../../core/validation/ValidationResult";
import HostRequest, {HostBindingType, HostRouteType} from "./model/HostRequest";
import HostResponse from "./model/HostResponse";
import Preloader from "../../core/components/preloader/Preloader";
import DeleteHostAction from "./actions/DeleteHostAction";
import Notification from "../../core/components/notification/Notification";
import {UnexpectedResponseError} from "../../core/apiclient/ApiResponse";
import ValidationResultConverter from "../../core/validation/ValidationResultConverter";
import ModalPreloader from "../../core/components/preloader/ModalPreloader";
import HostRoutes from "./components/HostRoutes";
import HostBindings from "./components/HostBindings";
import "./HostFormPage.css"

const DEFAULT_HOST: HostRequest = {
    enabled: true,
    default: false,
    domainNames: [""],
    bindings: [
        {
            ip: "0.0.0.0",
            port: 8080,
            type: HostBindingType.HTTP,
        }
    ],
    routes: [
        {
            priority: 0,
            type: HostRouteType.STATIC_RESPONSE,
            sourcePath: "/",
            response: {
                statusCode: 200,
                payload: "Hello world",
            }
        }
    ],
    featureSet: {
        websocketsSupport: true,
        http2Support: true,
        redirectHttpToHttps: false,
    },
}

interface HostFormPageState {
    formValues: HostRequest
    validationResult: ValidationResult
    loading: boolean
    notFound: boolean
}

export default class HostFormPage extends React.Component<any, HostFormPageState> {
    static contextType = AppShellContext
    context!: React.ContextType<typeof AppShellContext>

    private readonly service: HostService
    private readonly hostId?: string
    private readonly saveModal: ModalPreloader

    constructor(props: any, context: ShellOperations) {
        super(props, context);

        const hostId = routeParams().id
        this.hostId = hostId === "new" ? undefined : hostId
        this.service = new HostService()
        this.saveModal = new ModalPreloader()
        this.state = {
            validationResult: new ValidationResult(),
            loading: true,
            notFound: false,
            formValues: DEFAULT_HOST,
        }
    }

    private buildPayload(): HostRequest {
        const {formValues} = this.state

        const bindings = formValues.bindings.map(binding => {
            const certificateDetails = binding.certificateId as any
            const certificateId = certificateDetails?.id
            return {
                ...binding,
                certificateId,
            }
        })

        return {
            ...formValues,
            bindings,
        }
    }

    private convertToHostRequest(response: HostResponse): HostRequest {
        const {enabled, domainNames, routes, bindings, featureSet} = response
        return {
            enabled,
            domainNames,
            routes,
            bindings,
            featureSet,
            default: response.default,
        }
    }

    private async delete() {
        if (this.hostId === undefined)
            return

        return DeleteHostAction
            .execute(this.hostId)
            .then(() => navigateTo("/hosts"))
    }

    private submit() {
        const payload = this.buildPayload()
        this.saveModal.show("Hang on tight", "We're saving the host")

        const action = this.hostId === undefined
            ? this.service.create(payload)
            : this.service.updateById(this.hostId, payload)

        action
            .then(() => this.handleSuccess())
            .catch(error => this.handleError(error))
            .then(() => this.saveModal.close())
    }

    private handleSuccess() {
        Notification.success("Host saved", "The host was saved successfully")
        navigateTo("/hosts")
    }

    private handleError(error: Error) {
        if (error instanceof UnexpectedResponseError) {
            const validationResult = ValidationResultConverter.parse(error.response)
            if (validationResult != null)
                this.setState({ validationResult })
        }

        Notification.error(
            "That didn't work",
            "Please check the form to see if everything seems correct",
        )
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


        this.context.updateConfig({
            title: "Host details",
            subtitle: "Full details and configurations of the nginx's virtual host",
            actions,
        })
    }

    private renderForm() {
        const {validationResult, formValues} = this.state

        return (
            <Form<HostRequest>
                {...FormLayout.FormDefaults}
                onValuesChange={(_, formValues) => this.setState({formValues})}
                initialValues={formValues}
            >
                <h2 className="hosts-form-section-name">General</h2>
                <Flex className="hosts-form-inner-flex-container">
                    <Flex className="hosts-form-inner-flex-container-column hosts-form-expanded-label-size" style={{ maxWidth: "30%" }}>
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
                            name="default"
                            validateStatus={validationResult.getStatus("default")}
                            help={validationResult.getMessage("default")}
                            label="Default"
                            required
                        >
                            <Switch/>
                        </Form.Item>
                        <Form.Item
                            name={["featureSet", "websocketsSupport"]}
                            validateStatus={validationResult.getStatus("featureSet.websocketsSupport")}
                            help={validationResult.getMessage("featureSet.websocketsSupport")}
                            label="Websockets support"
                            required
                        >
                            <Switch/>
                        </Form.Item>
                        <Form.Item
                            name={["featureSet", "http2Support"]}
                            validateStatus={validationResult.getStatus("featureSet.http2Support")}
                            help={validationResult.getMessage("featureSet.http2Support")}
                            label="HTTP2 support"
                            required
                        >
                            <Switch/>
                        </Form.Item>
                        <Form.Item
                            name={["featureSet", "redirectHttpToHttps"]}
                            validateStatus={validationResult.getStatus("featureSet.redirectHttpToHttps")}
                            help={validationResult.getMessage("featureSet.redirectHttpToHttps")}
                            label="Redirect HTTP to HTTPS"
                            required
                        >
                            <Switch/>
                        </Form.Item>
                    </Flex>
                    <Flex className="hosts-form-inner-flex-container-column">
                        <DomainNamesList validationResult={validationResult} />
                    </Flex>
                </Flex>

                <h2 className="hosts-form-section-name">Routing</h2>
                <HostRoutes routes={formValues.routes} validationResult={validationResult}/>

                <h2 className="hosts-form-section-name">Bindings</h2>
                <HostBindings bindings={formValues.bindings} validationResult={validationResult}/>
            </Form>
        )
    }

    componentDidMount() {
        if (this.hostId === undefined) {
            this.setState({loading: false})
            this.updateShellConfig(true)
            return
        }

        this.service
            .getById(this.hostId!!)
            .then(hostDetails => {
                if (hostDetails === undefined)
                    this.setState({loading: false, notFound: true})
                else {
                    this.setState({loading: false, formValues: this.convertToHostRequest(hostDetails)})
                    this.updateShellConfig(true)
                }
            })

        this.updateShellConfig(false)
    }

    render() {
        const {loading, notFound} = this.state

        if (notFound)
            return <Empty description="Not found" />

        if (loading)
            return <Preloader loading />

        return this.renderForm()
    }
}
