package br.com.dillmann.nginxignition.application

import br.com.dillmann.nginxignition.api.ApiModule
import br.com.dillmann.nginxignition.certificate.acme.AcmeCertificateModule
import br.com.dillmann.nginxignition.certificate.custom.CustomCertificateModule
import br.com.dillmann.nginxignition.certificate.selfsigned.SelfSignedCertificateModule
import br.com.dillmann.nginxignition.core.CoreModule
import br.com.dillmann.nginxignition.core.common.lifecycle.LifecycleManager
import br.com.dillmann.nginxignition.core.common.log.logger
import br.com.dillmann.nginxignition.database.DatabaseModule
import br.com.dillmann.nginxignition.integration.docker.DockerIntegrationModule
import br.com.dillmann.nginxignition.integration.truenas.TrueNasIntegrationModule
import kotlinx.coroutines.runBlocking
import org.koin.core.Koin
import org.koin.core.KoinApplication
import org.koin.core.logger.Level

internal class Application {
    companion object {
        private val LOGGER = logger<Application>()
        lateinit var koin: Koin
            private set
    }

    private lateinit var koinContainer: KoinApplication

    suspend fun boot() {
        val startTime = System.currentTimeMillis()
        LOGGER.info("Welcome to nginx ignition")

        Runtime.getRuntime().addShutdownHook(Thread {
           runBlocking { stop() }
        })

        startKoin()
        koin.get<LifecycleManager>().fireStartupEvent()

        val timeTook = System.currentTimeMillis() - startTime
        LOGGER.info("Application initialized (took ${timeTook}ms)")
    }

    private suspend fun stop() {
        koin.get<LifecycleManager>().fireShutdownEvent()
        koinContainer.close()
    }

    private fun startKoin() {
        koinContainer = KoinApplication
            .init()
            .modules(
                CoreModule.initialize(),
                DatabaseModule.initialize(),
                AcmeCertificateModule.initialize(),
                CustomCertificateModule.initialize(),
                SelfSignedCertificateModule.initialize(),
                TrueNasIntegrationModule.initialize(),
                DockerIntegrationModule.initialize(),
                ApiModule.initialize(),
                ApplicationModule.initialize(),
            )
            .printLogger(level = Level.ERROR)

        koin = koinContainer.koin
    }
}
