import HostResponse from "./model/HostResponse";
import HostFormValues, {HostFormBinding, HostFormRoute, HostFormStaticResponse} from "./model/HostFormValues";
import HostRequest, {HostBinding, HostRoute, HostRouteStaticResponse} from "./model/HostRequest";
import CertificateService from "../certificate/CertificateService";

class HostConverter {
    private certificateService: CertificateService

    constructor() {
        this.certificateService = new CertificateService()
    }

    private notNull(value?: any) {
        return value !== undefined && value !== null
    }

    private staticResponseToFormValues(response: HostRouteStaticResponse): HostFormStaticResponse {
        const {statusCode, payload} = response
        const headers = response.headers !== undefined
            ? Object.entries(response.headers).map(([key, value]) => `${key}: ${value}`).join("\n")
            : ""

        return {
            headers,
            statusCode,
            payload,
        }
    }

    private routeToFormValues(route: HostRoute): HostFormRoute {
        const response = this.notNull(route.response)
            ? this.staticResponseToFormValues(route.response!!)
            : undefined

        return {
            ...route,
            response,
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
        const lines = headers
            .split("\n")
            .filter(line => line.trim().length > 0)
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
        const {statusCode, payload} = response
        const headers = this.notNull(response.headers)
            ? this.formValuesToHeaders(response.headers!!)
            : {}

        return {
            statusCode,
            payload,
            headers,
        }
    }

    private formValuesToRoute(route: HostFormRoute): HostRoute {
        const {priority, type, customSettings, targetUri, sourcePath} = route
        const response = this.notNull(route.response)
            ? this.formValuesToStaticResponse(route.response!!)
            : undefined

        return {
            priority,
            type,
            customSettings,
            targetUri,
            sourcePath,
            response,
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
        const {enabled, domainNames, featureSet, defaultServer} = response

        const routes =
            response.routes.map(route => this.routeToFormValues(route))
        const bindings =
            await Promise.all(response.bindings.map(binding => this.bindingToFormValues(binding)))

        return {
            enabled,
            domainNames,
            routes,
            bindings,
            featureSet,
            defaultServer,
        }
    }

    formValuesToRequest(formValues: HostFormValues): HostRequest {
        const {enabled, domainNames, featureSet, defaultServer} = formValues

        const routes = formValues.routes.map(route => this.formValuesToRoute(route))
        const bindings = formValues.bindings.map(binding => this.formValuesToBinding(binding))

        return {
            enabled,
            domainNames,
            featureSet,
            defaultServer,
            bindings,
            routes,
        }
    }
}

// eslint-disable-next-line import/no-anonymous-default-export
export default new HostConverter()
