package br.com.dillmann.nginxignition.core.nginx.log

import br.com.dillmann.nginxignition.core.common.configuration.ConfigurationProvider
import br.com.dillmann.nginxignition.core.common.log.logger
import br.com.dillmann.nginxignition.core.host.HostRepository
import br.com.dillmann.nginxignition.core.nginx.NginxProcessManager
import br.com.dillmann.nginxignition.core.settings.SettingsRepository
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import org.apache.commons.io.input.ReversedLinesFileReader
import java.io.File

internal class NginxLogRotator(
    private val configurationProvider: ConfigurationProvider,
    private val settingsRepository: SettingsRepository,
    private val hostRepository: HostRepository,
    private val processManager: NginxProcessManager,
) {
    private companion object {
        private val LOGGER = logger<NginxLogRotator>()
    }

    suspend fun rotate() {
        LOGGER.info("Starting log rotation")

        val basePath = configurationProvider.get("nginx-ignition.nginx.config-directory")
        val normalizedPath = basePath.dropLastWhile { it == '/' } + "/logs"
        val maximumLines = settingsRepository.get().logRotation.maximumLines

        rotate(maximumLines, normalizedPath, "main.log")
        val logFiles = hostRepository
            .findAllEnabled()
            .map { "host-${it.id}" }
            .flatMap { listOf("$it.access.log", "$it.error.log") } + "main.log"

        logFiles.forEach { rotate(maximumLines, normalizedPath, it) }
        processManager.sendReopenSignal()

        LOGGER.info("Log rotation finished with ${logFiles.size} files trimmed to up to $maximumLines lines")
    }

    private suspend fun rotate(maximumLines: Int, basePath: String, fileName: String) {
        try {
            val file = File(basePath, fileName)
            if (!file.exists() || !file.isFile) return

            val tail = readTail(file, maximumLines)
            if (tail.size < maximumLines)
                return

            val trimmedContent = tail.reversed().joinToString("\n").plus("\n")
            replaceContents(file, trimmedContent)
        } catch (ex: Exception) {
            LOGGER.warn("Unable to rotate log file {}", fileName, ex)
        }
    }

    private suspend fun replaceContents(file: File, contents: String) {
        withContext(Dispatchers.IO) {
            require(file.delete()) { "Unable to delete the file" }
            file.createNewFile()
            file.writer().use { it.write(contents) }
        }
    }

    private suspend fun readTail(file: File, size: Int): List<String> =
        withContext(Dispatchers.IO) {
            ReversedLinesFileReader
                .builder()
                .setFile(file)
                .setCharset(Charsets.UTF_8)
                .get()
                .use { it.readLines(size) }
        }
}
