import React from "react"
import { Button, Flex, Form, FormListFieldData, FormListOperation, InputNumber, Select, Space } from "antd"
import { DeleteOutlined, PlusOutlined } from "@ant-design/icons"
import FormLayout from "../../../core/components/form/FormLayout"
import ValidationResult from "../../../core/validation/ValidationResult"
import { I18n, i18n } from "../../../core/i18n/I18n"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"

const SECOND_ONWARDS_ACTION_ICON_STYLE = {
    marginLeft: 15,
    marginTop: 9,
}

const FIRST_ACTION_ICON_STYLE = {
    marginLeft: 15,
    alignItems: "start",
    marginTop: 40,
}

export interface CacheDurationsProps {
    validationResult: ValidationResult
}

export default class CacheDurations extends React.Component<CacheDurationsProps> {
    private renderEntry(field: FormListFieldData, operations: FormListOperation, index: number) {
        const { validationResult } = this.props
        const { name } = field

        return (
            <Flex key={field.key} align="start">
                <Form.Item
                    {...FormLayout.ExpandedLabeledItem}
                    layout="vertical"
                    name={[name, "statusCodes"]}
                    validateStatus={validationResult.getStatus(`durations[${index}].statusCodes`)}
                    help={validationResult.getMessage(`durations[${index}].statusCodes`)}
                    label={index === 0 ? <I18n id={MessageKey.FrontendCacheTabsGeneralStatusCodes} /> : undefined}
                    required
                    style={{ flex: 1 }}
                >
                    <Select
                        mode="tags"
                        placeholder={i18n(MessageKey.FrontendCacheComponentsDurationsPlaceholderStatusCodes)}
                        tokenSeparators={[",", " "]}
                    />
                </Form.Item>
                <Form.Item
                    {...FormLayout.ExpandedLabeledItem}
                    layout="vertical"
                    validateStatus={validationResult.getStatus(`durations[${index}].validTimeSeconds`)}
                    help={validationResult.getMessage(`durations[${index}].validTimeSeconds`)}
                    label={index === 0 ? <I18n id={MessageKey.FrontendCacheTabsGeneralValidTime} /> : undefined}
                    required
                    style={{ flex: 1, marginLeft: 10 }}
                >
                    <Space.Compact style={{ width: "100%" }}>
                        <Form.Item name={[name, "validTimeSeconds"]} noStyle>
                            <InputNumber min={0} style={{ width: "100%" }} />
                        </Form.Item>
                        <Space.Addon>
                            <I18n id={MessageKey.CommonUnitSeconds} />
                        </Space.Addon>
                    </Space.Compact>
                </Form.Item>

                <DeleteOutlined
                    onClick={() => operations.remove(index)}
                    style={index === 0 ? FIRST_ACTION_ICON_STYLE : SECOND_ONWARDS_ACTION_ICON_STYLE}
                />
            </Flex>
        )
    }

    private renderDurations(fields: FormListFieldData[], operations: FormListOperation) {
        const entries = fields.map((field, index) => this.renderEntry(field, operations, index))

        const addAction = (
            <Form.Item>
                <Button
                    type="dashed"
                    onClick={() => operations.add({ statusCodes: [], validTimeSeconds: 600 })}
                    icon={<PlusOutlined />}
                >
                    <I18n id={MessageKey.CommonAddDuration} />
                </Button>
            </Form.Item>
        )

        return [...entries, addAction]
    }

    render() {
        return (
            <Form.List name="durations">{(fields, operations) => this.renderDurations(fields, operations)}</Form.List>
        )
    }
}
