package br.com.dillmann.nginxignition.certificate.acme.dns.aws

import br.com.dillmann.nginxignition.certificate.acme.dns.DnsProvider
import br.com.dillmann.nginxignition.certificate.acme.letsencrypt.LetsEncryptDynamicFields
import kotlinx.coroutines.delay
import kotlinx.coroutines.future.await
import software.amazon.awssdk.auth.credentials.AwsBasicCredentials
import software.amazon.awssdk.regions.Region
import software.amazon.awssdk.services.route53.Route53AsyncClient
import software.amazon.awssdk.services.route53.model.*

internal class Route53DnsProvider: DnsProvider {
    companion object {
        private const val CHECK_CHANGE_STATUS_DELAY_MS = 250L
        private const val RECORD_TIME_TO_LIVE_SECONDS = 30L
        const val ID = "AWS_ROUTE53"
    }

    private data class RecordMetadata(
        val record: DnsProvider.ChallengeRecord,
        val hostedZoneId: String,
    )

    override val id = ID

    override suspend fun writeChallengeRecords(
        records: List<DnsProvider.ChallengeRecord>,
        dynamicFields: Map<String, Any?>,
    ) {
        val client = startClient(dynamicFields)
        val recordsMetadata = findHostedZoneIds(client, records)
        upsertRecords(client, recordsMetadata)
    }

    private fun startClient(dynamicFields: Map<String, Any?>): Route53AsyncClient {
        val accessKey = dynamicFields[LetsEncryptDynamicFields.AWS_ACCESS_KEY.id] as String
        val secretKey = dynamicFields[LetsEncryptDynamicFields.AWS_SECRET_KEY.id] as String
        System.setProperty("aws.region", Region.US_EAST_1.id())

        return Route53AsyncClient
            .builder()
            .credentialsProvider { AwsBasicCredentials.create(accessKey, secretKey) }
            .build()
    }

    private suspend fun findHostedZoneIds(
        client: Route53AsyncClient,
        records: List<DnsProvider.ChallengeRecord>,
    ): List<RecordMetadata> {
        val hostedZones = client
            .listHostedZones()
            .await()
            .hostedZones()
            .map { it.id() to it.name().dropLastWhile { it == '.' } }

        return records.map {
            RecordMetadata(
                record = it,
                hostedZoneId = findHostedZoneId(hostedZones, it.domainName),
            )
        }
    }

    private fun findHostedZoneId(hostedZones: List<Pair<String, String>>, domainName: String) =
        hostedZones
            .filter { (_, name) -> domainName.endsWith(name) }
            .maxBy { (_, name) -> name.length }
            .first

    private suspend fun upsertRecords(
        client: Route53AsyncClient,
        records: List<RecordMetadata>,
    ) {
        val changeIds = records
            .groupBy { it.hostedZoneId }
            .map { (hostedZoneId, items) -> upsertRecords(client, hostedZoneId, items.map { it.record }) }

        changeIds.forEach { changeId ->
            do {
                val status = client.getChange { it.id(changeId) }.await().changeInfo().status()
                if (status != ChangeStatus.INSYNC)
                    delay(CHECK_CHANGE_STATUS_DELAY_MS)
            } while (status != ChangeStatus.INSYNC)
        }
    }

    private suspend fun upsertRecords(
        client: Route53AsyncClient,
        hostedZoneId: String,
        records: List<DnsProvider.ChallengeRecord>,
    ): String {
        val changes = records
            .groupBy { it.domainName }
            .map { (domainName, tokens) ->
                val resourceRecords = tokens.map { (_, token) -> ResourceRecord.builder().value("\"$token\"").build() }
                val recordSet = ResourceRecordSet
                    .builder()
                    .name(domainName)
                    .type(RRType.TXT)
                    .resourceRecords(resourceRecords)
                    .ttl(RECORD_TIME_TO_LIVE_SECONDS)
                    .build()
                Change.builder().action(ChangeAction.UPSERT).resourceRecordSet(recordSet).build()
            }

        val batch = ChangeBatch.builder().changes(changes).build()
        return client
            .changeResourceRecordSets { it.hostedZoneId(hostedZoneId).changeBatch(batch) }
            .await()
            .changeInfo()
            .id()
    }
}
