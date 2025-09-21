import React from "react"
import ValidationResult from "../../../core/validation/ValidationResult"
import { Button, Flex, Form, FormListFieldData, FormListOperation, Input } from "antd"
import { DeleteOutlined, PlusOutlined } from "@ant-design/icons"
import FormLayout from "../../../core/components/form/FormLayout"
import { AccessListFormCredentialsDefaults } from "../AccessListFormDefaults"
import Password from "antd/es/input/Password"
import "./AccessListCredentials.css"

const ACTION_ICON_STYLE = {
    marginLeft: 15,
    alignItems: "start",
    marginTop: 37,
}

export interface AccessListCredentialsProps {
    validationResult: ValidationResult
}

export default class AccessListCredentials extends React.Component<AccessListCredentialsProps> {
    private renderEntry(field: FormListFieldData, operations: FormListOperation, index: number) {
        const { validationResult } = this.props
        const { name } = field

        return (
            <Flex className="access-list-credentials-container">
                <Form.Item
                    {...FormLayout.ExpandedLabeledItem}
                    className="access-list-credentials-username"
                    layout="vertical"
                    name={[name, "username"]}
                    validateStatus={validationResult.getStatus(`credentials[${index}].username`)}
                    help={validationResult.getMessage(`credentials[${index}].username`)}
                    label="Username"
                    required
                >
                    <Input />
                </Form.Item>
                <Form.Item
                    {...FormLayout.ExpandedLabeledItem}
                    className="access-list-credentials-password"
                    layout="vertical"
                    name={[name, "password"]}
                    validateStatus={validationResult.getStatus(`credentials[${index}].password`)}
                    help={validationResult.getMessage(`credentials[${index}].password`)}
                    label="Password"
                    required
                >
                    <Password />
                </Form.Item>

                <DeleteOutlined onClick={() => operations.remove(index)} style={ACTION_ICON_STYLE} />
            </Flex>
        )
    }

    private renderCredentials(fields: FormListFieldData[], operations: FormListOperation) {
        const entries = fields.map((field, index) => this.renderEntry(field, operations, index))

        const addAction = (
            <Form.Item>
                <Button
                    type="dashed"
                    onClick={() => operations.add({ ...AccessListFormCredentialsDefaults })}
                    icon={<PlusOutlined />}
                >
                    Add credentials
                </Button>
            </Form.Item>
        )

        return [...entries, addAction]
    }

    render() {
        return (
            <Form.List name="credentials">
                {(fields, operations) => this.renderCredentials(fields, operations)}
            </Form.List>
        )
    }
}
