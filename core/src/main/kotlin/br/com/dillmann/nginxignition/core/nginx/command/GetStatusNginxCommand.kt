package br.com.dillmann.nginxignition.core.nginx.command

fun interface GetStatusNginxCommand {
    suspend fun isRunning(): Boolean
}
