package br.com.dillmann.nginxignition.core.common.lifecycle

interface StartupCommand {
    val priority: Int
    suspend fun execute()
}
