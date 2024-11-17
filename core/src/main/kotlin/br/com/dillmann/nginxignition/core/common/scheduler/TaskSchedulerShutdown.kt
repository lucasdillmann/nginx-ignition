package br.com.dillmann.nginxignition.core.common.scheduler

import br.com.dillmann.nginxignition.core.common.lifecycle.ShutdownCommand

internal class TaskSchedulerShutdown: ShutdownCommand {
    override val priority = 100

    override suspend fun execute() {
        TaskScheduler.shutdown()
    }
}
