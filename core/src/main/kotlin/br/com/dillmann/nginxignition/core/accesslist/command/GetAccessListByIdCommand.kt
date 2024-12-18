package br.com.dillmann.nginxignition.core.accesslist.command

import br.com.dillmann.nginxignition.core.accesslist.AccessList
import java.util.*

fun interface GetAccessListByIdCommand {
    suspend fun getById(id: UUID): AccessList?
}
