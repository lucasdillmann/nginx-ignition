import { Form, FormItemProps, Input, Modal, Switch, Tabs } from "antd"
import FormLayout from "../../../core/components/form/FormLayout"
import TextArea from "antd/es/input/TextArea"
import React from "react"
import ValidationResult from "../../../core/validation/ValidationResult"
import { Link } from "react-router-dom"
import PaginatedSelect from "../../../core/components/select/PaginatedSelect"
import AccessListResponse from "../../accesslist/model/AccessListResponse"
import PageResponse from "../../../core/pagination/PageResponse"
import AccessListService from "../../accesslist/AccessListService"
import { HostRouteType } from "../model/HostRequest"
import { HostFormRoute } from "../model/HostFormValues"
import CacheResponse from "../../cache/model/CacheResponse"
import CacheService from "../../cache/CacheService"
import HostRouteConditionalConfig from "./HostRouteConditionalConfig"
import { I18n } from "../../../core/i18n/I18n"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"

const PROXY_ROUTE_TYPES: HostRouteType[] = [HostRouteType.INTEGRATION, HostRouteType.PROXY]
const ACCESS_LIST_SUPPORTED_ROUTE_TYPES: HostRouteType[] = [
    HostRouteType.INTEGRATION,
    HostRouteType.PROXY,
    HostRouteType.STATIC_FILES,
]

const ItemProps: FormItemProps = {
    labelCol: {
        sm: { span: 8 },
    },
    wrapperCol: {
        sm: { span: 16 },
    },
}

export interface HostRouteSettingsProps {
    open: boolean
    index: number
    route: HostFormRoute
    fieldPath: any
    onClose: () => void
    onCancel: () => void
    validationResult: ValidationResult
}

export default class HostRouteSettingsModal extends React.Component<HostRouteSettingsProps> {
    private readonly accessListService: AccessListService
    private readonly cacheService: CacheService

    constructor(props: HostRouteSettingsProps) {
        super(props)
        this.accessListService = new AccessListService()
        this.cacheService = new CacheService()
    }

