package br.com.dillmann.nginxignition.core.common.scheduler.lifecycle

import br.com.dillmann.nginxignition.core.common.lifecycle.StartupCommand
import br.com.dillmann.nginxignition.core.common.scheduler.TaskScheduler

internal class TaskSchedulerStartup(private val scheduler: TaskScheduler): StartupCommand {
    override suspend fun execute() {
        scheduler.startOrReload()
    }
}
