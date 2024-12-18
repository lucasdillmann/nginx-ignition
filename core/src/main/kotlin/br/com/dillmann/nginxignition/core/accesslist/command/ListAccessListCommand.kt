package br.com.dillmann.nginxignition.core.accesslist.command

import br.com.dillmann.nginxignition.core.accesslist.AccessList
import br.com.dillmann.nginxignition.core.common.pagination.Page

fun interface ListAccessListCommand {
    suspend fun getPage(pageSize: Int, pageNumber: Int, searchTerms: String?): Page<AccessList>
}
