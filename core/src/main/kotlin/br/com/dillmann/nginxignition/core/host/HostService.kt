package br.com.dillmann.nginxignition.core.host

import br.com.dillmann.nginxignition.core.common.pagination.Page
import br.com.dillmann.nginxignition.core.host.command.DeleteHostCommand
import br.com.dillmann.nginxignition.core.host.command.GetHostCommand
import br.com.dillmann.nginxignition.core.host.command.ListHostCommand
import br.com.dillmann.nginxignition.core.host.command.SaveHostCommand
import java.util.*

internal class HostService(
    private val repository: HostRepository,
    private val validator: HostValidator,
): SaveHostCommand, DeleteHostCommand, ListHostCommand, GetHostCommand {
    override suspend fun save(input: Host) {
        validator.validate(input)
        repository.save(input)
    }

    override suspend fun deleteById(id: UUID) {
        repository.deleteById(id)
    }

    override suspend fun list(pageSize: Int, pageNumber: Int): Page<Host> =
        repository.findPage(pageSize, pageNumber)

    override suspend fun getById(id: UUID): Host? =
        repository.findById(id)

    suspend fun getAll(): List<Host> =
        repository.findAll()
}
