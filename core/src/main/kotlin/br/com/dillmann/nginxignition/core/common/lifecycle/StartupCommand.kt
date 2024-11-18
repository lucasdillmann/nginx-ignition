package br.com.dillmann.nginxignition.core.common.lifecycle

interface StartupCommand {
    val priority: Int
        get() = 100

    suspend fun execute()
}
