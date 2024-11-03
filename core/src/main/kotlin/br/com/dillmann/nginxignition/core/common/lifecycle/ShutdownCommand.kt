package br.com.dillmann.nginxignition.core.common.lifecycle

interface ShutdownCommand {
    val priority: Int
    suspend fun execute()
}
