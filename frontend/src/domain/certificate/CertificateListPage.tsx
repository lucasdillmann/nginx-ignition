import React from "react";
import DataTable, {DataTableColumn} from "../../core/components/datatable/DataTable";
import {Flex} from "antd";
import {Link} from "react-router-dom";
import {DeleteOutlined, EyeOutlined, ReloadOutlined} from "@ant-design/icons";
import UserConfirmation from "../../core/components/confirmation/UserConfirmation";
import NotificationFacade from "../../core/components/notification/NotificationFacade";
import NginxReload from "../../core/components/nginx/NginxReload";
import PageResponse from "../../core/pagination/PageResponse";
import CertificateService from "./CertificateService";
import {CertificateResponse} from "./model/CertificateResponse";
import {UnexpectedResponseError} from "../../core/apiclient/ApiResponse";
import {RenewCertificateResponse} from "./model/RenewCertificateResponse";
import AvailableProviderResponse from "./model/AvailableProviderResponse";
import Preloader from "../../core/components/preloader/Preloader";
import TagGroup from "../../core/components/taggroup/TagGroup";

interface CertificateListPageState {
    loading: boolean
    providers: AvailableProviderResponse[]
}

export default class CertificateListPage extends React.Component<any, CertificateListPageState> {
    private readonly service: CertificateService
    private readonly table: React.RefObject<DataTable<CertificateResponse>>

    constructor(props: any) {
        super(props)
        this.service = new CertificateService()
        this.table = React.createRef()
        this.state = {
            loading: true,
            providers: []
        }
    }

    private translateProviderName(providerId: string): string {
        const {providers} = this.state
        return providers.find(provider => provider.uniqueId === providerId)?.name ?? providerId
    }

    private buildColumns(): DataTableColumn<CertificateResponse>[] {
        return [
            {
                id: "domainNames",
                description: "Domain names",
                renderer: (item) => <TagGroup values={item.domainNames} />,
            },
            {
                id: "provider",
                description: "Provider",
                renderer: (item) => this.translateProviderName(item.providerId),
                width: 250,
            },
            {
                id: "actions",
                description: "",
                renderer: (item) => (
                    <>
                        <Link to={`/certificates/${item.id}`}>
                            <EyeOutlined className="action-icon" />
                        </Link>
                        <Link to="" onClick={() => this.renewCertificate(item)}>
                            <ReloadOutlined className="action-icon" />
                        </Link>

                        <Link to="" onClick={() => this.deleteCertificate(item)}>
                            <DeleteOutlined className="action-icon" />
                        </Link>
                    </>
                ),
                fixed: true,
                width: 120,
            }
        ]

    }

    private async invokeCertificateRenew(certificate: CertificateResponse): Promise<void> {
        return this.service
            .renew(certificate.id)
            .then(() => NotificationFacade.success(
                `Certificate renewed`,
                `The certificate was renewed successfully`,
            ))
            .then(() => this.table.current?.refresh())
            .then(() => NginxReload.ask())
            .catch((error: UnexpectedResponseError<RenewCertificateResponse>) => NotificationFacade.error(
                `Unable to renew the certificate`,
                error.response.body?.errorReason ??
                `An unexpected error was found while trying to renew the certificate. Please try again later.`,
            ))
    }

    private renewCertificate(certificate: CertificateResponse) {
        return UserConfirmation.askWithCallback(
            `Do you really want to renew the certificate? Please be aware that this action can take several seconds.`,
            () => this.invokeCertificateRenew(certificate),
        )
    }

    private deleteCertificate(certificate: CertificateResponse) {
        UserConfirmation
            .ask("Do you really want to delete the certificate?")
            .then(() => this.service.delete(certificate.id))
            .then(() => NotificationFacade.success(
                `Certificate deleted`,
                `The certificate was deleted successfully`,
            ))
            .then(() => this.table.current?.refresh())
            .then(() => NginxReload.ask())
            .catch(() => NotificationFacade.error(
                `Unable to delete the certificate`,
                `An unexpected error was found while trying to delete the certificate. Please try again later.`,
            ))
    }

    private fetchData(pageSize: number, pageNumber: number): Promise<PageResponse<CertificateResponse>> {
        return this.service.list(pageSize, pageNumber)
    }

    componentDidMount() {
        this.service
            .availableProviders()
            .then(providers => this.setState({
                loading: false,
                providers,
            }))
    }

    render() {
        const {loading} = this.state
        if (loading)
            return (
                <Flex justify="center" align="center">
                    <Preloader loading />
                </Flex>
            )

        return (
            <DataTable
                ref={this.table}
                columns={this.buildColumns()}
                dataProvider={(pageSize, pageNumber) => this.fetchData(pageSize, pageNumber)}
                rowKey={item => item.id}
            />
        );
    }
}
