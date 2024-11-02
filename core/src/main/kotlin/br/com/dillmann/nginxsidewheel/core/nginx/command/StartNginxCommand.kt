package br.com.dillmann.nginxsidewheel.core.nginx.command

interface StartNginxCommand {
    suspend fun start()
}
