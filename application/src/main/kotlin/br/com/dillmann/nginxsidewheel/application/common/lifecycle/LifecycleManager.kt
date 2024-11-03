package br.com.dillmann.nginxsidewheel.application.common.lifecycle

import br.com.dillmann.nginxsidewheel.core.common.lifecycle.ShutdownCommand
import br.com.dillmann.nginxsidewheel.core.common.lifecycle.StartupCommand
import br.com.dillmann.nginxsidewheel.core.common.log.logger

class LifecycleManager(
    private val startupListeners: List<StartupCommand>,
    private val shutdownListeners: List<ShutdownCommand>,
) {
    private val logger = logger<LifecycleManager>()

    suspend fun fireStartupEvent() {
        startupListeners
            .sortedBy { it.priority }
            .forEach { it.execute() }
    }

    suspend fun fireShutdownEvent() {
        logger.info("Shutdown signal received. Starting graceful shutdown.")

        shutdownListeners
            .sortedBy { it.priority }
            .forEach {
                try {
                    it.execute()
                } catch (ex: Exception) {
                    logger.warn("Lifecycle listener failed with an exception", ex)
                }
            }
    }
}
