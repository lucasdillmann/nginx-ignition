import React from "react"
import { Empty, Flex } from "antd"
import AppShellContext from "../../core/components/shell/AppShellContext"
import Preloader from "../../core/components/preloader/Preloader"
import EmptyStates from "../../core/components/emptystate/EmptyStates"
import { isAccessGranted } from "../../core/components/accesscontrol/IsAccessGranted"
import { UserAccessLevel } from "../user/model/UserAccessLevel"
import MessageKey from "../../core/i18n/model/MessageKey.generated"
import { I18n, I18nMessage } from "../../core/i18n/I18n"
import NginxService from "../nginx/NginxService"
import NginxMetadata, { NginxSupportType } from "../nginx/model/NginxMetadata"
import HostService from "../host/HostService"
import StreamService from "../stream/StreamService"
import CertificateService from "../certificate/CertificateService"
import { CertificateResponse } from "../certificate/model/CertificateResponse"
import SettingsService from "../settings/SettingsService"
import SettingsDto from "../settings/model/SettingsDto"
import TrafficStatsService from "../trafficstats/TrafficStatsService"
import TrafficStatsResponse, { ZoneData } from "../trafficstats/model/TrafficStatsResponse"
import ZoneStatCards from "../trafficstats/components/ZoneStatCards"
import LogViewer from "../logs/components/LogViewer"
import LogLine from "../logs/model/LogLine"
import { Link } from "react-router-dom"
import TagGroup from "../../core/components/taggroup/TagGroup"
import CountCard from "./components/CountCard"
import HomeHeader from "./components/HomeHeader"
import NginxStatusCard from "./components/NginxStatusCard"
import "./HomePage.css"
import "../trafficstats/TrafficStatsPage.css"
import "../logs/components/LogViewer.css"

interface HomePageState {
    loading: boolean
    refreshToken: number
    metadata?: NginxMetadata
    nginxRunning?: boolean
    settings?: SettingsDto
    hostCount?: number
    streamCount?: number
    certificateCount?: number
    expiringCertificates: CertificateResponse[]
    errorLogs: LogLine[]
    stats?: TrafficStatsResponse
    error?: Error
}

export default class HomePage extends React.Component<object, HomePageState> {
    private readonly nginxService: NginxService
    private readonly hostService: HostService
    private readonly streamService: StreamService
    private readonly certificateService: CertificateService
    private readonly settingsService: SettingsService
    private readonly trafficStatsService: TrafficStatsService

    constructor(props: object) {
        super(props)
        this.nginxService = new NginxService()
        this.hostService = new HostService()
        this.streamService = new StreamService()
        this.certificateService = new CertificateService()
        this.settingsService = new SettingsService()
        this.trafficStatsService = new TrafficStatsService()
        this.state = {
            loading: true,
            refreshToken: 0,
            expiringCertificates: [],
            errorLogs: [],
        }
    }

    componentDidMount() {
        this.configureShell()
        this.fetchData()
    }

    private canViewNginxServer(): boolean {
        return isAccessGranted(UserAccessLevel.READ_ONLY, permissions => permissions.nginxServer)
    }

    private canViewHosts(): boolean {
        return isAccessGranted(UserAccessLevel.READ_ONLY, permissions => permissions.hosts)
    }

    private canViewStreams(): boolean {
        return isAccessGranted(UserAccessLevel.READ_ONLY, permissions => permissions.streams)
    }

    private canViewCertificates(): boolean {
        return isAccessGranted(UserAccessLevel.READ_ONLY, permissions => permissions.certificates)
    }

    private canViewLogs(): boolean {
        return isAccessGranted(UserAccessLevel.READ_ONLY, permissions => permissions.logs)
    }

    private canViewTrafficStats(): boolean {
        return isAccessGranted(UserAccessLevel.READ_ONLY, permissions => permissions.trafficStats)
    }

    private configureShell() {
        AppShellContext.get().updateConfig({
            noContainerPadding: true,
        })
    }

