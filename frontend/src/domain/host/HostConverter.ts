import HostResponse from "./model/HostResponse"
import HostFormValues, {
    HostFormBinding,
    HostFormRoute,
    HostFormRouteIntegration,
    HostFormStaticResponse,
} from "./model/HostFormValues"
import HostRequest, { HostBinding, HostRoute, HostRouteIntegration, HostRouteStaticResponse } from "./model/HostRequest"
import CertificateService from "../certificate/CertificateService"
import IntegrationService from "../integration/IntegrationService"
import AccessListService from "../accesslist/AccessListService"

class HostConverter {
    private readonly certificateService: CertificateService
    private readonly integrationService: IntegrationService
    private readonly accessListService: AccessListService

    constructor() {
        this.certificateService = new CertificateService()
        this.integrationService = new IntegrationService()
        this.accessListService = new AccessListService()
    }

    private notNull(value?: any) {
        return value !== undefined && value !== null
    }

    private staticResponseToFormValues(response: HostRouteStaticResponse): HostFormStaticResponse {
        const { statusCode, payload } = response
        const { headers } = response
        const joinedHeaders = this.notNull(headers)
                ? Object.entries(headers!).map(([key, value]) => `${key}: ${value}`).join("\n")
                : ""

        return {
            headers: joinedHeaders,
            statusCode,
            payload,
        }
    }

    private async integrationToFormValues(integration: HostRouteIntegration): Promise<HostFormRouteIntegration> {
        const { integrationId, optionId } = integration
        const option = await this.integrationService.getOptionById(integrationId, optionId)

        return {
            integrationId,
            option: option!!,
        }
    }

    private async routeToFormValues(route: HostRoute): Promise<HostFormRoute> {
        const accessList = this.notNull(route.accessListId)
            ? await this.accessListService.getById(route.accessListId!!)
            : undefined
        const response = this.notNull(route.response) ? this.staticResponseToFormValues(route.response!!) : undefined
        const integration = this.notNull(route.integration)
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
            integrationId: integration.integrationId,
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

    async responseToFormValues(response: HostResponse): Promise<HostFormValues> {
        const { enabled, domainNames, featureSet, defaultServer, useGlobalBindings, accessListId } = response

        const routes = response.routes.map(route => this.routeToFormValues(route))
        const bindings = await Promise.all(response.bindings.map(binding => this.bindingToFormValues(binding)))
        const accessList = this.notNull(accessListId) ? await this.accessListService.getById(accessListId!!) : undefined

        return {
            enabled,
            bindings,
            featureSet,
            defaultServer,
            useGlobalBindings,
            accessList,
            domainNames: domainNames ?? [""],
            routes: await Promise.all(routes),
        }
    }

    formValuesToRequest(formValues: HostFormValues): HostRequest {
        const { enabled, domainNames, featureSet, defaultServer, useGlobalBindings, accessList } = formValues

        const routes = formValues.routes.map(route => this.formValuesToRoute(route))
        const bindings = useGlobalBindings ? [] : formValues.bindings.map(binding => this.formValuesToBinding(binding))

        return {
            enabled,
            featureSet,
            defaultServer,
            useGlobalBindings,
            bindings,
            routes,
            accessListId: accessList?.id,
            domainNames: defaultServer ? [] : domainNames,
        }
    }
}

// eslint-disable-next-line import/no-anonymous-default-export
export default new HostConverter()
