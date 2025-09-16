import React, { ReactNode } from "react"
import AccessDeniedPage from "../../core/components/accesscontrol/AccessDeniedPage"
import { UserAccessLevel } from "../user/model/UserAccessLevel"
import { isAccessGranted } from "../../core/components/accesscontrol/IsAccessGranted"
import { Button } from "antd"
import ExportService from "./ExportService"

export default class ExportPage extends React.Component<any, any> {
    private readonly service: ExportService

    constructor(props: any) {
        super(props)
        this.service = new ExportService()
    }

    private renderPage(): ReactNode {
        return <Button onClick={() => this.service.downloadNginxConfigurationFiles()}>Download</Button>
    }

    render(): ReactNode {
        if (!isAccessGranted(UserAccessLevel.READ_ONLY, permissions => permissions.exportData)) {
            return <AccessDeniedPage />
        }

        return this.renderPage()
    }
}
