package br.com.dillmann.nginxignition.integration.docker

import br.com.dillmann.nginxignition.core.common.pagination.Page
import br.com.dillmann.nginxignition.core.integration.IntegrationAdapter
import com.github.dockerjava.api.model.Container
import com.github.dockerjava.api.model.ContainerPort
import java.net.URI

class DockerIntegrationAdapter: IntegrationAdapter {
    private companion object {
        private const val UNIQUE_ID_PORTION_SIZE_CHARS = 8
    }

    override val id = "DOCKER"
    override val name = "Docker"
    override val priority = 1
    override val description =
        "Enables easy pick of a Docker container with ports exposing a service as a target for your nginx " +
            "ignition's host routes."
    override val configurationFields =
        listOf(DockerDynamicFields.SOCKET_PATH, DockerDynamicFields.PROXY_URL)

    override suspend fun getAvailableOptions(
        parameters: Map<String, Any?>,
        pageNumber: Int,
        pageSize: Int,
        searchTerms: String?,
    ): Page<IntegrationAdapter.Option> =
        getAvailableOptions(parameters)
            .map { it.asIntegrationOption() }
            .filter { searchTerms == null || it.name.contains(searchTerms, ignoreCase = true) }
            .let { Page.of(it) }

    override suspend fun getAvailableOptionById(parameters: Map<String, Any?>, id: String): IntegrationAdapter.Option? =
        findById(id, parameters)?.asIntegrationOption()

    override suspend fun getOptionProxyUrl(id: String, parameters: Map<String, Any?>): String {
        val (container, port) = findById(id, parameters) ?: error("No container found with ID $id")
        val publicUrl = parameters[DockerDynamicFields.PROXY_URL.id] as? String
        val targetHost =
            publicUrl?.let { URI(it) }?.host
                ?: container.networkSettings?.networks?.values?.firstOrNull()?.ipAddress
                ?: error("No network or IP address found for the container with ID $id")

        val targetPort = if (publicUrl != null) port.publicPort else port.privatePort
        return "http://$targetHost:$targetPort"
    }

    private fun findById(id: String, parameters: Map<String, Any?>): Pair<Container, ContainerPort>? {
        val (containerId, portId) = id.split(":")
        return getAvailableOptions(parameters).firstOrNull { (container, port) ->
            containerId == container.id && portId == port.publicPort.toString()
        }
    }

    private fun Pair<Container, ContainerPort>.asIntegrationOption(): IntegrationAdapter.Option {
        val containerName = first.names.firstOrNull() ?: first.id.substring(0, UNIQUE_ID_PORTION_SIZE_CHARS)
        val normalizedName =
            if (containerName.startsWith("/")) containerName.removePrefix("/")
            else containerName

        return IntegrationAdapter.Option(
            id = "${first.id}:${second.publicPort}",
            name = "$normalizedName (${second.publicPort} HTTP)",
        )
    }

    private fun getAvailableOptions(parameters: Map<String, Any?>) =
        startClientAdapter(parameters)
            .listContainers()
            .flatMap { container -> container.ports.map { container to it } }
            .filter { (_, port) -> port.type == "tcp" }

    private fun startClientAdapter(parameters: Map<String, Any?>): DockerClientAdapter {
        val socketPath = parameters[DockerDynamicFields.SOCKET_PATH.id] as String
        return DockerClientAdapter(socketPath)
    }
}
