package br.com.dillmann.nginxignition.core.nginx.command

fun interface GetNginxMainLogsCommand {
    suspend fun getMainLogs(lines: Int): List<String>
}
