import HostResponse from "./model/HostResponse"
import HostFormValues, {
    HostFormBinding,
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
import CacheService from "../cache/CacheService"

class HostConverter {
    private readonly certificateService: CertificateService
    private readonly integrationService: IntegrationService
    private readonly accessListService: AccessListService
    private readonly cacheService: CacheService
    private readonly vpnService: VpnService

    constructor() {
        this.certificateService = new CertificateService()
        this.integrationService = new IntegrationService()
        this.accessListService = new AccessListService()
        this.cacheService = new CacheService()
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
        const accessListPromise = this.notNull(route.accessListId)
            ? this.accessListService.getById(route.accessListId!!)
            : Promise.resolve(undefined)
        const cachePromise = this.notNull(route.cacheId)
            ? this.cacheService.getById(route.cacheId!!)
            : Promise.resolve(undefined)
        const integrationPromise =
            this.notNull(route.integration) && route.type === HostRouteType.INTEGRATION
                ? this.integrationToFormValues(route.integration!!)
                : Promise.resolve(undefined)

        const response = this.notNull(route.response) ? this.staticResponseToFormValues(route.response!!) : undefined

        const [accessList, cache, integration] = await Promise.all([
            accessListPromise,
            cachePromise,
            integrationPromise,
        ])
        return {
            ...route,
            response,
            integration,
            accessList,
            cache,
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
        const {
            priority,
            enabled,
            type,
            settings,
            targetUri,
            sourcePath,
            accessList,
            redirectCode,
            sourceCode,
            cache,
        } = route
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
            cacheId: cache?.id,
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

    async responseToFormValues(response: HostResponse): Promise<HostFormValues> {
        const { enabled, domainNames, featureSet, defaultServer, useGlobalBindings, accessListId, cacheId } = response

        const routes = response.routes.map(route => this.routeToFormValues(route))
        const responseBindings = response.bindings ?? []
        const bindingsPromise = Promise.all(responseBindings.map(binding => this.bindingToFormValues(binding)))
        const vpnsPromise = Promise.all(response.vpns.map(vpn => this.vpnToFormValues(vpn)))
        const accessListPromise = this.notNull(accessListId)
            ? this.accessListService.getById(accessListId!!)
            : Promise.resolve(undefined)
        const cachePromise = this.notNull(cacheId) ? this.cacheService.getById(cacheId!!) : Promise.resolve(undefined)

        const [bindings, vpns, accessList, cache] = await Promise.all([
            bindingsPromise,
            vpnsPromise,
            accessListPromise,
            cachePromise,
        ])

        return {
            enabled,
            bindings,
            vpns,
            featureSet,
            defaultServer,
            useGlobalBindings,
            accessList,
            cache,
            domainNames: domainNames ?? [""],
            routes: await Promise.all(routes),
        }
    }

    formValuesToRequest(formValues: HostFormValues): HostRequest {
        const { enabled, domainNames, featureSet, defaultServer, useGlobalBindings, accessList, cache } = formValues

        const routes = formValues.routes.map(route => this.formValuesToRoute(route))
        const bindings = useGlobalBindings ? [] : formValues.bindings.map(binding => this.formValuesToBinding(binding))
        const vpns = formValues.vpns?.map(vpn => this.formValuesToVpn(vpn)) ?? []

        return {
            enabled,
            featureSet,
            defaultServer,
            useGlobalBindings,
            bindings,
            routes,
            vpns,
            accessListId: accessList?.id,
            cacheId: cache?.id,
            domainNames: defaultServer ? [] : domainNames,
        }
    }
}

export default new HostConverter()
