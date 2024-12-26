package br.com.dillmann.nginxignition.application.lifecycle

import br.com.dillmann.nginxignition.api.ApiModule
import br.com.dillmann.nginxignition.application.ApplicationModule
import br.com.dillmann.nginxignition.certificate.acme.AcmeCertificateModule
import br.com.dillmann.nginxignition.certificate.custom.CustomCertificateModule
import br.com.dillmann.nginxignition.certificate.selfsigned.SelfSignedCertificateModule
import br.com.dillmann.nginxignition.core.CoreModule
import br.com.dillmann.nginxignition.core.common.lifecycle.ShutdownCommand
import br.com.dillmann.nginxignition.core.common.lifecycle.StartupCommand
import br.com.dillmann.nginxignition.core.common.log.logger
import br.com.dillmann.nginxignition.database.DatabaseModule
import br.com.dillmann.nginxignition.integration.docker.DockerIntegrationModule
import br.com.dillmann.nginxignition.integration.truenas.TrueNasIntegrationModule
import kotlinx.coroutines.async
import kotlinx.coroutines.coroutineScope
import kotlinx.coroutines.runBlocking
import org.koin.core.Koin
import org.koin.core.KoinApplication
import org.koin.core.logger.Level

internal class ApplicationLifecycle {
    companion object {
        private val LOGGER = logger<ApplicationLifecycle>()
        lateinit var koin: Koin
            private set
    }

    private lateinit var koinContainer: KoinApplication

    suspend fun start() {
        val startTime = System.currentTimeMillis()
        LOGGER.info("Welcome to nginx ignition")

        Runtime.getRuntime().addShutdownHook(Thread {
           runBlocking { stop() }
        })

        startKoin()
        fireStartupEvents()

        val timeTook = System.currentTimeMillis() - startTime
        LOGGER.info("Application initialized (took ${timeTook}ms)")
    }

    private suspend fun stop() {
        fireShutdownEvents()
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

    private suspend fun fireStartupEvents() {
        koin.getAll<StartupCommand>()
            .sortedBy { it.priority }
            .forEach {
                if (it.async) it.executeAsync()
                else it.execute()
            }
    }

    private suspend fun StartupCommand.executeAsync() {
        val result = coroutineScope {
            async { execute() }
        }
        
        result.invokeOnCompletion { exception ->
            if (exception != null)
                LOGGER.warn("Startup command failed", exception)
        }
    }

    private suspend fun fireShutdownEvents() {
        koin.getAll<ShutdownCommand>()
            .sortedBy { it.priority }
            .forEach {
                try {
                    it.execute()
                } catch (ex: Exception) {
                    LOGGER.warn("Shutdown command failed with an exception", ex)
                }
            }
    }
}
