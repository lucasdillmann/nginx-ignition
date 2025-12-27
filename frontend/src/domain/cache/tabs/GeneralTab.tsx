import React from "react"
import { Flex, Form, Input, InputNumber, Select, Space, Switch } from "antd"
import { HttpMethod, UseStale } from "../model/CacheRequest"
import ValidationResult from "../../../core/validation/ValidationResult"
import CacheDurations from "../components/CacheDurations"

export interface GeneralTabProps {
    validationResult: ValidationResult
}

export default class GeneralTab extends React.Component<GeneralTabProps> {
    private formatStaleLabel(stale: UseStale): string {
        switch (stale) {
            case UseStale.ERROR:
                return "When an error happens"
            case UseStale.TIMEOUT:
                return "When timeout occurs"
            case UseStale.INVALID_HEADER:
                return "When an invalid header was sent"
            case UseStale.UPDATING:
                return "While updating the cache"
            case UseStale.HTTP_500:
                return "On HTTP 500"
            case UseStale.HTTP_502:
                return "On HTTP 502"
            case UseStale.HTTP_503:
                return "On HTTP 503"
            case UseStale.HTTP_504:
                return "On HTTP 504"
            case UseStale.HTTP_403:
                return "On HTTP 403"
            case UseStale.HTTP_404:
                return "On HTTP 404"
            case UseStale.HTTP_429:
                return "On HTTP 429"
            default:
                return stale
        }
    }

    render() {
        const { validationResult } = this.props

        return (
            <>
                <Flex className="cache-form-inner-flex-container">
                    <Flex className="cache-form-inner-flex-container-column cache-form-expanded-label-size">
                        <h2 className="cache-form-section-name">General</h2>
                        <p className="cache-form-section-help-text">
                            Main definitions and properties of the cache configuration.
                        </p>
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
                            help={validationResult.getMessage("storagePath")}
                            label="Storage path"
                        >
                            <Input />
                        </Form.Item>
                        <Form.Item
                            label="Maximum size"
                            validateStatus={validationResult.getStatus("maximumSizeMb")}
                            help={validationResult.getMessage("maximumSizeMb")}
                        >
                            <Space.Compact style={{ width: "100%" }}>
                                <Form.Item name="maximumSizeMb" noStyle>
                                    <InputNumber min={0} style={{ width: "100%" }} />
                                </Form.Item>
                                <Space.Addon>MB</Space.Addon>
                            </Space.Compact>
                        </Form.Item>
                        <Form.Item
                            label="Inactive time"
                            validateStatus={validationResult.getStatus("inactiveSeconds")}
                            help={
                                validationResult.getMessage("inactiveSeconds") ??
                                "How much time should pass before a cached response with no uses is deleted"
                            }
                        >
                            <Space.Compact style={{ width: "100%" }}>
                                <Form.Item name="inactiveSeconds" noStyle>
                                    <InputNumber min={0} style={{ width: "100%" }} />
                                </Form.Item>
                                <Space.Addon>seconds</Space.Addon>
                            </Space.Compact>
                        </Form.Item>
                    </Flex>

                    <Flex className="cache-form-inner-flex-container-column cache-form-expanded-label-size">
                        <h2 className="cache-form-section-name">Request</h2>
                        <p className="cache-form-section-help-text">Define which requests should be cached and how.</p>
                        <Form.Item
                            name="allowedMethods"
                            validateStatus={validationResult.getStatus("allowedMethods")}
                            help={validationResult.getMessage("allowedMethods")}
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
                            help={
                                validationResult.getMessage("minimumUsesBeforeCaching") ??
                                "The number of requests after which the response will be cached"
                            }
                            label="Minimum uses before caching"
                            required
                        >
                            <InputNumber min={1} style={{ width: "100%" }} />
                        </Form.Item>
                    </Flex>
                </Flex>

                <h2 className="cache-form-section-name" style={{ marginTop: 40 }}>
                    Stale & revalidation
                </h2>
                <p className="cache-form-section-help-text">
                    Define how the cache should behave when the content is stale and how nginx should handle
                    revalidation.
                </p>
                <Form.Item
                    name="backgroundUpdate"
                    validateStatus={validationResult.getStatus("backgroundUpdate")}
                    help={
                        validationResult.getMessage("backgroundUpdate") ??
                        "Allows starting a background subrequest to update an expired cache item, while a stale cached response is returned to the client"
                    }
                    label="Background update"
                    required
                >
                    <Switch />
                </Form.Item>
                <Form.Item
                    name="revalidate"
                    validateStatus={validationResult.getStatus("revalidate")}
                    help={
                        validationResult.getMessage("revalidate") ??
                        "Enables revalidation of expired cache items using conditional requests with 'If-Modified-Since' and 'If-None-Match' header fields"
                    }
                    label="Revalidate"
                    required
                >
                    <Switch />
                </Form.Item>
                <Form.Item
                    name="useStale"
                    validateStatus={validationResult.getStatus("useStale")}
                    help={
                        validationResult.getMessage("useStale") ??
                        "Determines in which cases a stale cached response can be used"
                    }
                    label="Use stale"
                    style={{ marginTop: 40 }}
                >
                    <Select mode="multiple">
                        {Object.values(UseStale).map(stale => (
                            <Select.Option key={stale} value={stale}>
                                {this.formatStaleLabel(stale)}
                            </Select.Option>
                        ))}
                    </Select>
                </Form.Item>

                <h2 className="cache-form-section-name">Durations</h2>
                <p className="cache-form-section-help-text">
                    Define how long the cache should be valid for different status codes.
                </p>
                <CacheDurations validationResult={validationResult} />
            </>
        )
    }
}
