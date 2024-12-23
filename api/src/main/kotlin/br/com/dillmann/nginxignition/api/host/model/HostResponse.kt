package br.com.dillmann.nginxignition.api.host.model

import br.com.dillmann.nginxignition.api.common.serialization.UuidString
import br.com.dillmann.nginxignition.core.host.Host
import kotlinx.serialization.Serializable

@Serializable
internal data class HostResponse(
    val id: UuidString,
    val enabled: Boolean,
    val defaultServer: Boolean,
    val useGlobalBindings: Boolean,
    val domainNames: List<String>?,
    val routes: List<Route>,
    val bindings: List<Binding>,
    val featureSet: FeatureSet,
    val accessListId: UuidString?,
) {
    @Serializable
    data class Route (
        val priority: Int,
        val type: Host.RouteType,
        val sourcePath: String,
        val settings: RouteSettings,
        val targetUri: String?,
        val redirectCode: Int?,
        val response: StaticResponse?,
        val integration: IntegrationConfig?,
        val accessListId: UuidString?,
        val sourceCode: RouteSourceCode?,
    )

    @Serializable
    data class RouteSourceCode(
        val language: Host.SourceCodeLanguage,
        val code: String,
        val mainFunction: String?,
    )

    @Serializable
    data class RouteSettings (
        val includeForwardHeaders: Boolean,
        val proxySslServerName: Boolean,
        val keepOriginalDomainName: Boolean,
        val custom: String?,
    )

    @Serializable
    data class IntegrationConfig(
        val integrationId: String,
        val optionId: String,
    )

    @Serializable
    data class StaticResponse(
        val statusCode: Int,
        val payload: String?,
        val headers: Map<String, String>,
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
        val certificateId: UuidString?,
    )
}
