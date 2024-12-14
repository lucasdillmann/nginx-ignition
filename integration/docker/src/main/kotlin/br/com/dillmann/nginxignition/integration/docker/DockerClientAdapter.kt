package br.com.dillmann.nginxignition.integration.docker

import com.github.dockerjava.api.model.Container
import com.github.dockerjava.core.DefaultDockerClientConfig
import com.github.dockerjava.core.DockerClientBuilder
import com.github.dockerjava.zerodep.ZerodepDockerHttpClient

internal class DockerClientAdapter(mode: DockerConnectionMode, host: String) {
    private companion object {
        private const val LIMIT = 1000
    }

    private val config = DefaultDockerClientConfig
        .createDefaultConfigBuilder()
        .withDockerHost(
            when (mode) {
                DockerConnectionMode.SOCKET -> "unix://$host"
                DockerConnectionMode.TCP -> host
            }
        )
        .build()

    private val httpClient = ZerodepDockerHttpClient
        .Builder()
        .dockerHost(config.dockerHost)
        .build()

    private val delegate = DockerClientBuilder
        .getInstance(config)
        .withDockerHttpClient(httpClient)
        .build()

    fun listContainers(): List<Container> =
        delegate
            .listContainersCmd()
            .withLimit(LIMIT)
            .withShowAll(true)
            .exec()
}
