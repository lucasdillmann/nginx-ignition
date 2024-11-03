package br.com.dillmann.nginxignition.core.nginx.command

interface ReloadNginxCommand {
    suspend fun reload()
}
