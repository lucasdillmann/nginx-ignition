import React from "react"
import { navigateTo, routeParams } from "../../core/components/router/AppRouter"
import CacheService from "./CacheService"
import { Form, FormInstance, Input, InputNumber, Select, Switch } from "antd"
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
import DeleteCacheAction from "./actions/DeleteCacheAction"
import ReloadNginxAction from "../nginx/actions/ReloadNginxAction"
import CacheRequest, { HttpMethod, UseStale } from "./model/CacheRequest"
import { UserAccessLevel } from "../user/model/UserAccessLevel"
import AccessControl from "../../core/components/accesscontrol/AccessControl"
import { isAccessGranted } from "../../core/components/accesscontrol/IsAccessGranted"
import { cacheFormDefaults } from "./CacheFormDefaults"
import CacheRules from "./components/CacheRules"
import CacheDurations from "./components/CacheDurations"
import "./CacheFormPage.css"

interface CacheFormState {
    formValues: CacheRequest
    validationResult: ValidationResult
    loading: boolean
    notFound: boolean
    error?: Error
}

export default class CacheFormPage extends React.Component<unknown, CacheFormState> {
    private readonly service: CacheService
    private readonly saveModal: ModalPreloader
    private readonly formRef: React.RefObject<FormInstance | null>
    private cacheId?: string

    constructor(props: any) {
        super(props)
        const cacheId = routeParams().id
        this.formRef = React.createRef()
        this.cacheId = cacheId === "new" ? undefined : cacheId
        this.service = new CacheService()
        this.saveModal = new ModalPreloader()
        this.state = {
            formValues: cacheFormDefaults(),
            validationResult: new ValidationResult(),
            loading: true,
            notFound: false,
        }
    }

    private submit() {
        const { formValues } = this.state
        this.saveModal.show("Hang on tight", "We're saving the cache configuration")
        this.setState({ validationResult: new ValidationResult() })

        const action =
            this.cacheId === undefined
                ? this.service.create(formValues).then(response => this.updateId(response.id))
                : this.service.updateById(this.cacheId, formValues)

        action.then(() => this.handleSuccess()).catch(error => this.handleError(error))
    }

    private updateId(id: string) {
        this.cacheId = id
        navigateTo(`/caches/${id}`, true)
        this.updateShellConfig(true)
    }

