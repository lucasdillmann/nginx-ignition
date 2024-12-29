package br.com.dillmann.nginxignition.core.nginx

import br.com.dillmann.nginxignition.core.common.log.logger
import br.com.dillmann.nginxignition.core.common.configuration.ConfigurationProvider
import br.com.dillmann.nginxignition.core.nginx.exception.NginxCommandException
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import java.io.File

internal class NginxProcessManager(configurationProvider: ConfigurationProvider) {
    private val configurationProvider = configurationProvider.withPrefix("nginx-ignition.nginx")
    private val logger = logger<NginxProcessManager>()

    suspend fun sendReloadSignal() {
        runCommand( "-s", "reload")
        logger.info("nginx reloaded")
    }

    suspend fun sendReopenSignal() {
        logger.info("Signaling nginx for log file reopen")
        runCommand( "-s", "reopen")
    }

    suspend fun sendStopSignal() {
        runCommand( "-s", "stop")
        logger.info("nginx stopped")
    }

    suspend fun start() {
        runCommand()
        logger.info("nginx started")
    }

    suspend fun currentPid(): Long? {
        val folder = configurationProvider.get("config-directory")

        return withContext(Dispatchers.IO) {
            val pidFile = File(folder, "nginx.pid")
            if (!pidFile.exists() || !pidFile.isFile)
                return@withContext null

            pidFile.readLines().first().toLong().takeIf { isPidAlive(it) }
        }
    }

    private suspend fun isPidAlive(pid: Long): Boolean =
        withContext(Dispatchers.IO) {
            ProcessBuilder()
                .command("kill", "-0", pid.toString())
                .start()
                .waitFor() == 0
        }

    private suspend fun runCommand(vararg extraArguments: String) {
        val binaryPath = configurationProvider.get("binary-path")
        val configDirectory = configurationProvider.get("config-directory")
        val arguments = arrayOf(binaryPath, "-c", "$configDirectory/config/nginx.conf", *extraArguments)
        val command = withContext(Dispatchers.IO) { ProcessBuilder().command(*arguments).start() }

        val exitCode = withContext(Dispatchers.IO) { command.waitFor() }
        if(exitCode != 0) {
            val output = withContext(Dispatchers.IO) { command.errorStream.readAllBytes().toString(Charsets.UTF_8) }
            throw NginxCommandException(
                command = arguments.joinToString(separator = " "),
                exitCode = exitCode,
                output = output,
            )
        }
    }
}
