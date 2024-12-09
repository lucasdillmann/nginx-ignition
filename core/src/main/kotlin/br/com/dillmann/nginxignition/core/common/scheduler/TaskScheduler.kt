package br.com.dillmann.nginxignition.core.common.scheduler

import br.com.dillmann.nginxignition.core.common.log.logger
import kotlinx.coroutines.runBlocking
import java.util.concurrent.Executors
import java.util.concurrent.ScheduledExecutorService
import java.util.concurrent.TimeUnit

internal class TaskScheduler(private val tasks: List<ScheduledTask>) {
    private companion object {
        private const val STOP_TIMEOUT_SECONDS = 30L
        private val logger = logger<TaskScheduler>()
    }

    private lateinit var delegate: ScheduledExecutorService

    fun shutdown() {
        logger.info("Stopping the task scheduler (${STOP_TIMEOUT_SECONDS}s timeout)")
        shutdownDelegateIfNeeded()
    }

    suspend fun startOrReload() {
        shutdownDelegateIfNeeded()
        delegate = Executors.newScheduledThreadPool(4)

        tasks.forEach { task ->
            val (enabled, unit, interval, initialDelay) = task.schedule()
            if (!enabled) return@forEach

            delegate.scheduleAtFixedRate(
                buildTaskProxy(task),
                initialDelay.toLong(),
                interval.toLong(),
                unit,
            )

            task.onScheduleStarted()
        }
    }

    private fun shutdownDelegateIfNeeded() {
        if (!::delegate.isInitialized)
            return

        delegate.shutdown()

        try {
            delegate.awaitTermination(STOP_TIMEOUT_SECONDS, TimeUnit.SECONDS)
        } catch (ex: Exception) {
            logger.warn("Task scheduler graceful shutdown failed", ex)
        }

        logger.info("Task scheduler stopped")
    }

    private fun buildTaskProxy(task: ScheduledTask): Runnable =
        Runnable {
            try {
                runBlocking { task.run() }
            } catch (ex: Exception) {
                logger.warn("Scheduled task failed with an exception", ex)
            }
        }
}
