import React from "react"
import { Flex, Form, Input, InputNumber, Select, Space, Switch } from "antd"
import { HttpMethod, UseStale } from "../model/CacheRequest"
import ValidationResult from "../../../core/validation/ValidationResult"
import CacheDurations from "../components/CacheDurations"
import { I18n } from "../../../core/i18n/I18n"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"

const CACHE_ALLOWED_METHOD_OPTIONS = Object.values(HttpMethod).map(method => ({
    value: method,
    label: method,
}))

const USE_STALE_OPTIONS_DATA = [
    { value: UseStale.ERROR, messageKey: MessageKey.FrontendCacheTabsGeneralStaleError },
    { value: UseStale.TIMEOUT, messageKey: MessageKey.FrontendCacheTabsGeneralStaleTimeout },
    { value: UseStale.INVALID_HEADER, messageKey: MessageKey.FrontendCacheTabsGeneralStaleInvalidHeader },
    { value: UseStale.UPDATING, messageKey: MessageKey.FrontendCacheTabsGeneralStaleUpdating },
    { value: UseStale.HTTP_500, messageKey: MessageKey.FrontendCacheTabsGeneralStaleHttp500 },
    { value: UseStale.HTTP_502, messageKey: MessageKey.FrontendCacheTabsGeneralStaleHttp502 },
    { value: UseStale.HTTP_503, messageKey: MessageKey.FrontendCacheTabsGeneralStaleHttp503 },
    { value: UseStale.HTTP_504, messageKey: MessageKey.FrontendCacheTabsGeneralStaleHttp504 },
    { value: UseStale.HTTP_403, messageKey: MessageKey.FrontendCacheTabsGeneralStaleHttp403 },
    { value: UseStale.HTTP_404, messageKey: MessageKey.FrontendCacheTabsGeneralStaleHttp404 },
    { value: UseStale.HTTP_429, messageKey: MessageKey.FrontendCacheTabsGeneralStaleHttp429 },
]

export interface GeneralTabProps {
    validationResult: ValidationResult
}

