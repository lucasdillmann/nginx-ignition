package br.com.dillmann.nginxignition.core.certificate.task

import br.com.dillmann.nginxignition.core.certificate.CertificateService
import br.com.dillmann.nginxignition.core.common.log.logger
import br.com.dillmann.nginxignition.core.common.scheduler.ScheduledTask
import br.com.dillmann.nginxignition.core.settings.Settings
import br.com.dillmann.nginxignition.core.settings.SettingsRepository
import java.util.concurrent.TimeUnit

internal class CertificateAutoRenewScheduledTask(
    private val settingsRepository: SettingsRepository,
    private val certificateService: CertificateService,
): ScheduledTask {
    override suspend fun schedule(): ScheduledTask.Schedule {
        val (enabled, unit, count) = settingsRepository.get().certificateAutoRenew
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

        logger<CertificateAutoRenewScheduledTask>()
            .info("Certificate auto-renew checks scheduled to execute every $minutes minutes")
    }

    override suspend fun run() {
        certificateService.renewAllDue()
    }
}
