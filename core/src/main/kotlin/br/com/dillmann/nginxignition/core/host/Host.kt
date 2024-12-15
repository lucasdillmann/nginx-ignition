package br.com.dillmann.nginxignition.core.host

import java.util.UUID

data class Host(
    val id: UUID,
    val enabled: Boolean,
    val defaultServer: Boolean,
    val useGlobalBindings: Boolean,
    val domainNames: List<String>?,
    val routes: List<Route>,
    val bindings: List<Binding>,
    val featureSet: FeatureSet,
) {
    data class Route (
        val id: UUID,
        val priority: Int,
        val type: RouteType,
        val sourcePath: String,
        val targetUri: String?,
        val redirectCode: Int?,
        val settings: RouteSettings,
        val response: StaticResponse?,
        val integration: IntegrationConfig?,
    )

    data class RouteSettings (
        val includeForwardHeaders: Boolean,
        val proxySslServerName: Boolean,
        val keepOriginalDomainName: Boolean,
        val forwardQueryParams: Boolean,
        val custom: String?,
    )

    enum class RouteType {
        PROXY,
        REDIRECT,
        STATIC_RESPONSE,
        INTEGRATION,
    }

    data class StaticResponse(
        val statusCode: Int,
        val payload: String?,
        val headers: Map<String, String>,
    )

    data class IntegrationConfig(
        val integrationId: String,
        val optionId: String,
    )

    data class FeatureSet(
        val websocketsSupport: Boolean,
        val http2Support: Boolean,
        val redirectHttpToHttps: Boolean,
    )

    data class Binding(
        val id: UUID,
        val type: BindingType,
        val ip: String,
        val port: Int,
        val certificateId: UUID?,
    )

    enum class BindingType {
        HTTP,
        HTTPS,
    }
}
