import React from "react"
import AppShellContext from "../../../core/components/shell/AppShellContext"
import { isAccessGranted } from "../../../core/components/accesscontrol/IsAccessGranted"
import { UserAccessLevel } from "../../user/model/UserAccessLevel"
import AccessDeniedPage from "../../../core/components/accesscontrol/AccessDeniedPage"

export default class ClientCertificateListPage extends React.Component {
    componentDidMount() {
        AppShellContext.get().updateConfig({
            title: "Client certificates",
            subtitle: "Relation of client certificates for use in the nginx's mutual authentication (mTLS)",
            actions: [],
        })
    }

    render() {
        if (!isAccessGranted(UserAccessLevel.READ_ONLY, permissions => permissions.clientCertificates)) {
            return <AccessDeniedPage />
        }

        return <></>
    }
}
