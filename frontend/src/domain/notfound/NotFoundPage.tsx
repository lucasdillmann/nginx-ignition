import React from "react"
import FullPageError from "../../core/components/error/FullPageError"
import AppContext from "../../core/components/context/AppContext"
import { Navigate } from "react-router-dom"
import { QuestionCircleFilled } from "@ant-design/icons"
import { buildLoginUrl } from "../../core/authentication/buildLoginUrl"
import MessageKey from "../../core/i18n/model/MessageKey.generated"

export default class NotFoundPage extends React.PureComponent {
    render() {
        const currentUser = AppContext.get().user
        if (currentUser === undefined) return <Navigate to={buildLoginUrl()} />

        return (
            <FullPageError
                title={MessageKey.GlobalErrorNotFoundTitle}
                message={MessageKey.GlobalErrorNotFoundDescription}
                icon={
                    <QuestionCircleFilled style={{ fontSize: 48, color: "var(--nginxIgnition-colorTextSecondary)" }} />
                }
            />
        )
    }
}
