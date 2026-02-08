import React from "react"
import AppShellContext, { ShellAction } from "../../core/components/shell/AppShellContext"
import { Flex, Form, FormInstance, Switch, Tooltip } from "antd"
import FormLayout from "../../core/components/form/FormLayout"
import DomainNamesList from "../../core/components/domainnames/DomainNamesList"
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
import If from "../../core/components/flowcontrol/If"
import CommonNotifications from "../../core/components/notification/CommonNotifications"
import EmptyStates from "../../core/components/emptystate/EmptyStates"
import PaginatedSelect from "../../core/components/select/PaginatedSelect"
import AccessListResponse from "../accesslist/model/AccessListResponse"
import PageResponse from "../../core/pagination/PageResponse"
import AccessListService from "../accesslist/AccessListService"
import { isAccessGranted } from "../../core/components/accesscontrol/IsAccessGranted"
import { UserAccessLevel } from "../user/model/UserAccessLevel"
import AccessDeniedPage from "../../core/components/accesscontrol/AccessDeniedPage"
import { hostFormValuesDefaults } from "./model/HostFormValuesDefaults"
import HostSupportWarning from "./components/HostSupportWarning"
import HostVpns from "./components/HostVpns"
import CacheService from "../cache/CacheService"
import CacheResponse from "../cache/model/CacheResponse"
import MessageKey from "../../core/i18n/model/MessageKey.generated"
import { i18n, I18n } from "../../core/i18n/I18n"
import NginxService from "../nginx/NginxService"
import NginxMetadata, { NginxSupportType } from "../nginx/model/NginxMetadata"

interface HostFormPageState {
    formValues: HostFormValues
    validationResult: ValidationResult
    loading: boolean
    notFound: boolean
    error?: Error
    metadata?: NginxMetadata
}

export default class HostFormPage extends React.Component<any, HostFormPageState> {
    private readonly hostService: HostService
    private readonly accessListService: AccessListService
    private readonly cacheService: CacheService
    private readonly nginxService: NginxService
    private readonly saveModal: ModalPreloader
    private readonly formRef: React.RefObject<FormInstance | null>
    private hostId?: string

    constructor(props: any) {
        super(props)

        const hostId = routeParams().id
        this.hostId = hostId === "new" ? undefined : hostId
        this.hostService = new HostService()
        this.accessListService = new AccessListService()
        this.cacheService = new CacheService()
        this.nginxService = new NginxService()
        this.saveModal = new ModalPreloader()
        this.formRef = React.createRef()
        this.state = {
            validationResult: new ValidationResult(),
            loading: true,
            notFound: false,
            formValues: hostFormValuesDefaults(),
        }
    }

    private async delete() {
        if (this.hostId === undefined) return

        return DeleteHostAction.execute(this.hostId).then(() => navigateTo("/hosts"))
    }

    private submit() {
        const payload = HostConverter.formValuesToRequest(this.state.formValues)
        this.saveModal.show(MessageKey.CommonHangOnTight, {
            id: MessageKey.CommonSavingType,
            params: { type: MessageKey.CommonHost },
        })
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
        Notification.success(
            { id: MessageKey.CommonTypeSaved, params: { type: MessageKey.CommonHost } },
            MessageKey.CommonSuccessMessage,
        )
        ReloadNginxAction.execute()
    }

    private handleError(error: Error) {
        if (error instanceof UnexpectedResponseError) {
            const validationResult = ValidationResultConverter.parse(error.response)
            if (validationResult != null) this.setState({ validationResult })
        }

        Notification.error(MessageKey.CommonThatDidntWork, MessageKey.CommonFormCheckMessage)
    }

