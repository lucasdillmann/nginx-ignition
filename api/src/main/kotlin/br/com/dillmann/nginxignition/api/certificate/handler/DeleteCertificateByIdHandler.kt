package br.com.dillmann.nginxignition.api.certificate.handler

import br.com.dillmann.nginxignition.core.certificate.command.DeleteCertificateCommand
import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.handler.UuidAwareRequestHandler
import java.util.*

internal class DeleteCertificateByIdHandler(
    private val deleteCommand: DeleteCertificateCommand,
): UuidAwareRequestHandler {
    override suspend fun handle(call: ApiCall, id: UUID) {
        deleteCommand.deleteById(id)
        call.respond(HttpStatus.NO_CONTENT)
    }
}
