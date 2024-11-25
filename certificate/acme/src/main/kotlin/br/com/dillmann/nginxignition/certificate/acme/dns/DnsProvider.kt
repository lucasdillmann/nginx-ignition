package br.com.dillmann.nginxignition.certificate.acme.dns

internal interface DnsProvider {
    data class ChallengeRecord(
        val domainName: String,
        val token: String,
    )

    val id: String

    suspend fun writeChallengeRecords(
        records: List<ChallengeRecord>,
        dynamicFields: Map<String, Any?>,
    )
}
