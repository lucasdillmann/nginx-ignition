package br.com.dillmann.nginxignition.integration.truenas

import br.com.dillmann.nginxignition.core.common.pagination.Page
import br.com.dillmann.nginxignition.core.integration.IntegrationAdapter
import br.com.dillmann.nginxignition.integration.truenas.client.TrueNasApiClient
import br.com.dillmann.nginxignition.integration.truenas.client.TrueNasAppDetailsResponse
import java.net.URI

class TrueNasIntegrationAdapter: IntegrationAdapter {
    override val id = "TRUENAS_SCALE"
    override val name = "TrueNAS Scale"
    override val description =
        "TrueNAS allows, alongside many other things, to run your favorite apps under Docker containers. With this " +
            "integration enabled, you will be able to easily pick any app exposing a service in your TrueNAS as a " +
            "target for your nginx ignition's host routes."
    override val configurationFields =
        listOf(
            TrueNasDynamicFields.URL,
            TrueNasDynamicFields.USERNAME,
            TrueNasDynamicFields.PASSWORD,
        )

    override suspend fun getAvailableOptions(
        parameters: Map<String, Any?>,
        pageNumber: Int,
        pageSize: Int,
    ): Page<IntegrationAdapter.Option> =
        getAvailableApps(parameters).flatMap(::buildOptions).let { Page.of(it) }

    override suspend fun getAvailableOptionById(parameters: Map<String, Any?>, id: String): IntegrationAdapter.Option? =
        runCatching {
            val (appId, containerPort) = id.split(":")
            getWorkloadPort(parameters, appId, containerPort.toInt())!!.let { (app, port) ->
                IntegrationAdapter.Option(
                    id = id,
                    name = "${app.name} (${port.hostPorts.first().hostPort} HTTP)",
                )
            }
        }.getOrNull()

    override suspend fun getOptionProxyUrl(id: String, parameters: Map<String, Any?>): String {
        val baseUrl = parameters[TrueNasDynamicFields.URL.id] as String
        val (appId, containerPort) = id.split(":")
        val (hostPort, hostIp) = getWorkloadPort(parameters, appId, containerPort.toInt())!!.second.hostPorts.first()
        val endpoint = if (hostIp == "0.0.0.0") URI(baseUrl).host else hostIp
        return "http://$endpoint:$hostPort"
    }

    private fun getWorkloadPort(
        parameters: Map<String, Any?>,
        appId: String,
        containerPort: Int,
    ): Pair<TrueNasAppDetailsResponse, TrueNasAppDetailsResponse.WorkloadPort>? {
        val app = getAvailableApps(parameters).find { it.id == appId } ?: return null
        val port = app.activeWorkloads.usedPorts.find { it.containerPort == containerPort } ?: return null
        return app to port
    }

    private fun buildOptions(app: TrueNasAppDetailsResponse): List<IntegrationAdapter.Option> =
        app
            .activeWorkloads
            .usedPorts
            .filter { it.protocol.equals("tcp", true) }
            .flatMap { buildOptions(app, it) }
            .sortedBy { it.name }

    private fun buildOptions(
        app: TrueNasAppDetailsResponse,
        port: TrueNasAppDetailsResponse.WorkloadPort,
    ): List<IntegrationAdapter.Option> =
        port.hostPorts.map {
            IntegrationAdapter.Option(
                id = "${app.id}:${port.containerPort}",
                name = "${app.name} (${it.hostPort} HTTP)"
            )
        }

    private fun getAvailableApps(parameters: Map<String, Any?>): List<TrueNasAppDetailsResponse> {
        val baseUrl = parameters[TrueNasDynamicFields.URL.id] as String
        val username = parameters[TrueNasDynamicFields.USERNAME.id] as String
        val password = parameters[TrueNasDynamicFields.PASSWORD.id] as String
        return TrueNasApiClient(baseUrl, username, password).getAvailableApps()
    }
}
