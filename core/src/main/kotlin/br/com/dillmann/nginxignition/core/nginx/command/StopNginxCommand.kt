package br.com.dillmann.nginxignition.core.nginx.command

interface StopNginxCommand {
    suspend fun stop()
}
