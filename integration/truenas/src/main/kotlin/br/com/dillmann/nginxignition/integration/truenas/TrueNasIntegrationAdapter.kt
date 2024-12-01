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
    ): Page<IntegrationAdapter.Option> {
        val baseUrl = parameters[TrueNasDynamicFields.URL.id] as String
        val username = parameters[TrueNasDynamicFields.USERNAME.id] as String
        val password = parameters[TrueNasDynamicFields.PASSWORD.id] as String
        val baseHost = URI(baseUrl).host

        return TrueNasApiClient(baseUrl, username, password)
            .getAvailableApps()
            .flatMap { buildOptions(it, baseHost) }
            .let { Page.of(it) }
    }

    private fun buildOptions(app: TrueNasAppDetailsResponse, baseHost: String): List<IntegrationAdapter.Option> =
        app
            .activeWorkloads
            .usedPorts
            .filter { it.protocol.equals("tcp", true) }
            .flatMap { buildOptions(app, it, baseHost) }
            .sortedBy { it.name }

    private fun buildOptions(
        app: TrueNasAppDetailsResponse,
        port: TrueNasAppDetailsResponse.WorkloadPort,
        baseHost: String
    ): List<IntegrationAdapter.Option> =
        port.hostPorts.map {
            val name = app.name
            val (hostPort, hostIp) = it
            val endpoint = if (hostIp == "0.0.0.0") baseHost else hostPort

            IntegrationAdapter.Option(
                id = "$endpoint:$hostPort:http",
                name = "$name ($hostPort HTTP)"
            )
        }

    override suspend fun getOptionProxyUrl(id: String, parameters: Map<String, Any?>): String {
        val (endpoint, hostPort, protocol) = id.split(":")
        return "$protocol://$endpoint:$hostPort"
    }
}
