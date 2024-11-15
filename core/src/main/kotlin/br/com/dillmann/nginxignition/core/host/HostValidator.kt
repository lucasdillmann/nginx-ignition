package br.com.dillmann.nginxignition.core.host

import br.com.dillmann.nginxignition.core.certificate.CertificateRepository
import br.com.dillmann.nginxignition.core.common.GlobalConstants.TLD_PATTERN
import br.com.dillmann.nginxignition.core.common.validation.ConsistencyException
import org.apache.commons.validator.routines.InetAddressValidator
import java.net.URI

private typealias ErrorCreator = (String, String) -> Unit

internal class HostValidator(
    private val hostRepository: HostRepository,
    private val certificateRepository: CertificateRepository,
) {
    suspend fun validate(host: Host) {
        val violations = mutableListOf<ConsistencyException.Violation>()
        val addError: ErrorCreator = { path, message ->
            violations += ConsistencyException.Violation(path, message)
        }

        validateDefaultFlag(host, addError)
        validateDomainNames(host, addError)
        validateRoutes(host, addError)
        validateBindings(host, addError)

        if (violations.isNotEmpty())
            throw ConsistencyException(violations)
    }

    private suspend fun validateDefaultFlag(host: Host, addError: ErrorCreator) {
        if (host.default) {
            val currentId = hostRepository.findDefault()?.id
            if (currentId != null && host.id != currentId)
                addError("default", "There's already another host marked as the default one")
        }
    }

    private fun validateDomainNames(host: Host, addError: ErrorCreator) {
        if (host.domainNames.isEmpty())
            addError("domainNames", "At least one domain name must be informed")

        host.domainNames.forEachIndexed { index, domainName ->
            if (!TLD_PATTERN.matcher(domainName).matches())
                addError("domainNames[$index]", "Value is not a valid domain name")
        }
    }

    private suspend fun validateBindings(host: Host, addError: ErrorCreator) {
        if (host.bindings.isEmpty())
            addError("bindings", "At least one binding must be informed")

        host.bindings.forEachIndexed { index, binding ->
            validateBinding(binding, index, addError)
        }
    }

    private suspend fun validateBinding(binding: Host.Binding, index: Int, addError: ErrorCreator) {
        if (!InetAddressValidator.getInstance().isValid(binding.ip))
            addError("bindings[$index].ip", "Not a valid IPv4 or IPv6 address")

        if (binding.port !in 1..65535)
            addError("bindings[$index].port", "Value must be between 1 and 65535")

        when {
            binding.type == Host.BindingType.HTTP && binding.certificateId != null ->
                addError("bindings[$index].certificateId", "Value cannot be informed for a HTTP binding")
            binding.type == Host.BindingType.HTTPS && binding.certificateId == null ->
                addError("bindings[$index].certificateId", "Value must be informed for a HTTPS binding")
            binding.type == Host.BindingType.HTTPS && !certificateRepository.existsById(binding.certificateId!!) ->
                addError("bindings[$index].certificateId", "No SSL certificate found with provided ID")
        }
    }

    private fun validateRoutes(host: Host, addError: ErrorCreator) {
        if (host.routes.isEmpty())
            addError("routes", "At least one route must be informed")

        host.routes
            .groupBy { it.priority }
            .filter { it.value.size > 1 }
            .forEach { (priority, _) ->
                addError("routes", "Priority $priority is duplicated in two or more routes")
            }

        host.routes.forEachIndexed { index, route ->
            validateRoute(route, index, addError)
        }
    }

    private fun validateRoute(route: Host.Route, index: Int, addError: ErrorCreator) {
        when (route.type) {
            Host.RouteType.PROXY -> validateProxyRoute(route, index, addError)
            Host.RouteType.REDIRECT -> validateRedirectRoute(route, index, addError)
            Host.RouteType.STATIC_RESPONSE -> validateStaticResponseRoute(route, index, addError)
        }
    }

    private fun validateProxyRoute(route: Host.Route, index: Int, addError: ErrorCreator) {
        if (route.targetUri.isNullOrBlank()) {
            addError(
                "routes[$index].targetUri",
                "Value is required when the type of the route is ${Host.RouteType.PROXY}",
            )
        } else {
            val parseResult = runCatching { URI(route.targetUri) }
            if (parseResult.isFailure)
                addError("routes[$index].targetUri", "Value is not a valid URI")
        }
    }

    private fun validateRedirectRoute(route: Host.Route, index: Int, addError: ErrorCreator) {
        if (route.targetUri.isNullOrBlank()) {
            addError(
                "routes[$index].targetUri",
                "Value is required when the type of the route is ${Host.RouteType.REDIRECT}",
            )
        } else {
            val parseResult = runCatching { URI(route.targetUri) }
            if (parseResult.isFailure)
                addError("routes[$index].targetUri", "Value is not a valid URI")
        }

        if (route.redirectCode !in 300..399)
            addError("routes[$index].redirectCode", "Value must be between 300 and 399")
    }

    private fun validateStaticResponseRoute(route: Host.Route, index: Int, addError: ErrorCreator) {
        if (route.response == null) {
            addError(
                "routes[$index].response",
                "Value is required when the type of the route is ${Host.RouteType.STATIC_RESPONSE}",
            )
            return
        }

        if (route.response.statusCode !in 100..599)
            addError("routes[$index].response.statusCode", "Value must be between 100 and 599")
    }
}
