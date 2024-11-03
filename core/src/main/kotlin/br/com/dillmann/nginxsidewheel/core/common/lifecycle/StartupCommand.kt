package br.com.dillmann.nginxsidewheel.core.common.lifecycle

interface StartupCommand {
    val priority: Int
    suspend fun execute()
}
