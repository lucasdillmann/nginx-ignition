import React from "react";
import DataTable, {DataTableColumn} from "../../core/components/datatable/DataTable";
import {Flex} from "antd";
import {Link} from "react-router-dom";
import {DeleteOutlined, EyeOutlined, ReloadOutlined} from "@ant-design/icons";
import UserConfirmation from "../../core/components/confirmation/UserConfirmation";
import Notification from "../../core/components/notification/Notification";
import NginxReload from "../../core/components/nginx/NginxReload";
import PageResponse from "../../core/pagination/PageResponse";
import CertificateService from "./CertificateService";
import {CertificateResponse} from "./model/CertificateResponse";
import {UnexpectedResponseError} from "../../core/apiclient/ApiResponse";
import {RenewCertificateResponse} from "./model/RenewCertificateResponse";
import AvailableProviderResponse from "./model/AvailableProviderResponse";
import Preloader from "../../core/components/preloader/Preloader";
import TagGroup from "../../core/components/taggroup/TagGroup";
import ShellAwareComponent, {ShellConfig} from "../../core/components/shell/ShellAwareComponent";

interface CertificateListPageState {
    loading: boolean
    providers: AvailableProviderResponse[]
}

export default class CertificateListPage extends ShellAwareComponent<any, CertificateListPageState> {
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
        return providers.find(provider => provider.id === providerId)?.name ?? providerId
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
            .then(() => Notification.success(
                `Certificate renewed`,
                `The certificate was renewed successfully`,
            ))
            .then(() => this.table.current?.refresh())
            .then(() => NginxReload.ask())
            .catch((error: UnexpectedResponseError<RenewCertificateResponse>) => Notification.error(
                `Unable to renew the certificate`,
                error.response.body?.errorReason ??
                `An unexpected error was found while trying to renew the certificate. Please try again later.`,
            ))
    }

    private renewCertificate(certificate: CertificateResponse) {
        return UserConfirmation.askWithCallback(
            `Renewing the certificate can take several seconds and is only recommended when something is wrong with it
            since, by default, nginx ignition will renew it automatically when it's close to expiring. Continue anyway?`,
            () => this.invokeCertificateRenew(certificate),
        )
    }

    private deleteCertificate(certificate: CertificateResponse) {
        UserConfirmation
            .ask("Do you really want to delete the certificate?")
            .then(() => this.service.delete(certificate.id))
            .then(() => Notification.success(
                `Certificate deleted`,
                `The certificate was deleted successfully`,
            ))
            .then(() => this.table.current?.refresh())
            .then(() => NginxReload.ask())
            .catch(() => Notification.error(
                `Unable to delete the certificate`,
                `An unexpected error was found while trying to delete the certificate. Please try again later.`,
            ))
    }

    private fetchData(pageSize: number, pageNumber: number): Promise<PageResponse<CertificateResponse>> {
        return this.service.list(pageSize, pageNumber)
    }

    shellConfig(): ShellConfig {
        return {
            title: "SSL certificates",
            subtitle: "Relation of issued SSL certificates for use in the nginx's virtual hosts",
            actions: [
                {
                    description: "Issue certificate",
                    onClick: "/certificates/new",
                },
            ],
        }
    }

    componentDidMount() {
        this.service
            .availableProviders()
            .then(providers => this.setState({
                loading: false,
                providers,
            }))
            .catch(() => Notification.error(
                "Unable to fetch the data",
                "We're unable to fetch the data at this moment. Please try again later.",
            ))
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
