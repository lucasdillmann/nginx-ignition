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

    override suspend fun issue(request: CertificateRequest): IssueCertificateCommand.Output {
        validator.validate(request)
        val provider = providers.first { it.uniqueId == request.providerId }
        val providerOutput = provider.issue(request)
        if (providerOutput.success)
            repository.save(providerOutput.certificate!!)

        return IssueCertificateCommand.Output(
            success = providerOutput.success,
            errorReason = providerOutput.errorReason,
            certificateId = providerOutput.certificate?.id,
        )
    }

    override suspend fun list(pageSize: Int, pageNumber: Int): Page<Certificate> =
        repository.findPage(pageSize, pageNumber)

    override suspend fun renewById(id: UUID): RenewCertificateCommand.Output {
        val certificate = repository.findById(id)
        if (certificate == null) {
            return RenewCertificateCommand.Output(
                success = false,
                errorReason = "No certificate found with ID $id",
            )
        }

        val provider = providers.first { it.uniqueId == certificate.providerId }
        val providerOutput = provider.renew(certificate)
        if (providerOutput.success)
            repository.save(providerOutput.certificate!!)

        return RenewCertificateCommand.Output(providerOutput.success, providerOutput.errorReason)
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
