package br.com.dillmann.nginxignition.integration.docker

import br.com.dillmann.nginxignition.core.common.dynamicfield.DynamicField

internal object DockerDynamicFields {
    val CONNECTION_MODE = DynamicField(
        id = "connectionMode",
        description = "Connection mode",
        priority = 1,
        required = true,
        type = DynamicField.Type.ENUM,
        enumOptions = DockerConnectionMode.entries.map { DynamicField.EnumOption(it.name, it.description) },
        defaultValue = DockerConnectionMode.SOCKET.name,
    )

    val SOCKET_PATH = DynamicField(
        id = "socketPath",
        description = "Socket path",
        priority = 2,
        required = true,
        type = DynamicField.Type.SINGLE_LINE_TEXT,
        defaultValue = "/var/run/docker.sock",
        condition = DynamicField.Condition(
            parentField = CONNECTION_MODE.id,
            value = DockerConnectionMode.SOCKET.name,
        ),
    )

    val HOST_URL = DynamicField(
        id = "hostUrl",
        description = "Host URL",
        priority = 2,
        required = true,
        helpText = "The URL to be used to connect to the Docker daemon, such as tcp://example.com:2375",
        type = DynamicField.Type.URL,
        condition = DynamicField.Condition(
            parentField = CONNECTION_MODE.id,
            value = DockerConnectionMode.TCP.name,
        ),
    )

    val PROXY_URL = DynamicField(
        id = "proxyUrl",
        description = "Apps URL",
        priority = 3,
        required = false,
        helpText = "The URL to be used when proxying a request to a Docker container. Use this if the apps are " +
            "exposed in another address that isn't the container IP (from the Docker network).",
        type = DynamicField.Type.URL,
    )
}