    private refreshData() {
        const { loading } = this.state
        if (loading) return

        this.setState(
            state => ({ error: undefined, refreshToken: state.refreshToken + 1 }),
            () => this.fetchData(),
        )
    }

    private filterExpiringCertificates(certificates: CertificateResponse[]): CertificateResponse[] {
        const now = new Date()
        const windowEnd = new Date(now.getTime() + 30 * 24 * 60 * 60 * 1000)

        return certificates
            .filter(certificate => {
                const validUntil = new Date(certificate.validUntil)
                return validUntil >= now && validUntil <= windowEnd
            })
            .sort((left, right) => new Date(left.validUntil).getTime() - new Date(right.validUntil).getTime())
    }

    private daysUntilExpiry(validUntil: string): number {
        const now = new Date()
        now.setHours(0, 0, 0, 0)
        const expiry = new Date(validUntil)
        expiry.setHours(0, 0, 0, 0)
        return Math.ceil((expiry.getTime() - now.getTime()) / (24 * 60 * 60 * 1000))
    }

    private overviewClassName(countCardsLength: number, canViewNginx: boolean): string {
        let className = `home-dashboard-overview-row home-dashboard-overview-row-count-${countCardsLength}`

        if (canViewNginx) className += " home-dashboard-overview-row-with-nginx"

        return className
    }

    private async fetchData() {
        try {
            const canNginxServer = this.canViewNginxServer()
            const canHosts = this.canViewHosts()
            const canStreams = this.canViewStreams()
            const canCertificates = this.canViewCertificates()
            const canLogs = this.canViewLogs()
            const canTrafficStats = this.canViewTrafficStats()
            const needsMetadata = canNginxServer || canTrafficStats
            const needsSettings = canLogs || canTrafficStats

            const [metadata, nginxRunning, settings, hostsPage, streamsPage, certificates] = await Promise.all([
                needsMetadata ? this.nginxService.getMetadata() : Promise.resolve(undefined),
                needsMetadata ? this.nginxService.isRunning() : Promise.resolve(undefined),
                needsSettings ? this.settingsService.get() : Promise.resolve(undefined),
                canHosts ? this.hostService.list(1, 0) : Promise.resolve(undefined),
                canStreams ? this.streamService.list(1, 0) : Promise.resolve(undefined),
                canCertificates ? this.certificateService.listAll() : Promise.resolve([]),
            ])

            let errorLogs: LogLine[] = []
            if (canLogs && settings?.nginx.logs.serverLogsEnabled) {
                errorLogs = await this.nginxService.logs(15, 0)
            }

            let stats: TrafficStatsResponse | undefined
            const statsSupported = metadata?.availableSupport.stats !== NginxSupportType.NONE
            const statsEnabled = metadata?.stats.enabled === true
            if (canTrafficStats && statsSupported && statsEnabled && nginxRunning) {
                stats = await this.trafficStatsService.getStats()
            }

            this.setState({
                loading: false,
                error: undefined,
                metadata,
                nginxRunning,
                settings,
                hostCount: hostsPage?.totalItems,
                streamCount: streamsPage?.totalItems,
                certificateCount: canCertificates ? certificates.length : undefined,
                expiringCertificates: canCertificates ? this.filterExpiringCertificates(certificates) : [],
                errorLogs,
                stats,
            })
        } catch (error) {
            this.setState({ loading: false, error: error as Error })
        }
    }

    private renderOverviewSection() {
        const { hostCount, streamCount, certificateCount } = this.state
        const countCards: React.ReactNode[] = []

        if (this.canViewHosts() && hostCount !== undefined) {
            countCards.push(<CountCard key="hosts" title={MessageKey.CommonHosts} count={hostCount} linkTo="/hosts" />)
        }

        if (this.canViewStreams() && streamCount !== undefined) {
            countCards.push(
                <CountCard key="streams" title={MessageKey.CommonStreams} count={streamCount} linkTo="/streams" />,
            )
        }

        if (this.canViewCertificates() && certificateCount !== undefined) {
            countCards.push(
                <CountCard
                    key="certificates"
                    title={MessageKey.CommonSslCertificates}
                    count={certificateCount}
                    linkTo="/certificates"
                />,
            )
        }

        const canViewNginx = this.canViewNginxServer()
        if (countCards.length === 0 && !canViewNginx) return null

        return (
            <Flex className={this.overviewClassName(countCards.length, canViewNginx)} align="stretch" wrap="wrap">
                {this.renderOverviewTotalsSection(countCards)}
                {this.renderOverviewNginxSection(canViewNginx)}
            </Flex>
        )
    }

