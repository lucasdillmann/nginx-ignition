package br.com.dillmann.nginxignition.core.nginx.configuration

import br.com.dillmann.nginxignition.core.common.log.logger
import br.com.dillmann.nginxignition.core.common.configuration.ConfigurationProvider
import br.com.dillmann.nginxignition.core.host.HostService
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import java.io.File

internal class NginxConfigurationFacade(
    private val hostService: HostService,
    private val fileProviders: List<NginxConfigurationFileProvider>,
    configurationProvider: ConfigurationProvider,
) {
    private val logger = logger<NginxConfigurationFacade>()
    private val configurationProvider = configurationProvider.withPrefix("nginx-ignition.nginx")

    suspend fun replaceConfigurationFiles() {
        val hosts = hostService.getAll()
        logger.info("Rebuilding nginx configuration files for ${hosts.size} hosts")

        val basePath = configurationProvider.get("config-directory")
        val normalizedPath = basePath.dropLastWhile { it == '/' }
        createMissingFolders(normalizedPath)
        emptyConfigFolder(normalizedPath)

        val configFiles = fileProviders.flatMap { it.provide(normalizedPath, hosts) }
        withContext(Dispatchers.IO) {
            configFiles.forEach { (name, contents) ->
                val configFile = File("$normalizedPath/config", name)
                configFile.createNewFile()
                configFile.writeText(contents)
            }
        }
    }

    private suspend fun createMissingFolders(basePath: String) {
        withContext(Dispatchers.IO) {
            listOf("logs", "config").forEach { folderName ->
                val folder = File(basePath, folderName)
                if (folder.exists() && !folder.isDirectory)
                    folder.delete()
                if (!folder.exists())
                    folder.mkdirs()
            }
        }
    }

    private suspend fun emptyConfigFolder(basePath: String) {
        withContext(Dispatchers.IO) {
            File(basePath, "config")
                .listFiles()
                ?.forEach { file -> file.deleteRecursively() }
        }
    }
}
