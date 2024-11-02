package br.com.dillmann.nginxsidewheel.application.controller.host.model

import br.com.dillmann.nginxsidewheel.application.common.serialization.UuidString
import br.com.dillmann.nginxsidewheel.core.host.Host
import kotlinx.serialization.Serializable

@Serializable
data class HostResponse(
    val id: UuidString,
    val default: Boolean,
    val enabled: Boolean,
    val domainNames: List<String>,
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
    )

    @Serializable
    data class StaticResponse(
        val statusCode: Int,
        val payload: String?,
        val headers: Map<String, List<String>>,
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
