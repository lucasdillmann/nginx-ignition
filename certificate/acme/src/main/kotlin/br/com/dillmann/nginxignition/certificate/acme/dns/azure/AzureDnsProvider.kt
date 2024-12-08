package br.com.dillmann.nginxignition.certificate.acme.dns.azure

import br.com.dillmann.nginxignition.certificate.acme.dns.DnsProvider
import br.com.dillmann.nginxignition.certificate.acme.letsencrypt.LetsEncryptDynamicFields
import com.azure.core.management.AzureEnvironment
import com.azure.core.management.profile.AzureProfile
import com.azure.identity.ClientSecretCredentialBuilder
import com.azure.resourcemanager.dns.DnsZoneManager
import com.azure.resourcemanager.dns.models.DnsZone
import com.azure.resourcemanager.dns.models.TxtRecord
import com.azure.resourcemanager.dns.models.TxtRecordSet
import kotlinx.coroutines.delay

internal class AzureDnsProvider: DnsProvider {
    enum class Environments(val description: String, val environment: AzureEnvironment) {
        DEFAULT("Azure (default)", AzureEnvironment.AZURE),
        CHINA("China", AzureEnvironment.AZURE_CHINA),
        US_GOVERNMENT("US Government", AzureEnvironment.AZURE_US_GOVERNMENT),
    }

    companion object {
        private const val PROPAGATION_WAIT_TIME_MILLIS = 10_000L
        const val ID = "AZURE_DNS"
    }

    override val id = ID

    override suspend fun writeChallengeRecords(
        records: List<DnsProvider.ChallengeRecord>,
        dynamicFields: Map<String, Any?>,
    ) {
        val manager = buildDnsManager(dynamicFields)
        val zones = manager.zones().list().toList()
        records
            .groupBy { findDnsZone(zones, it) }
            .map { (zone, newRecords) -> Triple(zone, zone.txtRecordSets().list().toList(), newRecords) }
            .forEach { (zone, currentRecords, newRecords) -> upsertRecords(zone, currentRecords, newRecords) }

        // Azure doesn't provide a way to wait for or check if the DNS changes where propagated/applied (the API
        // call returns before that). This is ugly, but it is the only way to for now to "wait" for changes to happen.
        delay(PROPAGATION_WAIT_TIME_MILLIS)
    }

    private fun upsertRecords(
        zone: DnsZone,
        currentRecords: List<TxtRecordSet>,
        newRecords: List<DnsProvider.ChallengeRecord>,
    ) {
        newRecords
            .groupBy { it.domainName }
            .forEach { (domainName, records) ->
                val normalizedDomainName = domainName.removeSuffix(".${zone.name()}")
                val values = records.map { it.token }
                val update = zone.update()
                val currentRecord = currentRecords.find { it.name() == normalizedDomainName }

                if (currentRecord != null) {
                    val updateOperation = update.updateTxtRecordSet(normalizedDomainName)
                    currentRecord.records().flatMap(TxtRecord::value).forEach(updateOperation::withoutText)
                    values.forEach(updateOperation::withText)
                } else {
                    val createOperation = update.defineTxtRecordSet(normalizedDomainName)
                    values.forEach(createOperation::withText)
                }

                update.apply()
            }
    }

    private fun findDnsZone(zones: List<DnsZone>, record: DnsProvider.ChallengeRecord): DnsZone =
        zones
            .filter { record.domainName.endsWith(it.name()) }
            .maxByOrNull { it.name().length }
            ?: error("No DNS zone found for ${record.domainName}")

    private fun buildDnsManager(dynamicFields: Map<String, Any?>): DnsZoneManager {
        val tenantId = dynamicFields[LetsEncryptDynamicFields.AZURE_TENANT_ID.id] as String
        val subscriptionId = dynamicFields[LetsEncryptDynamicFields.AZURE_SUBSCRIPTION_ID.id] as String
        val clientId = dynamicFields[LetsEncryptDynamicFields.AZURE_CLIENT_ID.id] as String
        val environmentId = dynamicFields[LetsEncryptDynamicFields.AZURE_ENVIRONMENT.id] as String
        val environment = Environments.valueOf(environmentId).environment
        val clientSecret = dynamicFields[LetsEncryptDynamicFields.AZURE_CLIENT_SECRET.id] as String
        val credentials = ClientSecretCredentialBuilder()
            .tenantId(tenantId)
            .clientId(clientId)
            .clientSecret(clientSecret)
            .build()
        val profile = AzureProfile(tenantId, subscriptionId, environment)
        return DnsZoneManager.authenticate(credentials, profile)
    }
}
