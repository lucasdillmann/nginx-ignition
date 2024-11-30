package br.com.dillmann.nginxignition.core.integration.command

import br.com.dillmann.nginxignition.core.common.pagination.Page

fun interface ListIntegrationOptionsCommand {
    data class Output(
        val id: String,
        val name: String,
    )

    suspend fun getIntegrationOptions(integrationId: String, pageNumber: Int, pageSize: Int): Page<Output>
}
