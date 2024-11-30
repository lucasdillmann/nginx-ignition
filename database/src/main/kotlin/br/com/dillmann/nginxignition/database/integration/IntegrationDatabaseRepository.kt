package br.com.dillmann.nginxignition.database.integration

import br.com.dillmann.nginxignition.core.integration.Integration
import br.com.dillmann.nginxignition.core.integration.IntegrationRepository
import br.com.dillmann.nginxignition.database.common.transaction.coTransaction
import br.com.dillmann.nginxignition.database.integration.mapping.IntegrationTable
import org.jetbrains.exposed.sql.upsert

internal class IntegrationDatabaseRepository(private val converter: IntegrationConverter): IntegrationRepository {
    override suspend fun findById(id: String): Integration? =
        coTransaction {
            val user = IntegrationTable
                .select(IntegrationTable.fields)
                .where { IntegrationTable.id eq id }
                .firstOrNull()
                ?: return@coTransaction null

            converter.toDomainModel(user)
        }

    override suspend fun save(integration: Integration) {
        coTransaction {
            IntegrationTable.upsert { converter.apply(integration, it) }
        }
    }
}
