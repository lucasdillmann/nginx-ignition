package br.com.dillmann.nginxignition.certificate.letsencrypt.dns

import br.com.dillmann.nginxignition.certificate.letsencrypt.DynamicFields
import kotlinx.coroutines.future.await
import org.shredzone.acme4j.challenge.Dns01Challenge
import software.amazon.awssdk.auth.credentials.AwsBasicCredentials
import software.amazon.awssdk.services.route53.Route53AsyncClient
import software.amazon.awssdk.services.route53.model.*

class Route53DnsProvider: DnsProvider {
    override val uniqueId = "ROUTE_53"

    override suspend fun writeChallengeRecord(authorization: Dns01Challenge, dynamicFields: Map<String, Any?>) {
        val client = startClient(dynamicFields)
        val domainName = authorization.location.host
        val hostedZoneId = findHostedZoneId(client, domainName) ?: error("No Route53 hosted zone found for $domainName")
        upsertRecord(client, hostedZoneId, domainName, authorization.authorization)
    }

    private fun startClient(dynamicFields: Map<String, Any?>): Route53AsyncClient {
        val accessKey = dynamicFields[DynamicFields.AWS_ACCESS_KEY.id] as String
        val secretKey = dynamicFields[DynamicFields.AWS_SECRET_KEY.id] as String
        return Route53AsyncClient
            .builder()
            .credentialsProvider { AwsBasicCredentials.create(accessKey, secretKey) }
            .build()
    }

    private suspend fun findHostedZoneId(client: Route53AsyncClient, domainName: String): String? {
        val response = client.listHostedZones().await()
        return response
            .hostedZones()
            .filter { domainName.endsWith(it.name()) }
            .maxByOrNull { it.name().length }
            ?.id()
    }

    private suspend fun upsertRecord(
        client: Route53AsyncClient,
        hostedZoneId: String,
        domainName: String,
        contents: String,
    ) {
        val record = ResourceRecord.builder().value(contents).build()
        val recordSet = ResourceRecordSet.builder().name(domainName).type(RRType.TXT).resourceRecords(record).build()
        val change = Change.builder().action(ChangeAction.UPSERT).resourceRecordSet(recordSet).build()
        val batch = ChangeBatch.builder().changes(change).build()
        client.changeResourceRecordSets { it.hostedZoneId(hostedZoneId).changeBatch(batch) }.await()
    }
}
