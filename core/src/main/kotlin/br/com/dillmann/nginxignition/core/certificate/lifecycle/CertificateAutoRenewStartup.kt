package br.com.dillmann.nginxignition.core.certificate.lifecycle

import br.com.dillmann.nginxignition.core.certificate.CertificateService
import br.com.dillmann.nginxignition.core.common.lifecycle.StartupCommand
import br.com.dillmann.nginxignition.core.common.log.logger
import br.com.dillmann.nginxignition.core.common.configuration.ConfigurationProvider
import br.com.dillmann.nginxignition.core.common.scheduler.TaskScheduler
import java.util.concurrent.TimeUnit

internal class CertificateAutoRenewStartup(
    private val configurationProvider: ConfigurationProvider,
    private val service: CertificateService,
): StartupCommand {
    override suspend fun execute() {
        val intervalMinutes =
            configurationProvider.get("nginx-ignition.certificate.auto-renew-interval-minutes").toLong()

        TaskScheduler.schedule(
            task = service::renewAllDue,
            initialDelay = intervalMinutes,
            interval = intervalMinutes,
            timeUnit = TimeUnit.MINUTES,
        )

        logger<CertificateAutoRenewStartup>()
            .info("Certificate auto-renew checks scheduled to execute every $intervalMinutes minutes")
    }
}
