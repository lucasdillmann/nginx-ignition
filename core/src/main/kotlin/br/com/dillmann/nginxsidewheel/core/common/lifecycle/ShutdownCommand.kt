package br.com.dillmann.nginxsidewheel.core.common.lifecycle

interface ShutdownCommand {
    val priority: Int
    suspend fun execute()
}