    private renderOverviewTotalsSection(countCards: React.ReactNode[]) {
        if (countCards.length === 0) return null

        return (
            <div className="home-dashboard-section home-dashboard-totals-section">
                <h3 className="home-dashboard-section-title">
                    <I18n id={MessageKey.FrontendHomeTotalsTitle} />
                </h3>
                <Flex className="traffic-stats-cards-row">{countCards}</Flex>
            </div>
        )
    }

    private renderOverviewNginxSection(canViewNginx: boolean) {
        if (!canViewNginx) return null

        const { refreshToken } = this.state

        return (
            <div className="home-dashboard-section home-dashboard-nginx-section">
                <h3 className="home-dashboard-section-title">
                    <I18n id={MessageKey.FrontendHomeNginxSectionTitle} />
                </h3>
                <div className="home-dashboard-nginx-slot">
                    <NginxStatusCard refreshToken={refreshToken} />
                </div>
            </div>
        )
    }

    private renderDashboardEmpty(message: I18nMessage) {
        return <Empty className="home-dashboard-empty" description={<I18n id={message} />} />
    }

    private renderTrafficEmptyState() {
        const { metadata, nginxRunning, stats } = this.state

        if (metadata?.availableSupport?.stats === NginxSupportType.NONE) {
            return this.renderDashboardEmpty(MessageKey.FrontendHomeTrafficUnsupported)
        }

        if (metadata && !metadata.stats.enabled) {
            return this.renderDashboardEmpty(MessageKey.FrontendHomeTrafficDisabled)
        }

        if (nginxRunning === false) {
            return this.renderDashboardEmpty(MessageKey.FrontendHomeTrafficOffline)
        }

        if (stats?.serverZones?.["*"]?.requestCounter === 0) {
            return this.renderDashboardEmpty(MessageKey.FrontendHomeTrafficNoData)
        }

        return null
    }

    private renderViewAllLink(to: string) {
        return (
            <Link to={to} className="home-dashboard-view-all-link">
                <I18n id={MessageKey.FrontendHomeViewAll} />
            </Link>
        )
    }

    private renderTrafficSectionHeader(emptyState: React.ReactNode | null) {
        return (
            <Flex className="home-dashboard-section-header">
                <h3 className="home-dashboard-section-title">
                    <I18n id={MessageKey.CommonTrafficStats} />
                </h3>
                {emptyState === null && this.renderViewAllLink("/traffic-stats")}
            </Flex>
        )
    }

    private renderTrafficSectionContent(emptyState: React.ReactNode | null, globalZone: ZoneData | undefined) {
        if (emptyState !== null) {
            return <div className="home-dashboard-panel home-dashboard-traffic-panel">{emptyState}</div>
        }

        if (globalZone === undefined) return null

        return (
            <ZoneStatCards
                requests={globalZone.requestCounter}
                inBytes={globalZone.inBytes}
                outBytes={globalZone.outBytes}
                avgResponseTime={globalZone.requestMsec}
            />
        )
    }

    private renderTrafficSection() {
        if (!this.canViewTrafficStats()) return null

        const { stats } = this.state
        const globalZone = stats?.serverZones?.["*"]
        const emptyState = this.renderTrafficEmptyState()

        return (
            <div className="home-dashboard-section">
                {this.renderTrafficSectionHeader(emptyState)}
                {this.renderTrafficSectionContent(emptyState, globalZone)}
            </div>
        )
    }

