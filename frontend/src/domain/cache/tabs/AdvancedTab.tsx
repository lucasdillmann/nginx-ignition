import React from "react"
import { Alert, Flex, Form, InputNumber, Space, Switch } from "antd"
import ValidationResult from "../../../core/validation/ValidationResult"
import CacheRules from "../components/CacheRules"
import { Link } from "react-router-dom"
import { I18n } from "../../../core/i18n/I18n"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"

export interface AdvancedTabProps {
    validationResult: ValidationResult
}

export default class AdvancedTab extends React.Component<AdvancedTabProps> {
    render() {
        const { validationResult } = this.props

        return (
            <>
                <Alert
                    message={<I18n id={MessageKey.CommonWarningProceedWithCaution} />}
                    description={<I18n id={MessageKey.CommonWarningProceedWithCautionDescription} />}
                    type="warning"
                    style={{ marginBottom: 20 }}
                    showIcon
                />
                <h2 className="cache-form-section-name" style={{ marginTop: 30 }}>
                    <I18n id={MessageKey.FrontendCacheTabsAdvancedConcurrencyLock} />
                </h2>
                <p className="cache-form-section-help-text">
                    <I18n id={MessageKey.FrontendCacheTabsAdvancedConcurrencyLockDescription} />
                </p>
                <Form.Item
                    name={["concurrencyLock", "enabled"]}
                    validateStatus={validationResult.getStatus("concurrencyLock.enabled")}
                    help={validationResult.getMessage("concurrencyLock.enabled")}
                    label={<I18n id={MessageKey.FrontendCacheTabsAdvancedEnableConcurrencyLock} />}
                    required
                >
                    <Switch />
                </Form.Item>
                <Form.Item
                    label={<I18n id={MessageKey.FrontendCacheTabsAdvancedLockTimeout} />}
                    validateStatus={validationResult.getStatus("concurrencyLock.timeoutSeconds")}
                    help={validationResult.getMessage("concurrencyLock.timeoutSeconds")}
                >
                    <Space.Compact style={{ width: "100%" }}>
                        <Form.Item name={["concurrencyLock", "timeoutSeconds"]} noStyle>
                            <InputNumber min={0} style={{ width: "100%" }} />
                        </Form.Item>
                        <Space.Addon>
                            <I18n id={MessageKey.CommonUnitSeconds} />
                        </Space.Addon>
                    </Space.Compact>
                </Form.Item>
                <Form.Item
                    label={<I18n id={MessageKey.FrontendCacheTabsAdvancedLockAge} />}
                    validateStatus={validationResult.getStatus("concurrencyLock.ageSeconds")}
                    help={
                        validationResult.getMessage("concurrencyLock.ageSeconds") ?? (
                            <I18n id={MessageKey.FrontendCacheTabsAdvancedLockAgeHelp} />
                        )
                    }
                >
                    <Space.Compact style={{ width: "100%" }}>
                        <Form.Item name={["concurrencyLock", "ageSeconds"]} noStyle>
                            <InputNumber min={0} style={{ width: "100%" }} />
                        </Form.Item>
                        <Space.Addon>
                            <I18n id={MessageKey.CommonUnitSeconds} />
                        </Space.Addon>
                    </Space.Compact>
                </Form.Item>

                <h2 className="cache-form-section-name">
                    <I18n id={MessageKey.FrontendCacheTabsAdvancedBypassNoCacheRules} />
                </h2>
                <p className="cache-form-section-help-text">
                    <I18n id={MessageKey.FrontendCacheTabsAdvancedRulesHelpPrefix} />{" "}
                    <Link
                        to="https://nginx.org/en/docs/http/ngx_http_proxy_module.html#proxy_cache_bypass"
                        target="_blank"
                    >
                        <I18n id={MessageKey.FrontendCacheTabsAdvancedNginxDocLink} />
                    </Link>{" "}
                    <I18n id={MessageKey.FrontendCacheTabsAdvancedRulesHelpSuffix} />
                </p>
                <Flex className="cache-form-inner-flex-container">
                    <Flex className="cache-form-inner-flex-container-column">
                        <h3 className="cache-form-section-name" style={{ fontSize: 16, marginTop: 0 }}>
                            <I18n id={MessageKey.FrontendCacheTabsAdvancedBypassRules} />
                        </h3>
                        <p className="cache-form-section-help-text" style={{ marginBottom: 20 }}>
                            <I18n id={MessageKey.FrontendCacheTabsAdvancedBypassRulesHelp} />
                        </p>
                        <CacheRules
                            name="bypassRules"
                            label={MessageKey.FrontendCacheTabsAdvancedBypassRule}
                            validationResult={validationResult}
                        />
                    </Flex>

                    <Flex className="cache-form-inner-flex-container-column">
                        <h3 className="cache-form-section-name" style={{ fontSize: 16, marginTop: 0 }}>
                            <I18n id={MessageKey.FrontendCacheTabsAdvancedNoCacheRules} />
                        </h3>
                        <p className="cache-form-section-help-text" style={{ marginBottom: 20 }}>
                            <I18n id={MessageKey.FrontendCacheTabsAdvancedNoCacheRulesHelp} />
                        </p>
                        <CacheRules
                            name="noCacheRules"
                            label={MessageKey.FrontendCacheTabsAdvancedNoCacheRule}
                            validationResult={validationResult}
                        />
                    </Flex>
                </Flex>
            </>
        )
    }
}
