package br.com.dillmann.nginxignition.core.integration

import br.com.dillmann.nginxignition.core.common.dynamicfield.DynamicField
import br.com.dillmann.nginxignition.core.common.pagination.Page

interface IntegrationAdapter {
    data class Option(
        val id: String,
        val name: String,
    )

    val id: String
    val name: String
    val description: String
    val configurationFields: List<DynamicField>

    suspend fun getAvailableOptions(parameters: Map<String, Any?>, pageNumber: Int, pageSize: Int): Page<Option>
    suspend fun getAvailableOptionById(parameters: Map<String, Any?>, id: String): Option?
    suspend fun getOptionProxyUrl(id: String, parameters: Map<String, Any?>): String
}
