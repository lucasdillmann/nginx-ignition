package br.com.dillmann.nginxignition.core.accesslist

import br.com.dillmann.nginxignition.core.accesslist.command.DeleteAccessListByIdCommand
import br.com.dillmann.nginxignition.core.accesslist.command.GetAccessListByIdCommand
import br.com.dillmann.nginxignition.core.accesslist.command.ListAccessListCommand
import br.com.dillmann.nginxignition.core.accesslist.command.SaveAccessListByCommand
import br.com.dillmann.nginxignition.core.common.pagination.Page
import java.util.*

internal class AccessListService(
    private val accessListRepository: AccessListRepository,
    private val validator: AccessListValidator,
): DeleteAccessListByIdCommand, GetAccessListByIdCommand, ListAccessListCommand, SaveAccessListByCommand {
    override suspend fun deleteById(id: UUID) {
        // TODO: Check if isn't in use before deleting
        accessListRepository.deleteById(id)
    }

    override suspend fun getById(id: UUID): AccessList? =
        accessListRepository.findById(id)

    override suspend fun getPage(pageNumber: Int, pageSize: Int, searchTerms: String?): Page<AccessList> =
        accessListRepository.findPage(pageNumber, pageSize, searchTerms)

    override suspend fun save(accessList: AccessList) {
        validator.validate(accessList)
        accessListRepository.save(accessList)
    }
}
