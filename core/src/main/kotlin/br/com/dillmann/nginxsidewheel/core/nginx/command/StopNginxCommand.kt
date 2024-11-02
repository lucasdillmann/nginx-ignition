package br.com.dillmann.nginxsidewheel.core.nginx.command

interface StopNginxCommand {
    suspend fun stop()
}
