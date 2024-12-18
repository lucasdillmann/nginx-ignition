package br.com.dillmann.nginxignition.core.accesslist.command

import br.com.dillmann.nginxignition.core.accesslist.AccessList

fun interface SaveAccessListByCommand {
    suspend fun save(accessList: AccessList)
}
