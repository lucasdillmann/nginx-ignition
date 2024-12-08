package br.com.dillmann.nginxignition.certificate.acme.dns.google

import br.com.dillmann.nginxignition.certificate.acme.dns.DnsProvider
import br.com.dillmann.nginxignition.certificate.acme.letsencrypt.LetsEncryptDynamicFields
import com.google.api.gax.paging.Page
import com.google.auth.oauth2.ServiceAccountCredentials
import com.google.cloud.dns.ChangeRequest
import com.google.cloud.dns.ChangeRequestInfo
import com.google.cloud.dns.Dns
import com.google.cloud.dns.DnsOptions
import com.google.cloud.dns.RecordSet
import com.google.cloud.dns.Zone
import java.lang.Thread.sleep
import java.util.concurrent.TimeUnit

internal class GoogleCloudDnsProvider : DnsProvider {
    companion object {
        private const val WAIT_FOR_COMPLETION_DELAY_MS = 100L
        private const val RECORD_TTL_SECONDS = 30
        private const val PAGE_SIZE = 100
        const val ID = "GOOGLE_CLOUD_DNS"
    }

    override val id = ID

    override suspend fun writeChallengeRecords(
        records: List<DnsProvider.ChallengeRecord>,
        dynamicFields: Map<String, Any?>,
    ) {
        val privateKey = dynamicFields[LetsEncryptDynamicFields.GOOGLE_CLOUD_PRIVATE_KEY.id] as String
        val credentials = ServiceAccountCredentials.fromStream(privateKey.byteInputStream())
        val options = DnsOptions.newBuilder().setCredentials(credentials).build()
        val client = DnsOptions.DefaultDnsFactory().create(options)

        val zones = fetchZones(client)
        records
            .groupBy { findDnsZone(zones, it) }
            .map { (zoneId, newRecords) -> Triple(zoneId, fetchTxtRecords(client, zoneId), newRecords) }
            .map { (zoneId, currentRecords, newRecords) -> updateRecords(client, zoneId, currentRecords, newRecords) }
            .forEach { (zoneId, changeId) -> waitForCompletion(client, zoneId, changeId) }
    }

    private fun waitForCompletion(client: Dns, zoneId: String, changeId: String) {
        do {
            val status = client.getChangeRequest(zoneId, changeId)
            if (status.status != ChangeRequestInfo.Status.DONE)
                sleep(WAIT_FOR_COMPLETION_DELAY_MS)
        } while (status.status != ChangeRequestInfo.Status.DONE)
    }

    private fun updateRecords(
        client: Dns,
        zoneId: String,
        currentRecords: List<RecordSet>,
        newRecords: List<DnsProvider.ChallengeRecord>,
    ): Pair<String, String> {
        val changeBuilder = ChangeRequest.newBuilder()
        newRecords
            .groupBy { it.domainName }
            .map { (domainName, records) ->
                val values = records.map { "\"${it.token}\"" }
                RecordSet
                    .newBuilder("$domainName.", RecordSet.Type.TXT)
                    .setRecords(values)
                    .setTtl(RECORD_TTL_SECONDS, TimeUnit.SECONDS)
                    .build()
            }
            .forEach { recordSet ->
                currentRecords.find { it.name == recordSet.name }?.let(changeBuilder::delete)
                changeBuilder.add(recordSet)
            }

        val batchChange = client.batch()
        val changeRequest = batchChange.applyChangeRequest(zoneId, changeBuilder.build())
        batchChange.submit()

        return zoneId to changeRequest.get().generatedId
    }

    private fun fetchTxtRecords(client: Dns, zoneId: String): List<RecordSet> {
        val pageSize = Dns.RecordSetListOption.pageSize(PAGE_SIZE)
        return fetchAll { nextPage ->
            val pageToken = nextPage?.let(Dns.RecordSetListOption::pageToken)
            val options = listOfNotNull(pageSize, pageToken).toTypedArray()
            client.listRecordSets(zoneId, *options)
        }.filter {
            it.type == RecordSet.Type.TXT
        }
    }

    private fun fetchZones(client: Dns): List<Zone> {
        val fields = Dns.ZoneListOption.fields(Dns.ZoneField.NAME, Dns.ZoneField.DNS_NAME, Dns.ZoneField.ZONE_ID)
        val pageSize = Dns.ZoneListOption.pageSize(PAGE_SIZE)
        return fetchAll { nextPage ->
            val pageToken = nextPage?.let(Dns.ZoneListOption::pageToken)
            val options = listOfNotNull(fields, pageSize, pageToken).toTypedArray()
            client.listZones(*options)
        }
    }

    private fun <T> fetchAll(pageProvider: (String?) -> Page<T>): List<T> {
        val output = mutableListOf<T>()
        var nextPage: String? = null

        do {
            val page = pageProvider(nextPage)
            output += page.values
            nextPage = page.nextPageToken
        } while (nextPage != null)

        return output
    }

    private fun findDnsZone(zones: List<Zone>, record: DnsProvider.ChallengeRecord): String =
        zones
            .filter {
                val normalizedDnsName =
                    if (it.dnsName.endsWith(".")) it.dnsName.removeSuffix(".")
                    else it.dnsName

                record.domainName.endsWith(normalizedDnsName)
            }
            .maxByOrNull { it.dnsName.length }
            ?.name
            ?: error("No DNS zone found for ${record.domainName}")
}
