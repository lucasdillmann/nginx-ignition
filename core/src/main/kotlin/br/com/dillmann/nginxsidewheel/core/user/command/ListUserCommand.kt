package br.com.dillmann.nginxsidewheel.core.user.command

import br.com.dillmann.nginxsidewheel.core.common.pagination.Page
import br.com.dillmann.nginxsidewheel.core.user.User

interface ListUserCommand {
    suspend fun list(pageSize: Int, pageNumber: Int): Page<User>
}
