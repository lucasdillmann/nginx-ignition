import React from "react"
import { Form, Input, Select, Switch, Upload } from "antd"
import ValidationResult from "../../validation/ValidationResult"
import TextArea from "antd/es/input/TextArea"
import { PlusOutlined } from "@ant-design/icons"
import Password from "antd/es/input/Password"
import DynamicField, { DynamicFieldType } from "../../dynamicfield/DynamicField"

export interface DynamicFieldProps {
    formValues: Record<string, any>
    validationResult: ValidationResult
    field: DynamicField
}

export default class DynamicInput extends React.Component<DynamicFieldProps> {
    private readonly qualifiedId: string

    constructor(props: DynamicFieldProps) {
        super(props)
        this.qualifiedId = `parameters.${props.field.id}`
    }

    private initialValue() {
        const { formValues, field } = this.props
        return formValues.parameters?.[field.id]
    }

    private evaluateConditions() {
        const { formValues, field } = this.props
        const { condition } = field
        if (condition === undefined || condition === null) return true

        const { parentField, value } = condition
        return formValues.parameters !== undefined && formValues.parameters[parentField] === value
    }

    private renderBoolean() {
        return <Switch value={this.initialValue()} />
    }

    private renderSingleLineText() {
        const {
            field: { sensitive },
        } = this.props
        if (sensitive) return <Password value={this.initialValue()} />
        else return <Input value={this.initialValue()} />
    }

    private renderMultiLineText() {
        return <TextArea rows={4} value={this.initialValue()} />
    }

    private renderEnum() {
        const { field } = this.props
        const options = field.enumOptions.map(option => ({
            value: option.id,
            label: option.description,
        }))

        return <Select value={this.initialValue()} options={options} showSearch />
    }

    private renderFileUpload() {
        return (
            <Upload type="drag" maxCount={1} beforeUpload={() => false}>
                <button style={{ border: 0, background: "none" }} type="button">
                    <PlusOutlined />
                    <div style={{ marginTop: 8 }}>Select file</div>
                </button>
            </Upload>
        )
    }

    render() {
        if (!this.evaluateConditions()) return undefined

        const { field, validationResult } = this.props
        let inputComponent: any

        switch (field.type) {
            case DynamicFieldType.BOOLEAN:
                inputComponent = this.renderBoolean()
                break
            case DynamicFieldType.EMAIL:
            case DynamicFieldType.SINGLE_LINE_TEXT:
            case DynamicFieldType.URL:
                inputComponent = this.renderSingleLineText()
                break
            case DynamicFieldType.MULTI_LINE_TEXT:
                inputComponent = this.renderMultiLineText()
                break
            case DynamicFieldType.FILE:
                inputComponent = this.renderFileUpload()
                break
            case DynamicFieldType.ENUM:
                inputComponent = this.renderEnum()
                break
        }

        return (
            <Form.Item
                name={["parameters", field.id]}
                validateStatus={validationResult.getStatus(this.qualifiedId)}
                help={validationResult.getMessage(this.qualifiedId) ?? field.helpText}
                label={field.description}
                required={field.required}
            >
                {inputComponent}
            </Form.Item>
        )
    }
}
