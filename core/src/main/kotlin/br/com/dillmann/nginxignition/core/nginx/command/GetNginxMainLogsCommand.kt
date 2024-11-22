package br.com.dillmann.nginxignition.core.nginx.command

interface GetNginxMainLogsCommand {
    suspend fun getMainLogs(lines: Int): List<String>
}
