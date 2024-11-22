package br.com.dillmann.nginxignition.core.nginx

import br.com.dillmann.nginxignition.core.nginx.command.*
import br.com.dillmann.nginxignition.core.nginx.configuration.NginxConfigurationFacade
import br.com.dillmann.nginxignition.core.nginx.log.NginxLogReader
import java.util.*

internal class NginxService(
    private val nginxConfiguration: NginxConfigurationFacade,
    private val processManager: NginxProcessManager,
    private val semaphore: NginxSemaphore,
    private val logReader: NginxLogReader,
): ReloadNginxCommand, StartNginxCommand, StopNginxCommand, GetStatusNginxCommand,
   GetNginxHostLogsCommand, GetNginxMainLogsCommand {

    override suspend fun reload() {
        semaphore.changeState(NginxSemaphore.State.RUNNING) {
            nginxConfiguration.replaceConfigurationFiles()
            processManager.sendReloadSignal()
        }
    }

    override suspend fun start() {
        if (semaphore.currentState == NginxSemaphore.State.RUNNING)
            return

        semaphore.changeState(NginxSemaphore.State.RUNNING) {
            nginxConfiguration.replaceConfigurationFiles()
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

    override suspend fun isRunning() =
        semaphore.currentState == NginxSemaphore.State.RUNNING

    override suspend fun getHostLogs(hostId: UUID, qualifier: String, lines: Int): List<String> =
        logReader.read("host-$hostId.$qualifier.log", lines)

    override suspend fun getMainLogs(lines: Int): List<String> =
        logReader.read("main.log", lines)
}
