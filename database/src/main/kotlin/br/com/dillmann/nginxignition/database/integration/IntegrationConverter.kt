package br.com.dillmann.nginxignition.database.integration

import br.com.dillmann.nginxignition.core.integration.Integration
import br.com.dillmann.nginxignition.database.common.json.toJsonObject
import br.com.dillmann.nginxignition.database.integration.mapping.IntegrationTable
import org.jetbrains.exposed.sql.ResultRow
import org.jetbrains.exposed.sql.statements.UpsertStatement

internal class IntegrationConverter {
    fun apply(integration: Integration, scope: UpsertStatement<out Any>) {
        with(IntegrationTable) {
            scope[id] = integration.id
            scope[enabled] = integration.enabled
            scope[parameters] = integration.parameters.toJsonObject().toString()
        }
    }

    fun toDomainModel(integration: ResultRow) =
        with(IntegrationTable) {
            Integration(
                id = integration[id],
                enabled = integration[enabled],
                parameters = integration[parameters].toJsonObject().toMap(),
            )
        }
}
