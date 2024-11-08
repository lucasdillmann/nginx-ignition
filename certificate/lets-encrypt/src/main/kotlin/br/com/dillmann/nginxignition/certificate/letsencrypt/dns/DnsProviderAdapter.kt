package br.com.dillmann.nginxignition.certificate.letsencrypt.dns

import org.shredzone.acme4j.challenge.Dns01Challenge

class DnsProviderAdapter(private val providers: List<DnsProvider>) {
    suspend fun writeChallengeRecord(
        providerId: String,
        authorization: Dns01Challenge,
        dynamicFields: Map<String, Any?>,
    ) {
        providers
            .first { provider -> provider.uniqueId == providerId }
            .writeChallengeRecord(authorization, dynamicFields)
    }
}
