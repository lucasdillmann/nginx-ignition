package br.com.dillmann.nginxignition.certificate.acme.dns.cloudflare

import br.com.dillmann.nginxignition.certificate.acme.dns.DnsProvider
import br.com.dillmann.nginxignition.certificate.acme.letsencrypt.LetsEncryptDynamicFields
import kotlinx.coroutines.delay

internal class CloudflareDnsProvider: DnsProvider {
    companion object {
        private const val PROPAGATION_WAIT_TIME_MILLIS = 10_000L
        const val ID = "CLOUDFLARE_DNS"
    }

    override val id = ID

    override suspend fun writeChallengeRecords(
        records: List<DnsProvider.ChallengeRecord>,
        dynamicFields: Map<String, Any?>,
    ) {
        val apiToken = dynamicFields[LetsEncryptDynamicFields.CLOUDFLARE_API_TOKEN.id] as String
        val client = CloudflareApiClient(apiToken)
        val zones = client.listZones()

        records
            .groupBy { findDnsZone(zones, it) }
            .map { (id, newRecords) -> Triple(id, client.listRecords(id), newRecords) }
            .forEach { (id, currentRecords, newRecords) -> updateRecords(client, id, currentRecords, newRecords) }

        // Cloudflare doesn't provide a way to wait for or check if the DNS changes where propagated/applied (the API
        // call returns before that). This is ugly, but it is the only way to for now to "wait" for changes to happen.
        delay(PROPAGATION_WAIT_TIME_MILLIS)
    }

    private fun findDnsZone(zones: List<CloudflareApi.Zone>, record: DnsProvider.ChallengeRecord): String =
        zones
            .filter { record.domainName.endsWith(it.name) }
            .maxByOrNull { it.name.length }
            ?.id
            ?: error("No DNS zone found for ${record.domainName}")

    private fun updateRecords(
        client: CloudflareApiClient,
        zoneId: String,
        currentRecords: List<CloudflareApi.DnsRecord>,
        newRecords: List<DnsProvider.ChallengeRecord>,
    ) {
        val currentRecordsPool = currentRecords.toMutableList()
        val (posts, puts) = newRecords.map { buildChange(currentRecordsPool, it) }.partition { it.id == null }
        val changes = CloudflareApi.BatchDnsRecordChange(posts, puts)
        client.batchRecordChange(zoneId, changes)
    }

    private fun buildChange(
        currentRecordsPool: MutableList<CloudflareApi.DnsRecord>,
        newRecord: DnsProvider.ChallengeRecord,
    ): CloudflareApi.DnsRecord {
        val currentRecord = currentRecordsPool.firstOrNull { it.name == newRecord.domainName && it.type == "TXT" }
        if (currentRecord != null)
            currentRecordsPool -= currentRecord

        return CloudflareApi.DnsRecord(
            id = currentRecord?.id,
            name = newRecord.domainName,
            content = "\"${newRecord.token}\"",
            type = "TXT",
        )
    }
}
