package br.com.dillmann.nginxignition.application.controller.certificate.handler

import br.com.dillmann.nginxignition.application.common.routing.RequestHandler
import br.com.dillmann.nginxignition.core.certificate.command.DeleteCertificateCommand
import io.ktor.http.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import java.util.*

class DeleteCertificateByIdHandler(
    private val deleteCommand: DeleteCertificateCommand,
): RequestHandler {
    override suspend fun handle(call: RoutingCall) {
        val certificateId = runCatching { call.request.pathVariables["id"].let(UUID::fromString) }.getOrNull()
        if (certificateId == null) {
            call.respond(HttpStatusCode.BadRequest)
            return
        }

        deleteCommand.deleteById(certificateId)
        call.respond(HttpStatusCode.NoContent)
    }
}
