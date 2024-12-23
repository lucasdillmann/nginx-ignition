package br.com.dillmann.nginxignition.core.host

import br.com.dillmann.nginxignition.core.certificate.CertificateRepository
import br.com.dillmann.nginxignition.core.common.GlobalConstants.TLD_PATTERN
import br.com.dillmann.nginxignition.core.common.validation.ConsistencyValidator
import br.com.dillmann.nginxignition.core.common.validation.ErrorCreator
import org.apache.commons.validator.routines.InetAddressValidator
import java.net.URI

internal class HostValidator(
    private val hostRepository: HostRepository,
    private val certificateRepository: CertificateRepository,
): ConsistencyValidator() {
    private companion object {
        private const val BINDINGS_PATH = "bindings"
        private const val MINIMUM_PORT = 1
        private const val MAXIMUM_PORT = 65535
        private const val MINIMUM_REDIRECT_STATUS_CODE = 300
        private const val MAXIMUM_REDIRECT_STATUS_CODE = 399
        private const val MINIMUM_STATUS_CODE = 100
        private const val MAXIMUM_STATUS_CODE = 599
    }

    suspend fun validate(host: Host) {
        withValidationScope { addError ->
            validateDefaultFlag(host, addError)
            validateDomainNames(host, addError)
            validateRoutes(host, addError)
            validateBindings(host, addError)
        }
    }

    private suspend fun validateDefaultFlag(host: Host, addError: ErrorCreator) {
        if (!host.defaultServer) return

        val currentId = hostRepository.findDefault()?.id
        if (currentId != null && host.id != currentId)
            addError("defaultServer", "There's already another host marked as the default one")

        if (!host.domainNames.isNullOrEmpty())
            addError("domainNames", "Must be empty when the host is the default one")
    }

    private fun validateDomainNames(host: Host, addError: ErrorCreator) {
        if (host.domainNames.isNullOrEmpty() && !host.defaultServer)
            addError("domainNames", "At least one domain name must be informed")

        host.domainNames?.forEachIndexed { index, domainName ->
            if (!TLD_PATTERN.matcher(domainName).matches())
                addError("domainNames[$index]", "Value is not a valid domain name")
        }
    }

    private suspend fun validateBindings(host: Host, addError: ErrorCreator) {
        if (host.useGlobalBindings && host.bindings.isNotEmpty())
            addError(BINDINGS_PATH, "Must be empty when using global bindings")

        if (!host.useGlobalBindings) {
            if (host.bindings.isEmpty())
                addError(BINDINGS_PATH, "At least one binding must be informed")

            host.bindings.forEachIndexed { index, binding ->
                validateBinding(BINDINGS_PATH, binding, index, addError)
            }
        }
    }

    suspend fun validateBinding(pathPrefix: String, binding: Host.Binding, index: Int, addError: ErrorCreator) {
        if (!InetAddressValidator.getInstance().isValid(binding.ip))
            addError("$pathPrefix[$index].ip", "Not a valid IPv4 or IPv6 address")

        if (binding.port !in MINIMUM_PORT..MAXIMUM_PORT)
            addError("$pathPrefix[$index].port", "Value must be between $MINIMUM_PORT and $MAXIMUM_PORT")

        val certificateIdField = "$pathPrefix[$index].certificateId"
        when {
            binding.type == Host.BindingType.HTTP && binding.certificateId != null ->
                addError(certificateIdField, "Value cannot be informed for a HTTP binding")
            binding.type == Host.BindingType.HTTPS && binding.certificateId == null ->
                addError(certificateIdField, "Value must be informed for a HTTPS binding")
            binding.type == Host.BindingType.HTTPS && !certificateRepository.existsById(binding.certificateId!!) ->
                addError(certificateIdField, "No SSL certificate found with provided ID")
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

        val distinctPaths = mutableSetOf<String>()
        host.routes.forEachIndexed { index, route ->
            validateRoute(route, index, addError, distinctPaths)
        }
    }

    private fun validateRoute(
        route: Host.Route,
        index: Int,
        addError: ErrorCreator,
        distinctPaths: MutableSet<String>,
    ) {
        if (!distinctPaths.add(route.sourcePath))
            addError("routes[$index].sourcePath", "Source path was already used in another route")

        when (route.type) {
            Host.RouteType.PROXY -> validateProxyRoute(route, index, addError)
            Host.RouteType.REDIRECT -> validateRedirectRoute(route, index, addError)
            Host.RouteType.STATIC_RESPONSE -> validateStaticResponseRoute(route, index, addError)
            Host.RouteType.INTEGRATION -> validateIntegrationRoute(route, index, addError)
            Host.RouteType.SOURCE_CODE -> validateSourceCodeRoute(route, index, addError)
        }
    }

    private fun validateProxyRoute(route: Host.Route, index: Int, addError: ErrorCreator) {
        @Suppress("StringLiteralDuplication")
        val targetUriField = "routes[$index].targetUri"
        if (route.targetUri.isNullOrBlank()) {
            addError(
                targetUriField,
                "Value is required when the type of the route is proxy",
            )
        } else {
            val parseResult = runCatching { requireNotNull(URI(route.targetUri).host) }
            if (parseResult.isFailure)
                addError(targetUriField, "Value is not a valid URL")
        }
    }

    private fun validateRedirectRoute(route: Host.Route, index: Int, addError: ErrorCreator) {
        if (route.targetUri.isNullOrBlank()) {
            addError(
                "routes[$index].targetUri",
                "Value is required when the type of the route is redirect",
            )
        } else {
            val parseResult = runCatching { requireNotNull(URI(route.targetUri).host) }
            if (parseResult.isFailure)
                addError("routes[$index].targetUri", "Value is not a valid URI")
        }

        if (route.redirectCode !in MINIMUM_REDIRECT_STATUS_CODE..MAXIMUM_REDIRECT_STATUS_CODE)
            addError(
                "routes[$index].redirectCode",
                "Value must be between $MINIMUM_REDIRECT_STATUS_CODE and $MAXIMUM_REDIRECT_STATUS_CODE",
            )
    }

    private fun validateStaticResponseRoute(route: Host.Route, index: Int, addError: ErrorCreator) {
        if (route.response?.statusCode !in MINIMUM_STATUS_CODE..MAXIMUM_STATUS_CODE)
            addError(
                "routes[$index].response.statusCode",
                "Value must be between $MINIMUM_STATUS_CODE and $MAXIMUM_STATUS_CODE",
            )
    }

    private fun validateIntegrationRoute(route: Host.Route, index: Int, addError: ErrorCreator) {
        if (route.integration?.integrationId.isNullOrBlank())
            addError(
                "routes[$index].integration.integrationId",
                "Value is required when the type of the route is integration",
            )

        if (route.integration?.optionId.isNullOrBlank())
            addError(
                "routes[$index].integration.optionId",
                "Value is required when the type of the route is integration",
            )
    }

    private fun validateSourceCodeRoute(route: Host.Route, index: Int, addError: ErrorCreator) {
        if (route.sourceCode?.code.isNullOrBlank())
            addError(
                "routes[$index].sourceCode.code",
                "Value is required when the type of the route is source code",
            )

        if (route.sourceCode?.language == null)
            addError(
                "routes[$index].sourceCode.language",
                "Value is required when the type of the route is source code",
            )

        if (route.sourceCode?.language == Host.SourceCodeLanguage.JAVASCRIPT &&
            route.sourceCode.mainFunction.isNullOrBlank())
            addError(
                "routes[$index].sourceCode.mainFunction",
                "Value is required when the language is JavaScript",
            )
    }
}
