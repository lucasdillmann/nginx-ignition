import { ProDescriptionsProps } from "@ant-design/pro-components"
import FormLayout from "../form/FormLayout"
import "./DescriptionLayout.css"

const Defaults: ProDescriptionsProps = {
    colon: false,
    column: 2,
    layout: "horizontal",
    size: "medium",
    classNames: {
        label: "description-layout-label",
    },
    formProps: FormLayout.FormDefaults,
}

export default Object.freeze({
    Defaults,
})
