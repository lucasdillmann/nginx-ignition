package br.com.dillmann.nginxignition.database.certificate.lifecycle

import br.com.dillmann.nginxignition.core.common.lifecycle.StartupCommand
import br.com.dillmann.nginxignition.database.common.transaction.coTransaction
import br.com.dillmann.nginxignition.database.certificate.mapping.CertificateTable
import org.jetbrains.exposed.sql.SchemaUtils

class CertificateMigrations: StartupCommand {
    override val priority = 100

    override suspend fun execute() {
        coTransaction {
            SchemaUtils.createMissingTablesAndColumns(CertificateTable)
        }
    }
}
