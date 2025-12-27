import React from "react"
import { Button, Flex, Form, FormListFieldData, FormListOperation, Input } from "antd"
import { DeleteOutlined, PlusOutlined } from "@ant-design/icons"
import FormLayout from "../../../core/components/form/FormLayout"
import ValidationResult from "../../../core/validation/ValidationResult"

const ACTION_ICON_STYLE = {
    marginLeft: 15,
    marginTop: 7,
}

export interface CacheRulesProps {
    name: string
    label: string
    validationResult: ValidationResult
}

export default class CacheRules extends React.Component<CacheRulesProps> {
    private renderEntry(field: FormListFieldData, operations: FormListOperation, index: number) {
        const { validationResult, name: listName } = this.props
        const { name } = field

        return (
            <Flex key={field.key} align="start">
                <Form.Item
                    {...FormLayout.ExpandedLabeledItem}
                    layout="vertical"
                    name={name}
                    validateStatus={validationResult.getStatus(`${listName}[${index}]`)}
                    help={validationResult.getMessage(`${listName}[${index}]`)}
                    required
                    style={{ flex: 1 }}
                >
                    <Input placeholder="e.g. $http_cache_control ~* no-cache" />
                </Form.Item>

                <DeleteOutlined onClick={() => operations.remove(index)} style={ACTION_ICON_STYLE} />
            </Flex>
        )
    }

    private renderRules(fields: FormListFieldData[], operations: FormListOperation) {
        const { label } = this.props
        const entries = fields.map((field, index) => this.renderEntry(field, operations, index))

        const addAction = (
            <Form.Item>
                <Button type="dashed" onClick={() => operations.add("")} icon={<PlusOutlined />}>
                    Add {label.toLowerCase()}
                </Button>
            </Form.Item>
        )

        return [...entries, addAction]
    }

    render() {
        const { name } = this.props
        return <Form.List name={name}>{(fields, operations) => this.renderRules(fields, operations)}</Form.List>
    }
}
