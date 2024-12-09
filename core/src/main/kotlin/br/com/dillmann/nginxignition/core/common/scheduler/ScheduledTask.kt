package br.com.dillmann.nginxignition.core.common.scheduler

import java.util.concurrent.TimeUnit

interface ScheduledTask {
    data class Schedule(
        val enabled: Boolean,
        val unit: TimeUnit,
        val interval: Int,
        val initialDelay: Int,
    )

    suspend fun run()
    suspend fun schedule(): Schedule
    suspend fun onScheduleStarted()
}
