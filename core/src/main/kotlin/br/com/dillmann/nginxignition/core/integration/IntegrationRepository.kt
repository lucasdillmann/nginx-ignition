package br.com.dillmann.nginxignition.core.integration

interface IntegrationRepository {
    suspend fun findById(id: String): Integration?
    suspend fun save(integration: Integration)
}
