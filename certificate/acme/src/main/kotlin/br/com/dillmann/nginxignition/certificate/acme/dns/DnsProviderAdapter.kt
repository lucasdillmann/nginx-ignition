package br.com.dillmann.nginxignition.certificate.acme.dns

internal class DnsProviderAdapter(private val providers: List<DnsProvider>) {
    suspend fun writeChallengeRecord(
        providerId: String,
        challengeRecords: List<DnsProvider.ChallengeRecord>,
        dynamicFields: Map<String, Any?>,
    ) {
        providers
            .first { provider -> provider.id == providerId }
            .writeChallengeRecords(challengeRecords, dynamicFields)
    }
}
