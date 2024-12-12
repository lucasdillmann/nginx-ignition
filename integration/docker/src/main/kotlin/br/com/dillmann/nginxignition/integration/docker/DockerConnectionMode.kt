package br.com.dillmann.nginxignition.integration.docker

internal enum class DockerConnectionMode(val description: String) {
    SOCKET("Socket"),
    TCP("TCP"),
}
