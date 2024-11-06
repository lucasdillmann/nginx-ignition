package br.com.dillmann.nginxignition.core.certificate

import br.com.dillmann.nginxignition.core.certificate.command.*
import br.com.dillmann.nginxignition.core.certificate.model.AvailableCertificateProvider
import br.com.dillmann.nginxignition.core.certificate.provider.CertificateProvider
import br.com.dillmann.nginxignition.core.certificate.provider.CertificateRequest
import br.com.dillmann.nginxignition.core.common.pagination.Page
import java.util.*

internal class CertificateService(
    private val repository: CertificateRepository,
    private val validator: CertificateValidator,
    private val providers: List<CertificateProvider>,
): DeleteCertificateCommand, GetCertificateCommand, IssueCertificateCommand,
   ListCertificateCommand, RenewCertificateCommand, GetAvailableProvidersCommand {
    override suspend fun deleteById(id: UUID) {
        repository.deleteById(id)
    }

    override suspend fun getById(id: UUID): Certificate? =
        repository.findById(id)

    override suspend fun issue(request: CertificateRequest): UUID {
        validator.validate(request)
        val provider = providers.first { it.uniqueId == request.providerId }
        val certificate = provider.issue(request)
        repository.save(certificate)
        return certificate.id
    }

    override suspend fun list(pageSize: Int, pageNumber: Int): Page<Certificate> =
        repository.findPage(pageSize, pageNumber)

    override suspend fun renewById(id: UUID) {
        val certificate = repository.findById(id)
        require(certificate != null) { "No certificate found with ID $id" }

        val provider = providers.first { it.uniqueId == certificate.providerId }
        val updatedCertificate = provider.renew(certificate)
        repository.save(updatedCertificate)
    }

    override suspend fun getAvailableProviders(): List<AvailableCertificateProvider> =
        providers.map {
            AvailableCertificateProvider(
                name = it.name,
                uniqueId = it.uniqueId,
                dynamicFields = it.dynamicFields,
            )
        }
}
