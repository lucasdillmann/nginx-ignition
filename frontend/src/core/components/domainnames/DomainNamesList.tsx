import React from "react"
import { Button, Flex, Form, FormListFieldData, FormListOperation, Input } from "antd"
import FormLayout from "../form/FormLayout"
import If from "../flowcontrol/If"
import { PlusOutlined, DeleteOutlined } from "@ant-design/icons"
import ValidationResult from "../../validation/ValidationResult"
import { i18n, I18n } from "../../i18n/I18n"
import MessageKey from "../../i18n/model/MessageKey.generated"

export interface DomainNamesListProps {
    pathPrefix?: { merged: string; name: any[] }
    validationResult: ValidationResult
    expandedLabelSize?: boolean
    className?: string
    disableTitle?: boolean
}

export default class DomainNamesList extends React.PureComponent<DomainNamesListProps> {
    private renderFields(fields: FormListFieldData[], operations: FormListOperation) {
        const { validationResult, expandedLabelSize, className, disableTitle } = this.props
        const layout = expandedLabelSize === true ? FormLayout.ExpandedUnlabeledItem : FormLayout.UnlabeledItem
        const pathPrefix = this.props.pathPrefix === undefined ? "" : this.props.pathPrefix.merged + "."

        const domainNameFields = fields.map((field, index) => (
            <Form.Item
                {...(index > 0 && !disableTitle ? layout : undefined)}
                label={index === 0 && !disableTitle ? <I18n id={MessageKey.CommonDomainNames} /> : ""}
                key={field.key}
                className={className}
                required
            >
                <Flex>
                    <Form.Item
                        {...field}
                        validateStatus={validationResult.getStatus(`${pathPrefix}domainNames[${index}]`)}
                        help={validationResult.getMessage(`${pathPrefix}domainNames[${index}]`)}
                        style={{ marginBottom: 0, width: "100%" }}
                    >
                        <Input placeholder={i18n(MessageKey.FrontendComponentsDomainnamesPlaceholder)} />
                    </Form.Item>
                    <If condition={fields.length > 1}>
                        <DeleteOutlined
                            style={{ marginLeft: 15, alignItems: "start", marginTop: 7 }}
                            onClick={() => operations.remove(field.name)}
                        />
                    </If>
                </Flex>
            </Form.Item>
        ))

        const addAction = (
            <Form.Item {...layout} key="add-domain" className={className}>
                <Button type="dashed" onClick={() => operations.add()} icon={<PlusOutlined />}>
                    <I18n id={MessageKey.FrontendHostRouteAddDomain} />
                </Button>
            </Form.Item>
        )

        return [...domainNameFields, addAction]
    }

    render() {
        const { pathPrefix } = this.props
        const path = pathPrefix === undefined ? [] : pathPrefix.name

        return (
            <Form.List name={[...path, "domainNames"]}>
                {(fields, operations) => this.renderFields(fields, operations)}
            </Form.List>
        )
    }
}
