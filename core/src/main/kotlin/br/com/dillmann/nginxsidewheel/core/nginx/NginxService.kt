package br.com.dillmann.nginxsidewheel.core.nginx

import br.com.dillmann.nginxsidewheel.core.nginx.command.ReloadNginxCommand
import br.com.dillmann.nginxsidewheel.core.nginx.command.StartNginxCommand
import br.com.dillmann.nginxsidewheel.core.nginx.command.StopNginxCommand

internal class NginxService(
    private val configurationFiles: NginxConfigurationFiles,
    private val processManager: NginxProcessManager,
    private val semaphore: NginxSemaphore,
): ReloadNginxCommand, StartNginxCommand, StopNginxCommand {
    override suspend fun reload() {
        semaphore.changeState(NginxSemaphore.State.RUNNING) {
            configurationFiles.generate()
            processManager.sendReloadSignal()
        }
    }

    override suspend fun start() {
        if (semaphore.currentState == NginxSemaphore.State.RUNNING)
            return

        semaphore.changeState(NginxSemaphore.State.RUNNING) {
            configurationFiles.generate()
            processManager.start()
        }
    }

    override suspend fun stop() {
        if (semaphore.currentState == NginxSemaphore.State.STOPPED)
            return

        semaphore.changeState(NginxSemaphore.State.STOPPED) {
            processManager.sendStopSignal()
        }
    }
}
