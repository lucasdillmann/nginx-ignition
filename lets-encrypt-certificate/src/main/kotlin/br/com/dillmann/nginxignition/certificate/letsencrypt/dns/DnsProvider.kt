package br.com.dillmann.nginxignition.certificate.letsencrypt.dns

import org.shredzone.acme4j.challenge.Dns01Challenge

interface DnsProvider {
    val uniqueId: String
    suspend fun writeChallengeRecord(authorization: Dns01Challenge, dynamicFields: Map<String, Any?>)
}
