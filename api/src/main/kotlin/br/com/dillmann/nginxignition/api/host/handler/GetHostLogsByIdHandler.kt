package br.com.dillmann.nginxignition.api.host.handler

import br.com.dillmann.nginxignition.api.common.logs.LogRequestHandler
import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.handler.UuidAwareRequestHandler
import br.com.dillmann.nginxignition.api.common.request.respond
import br.com.dillmann.nginxignition.core.host.command.HostExistsByIdCommand
import br.com.dillmann.nginxignition.core.nginx.command.GetNginxHostLogsCommand
import java.util.*

internal class GetHostLogsByIdHandler(
    private val existsCommand: HostExistsByIdCommand,
    private val getLogsCommand: GetNginxHostLogsCommand,
): UuidAwareRequestHandler, LogRequestHandler {
    private companion object {
        private val ALLOWED_QUALIFIERS = listOf("access", "error")
    }

    override suspend fun handle(call: ApiCall, id: UUID) {
        val exists = existsCommand.existsById(id)
        val qualifier = call.pathVariables()["qualifier"] ?: ""
        if (!exists || qualifier !in ALLOWED_QUALIFIERS) {
            call.respond(HttpStatus.NOT_FOUND)
            return
        }

        val lineCount = parseLinesAmount(call)
        if (lineCount == null) {
            sendLinesAmountErrorResponse(call)
            return
        }

        val logs = getLogsCommand.getHostLogs(id, qualifier, lineCount)
        call.respond(HttpStatus.OK, logs)
    }
}
