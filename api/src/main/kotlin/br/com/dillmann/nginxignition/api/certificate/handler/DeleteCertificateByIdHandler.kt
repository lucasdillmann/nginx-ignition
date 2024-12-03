package br.com.dillmann.nginxignition.api.certificate.handler

import br.com.dillmann.nginxignition.core.certificate.command.DeleteCertificateCommand
import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.handler.UuidAwareRequestHandler
import br.com.dillmann.nginxignition.api.common.request.respond
import java.util.*

internal class DeleteCertificateByIdHandler(
    private val deleteCommand: DeleteCertificateCommand,
): UuidAwareRequestHandler {
    override suspend fun handle(call: ApiCall, id: UUID) {
        val (deleted, reason) = deleteCommand.deleteById(id)
        if (!deleted) {
            val payload = mapOf("message" to reason)
            call.respond(HttpStatus.PRECONDITION_FAILED, payload)
            return
        }

        call.respond(HttpStatus.NO_CONTENT)
    }
}
