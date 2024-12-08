package br.com.dillmann.nginxignition.core.common.scheduler.lifecycle

import br.com.dillmann.nginxignition.core.common.lifecycle.ShutdownCommand
import br.com.dillmann.nginxignition.core.common.scheduler.TaskScheduler

internal class TaskSchedulerShutdown(private val scheduler: TaskScheduler): ShutdownCommand {
    override suspend fun execute() {
        scheduler.shutdown()
    }
}
