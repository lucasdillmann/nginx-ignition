package br.com.dillmann.nginxignition.database.host

import br.com.dillmann.nginxignition.core.host.Host
import br.com.dillmann.nginxignition.database.host.mapping.HostBindingTable
import br.com.dillmann.nginxignition.database.host.mapping.HostRouteTable
import br.com.dillmann.nginxignition.database.host.mapping.HostTable
import kotlinx.serialization.encodeToString
import kotlinx.serialization.json.Json
import org.jetbrains.exposed.sql.ResultRow
import org.jetbrains.exposed.sql.statements.InsertStatement
import java.util.*

internal class HostConverter {
    fun apply(host: Host, scope: InsertStatement<out Any>) {
        with(HostTable) {
            scope[id] = host.id
            scope[default] = host.default
            scope[enabled] = host.enabled
            scope[domainNames] = Json.encodeToString(host.domainNames)
            scope[websocketSupport] = host.featureSet.websocketsSupport
            scope[http2Support] = host.featureSet.http2Support
            scope[redirectHttpToHttps] = host.featureSet.redirectHttpToHttps
        }
    }

    fun apply(parentId: UUID, route: Host.Route, scope: InsertStatement<out Any>) {
        with(HostRouteTable) {
            scope[id] = route.id
            scope[priority] = route.priority
            scope[hostId] = parentId
            scope[type] = route.type.name
            scope[sourcePath] = route.sourcePath
            scope[targetUri] = route.targetUri
            scope[customSettings] = route.customSettings
            scope[redirectCode] = route.redirectCode
            scope[staticResponseCode] = route.response?.statusCode
            scope[staticResponsePayload] = route.response?.payload
            scope[staticResponseHeaders] = route.response?.headers?.let(Json::encodeToString)
        }
    }

    fun apply(parentId: UUID, binding: Host.Binding, scope: InsertStatement<out Any>) {
        with(HostBindingTable) {
            scope[id] = binding.id
            scope[hostId] = parentId
            scope[type] = binding.type.name
            scope[ip] = binding.ip
            scope[port] = binding.port
            scope[certificateId] = binding.certificateId
        }
    }

    fun toHost(host: ResultRow, bindings: List<ResultRow>, routes: List<ResultRow>) =
        Host(
            id = host[HostTable.id],
            default = host[HostTable.default],
            enabled = host[HostTable.enabled],
            domainNames = Json.decodeFromString(host[HostTable.domainNames]),
            routes = routes.map(::toRoute),
            bindings = bindings.map(::toBinding),
            featureSet = Host.FeatureSet(
                websocketsSupport = host[HostTable.websocketSupport],
                http2Support = host[HostTable.http2Support],
                redirectHttpToHttps = host[HostTable.redirectHttpToHttps],
            )
        )

    private fun toRoute(route: ResultRow): Host.Route {
        val response =
            if (route[HostRouteTable.staticResponseCode] == null) null
            else Host.StaticResponse(
                statusCode = route[HostRouteTable.staticResponseCode]!!,
                payload = route[HostRouteTable.staticResponsePayload]!!,
                headers = Json.decodeFromString(route[HostRouteTable.staticResponseHeaders]!!),
            )

        return Host.Route(
            id = route[HostRouteTable.id],
            priority = route[HostRouteTable.priority],
            type = Host.RouteType.valueOf(route[HostRouteTable.type]),
            customSettings = route[HostRouteTable.customSettings],
            sourcePath = route[HostRouteTable.sourcePath],
            targetUri = route[HostRouteTable.targetUri],
            redirectCode = route[HostRouteTable.redirectCode],
            response = response,
        )
    }

    private fun toBinding(binding: ResultRow) =
        Host.Binding(
            id = binding[HostBindingTable.id],
            type = Host.BindingType.valueOf(binding[HostBindingTable.type]),
            ip = binding[HostBindingTable.ip],
            port = binding[HostBindingTable.port],
            certificateId = binding[HostBindingTable.certificateId],
        )
}