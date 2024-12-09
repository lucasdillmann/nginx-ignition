package br.com.dillmann.nginxignition.core.nginx.task

import br.com.dillmann.nginxignition.core.common.log.logger
import br.com.dillmann.nginxignition.core.common.scheduler.ScheduledTask
import br.com.dillmann.nginxignition.core.nginx.NginxService
import br.com.dillmann.nginxignition.core.settings.Settings
import br.com.dillmann.nginxignition.core.settings.SettingsRepository
import java.util.concurrent.TimeUnit

internal class LogRotationScheduledTask(
    private val service: NginxService,
    private val settingsRepository: SettingsRepository,
): ScheduledTask {
    override suspend fun run() {
        service.rotateLogs()
    }

    override suspend fun schedule(): ScheduledTask.Schedule {
        val (enabled, _, unit, count) = settingsRepository.get().logRotation
        val timeUnit =
            when (unit) {
                Settings.TimeUnit.MINUTES -> TimeUnit.MINUTES
                Settings.TimeUnit.HOURS -> TimeUnit.HOURS
                Settings.TimeUnit.DAYS -> TimeUnit.DAYS
            }

        return ScheduledTask.Schedule(enabled, timeUnit, count, count)
    }

    override suspend fun onScheduleStarted() {
        val (_, unit, count) = schedule()
        val minutes = unit.toMinutes(count.toLong())
        logger<LogRotationScheduledTask>().info("Log rotation scheduled to execute every $minutes minutes")
    }
}
