package br.com.dillmann.nginxignition.certificate.letsencrypt.dns

class DnsProviderAdapter(private val providers: List<DnsProvider>) {
    suspend fun writeChallengeRecord(
        providerId: String,
        challengeRecords: List<DnsProvider.ChallengeRecord>,
        dynamicFields: Map<String, Any?>,
    ) {
        providers
            .first { provider -> provider.uniqueId == providerId }
            .writeChallengeRecords(challengeRecords, dynamicFields)
    }
}
