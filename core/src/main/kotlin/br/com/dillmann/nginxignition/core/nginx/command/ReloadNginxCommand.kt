package br.com.dillmann.nginxignition.core.nginx.command

fun interface ReloadNginxCommand {
    suspend fun reload()
}
