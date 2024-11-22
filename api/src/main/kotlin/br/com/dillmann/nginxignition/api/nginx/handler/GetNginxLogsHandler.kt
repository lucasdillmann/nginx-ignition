package br.com.dillmann.nginxignition.api.nginx.handler

import br.com.dillmann.nginxignition.api.common.logs.LogRequestHandler
import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.handler.RequestHandler
import br.com.dillmann.nginxignition.api.common.request.respond
import br.com.dillmann.nginxignition.core.nginx.command.GetNginxMainLogsCommand

internal class GetNginxLogsHandler(
    private val getLogsCommand: GetNginxMainLogsCommand,
): RequestHandler, LogRequestHandler {
    override suspend fun handle(call: ApiCall) {
        val lineCount = parseLinesAmount(call)
        if (lineCount == null) {
            sendLinesAmountErrorResponse(call)
            return
        }

        val logs = getLogsCommand.getMainLogs(lineCount)
        call.respond(HttpStatus.OK, logs)
    }
}
