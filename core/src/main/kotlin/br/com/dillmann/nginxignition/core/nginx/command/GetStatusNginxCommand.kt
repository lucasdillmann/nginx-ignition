package br.com.dillmann.nginxignition.core.nginx.command

interface GetStatusNginxCommand {
    suspend fun isRunning(): Boolean
}
