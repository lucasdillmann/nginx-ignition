package br.com.dillmann.nginxignition.integration.docker

import br.com.dillmann.nginxignition.core.common.dynamicfield.DynamicField

internal object DockerDynamicFields {
    val SOCKET_PATH = DynamicField(
        id = "socketPath",
        description = "Socket path",
        priority = 1,
        required = true,
        helpText = "Path to the Docker socket file",
        type = DynamicField.Type.SINGLE_LINE_TEXT,
        defaultValue = "/var/run/docker.sock",
    )

    val PROXY_URL = DynamicField(
        id = "proxyUrl",
        description = "Apps URL",
        priority = 2,
        required = false,
        helpText = "The URL to be used when proxying a request to a Docker container. Use this if the apps are " +
            "exposed in another address that isn't the container IP (from the Docker network).",
        type = DynamicField.Type.URL,
    )
}
