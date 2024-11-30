package br.com.dillmann.nginxignition.core.nginx.command

fun interface StartNginxCommand {
    suspend fun start()
}
