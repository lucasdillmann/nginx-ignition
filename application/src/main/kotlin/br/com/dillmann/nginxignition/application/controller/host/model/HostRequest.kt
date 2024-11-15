package br.com.dillmann.nginxignition.application.controller.host.model

import br.com.dillmann.nginxignition.application.common.serialization.UuidString
import br.com.dillmann.nginxignition.core.host.Host
import kotlinx.serialization.Serializable

@Serializable
data class HostRequest(
    val default: Boolean,
    val enabled: Boolean,
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
        val payload: String?,
        val headers: Map<String, List<String>> = emptyMap(),
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
