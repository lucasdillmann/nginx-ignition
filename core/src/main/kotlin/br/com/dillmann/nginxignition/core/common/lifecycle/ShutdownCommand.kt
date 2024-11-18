package br.com.dillmann.nginxignition.core.common.lifecycle

interface ShutdownCommand {
    val priority: Int
        get() = 100

    suspend fun execute()
}
