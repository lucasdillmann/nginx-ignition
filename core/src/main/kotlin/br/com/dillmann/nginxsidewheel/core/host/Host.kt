package br.com.dillmann.nginxsidewheel.core.host

import java.util.UUID

data class Host(
    val id: UUID,
    val default: Boolean,
    val enabled: Boolean,
    val domainNames: List<String>,
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
        val customSettings: String?,
        val response: StaticResponse?,
    )

    enum class RouteType {
        PROXY,
        REDIRECT,
        STATIC_RESPONSE,
    }

    data class StaticResponse(
        val statusCode: Int,
        val payload: String?,
        val headers: Map<String, List<String>>,
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
