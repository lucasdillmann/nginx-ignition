package br.com.dillmann.nginxignition.core.common.scheduler

import java.util.concurrent.TimeUnit

interface ScheduledTask {
    data class Schedule(
        val unit: TimeUnit,
        val interval: Long,
        val initialDelay: Long,
    )

    fun schedule(): Schedule

    fun onScheduleStarted() {
        // no-op by default
    }

    suspend fun run()
}
