package br.com.dillmann.nginxignition.core.user.command

interface GetUserCountCommand {
    suspend fun count(): Long
}
