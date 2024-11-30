package br.com.dillmann.nginxignition.integration.truenas

import br.com.dillmann.nginxignition.core.common.dynamicfield.DynamicField
import br.com.dillmann.nginxignition.core.common.pagination.Page
import br.com.dillmann.nginxignition.core.integration.IntegrationAdapter

class TrueNasIntegrationAdapter: IntegrationAdapter {
    override val id = "TRUENAS_SCALE"
    override val imageId = null // TODO
    override val name = "TrueNAS Scale"
    override val description = "Lorem ipsum dolor sit amet" // TODO
    override val configurationFields = listOf(
        DynamicField(
            id = "url",
            description = "URL",
            priority = 1,
            required = true,
            type = DynamicField.Type.SINGLE_LINE_TEXT,
        ),
        DynamicField(
            id = "apiToken",
            description = "API token",
            priority = 2,
            required = true,
            sensitive = true,
            type = DynamicField.Type.SINGLE_LINE_TEXT,
        ),
    )

    override suspend fun getAvailableOptions(
        parameters: Map<String, Any?>,
        pageNumber: Int,
        pageSize: Int,
    ): Page<IntegrationAdapter.Option> =
        // TODO: Fetch actual values (remove mocked response)
        Page.of(
            IntegrationAdapter.Option(
                id = "1",
                name = "Lorem",
            ),
            IntegrationAdapter.Option(
                id = "2",
                name = "Ipsum",
            )
        )

    override suspend fun getOptionProxyUrl(id: String, parameters: Map<String, Any?>): String =
        // TODO: Fetch actual values (remove mocked response)
        "https://dillmann.dev"
}
