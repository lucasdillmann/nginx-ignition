package br.com.dillmann.nginxignition.core.nginx.configuration.provider

import br.com.dillmann.nginxignition.core.certificate.CertificateRepository
import br.com.dillmann.nginxignition.core.common.extensions.decodeBase64
import br.com.dillmann.nginxignition.core.host.Host
import br.com.dillmann.nginxignition.core.nginx.configuration.NginxConfigurationFileProvider
import java.util.*

internal class HostCertificateFileProvider(
    private val certificateRepository: CertificateRepository,
): NginxConfigurationFileProvider {
    private companion object {
        private const val LINE_LENGTH = 64
        private const val LINE_SEPARATOR = "\n"
        private const val PRIVATE_KEY_HEADER = "-----BEGIN PRIVATE KEY-----"
        private const val PRIVATE_KEY_FOOTER = "-----END PRIVATE KEY-----"
        private const val PUBLIC_KEY_HEADER = "-----BEGIN CERTIFICATE-----"
        private const val PUBLIC_KEY_FOOTER = "-----END CERTIFICATE-----"
    }

    override suspend fun provide(basePath: String, hosts: List<Host>): List<NginxConfigurationFileProvider.Output> =
        hosts
            .flatMap { it.bindings }
            .filter { it.type == Host.BindingType.HTTPS }
            .mapNotNull { it.certificateId }
            .distinct()
            .map { buildCertificateFile(it) }

    private suspend fun buildCertificateFile(certificateId: UUID, ): NginxConfigurationFileProvider.Output {
        val certificate = certificateRepository.findById(certificateId)
            ?: error("No certificate found with ID $certificateId")

        val certificateChain = certificate
            .certificationChain
            .joinToString(separator = LINE_SEPARATOR) { convertToPemEncodedString(it, null) }
        val mainContents = convertToPemEncodedString(certificate.publicKey, certificate.privateKey)

        return NginxConfigurationFileProvider.Output(
            name = "certificate-$certificateId.pem",
            contents = "$certificateChain$mainContents"
        )
    }

    private fun convertToPemEncodedString(
        publicKey: String,
        privateKey: String?,
    ): String {
        val encoder = Base64.getMimeEncoder(LINE_LENGTH, LINE_SEPARATOR.encodeToByteArray())
        val publicKeyBytes = publicKey.decodeBase64()
        val buffer = StringBuffer()

        buffer.appendLine(PUBLIC_KEY_HEADER)
        buffer.appendLine(encoder.encodeToString(publicKeyBytes))
        buffer.appendLine(PUBLIC_KEY_FOOTER)

        if (privateKey != null) {
            val privateKeyBytes = privateKey.decodeBase64()
            buffer.appendLine(PRIVATE_KEY_HEADER)
            buffer.appendLine(encoder.encodeToString(privateKeyBytes))
            buffer.appendLine(PRIVATE_KEY_FOOTER)
        }

        return buffer.toString()
    }
}