    private handleSuccess() {
        this.saveModal.close()
        Notification.success("Cache configuration saved", "The cache configuration was saved successfully")
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

    private handleChange(newValues: CacheRequest) {
        this.setState({
            formValues: {
                ...newValues,
            },
        })
    }

    private renderForm() {
        const { formValues, validationResult } = this.state

        return (
            <Form<CacheRequest>
                {...FormLayout.FormDefaults}
                ref={this.formRef}
                onValuesChange={(_, formValues) => this.handleChange(formValues)}
                initialValues={formValues}
            >
                <h2 className="cache-form-section-name">General</h2>
                <p className="cache-form-section-help-text">Main definitions and properties of the cache configuration.</p>
                <Form.Item
                    name="name"
                    validateStatus={validationResult.getStatus("name")}
                    help={validationResult.getMessage("name")}
                    label="Name"
                    required
                >
                    <Input />
                </Form.Item>
                <Form.Item
                    name="storagePath"
                    validateStatus={validationResult.getStatus("storagePath")}
                    help={validationResult.getMessage("storagePath") ?? "The directory where the cached files will be stored"}
                    label="Storage path"
                >
                    <Input />
                </Form.Item>
                <Form.Item
                    name="inactiveSeconds"
                    validateStatus={validationResult.getStatus("inactiveSeconds")}
                    help={validationResult.getMessage("inactiveSeconds") ?? "Cached data that are not accessed during the time specified by the inactive parameter are removed from the cache"}
                    label="Inactive (seconds)"
                >
                    <InputNumber min={0} style={{ width: "100%" }} />
                </Form.Item>
                <Form.Item
                    name="maximumSizeMb"
                    validateStatus={validationResult.getStatus("maximumSizeMb")}
                    help={validationResult.getMessage("maximumSizeMb") ?? "The maximum size of the cache in megabytes"}
                    label="Maximum size (MB)"
                >
                    <InputNumber min={0} style={{ width: "100%" }} />
                </Form.Item>

                <h2 className="cache-form-section-name">Request</h2>
                <p className="cache-form-section-help-text">Define which requests should be cached and how.</p>
                <Form.Item
                    name="allowedMethods"
                    validateStatus={validationResult.getStatus("allowedMethods")}
                    help={validationResult.getMessage("allowedMethods") ?? "The HTTP methods for which the response will be cached"}
                    label="Allowed methods"
                >
                    <Select mode="multiple">
                        {Object.values(HttpMethod).map(method => (
                            <Select.Option key={method} value={method}>
                                {method}
                            </Select.Option>
                        ))}
                    </Select>
                </Form.Item>
                <Form.Item
                    name="minimumUsesBeforeCaching"
                    validateStatus={validationResult.getStatus("minimumUsesBeforeCaching")}
                    help={validationResult.getMessage("minimumUsesBeforeCaching") ?? "The number of requests after which the response will be cached"}
                    label="Minimum uses before caching"
                    required
                >
                    <InputNumber min={1} style={{ width: "100%" }} />
                </Form.Item>

                <h2 className="cache-form-section-name">Stale & Revalidation</h2>
                <p className="cache-form-section-help-text">Define how the cache should behave when the content is stale.</p>
                <Form.Item
                    name="useStale"
                    validateStatus={validationResult.getStatus("useStale")}
                    help={validationResult.getMessage("useStale") ?? "Determines in which cases a stale cached response can be used"}
                    label="Use stale"
                >
                    <Select mode="multiple">
                        {Object.values(UseStale).map(stale => (
                            <Select.Option key={stale} value={stale}>
                                {stale}
                            </Select.Option>
                        ))}
                    </Select>
                </Form.Item>
                <Form.Item
                    name="backgroundUpdate"
                    validateStatus={validationResult.getStatus("backgroundUpdate")}
                    help={validationResult.getMessage("backgroundUpdate") ?? "Allows starting a background subrequest to update an expired cache item, while a stale cached response is returned to the client"}
                    label="Background update"
                    required
                >
                    <Switch />
                </Form.Item>
                <Form.Item
                    name="revalidate"
                    validateStatus={validationResult.getStatus("revalidate")}
                    help={validationResult.getMessage("revalidate") ?? "Enables revalidation of expired cache items using conditional requests with 'If-Modified-Since' and 'If-None-Match' header fields"}
                    label="Revalidate"
                    required
                >
                    <Switch />
                </Form.Item>

                <h2 className="cache-form-section-name">Concurrency Lock</h2>
                <p className="cache-form-section-help-text">Settings to prevent multiple simultaneous requests from populating the cache at the same time.</p>
                <Form.Item
                    name={["concurrencyLock", "enabled"]}
                    validateStatus={validationResult.getStatus("concurrencyLock.enabled")}
                    help={validationResult.getMessage("concurrencyLock.enabled") ?? "When enabled, only one request at a time will be allowed to populate a new cache element"}
                    label="Enable concurrency lock"
                    required
                >
                    <Switch />
                </Form.Item>
                <Form.Item
                    name={["concurrencyLock", "timeoutSeconds"]}
                    validateStatus={validationResult.getStatus("concurrencyLock.timeoutSeconds")}
                    help={validationResult.getMessage("concurrencyLock.timeoutSeconds") ?? "Sets the timeout for the concurrency lock"}
                    label="Lock timeout (seconds)"
                >
                    <InputNumber min={0} style={{ width: "100%" }} />
                </Form.Item>
                <Form.Item
                    name={["concurrencyLock", "ageSeconds"]}
                    validateStatus={validationResult.getStatus("concurrencyLock.ageSeconds")}
                    help={validationResult.getMessage("concurrencyLock.ageSeconds") ?? "If the last request that was populating the cache hasn't finished after this time, another request may be allowed to populate it"}
                    label="Lock age (seconds)"
                >
                    <InputNumber min={0} style={{ width: "100%" }} />
                </Form.Item>

                <h2 className="cache-form-section-name">Bypass & No-Cache Rules</h2>
                <p className="cache-form-section-help-text">Rules to determine when to bypass or not store the cache.</p>
                <h3 className="cache-form-section-name" style={{ fontSize: 16, marginTop: 20 }}>Bypass Rules</h3>
                <p className="cache-form-section-help-text" style={{ marginBottom: 10 }}>Defines conditions under which the response will not be taken from a cache.</p>
                <CacheRules name="bypassRules" label="Bypass Rule" validationResult={validationResult} />

                <h3 className="cache-form-section-name" style={{ fontSize: 16, marginTop: 20 }}>No-Cache Rules</h3>
                <p className="cache-form-section-help-text" style={{ marginBottom: 10 }}>Defines conditions under which the response will not be saved to a cache.</p>
                <CacheRules name="noCacheRules" label="No-Cache Rule" validationResult={validationResult} />

                <h2 className="cache-form-section-name">Durations</h2>
                <p className="cache-form-section-help-text">Define how long the cache should be valid for different status codes.</p>
                <CacheDurations validationResult={validationResult} />
            </Form>
        )
    }

    private async delete() {
        if (this.cacheId === undefined) return

        return DeleteCacheAction.execute(this.cacheId).then(() => navigateTo("/caches"))
    }

    private updateShellConfig(enableActions: boolean) {
        if (!isAccessGranted(UserAccessLevel.READ_WRITE, permissions => permissions.caches)) {
            enableActions = false
        }

        const actions: ShellAction[] = [
            {
                description: "Save",
                disabled: !enableActions,
                onClick: () => this.submit(),
            },
        ]

        if (this.cacheId !== undefined)
            actions.unshift({
                description: "Delete",
                disabled: !enableActions,
                color: "danger",
                onClick: () => this.delete(),
            })

        AppShellContext.get().updateConfig({
            title: "Cache configuration details",
            subtitle: "Full details of the cache configuration",
            actions,
        })
    }

    componentDidMount() {
        if (this.cacheId === undefined) {
            this.setState({ loading: false })
            this.updateShellConfig(true)
            return
        }

        this.service
            .getById(this.cacheId!!)
            .then(cache => {
                if (cache === undefined) this.setState({ loading: false, notFound: true })
                else {
                    this.setState({ loading: false, formValues: cache })
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

        return (
            <AccessControl
                requiredAccessLevel={UserAccessLevel.READ_ONLY}
                permissionResolver={permissions => permissions.caches}
            >
                {this.renderForm()}
            </AccessControl>
        )
    }
}
