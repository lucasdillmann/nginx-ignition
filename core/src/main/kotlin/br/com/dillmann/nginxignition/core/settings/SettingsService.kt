package br.com.dillmann.nginxignition.core.settings

import br.com.dillmann.nginxignition.core.common.scheduler.TaskScheduler
import br.com.dillmann.nginxignition.core.settings.command.GetSettingsCommand
import br.com.dillmann.nginxignition.core.settings.command.SaveSettingsCommand

internal class SettingsService(
    private val repository: SettingsRepository,
    private val scheduler: TaskScheduler,
    private val validator: SettingsValidator,
): GetSettingsCommand, SaveSettingsCommand {
    override suspend fun get(): Settings =
        repository.get()

    override suspend fun save(settings: Settings) {
        validator.validate(settings)
        repository.save(settings)
        scheduler.startOrReload()
    }
}
