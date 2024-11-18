package br.com.dillmann.nginxignition.core.common.scheduler

import br.com.dillmann.nginxignition.core.common.lifecycle.ShutdownCommand

internal class TaskSchedulerShutdown: ShutdownCommand {
    override suspend fun execute() {
        TaskScheduler.shutdown()
    }
}
