package br.com.dillmann.nginxignition.api.host.model

import br.com.dillmann.nginxignition.api.common.serialization.UuidString
import br.com.dillmann.nginxignition.core.host.Host
import kotlinx.serialization.Serializable

@Serializable
internal data class HostResponse(
    val id: UuidString,
    val enabled: Boolean,
    val defaultServer: Boolean,
    val domainNames: List<String>?,
    val routes: List<Route>,
    val bindings: List<Binding>,
    val featureSet: FeatureSet,
) {
    @Serializable
    data class Route (
        val priority: Int,
        val type: Host.RouteType,
        val sourcePath: String,
        val targetUri: String?,
        val customSettings: String?,
        val response: StaticResponse?,
        val integration: IntegrationConfig?,
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
