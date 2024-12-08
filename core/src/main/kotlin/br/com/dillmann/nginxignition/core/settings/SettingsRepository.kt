package br.com.dillmann.nginxignition.core.settings

interface SettingsRepository {
    suspend fun save(settings: Settings)
    suspend fun get(): Settings
}
