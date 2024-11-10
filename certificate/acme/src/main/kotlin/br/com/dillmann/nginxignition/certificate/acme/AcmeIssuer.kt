package br.com.dillmann.nginxignition.certificate.acme

import br.com.dillmann.nginxignition.certificate.acme.dns.DnsProvider
import org.shredzone.acme4j.AccountBuilder
import org.shredzone.acme4j.Authorization
import org.shredzone.acme4j.Certificate
import org.shredzone.acme4j.Order
import org.shredzone.acme4j.Session
import org.shredzone.acme4j.Status
import org.shredzone.acme4j.challenge.Dns01Challenge
import java.security.KeyPair
import java.time.Duration
import kotlin.jvm.optionals.getOrNull

internal class AcmeIssuer {
    data class Context(
        val userKeys: KeyPair,
        val userMail: String,
        val domainKeys: KeyPair,
        val domainNames: List<String>,
        val timeout: Duration,
        val createDnsRecordAction: suspend (List<DnsProvider.ChallengeRecord>) -> Unit,
        val issuerUrl: String,
    )

    data class Output(
        val success: Boolean,
        val exception: Throwable? = null,
        val certificate: Certificate? = null,
    )

    suspend fun issue(context: Context): Output =
        try {
            execute(context)
        } catch (ex: Exception) {
            Output(success = false, exception = ex)
        }

    private suspend fun execute(context: Context): Output {
        val session = Session(context.issuerUrl)
        val account = AccountBuilder()
            .agreeToTermsOfService()
            .useKeyPair(context.userKeys)
            .addEmail(context.userMail)
            .create(session)

        val order = account.newOrder().domains(context.domainNames).create()
        val challengeRecords = order
            .authorizations
            .map {
                val domainName = "_acme-challenge.${it.identifier.domain}"
                val token = it.dnsChallenge().digest

                DnsProvider.ChallengeRecord(
                    domainName = domainName,
                    token = token,
                )
            }

        context.createDnsRecordAction(challengeRecords)
        order.authorizations.forEach { it.dnsChallenge().trigger() }

        order.waitForCompletion(context.timeout).requireSuccess(order)
        order.execute(context.domainKeys)
        order.waitForCompletion(context.timeout).requireSuccess(order)

        return Output(success = true, certificate = order.certificate)
    }

    private fun Status.requireSuccess(order: Order) {
        if (this != Status.VALID) {
            val reason = order.error.map { it.toString() }.getOrNull()
                ?: order.authorizations.mapNotNull { it.dnsChallenge().error.getOrNull() }.joinToString().takeIf { it.isNotBlank() }
                ?: "unknown error"
            error("Certificate failed to be issued: $reason")
        }
    }

    private fun Authorization.dnsChallenge() =
        findChallenge(Dns01Challenge::class.java).get()
}
