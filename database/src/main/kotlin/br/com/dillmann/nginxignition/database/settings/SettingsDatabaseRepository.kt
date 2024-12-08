package br.com.dillmann.nginxignition.database.settings

import br.com.dillmann.nginxignition.core.settings.Settings
import br.com.dillmann.nginxignition.core.settings.SettingsRepository
import br.com.dillmann.nginxignition.database.common.transaction.coTransaction
import br.com.dillmann.nginxignition.database.settings.mapping.SettingsCertificateAutoRenewTable
import br.com.dillmann.nginxignition.database.settings.mapping.SettingsGlobalBindingTable
import br.com.dillmann.nginxignition.database.settings.mapping.SettingsLogRotationTable
import br.com.dillmann.nginxignition.database.settings.mapping.SettingsNginxTable
import org.jetbrains.exposed.sql.deleteAll
import org.jetbrains.exposed.sql.insert
import org.jetbrains.exposed.sql.selectAll
import org.jetbrains.exposed.sql.update

internal class SettingsDatabaseRepository(private val converter: SettingsConverter): SettingsRepository {
    override suspend fun save(settings: Settings) {
        coTransaction {
            SettingsCertificateAutoRenewTable.update { converter.apply(settings.certificateAutoRenew, it) }
            SettingsLogRotationTable.update { converter.apply(settings.logRotation, it) }
            SettingsNginxTable.update { converter.apply(settings.nginx, it) }

            SettingsGlobalBindingTable.deleteAll()
            settings.globalBindings.forEach { binding ->
                SettingsGlobalBindingTable.insert { converter.apply(binding, it) }
            }
        }
    }

    override suspend fun get(): Settings =
        coTransaction {
            val certificateAutoRenew =
                SettingsCertificateAutoRenewTable.selectAll().first().let(converter::toCertificateAutoRenewSettings)
            val logRotation =
                SettingsLogRotationTable.selectAll().first().let(converter::toLogRotationSettings)
            val nginx =
                SettingsNginxTable.selectAll().first().let(converter::toNginxSettings)
            val bindings =
                SettingsGlobalBindingTable.selectAll().map(converter::toGlobalHostBinding)

            Settings(nginx, logRotation, certificateAutoRenew, bindings)
        }
}
