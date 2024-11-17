package br.com.dillmann.nginxignition.core.common

import br.com.dillmann.nginxignition.core.common.lifecycle.ShutdownCommand
import br.com.dillmann.nginxignition.core.common.scheduler.TaskSchedulerShutdown
import org.koin.core.module.Module
import org.koin.dsl.bind

internal fun Module.commonBeans() {
    single { TaskSchedulerShutdown() } bind ShutdownCommand::class
}
