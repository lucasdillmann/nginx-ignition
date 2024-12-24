import React from "react"
import DataTable, { DataTableColumn } from "../../core/components/datatable/DataTable"
import { Flex } from "antd"
import { Link } from "react-router-dom"
import { DeleteOutlined, EyeOutlined, ReloadOutlined } from "@ant-design/icons"
import PageResponse from "../../core/pagination/PageResponse"
import CertificateService from "./CertificateService"
import { CertificateResponse } from "./model/CertificateResponse"
import AvailableProviderResponse from "./model/AvailableProviderResponse"
import Preloader from "../../core/components/preloader/Preloader"
import TagGroup from "../../core/components/taggroup/TagGroup"
import RenewCertificateAction from "./actions/RenewCertificateAction"
import DeleteCertificateAction from "./actions/DeleteCertificateAction"
import AppShellContext from "../../core/components/shell/AppShellContext"
import CommonNotifications from "../../core/components/notification/CommonNotifications"
import EmptyStates from "../../core/components/emptystate/EmptyStates"

interface CertificateListPageState {
    loading: boolean
    providers: AvailableProviderResponse[]
    error?: Error
}

export default class CertificateListPage extends React.Component<any, CertificateListPageState> {
    private readonly service: CertificateService
    private readonly table: React.RefObject<DataTable<CertificateResponse> | null>

    constructor(props: any) {
        super(props)
        this.service = new CertificateService()
        this.table = React.createRef()
        this.state = {
            loading: true,
            providers: [],
        }
    }

    private translateProviderName(providerId: string): string {
        const { providers } = this.state
        return providers.find(provider => provider.id === providerId)?.name ?? providerId
    }

    private buildColumns(): DataTableColumn<CertificateResponse>[] {
        return [
            {
                id: "domainNames",
                description: "Domain names",
                renderer: item => <TagGroup values={item.domainNames} />,
            },
            {
                id: "provider",
                description: "Provider",
                renderer: item => this.translateProviderName(item.providerId),
                width: 250,
            },
            {
                id: "actions",
                description: "",
                renderer: item => (
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
                width: 120,
            },
        ]
    }

    private async renewCertificate(certificate: CertificateResponse) {
        return RenewCertificateAction.execute(certificate.id).then(() => this.table.current?.refresh())
    }

    private async deleteCertificate(certificate: CertificateResponse) {
        return DeleteCertificateAction.execute(certificate.id).then(() => this.table.current?.refresh())
    }

    private fetchData(
        pageSize: number,
        pageNumber: number,
        searchTerms?: string,
    ): Promise<PageResponse<CertificateResponse>> {
        return this.service.list(pageSize, pageNumber, searchTerms)
    }

    componentDidMount() {
        this.service
            .availableProviders()
            .then(providers =>
                this.setState({
                    loading: false,
                    providers,
                }),
            )
            .catch(error => {
                CommonNotifications.failedToFetch()
                this.setState({ loading: false, error })
            })

        AppShellContext.get().updateConfig({
            title: "SSL certificates",
            subtitle: "Relation of issued SSL certificates for use in the nginx's virtual hosts",
            actions: [
                {
                    description: "Issue certificate",
                    onClick: "/certificates/new",
                },
            ],
        })
    }

    render() {
        const { loading, error } = this.state
        if (loading)
            return (
                <Flex justify="center" align="center">
                    <Preloader loading />
                </Flex>
            )

        if (error !== undefined) return EmptyStates.FailedToFetch

        return (
            <DataTable
                ref={this.table}
                columns={this.buildColumns()}
                dataProvider={(pageSize, pageNumber, searchTerms) => this.fetchData(pageSize, pageNumber, searchTerms)}
                rowKey={item => item.id}
            />
        )
    }
}
