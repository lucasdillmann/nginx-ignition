package br.com.dillmann.nginxignition.core.user.command

fun interface GetUserCountCommand {
    suspend fun count(): Long
}
