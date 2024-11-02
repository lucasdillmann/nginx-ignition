package br.com.dillmann.nginxsidewheel.core.host.command

import br.com.dillmann.nginxsidewheel.core.common.pagination.Page
import br.com.dillmann.nginxsidewheel.core.host.Host

interface ListHostCommand {
    suspend fun list(pageSize: Int, pageNumber: Int): Page<Host>
}