    private updateShellConfig(enableActions: boolean) {
        if (!isAccessGranted(UserAccessLevel.READ_WRITE, permissions => permissions.hosts)) {
            enableActions = false
        }

        const actions: ShellAction[] = [
            {
                description: MessageKey.CommonSave,
                disabled: !enableActions,
                onClick: () => this.submit(),
            },
        ]

        if (this.hostId !== undefined)
            actions.unshift({
                description: MessageKey.CommonDelete,
                disabled: !enableActions,
                color: "danger",
                onClick: () => this.delete(),
            })

        AppShellContext.get().updateConfig({
            title: MessageKey.FrontendHostFormTitle,
            subtitle: MessageKey.FrontendHostFormSubtitle,
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
            newDomainNames = hostFormValuesDefaults().domainNames
        }

        const updatedData: HostFormValues = {
            ...host,
            routes: sortedRoutes,
            bindings: injectBindingNeeded ? hostFormValuesDefaults().bindings : bindings,
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

    private fetchCaches(
        pageSize: number,
        pageNumber: number,
        searchTerms?: string,
    ): Promise<PageResponse<CacheResponse>> {
        return this.cacheService.list(pageSize, pageNumber, searchTerms)
    }

    private renderStatsSwitch() {
        const { metadata, validationResult } = this.state

        const statsUnsupported = metadata !== undefined && metadata.availableSupport.stats === NginxSupportType.NONE
        if (statsUnsupported) return null

        const statsDisabled = metadata !== undefined && !metadata.stats.enabled
        const statsAllHosts = metadata !== undefined && metadata.stats.allHosts

        let switchDisabled = false
        let switchChecked: boolean | undefined = undefined
        let tooltipText: string | undefined = undefined

        if (statsDisabled) {
            switchDisabled = true
            switchChecked = false
            tooltipText = i18n(MessageKey.FrontendHostFormStatsDisabledTooltip)
        } else if (statsAllHosts) {
            switchDisabled = true
            switchChecked = true
            tooltipText = i18n(MessageKey.FrontendHostFormStatsAllHostsTooltip)
        }

        const switchElement = (
            <Form.Item
                name="statsEnabled"
                validateStatus={validationResult.getStatus("statsEnabled")}
                help={validationResult.getMessage("statsEnabled")}
                label={<I18n id={MessageKey.FrontendHostFormStatsEnabled} />}
                required
            >
                <Switch disabled={switchDisabled} checked={switchChecked} />
            </Form.Item>
        )

        if (tooltipText !== undefined) return <Tooltip title={tooltipText}>{switchElement}</Tooltip>

        return switchElement
    }

    private renderForm() {
        const { validationResult, formValues } = this.state

        return (
            <Form<HostFormValues>
                {...FormLayout.FormDefaults}
                ref={this.formRef}
                onValuesChange={(_, formValues) => this.handleChange(formValues)}
                initialValues={formValues}
            >
                <HostSupportWarning />

                <h2 className="hosts-form-section-name">
                    <I18n id={MessageKey.CommonGeneral} />
                </h2>
                <p className="hosts-form-section-help-text">
                    <I18n id={MessageKey.FrontendHostFormSectionGeneralHelp} />
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
                            label={<I18n id={MessageKey.CommonEnabled} />}
                            required
                        >
                            <Switch />
                        </Form.Item>
                        <Form.Item
                            name="defaultServer"
                            validateStatus={validationResult.getStatus("defaultServer")}
                            help={validationResult.getMessage("defaultServer")}
                            label={<I18n id={MessageKey.CommonDefaultServerLabel} />}
                            required
                        >
                            <Switch />
                        </Form.Item>
                        <Form.Item
                            name={["featureSet", "websocketsSupport"]}
                            validateStatus={validationResult.getStatus("featureSet.websocketsSupport")}
                            help={validationResult.getMessage("featureSet.websocketsSupport")}
                            label={<I18n id={MessageKey.FrontendHostFormWebsocketsSupport} />}
                            required
                        >
                            <Switch />
                        </Form.Item>
                        <Form.Item
                            name={["featureSet", "http2Support"]}
                            validateStatus={validationResult.getStatus("featureSet.http2Support")}
                            help={validationResult.getMessage("featureSet.http2Support")}
                            label={<I18n id={MessageKey.FrontendHostFormHttp2Support} />}
                            required
                        >
                            <Switch />
                        </Form.Item>
                        <Form.Item
                            name={["featureSet", "redirectHttpToHttps"]}
                            validateStatus={validationResult.getStatus("featureSet.redirectHttpToHttps")}
                            help={validationResult.getMessage("featureSet.redirectHttpToHttps")}
                            label={<I18n id={MessageKey.FrontendHostFormRedirectHttpToHttps} />}
                            required
                        >
                            <Switch />
                        </Form.Item>
                        {this.renderStatsSwitch()}
                    </Flex>
                    <Flex className="hosts-form-inner-flex-container-column">
                        <Form.Item
                            name="cache"
                            validateStatus={validationResult.getStatus("cacheId")}
                            help={validationResult.getMessage("cacheId")}
                            label={<I18n id={MessageKey.CommonCache} />}
                        >
                            <PaginatedSelect<CacheResponse>
                                itemDescription={item => item?.name}
                                itemKey={item => item?.id}
                                pageProvider={(pageSize, pageNumber, searchTerms) =>
                                    this.fetchCaches(pageSize, pageNumber, searchTerms)
                                }
                                allowEmpty
                            />
                        </Form.Item>
                        <Form.Item
                            name="accessList"
                            validateStatus={validationResult.getStatus("accessListId")}
                            help={validationResult.getMessage("accessListId")}
                            label={<I18n id={MessageKey.CommonAccessList} />}
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
                            <Form.Item label={<I18n id={MessageKey.CommonDomainNames} />} required>
                                <Flex>
                                    <I18n id={MessageKey.FrontendHostFormDefaultServerUnavailable} />
                                </Flex>
                            </Form.Item>
                        </If>
                        <DomainNamesList
                            validationResult={validationResult}
                            className={formValues.defaultServer ? "hosts-form-invisible-input" : undefined}
                        />
                    </Flex>
                </Flex>

                <h2 className="hosts-form-section-name">
                    <I18n id={MessageKey.FrontendHostFormSectionRouting} />
                </h2>
                <p className="hosts-form-section-help-text">
                    <I18n id={MessageKey.FrontendHostFormSectionRoutingHelp} />
                </p>
                <HostRoutes
                    routes={formValues.routes}
                    validationResult={validationResult}
                    onRouteRemove={index => this.removeRoute(index)}
                    onChange={() => this.formRef.current?.resetFields()}
                />

                <Flex style={{ marginTop: 50 }}>
                    <Flex style={{ flexGrow: 1 }} vertical>
                        <h2 className="hosts-form-section-name">
                            <I18n id={MessageKey.FrontendHostFormSectionStandardBindings} />
                        </h2>
                        <p className="hosts-form-section-help-text">
                            <I18n id={MessageKey.FrontendHostFormSectionStandardBindingsHelp} />
                        </p>
                    </Flex>
                    <Flex>
                        <Form.Item
                            name="useGlobalBindings"
                            validateStatus={validationResult.getStatus("useGlobalBindings")}
                            help={validationResult.getMessage("useGlobalBindings")}
                            label={<I18n id={MessageKey.FrontendHostFormUseGlobalBindings} />}
                            wrapperCol={{ style: { flexGrow: 0, minWidth: 65 } }}
                            labelCol={{ style: { order: 2, flexGrow: 1 } }}
                            required
                        >
                            <Switch />
                        </Form.Item>
                    </Flex>
                </Flex>
                <HostBindings
                    className={formValues.useGlobalBindings ? "hosts-form-invisible-input" : undefined}
                    pathPrefix="bindings"
                    bindings={formValues.bindings}
                    validationResult={validationResult}
                />

                <h2 className="hosts-form-section-name">
                    <I18n id={MessageKey.FrontendHostFormSectionVpnBindings} />
                </h2>
                <p className="hosts-form-section-help-text">
                    <I18n id={MessageKey.FrontendHostFormSectionVpnBindingsHelp} />
                </p>
                <HostVpns vpns={formValues.vpns} validationResult={validationResult} />
            </Form>
        )
    }

    componentDidMount() {
        this.updateShellConfig(false)
        const metadataPromise = this.nginxService.getMetadata()

        const copyFrom = queryParams().copyFrom as string | undefined
        if (this.hostId === undefined && copyFrom === undefined) {
            metadataPromise.then(metadata => this.setState({ metadata, loading: false }))
            this.updateShellConfig(true)
            return
        }

        const hostPromisse = this.hostService
            .getById((this.hostId ?? copyFrom)!!)
            .then(response => (response === undefined ? undefined : HostConverter.responseToFormValues(response)))

        Promise.all([metadataPromise, hostPromisse])
            .then(([metadata, formValues]) => {
                if (formValues === undefined) {
                    this.setState({ metadata, loading: false, notFound: true })
                    return
                }

                if (copyFrom !== undefined)
                    Notification.success(
                        MessageKey.FrontendHostValuesCopied,
                        MessageKey.FrontendHostValuesCopiedFullDescription,
                    )

                this.setState({ formValues, metadata, loading: false })
                this.updateShellConfig(true)
            })
            .catch(error => {
                CommonNotifications.failedToFetch()
                this.setState({ loading: false, error })
            })
    }

    render() {
        if (!isAccessGranted(UserAccessLevel.READ_ONLY, permissions => permissions.hosts)) {
            return <AccessDeniedPage />
        }

        const { loading, notFound, error } = this.state

        if (error !== undefined) return EmptyStates.FailedToFetch
        if (notFound) return EmptyStates.NotFound
        if (loading) return <Preloader loading />

        return this.renderForm()
    }
}
