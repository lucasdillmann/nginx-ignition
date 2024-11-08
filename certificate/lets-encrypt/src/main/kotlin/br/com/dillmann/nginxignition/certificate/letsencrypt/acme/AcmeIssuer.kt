package br.com.dillmann.nginxignition.certificate.letsencrypt.acme

import org.shredzone.acme4j.AccountBuilder
import org.shredzone.acme4j.Certificate
import org.shredzone.acme4j.Session
import org.shredzone.acme4j.Status
import org.shredzone.acme4j.challenge.Dns01Challenge
import java.security.KeyPair
import java.time.Duration
import kotlin.jvm.optionals.getOrNull

class AcmeIssuer {
    data class Context(
        val userKeys: KeyPair,
        val userMail: String,
        val domainKeys: KeyPair,
        val domainNames: List<String>,
        val timeout: Duration,
        val createDnsRecordAction: suspend (Dns01Challenge) -> Unit,
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
        order.authorizations.forEach {
            val dnsChallenge = it
                .findChallenge(Dns01Challenge::class.java)
                .getOrNull()
                ?: error("DNS challenge not found")
            context.createDnsRecordAction(dnsChallenge)
            dnsChallenge.trigger()
        }

        order.waitForCompletion(context.timeout)
        order.execute(context.domainKeys)

        val status = order.waitForCompletion(context.timeout)
        if (status != Status.VALID) {
            val reason = order.error.map { it.toString() }.getOrNull() ?: "unknown error"
            error("Certificate failed to be issued: $reason")
        }

        return Output(success = true, certificate = order.certificate)
    }
}
