package br.com.dillmann.nginxignition.core.user.command

import br.com.dillmann.nginxignition.core.common.pagination.Page
import br.com.dillmann.nginxignition.core.user.User

fun interface ListUserCommand {
    suspend fun list(pageSize: Int, pageNumber: Int): Page<User>
}
