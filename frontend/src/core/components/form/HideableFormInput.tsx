import React from "react"
import FormItem from "antd/es/form/FormItem"
import If from "../flowcontrol/If"
import { I18n, I18nMessage } from "../../i18n/I18n"

export interface HideableFormInputProps {
    hidden: boolean
    reason: I18nMessage
    children: any
}

export default class HideableFormInput extends React.PureComponent<HideableFormInputProps> {
    render() {
        const { hidden, reason, children } = this.props
        const targetChildren = React.cloneElement(children, { hidden })

        return (
            <>
                <If condition={hidden}>
                    <FormItem {...children.props} hidden={false} help={undefined}>
                        <span style={{ color: "gray" }}>
                            <I18n id={reason} />
                        </span>
                    </FormItem>
                </If>

                {targetChildren}
            </>
        )
    }
}
