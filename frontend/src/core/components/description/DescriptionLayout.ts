import {ProDescriptionsProps} from "@ant-design/pro-components";
import FormLayout from "../form/FormLayout";

const Defaults: ProDescriptionsProps = {
    colon: false,
    column: 2,
    layout: "horizontal",
    size: "middle",
    labelStyle: {
        width: 250,
    },
    formProps: FormLayout.FormDefaults,
}

export default Object.freeze({
    Defaults,
})
