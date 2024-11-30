package br.com.dillmann.nginxignition.core.nginx.command

fun interface StopNginxCommand {
    suspend fun stop()
}
