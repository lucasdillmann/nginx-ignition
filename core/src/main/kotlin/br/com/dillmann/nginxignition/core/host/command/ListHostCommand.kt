package br.com.dillmann.nginxignition.core.host.command

import br.com.dillmann.nginxignition.core.common.pagination.Page
import br.com.dillmann.nginxignition.core.host.Host

interface ListHostCommand {
    suspend fun list(pageSize: Int, pageNumber: Int): Page<Host>
}
