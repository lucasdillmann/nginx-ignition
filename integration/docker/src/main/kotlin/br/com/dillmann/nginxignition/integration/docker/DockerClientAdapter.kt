package br.com.dillmann.nginxignition.integration.docker

import com.github.dockerjava.api.model.Container
import com.github.dockerjava.core.DefaultDockerClientConfig
import com.github.dockerjava.core.DockerClientBuilder
import com.github.dockerjava.zerodep.ZerodepDockerHttpClient

class DockerClientAdapter(socketPath: String) {
    private val config = DefaultDockerClientConfig
        .createDefaultConfigBuilder()
        .withDockerHost("unix://$socketPath")
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
            .withShowAll(true)
            .exec()
}
