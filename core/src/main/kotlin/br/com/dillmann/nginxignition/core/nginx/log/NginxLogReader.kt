package br.com.dillmann.nginxignition.core.nginx.log

import br.com.dillmann.nginxignition.core.common.configuration.ConfigurationProvider
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import org.apache.commons.io.input.ReversedLinesFileReader
import java.io.File

internal class NginxLogReader(private val configurationProvider: ConfigurationProvider) {
    suspend fun read(fileName: String, tailSize: Int): List<String> {
        val basePath = configurationProvider.get("nginx-ignition.nginx.config-directory")
        val normalizedPath = basePath.dropLastWhile { it == '/' } + "/logs"

        val file = File(normalizedPath, fileName)
        if (!file.exists() || file.isDirectory || file.parent != normalizedPath)
            return emptyList()

        return withContext(Dispatchers.IO) {
            ReversedLinesFileReader
                .builder()
                .setFile(file)
                .setCharset(Charsets.UTF_8)
                .get()
                .use { it.readLines(tailSize) }
        }
    }
}
