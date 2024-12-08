package br.com.dillmann.nginxignition.core.certificate.task

import br.com.dillmann.nginxignition.core.certificate.CertificateService
import br.com.dillmann.nginxignition.core.common.log.logger
import br.com.dillmann.nginxignition.core.common.configuration.ConfigurationProvider
import br.com.dillmann.nginxignition.core.common.scheduler.ScheduledTask
import java.util.concurrent.TimeUnit

internal class CertificateAutoRenewScheduledTask(
    private val configurationProvider: ConfigurationProvider,
    private val service: CertificateService,
): ScheduledTask {
    override fun schedule(): ScheduledTask.Schedule {
        val intervalMinutes =
            configurationProvider.get("nginx-ignition.certificate.auto-renew-interval-minutes").toLong()

        return ScheduledTask.Schedule(
            initialDelay = intervalMinutes,
            interval = intervalMinutes,
            unit = TimeUnit.MINUTES,
        )
    }

    override fun onScheduleStarted() {
        val intervalMinutes =
            configurationProvider.get("nginx-ignition.certificate.auto-renew-interval-minutes").toLong()

        logger<CertificateAutoRenewScheduledTask>()
            .info("Certificate auto-renew checks scheduled to execute every $intervalMinutes minutes")
    }

    override suspend fun run() {
        service.renewAllDue()
    }
}
