import React from "react"
import { Alert, Flex, Form, InputNumber, Space, Switch } from "antd"
import ValidationResult from "../../../core/validation/ValidationResult"
import CacheRules from "../components/CacheRules"
import { Link } from "react-router-dom"

export interface AdvancedTabProps {
    validationResult: ValidationResult
}

export default class AdvancedTab extends React.Component<AdvancedTabProps> {
    render() {
        const { validationResult } = this.props

        return (
            <>
                <Alert
                    message="Proceed with caution"
                    description={`
                        These settings can break your nginx server or make it misbehave. Using the default values will 
                        work just fine for almost all use cases. 
                    `}
                    type="warning"
                    style={{ marginBottom: 20 }}
                    showIcon
                />
                <h2 className="cache-form-section-name" style={{ marginTop: 30 }}>
                    Concurrency lock
                </h2>
                <p className="cache-form-section-help-text">
                    Settings to prevent multiple simultaneous requests from populating the cache at the same time.
                </p>
                <Form.Item
                    name={["concurrencyLock", "enabled"]}
                    validateStatus={validationResult.getStatus("concurrencyLock.enabled")}
                    help={
                        validationResult.getMessage("concurrencyLock.enabled") ??
                        "When enabled, only one request at a time will be allowed to populate a new cache element"
                    }
                    label="Enable concurrency lock"
                    required
                >
                    <Switch />
                </Form.Item>
                <Form.Item
                    label="Lock timeout"
                    validateStatus={validationResult.getStatus("concurrencyLock.timeoutSeconds")}
                    help={
                        validationResult.getMessage("concurrencyLock.timeoutSeconds") ??
                        "Sets the timeout for the concurrency lock"
                    }
                >
                    <Space.Compact style={{ width: "100%" }}>
                        <Form.Item name={["concurrencyLock", "timeoutSeconds"]} noStyle>
                            <InputNumber min={0} style={{ width: "100%" }} />
                        </Form.Item>
                        <Space.Addon>seconds</Space.Addon>
                    </Space.Compact>
                </Form.Item>
                <Form.Item
                    label="Lock age"
                    validateStatus={validationResult.getStatus("concurrencyLock.ageSeconds")}
                    help={
                        validationResult.getMessage("concurrencyLock.ageSeconds") ??
                        "If the last request that was populating the cache hasn't finished after this time, another request may be allowed to populate it"
                    }
                >
                    <Space.Compact style={{ width: "100%" }}>
                        <Form.Item name={["concurrencyLock", "ageSeconds"]} noStyle>
                            <InputNumber min={0} style={{ width: "100%" }} />
                        </Form.Item>
                        <Space.Addon>seconds</Space.Addon>
                    </Space.Compact>
                </Form.Item>

                <h2 className="cache-form-section-name">Bypass & no-cache rules</h2>
                <p className="cache-form-section-help-text">
                    Rules to determine when to bypass or not store the cache. Any instruction set here will be placed in
                    the nginx configuration files as-is, please check the{" "}
                    <Link
                        to="https://nginx.org/en/docs/http/ngx_http_proxy_module.html#proxy_cache_bypass"
                        target="_blank"
                    >
                        nginx documentation
                    </Link>{" "}
                    for more information on the syntax of these rules.
                </p>
                <Flex className="cache-form-inner-flex-container">
                    <Flex className="cache-form-inner-flex-container-column">
                        <h3 className="cache-form-section-name" style={{ fontSize: 16, marginTop: 0 }}>
                            Bypass rules
                        </h3>
                        <p className="cache-form-section-help-text" style={{ marginBottom: 20 }}>
                            Defines conditions under which the response will not be taken from a cache.
                        </p>
                        <CacheRules name="bypassRules" label="Bypass rule" validationResult={validationResult} />
                    </Flex>

                    <Flex className="cache-form-inner-flex-container-column">
                        <h3 className="cache-form-section-name" style={{ fontSize: 16, marginTop: 0 }}>
                            No-cache rules
                        </h3>
                        <p className="cache-form-section-help-text" style={{ marginBottom: 20 }}>
                            Defines conditions under which the response will not be saved to a cache.
                        </p>
                        <CacheRules name="noCacheRules" label="No-cache rule" validationResult={validationResult} />
                    </Flex>
                </Flex>
            </>
        )
    }
}