    private renderAdvancedTab() {
        const { index, validationResult, fieldPath } = this.props
        return (
            <>
                <p>
                    <I18n id={MessageKey.FrontendHostComponentsHostroutesettingsCustomSettingsDescription1} />
                </p>
                <p>
                    <I18n id={MessageKey.FrontendHostComponentsHostroutesettingsCustomSettingsDescription2} />{" "}
                    <Link
                        to="https://nginx.org/en/docs/http/ngx_http_core_module.html#location"
                        target="_blank"
                        rel="noreferrer"
                    >
                        <I18n id={MessageKey.CommonNginxDocLink} />
                    </Link>{" "}
                    <I18n id={MessageKey.FrontendHostComponentsHostroutesettingsCustomSettingsDescription3} />
                </p>

                <Form.Item
                    {...FormLayout.ExpandedLabeledItem}
                    className="host-form-route-custom-settings"
                    name={[fieldPath, "settings", "custom"]}
                    validateStatus={validationResult.getStatus(`routes[${index}].settings.custom`)}
                    help={validationResult.getMessage(`routes[${index}].settings.custom`)}
                    label={<I18n id={MessageKey.CommonCustomSettings} />}
                    required
                >
                    <TextArea rows={10} />
                </Form.Item>
            </>
        )
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

    private renderMainTab() {
        const { index, validationResult, fieldPath, route } = this.props

        return (
            <>
                <Form.Item
                    {...ItemProps}
                    name={[fieldPath, "cache"]}
                    validateStatus={validationResult.getStatus(`routes[${index}].cacheId`)}
                    help={validationResult.getMessage(`routes[${index}].cacheId`)}
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
                <HostRouteConditionalConfig route={route} types={ACCESS_LIST_SUPPORTED_ROUTE_TYPES}>
                    <Form.Item
                        {...ItemProps}
                        name={[fieldPath, "accessList"]}
                        validateStatus={validationResult.getStatus(`routes[${index}].accessListId`)}
                        help={validationResult.getMessage(`routes[${index}].accessListId`)}
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
                </HostRouteConditionalConfig>
                <HostRouteConditionalConfig route={route} types={[HostRouteType.STATIC_FILES]}>
                    <Form.Item
                        {...ItemProps}
                        name={[fieldPath, "settings", "directoryListingEnabled"]}
                        label={<I18n id={MessageKey.FrontendHostComponentsHostroutesettingsEnableDirectoryListing} />}
                        validateStatus={validationResult.getStatus(`routes[${index}].settings.directoryListingEnabled`)}
                        help={
                            validationResult.getMessage(`routes[${index}].settings.directoryListingEnabled`) ?? (
                                <I18n
                                    id={MessageKey.FrontendHostComponentsHostroutesettingsEnableDirectoryListingHelp}
                                />
                            )
                        }
                        required
                    >
                        <Switch />
                    </Form.Item>
                </HostRouteConditionalConfig>
                <HostRouteConditionalConfig route={route} types={[HostRouteType.STATIC_FILES]}>
                    <Form.Item
                        {...ItemProps}
                        name={[fieldPath, "settings", "indexFile"]}
                        label={<I18n id={MessageKey.FrontendHostComponentsHostroutesettingsIndexFile} />}
                        validateStatus={validationResult.getStatus(`routes[${index}].settings.indexFile`)}
                        help={
                            validationResult.getMessage(`routes[${index}].settings.indexFile`) ?? (
                                <I18n id={MessageKey.FrontendHostComponentsHostroutesettingsIndexFileHelp} />
                            )
                        }
                    >
                        <Input />
                    </Form.Item>
                </HostRouteConditionalConfig>
                <HostRouteConditionalConfig route={route} types={PROXY_ROUTE_TYPES}>
                    <Form.Item
                        {...ItemProps}
                        name={[fieldPath, "settings", "keepOriginalDomainName"]}
                        label={<I18n id={MessageKey.FrontendHostComponentsHostroutesettingsKeepOriginalDomainName} />}
                        validateStatus={validationResult.getStatus(`routes[${index}].settings.keepOriginalDomainName`)}
                        help={
                            validationResult.getMessage(`routes[${index}].settings.keepOriginalDomainName`) ?? (
                                <I18n
                                    id={MessageKey.FrontendHostComponentsHostroutesettingsKeepOriginalDomainNameHelp}
                                />
                            )
                        }
                        required
                    >
                        <Switch />
                    </Form.Item>
                </HostRouteConditionalConfig>
                <HostRouteConditionalConfig route={route} types={PROXY_ROUTE_TYPES}>
                    <Form.Item
                        {...ItemProps}
                        name={[fieldPath, "settings", "proxySslServerName"]}
                        label={<I18n id={MessageKey.FrontendHostComponentsHostroutesettingsProxySslServerName} />}
                        validateStatus={validationResult.getStatus(`routes[${index}].settings.proxySslServerName`)}
                        help={
                            validationResult.getMessage(`routes[${index}].settings.proxySslServerName`) ?? (
                                <I18n id={MessageKey.FrontendHostComponentsHostroutesettingsProxySslServerNameHelp} />
                            )
                        }
                        required
                    >
                        <Switch />
                    </Form.Item>
                </HostRouteConditionalConfig>
                <HostRouteConditionalConfig route={route} types={PROXY_ROUTE_TYPES}>
                    <Form.Item
                        {...ItemProps}
                        name={[fieldPath, "settings", "includeForwardHeaders"]}
                        label={<I18n id={MessageKey.FrontendHostComponentsHostroutesettingsIncludeForwardHeaders} />}
                        validateStatus={validationResult.getStatus(`routes[${index}].settings.includeForwardHeaders`)}
                        help={
                            validationResult.getMessage(`routes[${index}].settings.includeForwardHeaders`) ?? (
                                <I18n
                                    id={MessageKey.FrontendHostComponentsHostroutesettingsIncludeForwardHeadersHelp}
                                />
                            )
                        }
                        required
                    >
                        <Switch />
                    </Form.Item>
                </HostRouteConditionalConfig>
            </>
        )
    }

    private tabsDefinitions() {
        return [
            {
                key: "main",
                label: <I18n id={MessageKey.FrontendHostComponentsHostroutesettingsTabMain} />,
                children: this.renderMainTab(),
            },
            {
                key: "advanced",
                label: <I18n id={MessageKey.CommonAdvanced} />,
                children: this.renderAdvancedTab(),
            },
        ]
    }

    render() {
        const { open, onClose, onCancel } = this.props
        return (
            <Modal
                title={<I18n id={MessageKey.FrontendHostComponentsHostroutesettingsTitle} />}
                open={open}
                afterClose={onClose}
                onCancel={onCancel}
                footer={null}
                width={750}
            >
                <Tabs defaultActiveKey="1" items={this.tabsDefinitions()} />
            </Modal>
        )
    }
}
