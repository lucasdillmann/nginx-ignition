package br.com.dillmann.nginxignition.core.nginx.command

interface StartNginxCommand {
    suspend fun start()
}
