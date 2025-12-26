import HostResponse from "./model/HostResponse"
import HostFormValues, {
    HostFormBinding,
    HostFormGlobalBindingCertificateOverride,
    HostFormRoute,
    HostFormRouteIntegration,
    HostFormStaticResponse,
    HostFormVpn,
} from "./model/HostFormValues"
import HostRequest, {
    HostBinding,
    HostRoute,
    HostRouteIntegration,
    HostRouteStaticResponse,
    HostRouteType,
    HostVpn,
} from "./model/HostRequest"
import CertificateService from "../certificate/CertificateService"
import IntegrationService from "../integration/IntegrationService"
import AccessListService from "../accesslist/AccessListService"
import VpnService from "../vpn/VpnService"

class HostConverter {
    private readonly certificateService: CertificateService
    private readonly integrationService: IntegrationService
    private readonly accessListService: AccessListService
    private readonly vpnService: VpnService

    constructor() {
        this.certificateService = new CertificateService()
        this.integrationService = new IntegrationService()
        this.accessListService = new AccessListService()
        this.vpnService = new VpnService()
    }

    private notNull(value?: any) {
        return value !== undefined && value !== null
    }

    private staticResponseToFormValues(response: HostRouteStaticResponse): HostFormStaticResponse {
        const { statusCode, payload } = response
        const { headers } = response
        const joinedHeaders = this.notNull(headers)
            ? Object.entries(headers!)
                  .map(([key, value]) => `${key}: ${value}`)
                  .join("\n")
            : ""

        return {
            headers: joinedHeaders,
            statusCode,
            payload,
        }
    }

    private async integrationToFormValues(data: HostRouteIntegration): Promise<HostFormRouteIntegration> {
        const { integrationId, optionId } = data
        const integration = await this.integrationService.getById(integrationId)
        const option = await this.integrationService.getOptionById(integrationId, optionId)

        return {
            integration: integration!!,
            option: option!!,
        }
    }

    private async routeToFormValues(route: HostRoute): Promise<HostFormRoute> {
        const accessList = this.notNull(route.accessListId)
            ? await this.accessListService.getById(route.accessListId!!)
            : undefined
        const response = this.notNull(route.response) ? this.staticResponseToFormValues(route.response!!) : undefined
        const integration =
            this.notNull(route.integration) && route.type === HostRouteType.INTEGRATION
                ? await this.integrationToFormValues(route.integration!!)
                : undefined

        return {
            ...route,
            response,
            integration,
            accessList,
        }
    }

    private async bindingToFormValues(binding: HostBinding): Promise<HostFormBinding> {
        const certificate = this.notNull(binding.certificateId)
            ? await this.certificateService.getById(binding.certificateId!!)
            : undefined

        return {
            ...binding,
            certificate,
        }
    }

    private async vpnToFormValues(data: HostVpn): Promise<HostFormVpn> {
        const vpn = await this.vpnService.getById(data.vpnId!!)

        return {
            ...data,
            vpn: vpn!!,
        }
    }

    private formValuesToHeaders(headers: string): Record<string, string> {
        const lines = headers.split("\n").filter(line => line.trim().length > 0)
        const pairs = lines
            .map(line => line.split(":"))
            .filter(line => line.length >= 2)
            .map(([key, ...remaining]) => [key, remaining.join(":")])

        const output: Record<string, string> = {}
        for (const [key, value] of pairs) {
            output[key.trim()] = value.trim()
        }
        return output
    }

    private formValuesToStaticResponse(response: HostFormStaticResponse): HostRouteStaticResponse {
        const { statusCode, payload } = response
        const headers = this.notNull(response.headers) ? this.formValuesToHeaders(response.headers!!) : {}

        return {
            statusCode,
            payload,
            headers,
        }
    }

    private formValuesToIntegration(integration: HostFormRouteIntegration): HostRouteIntegration {
        return {
            integrationId: integration.integration.id,
            optionId: integration.option.id,
        }
    }

