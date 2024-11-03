package br.com.dillmann.nginxsidewheel.core.nginx.command

interface GetStatusNginxCommand {
    suspend fun isRunning(): Boolean
}
