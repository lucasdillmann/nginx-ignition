package br.com.dillmann.nginxignition.core.common

import br.com.dillmann.nginxignition.core.common.lifecycle.LifecycleManager
import br.com.dillmann.nginxignition.core.common.lifecycle.ShutdownCommand
import br.com.dillmann.nginxignition.core.common.lifecycle.StartupCommand
import br.com.dillmann.nginxignition.core.common.scheduler.TaskScheduler
import br.com.dillmann.nginxignition.core.common.scheduler.lifecycle.TaskSchedulerShutdown
import br.com.dillmann.nginxignition.core.common.scheduler.lifecycle.TaskSchedulerStartup
import org.koin.core.module.Module
import org.koin.dsl.bind

internal fun Module.commonBeans() {
    single { TaskScheduler(getAll()) }
    single { TaskSchedulerShutdown(get()) } bind ShutdownCommand::class
    single { TaskSchedulerStartup(get()) } bind StartupCommand::class
    single { LifecycleManager(getAll(), getAll()) }
}
