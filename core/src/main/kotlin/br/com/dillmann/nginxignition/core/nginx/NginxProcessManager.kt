package br.com.dillmann.nginxignition.core.nginx

import br.com.dillmann.nginxignition.core.common.log.logger
import br.com.dillmann.nginxignition.core.common.configuration.ConfigurationProvider
import br.com.dillmann.nginxignition.core.nginx.exception.NginxCommandException
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext

internal class NginxProcessManager(configurationProvider: ConfigurationProvider) {
    private val configurationProvider = configurationProvider.withPrefix("nginx-ignition.nginx")
    private val logger = logger<NginxProcessManager>()

    suspend fun sendReloadSignal() {
        logger.info("Reloading nginx configuration")
        runCommand( "-s", "reload")
    }

    suspend fun sendStopSignal() {
        logger.info("Stopping nginx")
        runCommand( "-s", "stop")
    }

    suspend fun start() {
        logger.info("Starting nginx")
        runCommand()
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
