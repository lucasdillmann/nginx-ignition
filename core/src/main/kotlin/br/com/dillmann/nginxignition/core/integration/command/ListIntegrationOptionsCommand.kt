package br.com.dillmann.nginxignition.core.integration.command

import br.com.dillmann.nginxignition.core.common.pagination.Page
import br.com.dillmann.nginxignition.core.integration.model.IntegrationOption

fun interface ListIntegrationOptionsCommand {
    suspend fun getIntegrationOptions(
        integrationId: String,
        pageNumber: Int,
        pageSize: Int,
        searchTerms: String?,
    ): Page<IntegrationOption>
}