    private renderCertificateExpiryMessage(days: number) {
        if (days <= 0) return <I18n id={MessageKey.FrontendHomeExpiresToday} />

        return <I18n id={MessageKey.FrontendHomeExpiresInDays} params={{ days }} />
    }

    private renderExpiringCertificateItem(certificate: CertificateResponse) {
        const days = this.daysUntilExpiry(certificate.validUntil)
        const expiryMessage = this.renderCertificateExpiryMessage(days)

        return (
            <Flex key={certificate.id} className="home-dashboard-cert-item" align="center">
                <Link to={`/certificates/${certificate.id}`} className="home-dashboard-cert-domains">
                    <TagGroup values={certificate.domainNames} maximumSize={1} />
                </Link>
                <span className="home-dashboard-cert-expiry">{expiryMessage}</span>
            </Flex>
        )
    }

    private renderExpiringCertificatesContent() {
        const { expiringCertificates } = this.state

        if (expiringCertificates.length === 0) {
            return this.renderDashboardEmpty(MessageKey.FrontendHomeNoCertificatesExpiring)
        }

        return (
            <Flex className="home-dashboard-cert-list" vertical>
                {expiringCertificates.map(certificate => this.renderExpiringCertificateItem(certificate))}
            </Flex>
        )
    }

    private renderExpiringCertificatesPanel() {
        if (!this.canViewCertificates()) return null

        return (
            <div className="home-dashboard-details-column">
                <Flex className="home-dashboard-section-header">
                    <h3 className="home-dashboard-section-title">
                        <I18n id={MessageKey.FrontendHomeCertificatesExpiringTitle} />
                    </h3>
                    {this.renderViewAllLink("/certificates")}
                </Flex>
                <div className="home-dashboard-panel">{this.renderExpiringCertificatesContent()}</div>
            </div>
        )
    }

    private renderRecentErrorsContent() {
        const { settings, errorLogs } = this.state
        const serverLogsEnabled = settings?.nginx.logs.serverLogsEnabled === true

        if (!serverLogsEnabled) {
            return this.renderDashboardEmpty(MessageKey.FrontendHomeRecentErrorsDisabled)
        }

        if (errorLogs.length === 0) {
            return this.renderDashboardEmpty(MessageKey.FrontendHomeRecentErrorsEmpty)
        }

        return <LogViewer lines={errorLogs} />
    }

    private renderRecentErrorsBody(content: React.ReactNode) {
        const serverLogsEnabled = this.state.settings?.nginx.logs.serverLogsEnabled === true

        if (serverLogsEnabled) {
            return <div className="home-dashboard-log-content">{content}</div>
        }

        return <div className="home-dashboard-panel">{content}</div>
    }

    private renderRecentErrorsPanel() {
        if (!this.canViewLogs()) return null

        const content = this.renderRecentErrorsContent()

        return (
            <div className="home-dashboard-details-column">
                <Flex className="home-dashboard-section-header">
                    <h3 className="home-dashboard-section-title">
                        <I18n id={MessageKey.FrontendHomeRecentErrorsTitle} />
                    </h3>
                    {this.renderViewAllLink("/logs")}
                </Flex>
                {this.renderRecentErrorsBody(content)}
            </div>
        )
    }

    private renderDetailsSection() {
        const expiringPanel = this.renderExpiringCertificatesPanel()
        const errorsPanel = this.renderRecentErrorsPanel()

        if (expiringPanel === null && errorsPanel === null) return null

        return (
            <Flex className="home-dashboard-split-row">
                {expiringPanel}
                {errorsPanel}
            </Flex>
        )
    }

    render() {
        const { loading, error } = this.state

        if (error !== undefined) return EmptyStates.FailedToFetch

        if (loading) return <Preloader loading />

        const { metadata } = this.state

        return (
            <>
                <HomeHeader metadata={metadata} onRefresh={() => this.refreshData()} />
                <Flex className="home-dashboard-container" vertical>
                    {this.renderOverviewSection()}
                    {this.renderTrafficSection()}
                    {this.renderDetailsSection()}
                </Flex>
            </>
        )
    }
}