    private formValuesToRoute(route: HostFormRoute): HostRoute {
        const { priority, enabled, type, settings, targetUri, sourcePath, accessList, redirectCode, sourceCode } = route
        const response = this.notNull(route.response) ? this.formValuesToStaticResponse(route.response!!) : undefined
        const integration = this.notNull(route.integration)
            ? this.formValuesToIntegration(route.integration!!)
            : undefined

        return {
            priority,
            enabled,
            type,
            settings,
            targetUri,
            sourcePath,
            response,
            integration,
            redirectCode,
            sourceCode,
            accessListId: accessList?.id,
        }
    }

    private formValuesToBinding(binding: HostFormBinding): HostBinding {
        const output = {
            ...binding,
            certificateId: binding.certificate?.id,
        }

        delete output.certificate
        return output
    }

    private formValuesToVpn(vpn: HostFormVpn): HostVpn {
        const data = {
            vpnId: vpn.vpn?.id,
            name: vpn.name,
            host: vpn.host,
        }

        if (!vpn.host || vpn.host.trim() === "") {
            data.host = undefined
        }

        return data
    }

    private async certificateOverridesToFormValues(
        overrides: Record<string, string | null> | undefined,
        globalBindings: HostBinding[] | undefined,
    ): Promise<HostFormGlobalBindingCertificateOverride[]> {
        if (!globalBindings || globalBindings.length === 0) {
            return []
        }

        const result: HostFormGlobalBindingCertificateOverride[] = []

        for (const binding of globalBindings) {
            if (!binding.id) continue

            const certificateId = overrides?.[binding.id] ?? null
            const certificate = certificateId ? await this.certificateService.getById(certificateId) : undefined

            result.push({
                bindingId: binding.id,
                certificate,
            })
        }

        return result
    }

    private formValuesToCertificateOverrides(
        overrides: HostFormGlobalBindingCertificateOverride[],
    ): Record<string, string | null> | undefined {
        if (!overrides || overrides.length === 0) {
            return undefined
        }

        const result: Record<string, string | null> = {}

        for (const override of overrides) {
            if (!override?.bindingId) continue

            const certificateId = override.certificate?.id ?? null
            result[override.bindingId] = certificateId
        }

        return Object.keys(result).length > 0 ? result : undefined
    }

    async responseToFormValues(response: HostResponse): Promise<HostFormValues> {
        const {
            enabled,
            domainNames,
            featureSet,
            defaultServer,
            useGlobalBindings,
            accessListId,
            globalBindingCertificateOverrides,
            globalBindings,
        } = response

        const routes = response.routes.map(route => this.routeToFormValues(route))
        const responseBindings = response.bindings ?? []
        const bindings = await Promise.all(responseBindings.map(binding => this.bindingToFormValues(binding)))
        const vpns = await Promise.all(response.vpns.map(vpn => this.vpnToFormValues(vpn)))
        const accessList = this.notNull(accessListId) ? await this.accessListService.getById(accessListId!!) : undefined

        const certificateOverrides = await this.certificateOverridesToFormValues(
            globalBindingCertificateOverrides,
            globalBindings,
        )

        return {
            enabled,
            bindings,
            vpns,
            featureSet,
            defaultServer,
            useGlobalBindings,
            globalBindingCertificateOverrides: certificateOverrides,
            accessList,
            domainNames: domainNames ?? [""],
            routes: await Promise.all(routes),
        }
    }

    formValuesToRequest(formValues: HostFormValues): HostRequest {
        const {
            enabled,
            domainNames,
            featureSet,
            defaultServer,
            useGlobalBindings,
            accessList,
            globalBindingCertificateOverrides,
        } = formValues

        const routes = formValues.routes.map(route => this.formValuesToRoute(route))
        const bindings = useGlobalBindings ? [] : formValues.bindings.map(binding => this.formValuesToBinding(binding))
        const vpns = formValues.vpns?.map(vpn => this.formValuesToVpn(vpn)) ?? []

        const overrides = this.formValuesToCertificateOverrides(globalBindingCertificateOverrides)

        return {
            enabled,
            featureSet,
            defaultServer,
            useGlobalBindings,
            globalBindingCertificateOverrides: overrides,
            bindings,
            routes,
            vpns,
            accessListId: accessList?.id,
            domainNames: defaultServer ? [] : domainNames,
        }
    }
}

export default new HostConverter()
