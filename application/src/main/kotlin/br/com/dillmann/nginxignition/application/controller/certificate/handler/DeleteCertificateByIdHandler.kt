package br.com.dillmann.nginxignition.application.controller.certificate.handler

import br.com.dillmann.nginxignition.application.common.routing.template.IdAwareRequestHandler
import br.com.dillmann.nginxignition.core.certificate.command.DeleteCertificateCommand
import io.ktor.http.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import java.util.*

class DeleteCertificateByIdHandler(
    private val deleteCommand: DeleteCertificateCommand,
): IdAwareRequestHandler {
    override suspend fun handle(call: RoutingCall, id: UUID) {
        deleteCommand.deleteById(id)
        call.respond(HttpStatusCode.NoContent)
    }
}
