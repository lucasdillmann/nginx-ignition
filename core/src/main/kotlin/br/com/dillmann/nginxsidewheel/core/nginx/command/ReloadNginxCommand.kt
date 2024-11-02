package br.com.dillmann.nginxsidewheel.core.nginx.command

interface ReloadNginxCommand {
    suspend fun reload()
}
