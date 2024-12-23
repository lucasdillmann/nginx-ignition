package br.com.dillmann.nginxignition.api.host.model

import br.com.dillmann.nginxignition.api.common.serialization.UuidString
import br.com.dillmann.nginxignition.core.host.Host
import kotlinx.serialization.Serializable

@Serializable
internal data class HostRequest(
    val enabled: Boolean,
    val defaultServer: Boolean,
    val useGlobalBindings: Boolean,
    val domainNames: List<String> = emptyList(),
    val routes: List<Route> = emptyList(),
    val bindings: List<Binding> = emptyList(),
    val featureSet: FeatureSet,
    val accessListId: UuidString? = null,
) {
    @Serializable
    data class Route (
        val priority: Int,
        val type: Host.RouteType,
        val sourcePath: String,
        val settings: RouteSettings,
        val targetUri: String? = null,
        val redirectCode: Int? = null,
        val response: StaticResponse? = null,
        val integration: IntegrationConfig? = null,
        val accessListId: UuidString? = null,
        val sourceCode: RouteSourceCode? = null,
    )

    @Serializable
    data class RouteSourceCode(
        val language: Host.SourceCodeLanguage,
        val code: String,
        val mainFunction: String? = null,
    )

    @Serializable
    data class RouteSettings (
        val includeForwardHeaders: Boolean,
        val proxySslServerName: Boolean,
        val keepOriginalDomainName: Boolean,
        val custom: String? = null,
    )

    @Serializable
    data class IntegrationConfig(
        val integrationId: String,
        val optionId: String,
    )

    @Serializable
    data class StaticResponse(
        val statusCode: Int,
        val payload: String? = null,
        val headers: Map<String, String> = emptyMap(),
    )

    @Serializable
    data class FeatureSet(
        val websocketsSupport: Boolean,
        val http2Support: Boolean,
        val redirectHttpToHttps: Boolean,
    )

    @Serializable
    data class Binding(
        val type: Host.BindingType,
        val ip: String,
        val port: Int,
        val certificateId: UuidString? = null,
    )
}