export default class GeneralTab extends React.Component<GeneralTabProps> {
    render() {
        const { validationResult } = this.props

        return (
            <>
                <Flex className="cache-form-inner-flex-container">
                    <Flex className="cache-form-inner-flex-container-column cache-form-expanded-label-size">
                        <h2 className="cache-form-section-name">
                            <I18n id={MessageKey.CommonGeneral} />
                        </h2>
                        <p className="cache-form-section-help-text">
                            <I18n id={MessageKey.FrontendCacheTabsGeneralSectionGeneralHelp} />
                        </p>
                        <Form.Item
                            name="cacheStatusResponseHeaderEnabled"
                            validateStatus={validationResult.getStatus("cacheStatusResponseHeaderEnabled")}
                            help={validationResult.getMessage("cacheStatusResponseHeaderEnabled")}
                            label={<I18n id={MessageKey.FrontendCacheTabsGeneralCacheStatusResponseHeader} />}
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
                            name="storagePath"
                            validateStatus={validationResult.getStatus("storagePath")}
                            help={validationResult.getMessage("storagePath")}
                            label={<I18n id={MessageKey.FrontendCacheTabsGeneralStoragePath} />}
                        >
                            <Input />
                        </Form.Item>
                        <Form.Item
                            label={<I18n id={MessageKey.CommonMaximumSize} />}
                            validateStatus={validationResult.getStatus("maximumSizeMb")}
                            help={validationResult.getMessage("maximumSizeMb")}
                        >
                            <Space.Compact style={{ width: "100%" }}>
                                <Form.Item name="maximumSizeMb" noStyle>
                                    <InputNumber min={0} style={{ width: "100%" }} />
                                </Form.Item>
                                <Space.Addon>
                                    <I18n id={MessageKey.CommonUnitMb} />
                                </Space.Addon>
                            </Space.Compact>
                        </Form.Item>
                        <Form.Item
                            label={<I18n id={MessageKey.FrontendCacheTabsGeneralInactiveTime} />}
                            validateStatus={validationResult.getStatus("inactiveSeconds")}
                            help={
                                validationResult.getMessage("inactiveSeconds") ?? (
                                    <I18n id={MessageKey.FrontendCacheTabsGeneralInactiveTimeHelp} />
                                )
                            }
                        >
                            <Space.Compact style={{ width: "100%" }}>
                                <Form.Item name="inactiveSeconds" noStyle>
                                    <InputNumber min={0} style={{ width: "100%" }} />
                                </Form.Item>
                                <Space.Addon>
                                    <I18n id={MessageKey.CommonUnitSeconds} />
                                </Space.Addon>
                            </Space.Compact>
                        </Form.Item>
                    </Flex>

                    <Flex className="cache-form-inner-flex-container-column cache-form-expanded-label-size">
                        <h2 className="cache-form-section-name">
                            <I18n id={MessageKey.FrontendCacheTabsGeneralRequestMatching} />
                        </h2>
                        <p className="cache-form-section-help-text">
                            <I18n id={MessageKey.FrontendCacheTabsGeneralRequestMatchingHelp} />
                        </p>
                        <Form.Item
                            name="allowedMethods"
                            validateStatus={validationResult.getStatus("allowedMethods")}
                            help={validationResult.getMessage("allowedMethods")}
                            label={<I18n id={MessageKey.FrontendCacheTabsGeneralAllowedMethods} />}
                        >
                            <Select mode="multiple" options={CACHE_ALLOWED_METHOD_OPTIONS} />
                        </Form.Item>
                        <Form.Item
                            name="minimumUsesBeforeCaching"
                            validateStatus={validationResult.getStatus("minimumUsesBeforeCaching")}
                            help={
                                validationResult.getMessage("minimumUsesBeforeCaching") ?? (
                                    <I18n id={MessageKey.FrontendCacheTabsGeneralMinimumUsesBeforeCachingHelp} />
                                )
                            }
                            label={<I18n id={MessageKey.FrontendCacheTabsGeneralMinimumUsesBeforeCaching} />}
                            required
                        >
                            <InputNumber min={1} style={{ width: "100%" }} />
                        </Form.Item>
                        <Form.Item
                            name="fileExtensions"
                            validateStatus={validationResult.getStatus("fileExtensions")}
                            help={
                                validationResult.getMessage("fileExtensions") ?? (
                                    <I18n id={MessageKey.FrontendCacheTabsGeneralFileExtensionsHelp} />
                                )
                            }
                            label={<I18n id={MessageKey.CommonFileExtensions} />}
                        >
                            <Select mode="tags" />
                        </Form.Item>
                        <Form.Item
                            name="ignoreUpstreamCacheHeaders"
                            validateStatus={validationResult.getStatus("ignoreUpstreamCacheHeaders")}
                            help={
                                validationResult.getMessage("ignoreUpstreamCacheHeaders") ?? (
                                    <I18n id={MessageKey.FrontendCacheTabsGeneralIgnoreUpstreamCacheHeadersHelp} />
                                )
                            }
                            label={<I18n id={MessageKey.FrontendCacheTabsGeneralIgnoreUpstreamCacheHeaders} />}
                            required
                        >
                            <Switch />
                        </Form.Item>
                    </Flex>
                </Flex>

                <h2 className="cache-form-section-name" style={{ marginTop: 40 }}>
                    <I18n id={MessageKey.FrontendCacheTabsGeneralStaleContents} />
                </h2>
                <p className="cache-form-section-help-text">
                    <I18n id={MessageKey.FrontendCacheTabsGeneralStaleContentsHelp} />
                </p>
                <Form.Item
                    name="backgroundUpdate"
                    validateStatus={validationResult.getStatus("backgroundUpdate")}
                    help={
                        validationResult.getMessage("backgroundUpdate") ?? (
                            <I18n id={MessageKey.FrontendCacheTabsGeneralBackgroundUpdateHelp} />
                        )
                    }
                    label={<I18n id={MessageKey.FrontendCacheTabsGeneralBackgroundUpdate} />}
                    required
                >
                    <Switch />
                </Form.Item>
                <Form.Item
                    name="revalidate"
                    validateStatus={validationResult.getStatus("revalidate")}
                    help={
                        validationResult.getMessage("revalidate") ?? (
                            <I18n id={MessageKey.FrontendCacheTabsGeneralRevalidateHelp} />
                        )
                    }
                    label={<I18n id={MessageKey.FrontendCacheTabsGeneralRevalidate} />}
                    required
                >
                    <Switch />
                </Form.Item>
                <Form.Item
                    name="useStale"
                    validateStatus={validationResult.getStatus("useStale")}
                    help={
                        validationResult.getMessage("useStale") ?? (
                            <I18n id={MessageKey.FrontendCacheTabsGeneralUseStaleContentsHelp} />
                        )
                    }
                    label={<I18n id={MessageKey.FrontendCacheTabsGeneralUseStaleContents} />}
                    style={{ marginTop: 40 }}
                >
                    <Select
                        mode="multiple"
                        options={USE_STALE_OPTIONS_DATA.map(item => ({
                            value: item.value,
                            label: <I18n id={item.messageKey} />,
                        }))}
                    />
                </Form.Item>

                <h2 className="cache-form-section-name">
                    <I18n id={MessageKey.FrontendCacheTabsGeneralSectionCacheContentExpiration} />
                </h2>
                <p className="cache-form-section-help-text">
                    <I18n id={MessageKey.FrontendCacheTabsGeneralSectionCacheContentExpirationHelp} />
                </p>
                <CacheDurations validationResult={validationResult} />
            </>
        )
    }
}
