package br.com.dillmann.nginxignition.core.host

import br.com.dillmann.nginxignition.core.common.pagination.Page
import br.com.dillmann.nginxignition.core.host.command.*
import java.util.*

internal class HostService(
    private val repository: HostRepository,
    private val validator: HostValidator,
): SaveHostCommand, DeleteHostCommand, ListHostCommand, GetHostCommand, HostExistsByIdCommand {
    override suspend fun save(input: Host) {
        validator.validate(input)
        repository.save(input)
    }

    override suspend fun deleteById(id: UUID) {
        repository.deleteById(id)
    }

    override suspend fun list(pageSize: Int, pageNumber: Int, searchTerms: String?): Page<Host> =
        repository.findPage(pageSize, pageNumber, searchTerms)

    override suspend fun getById(id: UUID): Host? =
        repository.findById(id)

    suspend fun getAllEnabled(): List<Host> =
        repository.findAllEnabled()

    override suspend fun existsById(id: UUID) =
        repository.existsById(id)
}
