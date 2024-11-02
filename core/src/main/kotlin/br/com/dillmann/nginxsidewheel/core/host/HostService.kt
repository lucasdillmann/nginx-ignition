package br.com.dillmann.nginxsidewheel.core.host

import br.com.dillmann.nginxsidewheel.core.common.pagination.Page
import br.com.dillmann.nginxsidewheel.core.host.command.DeleteHostCommand
import br.com.dillmann.nginxsidewheel.core.host.command.GetHostCommand
import br.com.dillmann.nginxsidewheel.core.host.command.ListHostCommand
import br.com.dillmann.nginxsidewheel.core.host.command.SaveHostCommand
import java.util.*

internal class HostService(
    private val repository: HostRepository,
): SaveHostCommand, DeleteHostCommand, ListHostCommand, GetHostCommand {
    override suspend fun save(input: Host) {
        HostValidator.validate(input)
        repository.save(input)
    }

    override suspend fun deleteById(id: UUID) {
        repository.deleteById(id)
    }

    override suspend fun list(pageSize: Int, pageNumber: Int): Page<Host> =
        repository.findAll(pageSize, pageNumber)

    override suspend fun getById(id: UUID): Host? =
        repository.findById(id)
}
