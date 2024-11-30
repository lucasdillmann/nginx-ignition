import { FormItemProps, FormProps } from "antd"

const LabeledItem: FormItemProps = {
    labelCol: {
        xs: { span: 24 },
        sm: { span: 5 },
    },
    wrapperCol: {
        xs: { span: 24 },
        sm: { span: 20 },
    },
}

const ExpandedLabeledItem: FormItemProps = {
    labelCol: {
        flex: "auto",
    },
    wrapperCol: {
        flex: "auto",
    },
}

const FormDefaults: FormProps = {
    labelCol: LabeledItem.labelCol,
    wrapperCol: LabeledItem.wrapperCol,
    layout: "horizontal",
    requiredMark: "optional",
    colon: false,
    labelAlign: "left",
    labelWrap: true,
}

const UnlabeledItem: FormItemProps = {
    wrapperCol: {
        xs: { span: 24, offset: 0 },
        sm: { span: 20, offset: 5 },
    },
}

const ExpandedUnlabeledItem: FormItemProps = {
    wrapperCol: {
        flex: "auto",
    },
}

export default Object.freeze({
    FormDefaults,
    LabeledItem,
    ExpandedLabeledItem,
    UnlabeledItem,
    ExpandedUnlabeledItem,
})
