package br.com.dillmann.nginxsidewheel.core.common.startup

interface StartupCommand {
    val priority: Int
    suspend fun execute()
}
