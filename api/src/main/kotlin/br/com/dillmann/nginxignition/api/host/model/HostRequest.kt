package br.com.dillmann.nginxignition.api.host.model

import br.com.dillmann.nginxignition.api.common.serialization.UuidString
import br.com.dillmann.nginxignition.core.host.Host
import kotlinx.serialization.Serializable

@Serializable
internal data class HostRequest(
    val enabled: Boolean,
    val defaultServer: Boolean,
    val domainNames: List<String> = emptyList(),
    val routes: List<Route> = emptyList(),
    val bindings: List<Binding> = emptyList(),
    val featureSet: FeatureSet,
) {
    @Serializable
    data class Route (
        val priority: Int,
        val type: Host.RouteType,
        val sourcePath: String,
        val targetUri: String? = null,
        val customSettings: String? = null,
        val redirectCode: Int? = null,
        val response: StaticResponse? = null,
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
