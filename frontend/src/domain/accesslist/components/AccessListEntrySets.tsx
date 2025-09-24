import React from "react"
import ValidationResult from "../../../core/validation/ValidationResult"
import { Button, Flex, Form, FormListFieldData, FormListOperation, Select } from "antd"
import { ArrowDownOutlined, ArrowUpOutlined, DeleteOutlined, PlusOutlined } from "@ant-design/icons"
import FormLayout from "../../../core/components/form/FormLayout"
import TextArea from "antd/es/input/TextArea"
import { AccessListOutcome } from "../model/AccessListRequest"
import "./AccessListEntrySets.css"
import { AccessListEntrySetFormValues } from "../model/AccessListFormValues"
import { accessListFormEntryDefaults } from "../AccessListFormDefaults"

const ACTION_ICON_STYLE = {
    marginLeft: 15,
    alignItems: "start",
    marginTop: 37,
}

const DISABLED_ACTION_ICON_STYLE = {
    ...ACTION_ICON_STYLE,
    color: "var(--nginxIgnition-colorTextSecondary)",
}

const ENABLED_ACTION_ICON_STYLE = {
    ...ACTION_ICON_STYLE,
    color: "var(--nginxIgnition-colorText)",
}

export interface AccessListEntrySetsProps {
    entrySets: AccessListEntrySetFormValues[]
    validationResult: ValidationResult
    onRemove: (index: number) => void
}

export default class AccessListEntrySets extends React.Component<AccessListEntrySetsProps> {
    private moveEntry(operations: FormListOperation, index: number, offset: number) {
        const { entrySets } = this.props

        if (index === 0 && offset < 0) return
        if (index === entrySets.length - 1 && offset > 0) return

        const currentPosition = entrySets[index]
        const newPosition = entrySets[index + offset]

        currentPosition.priority = currentPosition.priority + offset
        newPosition.priority = newPosition.priority - offset

        const indexToRemove = offset > 0 ? index : index + offset
        operations.remove(indexToRemove)
        operations.remove(indexToRemove)
        operations.add(currentPosition)
        operations.add(newPosition)
    }

    private removeEntry(index: number) {
        const { onRemove } = this.props
        onRemove(index)
    }

    private renderEntry(field: FormListFieldData, operations: FormListOperation, index: number) {
        const { validationResult, entrySets } = this.props
        const { name } = field

        return (
            <Flex className="access-list-entry-container">
                <Form.Item
                    {...FormLayout.ExpandedLabeledItem}
                    className="access-list-entry-outcome"
                    layout="vertical"
                    name={[name, "outcome"]}
                    validateStatus={validationResult.getStatus(`entries[${index}].outcome`)}
                    help={validationResult.getMessage(`entries[${index}].outcome`)}
                    label="Outcome"
                    required
                >
                    <Select>
                        <Select.Option value={AccessListOutcome.DENY}>Deny access</Select.Option>
                        <Select.Option value={AccessListOutcome.ALLOW}>Allow access</Select.Option>
                    </Select>
                </Form.Item>
                <Form.Item
                    {...FormLayout.ExpandedLabeledItem}
                    className="access-list-entry-ip-list"
                    layout="vertical"
                    name={[name, "sourceAddresses"]}
                    validateStatus={validationResult.getStatus(`entries[${index}].sourceAddresses`)}
                    help={validationResult.getMessage(`entries[${index}].sourceAddresses`)}
                    label="IP addresses or ranges"
                    required
                >
                    <TextArea rows={4} />
                </Form.Item>

                <ArrowUpOutlined
                    onClick={() => this.moveEntry(operations, index, -1)}
                    style={index === 0 ? DISABLED_ACTION_ICON_STYLE : ENABLED_ACTION_ICON_STYLE}
                />
                <ArrowDownOutlined
                    onClick={() => this.moveEntry(operations, index, 1)}
                    style={index === entrySets.length - 1 ? DISABLED_ACTION_ICON_STYLE : ENABLED_ACTION_ICON_STYLE}
                />

                <DeleteOutlined onClick={() => this.removeEntry(index)} style={ACTION_ICON_STYLE} />
            </Flex>
        )
    }

    private renderEntries(fields: FormListFieldData[], operations: FormListOperation) {
        const entries = fields.map((field, index) => this.renderEntry(field, operations, index))

        const addAction = (
            <Form.Item>
                <Button
                    type="dashed"
                    onClick={() =>
                        operations.add({
                            ...accessListFormEntryDefaults(),
                            priority: fields.length,
                        })
                    }
                    icon={<PlusOutlined />}
                >
                    Add IP address list
                </Button>
            </Form.Item>
        )

        return [...entries, addAction]
    }

    render() {
        return <Form.List name="entries">{(fields, operations) => this.renderEntries(fields, operations)}</Form.List>
    }
}
