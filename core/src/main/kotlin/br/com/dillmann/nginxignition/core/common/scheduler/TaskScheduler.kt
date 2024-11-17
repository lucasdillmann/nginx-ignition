package br.com.dillmann.nginxignition.core.common.scheduler

import br.com.dillmann.nginxignition.core.common.log.logger
import kotlinx.coroutines.runBlocking
import java.util.concurrent.Executors
import java.util.concurrent.TimeUnit

object TaskScheduler {
    private val logger = logger<TaskScheduler>()
    private val delegate = Executors.newScheduledThreadPool(4)

    fun shutdown() {
        logger.info("Stopping the task scheduler (30s timeout)")
        delegate.shutdown()

        try {
            delegate.awaitTermination(30, TimeUnit.SECONDS)
        } catch (ex: Exception) {
            logger.warn("Task scheduler graceful shutdown failed", ex)
        }
    }

    fun schedule(
        task: suspend () -> Unit,
        timeUnit: TimeUnit,
        interval: Long,
        initialDelay: Long,
    ) {
        delegate.scheduleAtFixedRate(
            buildTaskProxy(task),
            initialDelay,
            interval,
            timeUnit,
        )
    }

    private fun buildTaskProxy(task: suspend () -> Unit): Runnable =
        Runnable {
            try {
                runBlocking { task() }
            } catch (ex: Exception) {
                logger.warn("Scheduled task failed with an exception", ex)
            }
        }
}
