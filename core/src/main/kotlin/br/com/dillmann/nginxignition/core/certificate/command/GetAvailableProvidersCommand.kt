package br.com.dillmann.nginxignition.core.certificate.command

import br.com.dillmann.nginxignition.core.certificate.model.AvailableCertificateProvider

fun interface GetAvailableProvidersCommand {
    suspend fun getAvailableProviders(): List<AvailableCertificateProvider>
}
