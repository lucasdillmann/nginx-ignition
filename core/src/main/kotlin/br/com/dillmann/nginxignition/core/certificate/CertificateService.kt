package br.com.dillmann.nginxignition.core.certificate

import br.com.dillmann.nginxignition.core.certificate.command.*
import br.com.dillmann.nginxignition.core.certificate.model.AvailableCertificateProvider
import br.com.dillmann.nginxignition.core.certificate.provider.CertificateProvider
import br.com.dillmann.nginxignition.core.certificate.provider.CertificateRequest
import br.com.dillmann.nginxignition.core.common.dynamicfield.DynamicFields
import br.com.dillmann.nginxignition.core.common.log.logger
import br.com.dillmann.nginxignition.core.common.pagination.Page
import br.com.dillmann.nginxignition.core.host.HostRepository
import br.com.dillmann.nginxignition.core.nginx.NginxService
import java.util.*

internal class CertificateService(
    private val repository: CertificateRepository,
    private val validator: CertificateValidator,
    private val providers: List<CertificateProvider>,
    private val nginxService: NginxService,
    private val hostRepository: HostRepository,
) : DeleteCertificateCommand, GetCertificateCommand, IssueCertificateCommand,
    ListCertificateCommand, RenewCertificateCommand, GetAvailableProvidersCommand {
    private val logger = logger<CertificateService>()

    override suspend fun deleteById(id: UUID): DeleteCertificateCommand.Output {
        if (hostRepository.existsByCertificateId(id)) {
            return DeleteCertificateCommand.Output(
                deleted = false,
                reason = "Certificate is being used by at least one host. Please update them and try again.",
            )
        }

        repository.deleteById(id)
        return DeleteCertificateCommand.Output(
            deleted = true,
            reason = "Certificate deleted successfully",
        )
    }

    override suspend fun getById(id: UUID): Certificate? =
        repository.findById(id)?.let(::removeSensitiveParameters)

    override suspend fun issue(request: CertificateRequest): IssueCertificateCommand.Output {
        validator.validate(request)
        val provider = providers.first { it.id == request.providerId }
        val providerResult = runCatching { provider.issue(request) }
        if (providerResult.isFailure)
            return IssueCertificateCommand.Output(
                success = false,
                errorReason = providerResult.exceptionOrNull()?.message,
            )

        val providerOutput = providerResult.getOrThrow()
        if (providerOutput.success)
            repository.save(providerOutput.certificate!!)

        return IssueCertificateCommand.Output(
            success = providerOutput.success,
            errorReason = providerOutput.errorReason,
            certificateId = providerOutput.certificate?.id,
        )
    }

    override suspend fun list(pageSize: Int, pageNumber: Int): Page<Certificate> =
        repository
            .findPage(pageSize, pageNumber)
            .let {
                val updatedContents = it.contents.map(::removeSensitiveParameters)
                it.copy(contents = updatedContents)
            }

    override suspend fun renewById(id: UUID): RenewCertificateCommand.Output {
        val certificate = repository.findById(id)
        if (certificate == null) {
            return RenewCertificateCommand.Output(
                success = false,
                errorReason = "No certificate found with ID $id",
            )
        }

        val providerOutput = renew(certificate)
        return RenewCertificateCommand.Output(providerOutput.success, providerOutput.errorReason)
    }

    private suspend fun renew(certificate: Certificate): CertificateProvider.Output {
        val provider = providers.first { it.id == certificate.providerId }
        val providerOutput = provider.renew(certificate)
        if (providerOutput.success)
            repository.save(providerOutput.certificate!!)

        return providerOutput.copy(
            certificate = providerOutput.certificate?.let(::removeSensitiveParameters),
        )
    }

    override suspend fun getAvailableProviders(): List<AvailableCertificateProvider> =
        providers.map {
            AvailableCertificateProvider(
                name = it.name,
                id = it.id,
                dynamicFields = it.dynamicFields,
            )
        }

    suspend fun renewAllDue() {
        val certificatesToRenew = repository.findAllDueToRenew()
        if (certificatesToRenew.isEmpty()) {
            logger.info("Certificates auto-renew triggered, but no certificates are due to be renewed yet")
            return
        }

        logger.info("${certificatesToRenew.size} certificates due to be renewed")
        certificatesToRenew.forEach { certificate ->
            val providerOutput = renew(certificate)
            if (providerOutput.success)
                logger.info("Certificate ${certificate.id} renewed successfully")
            else
                logger.warn("Certificate ${certificate.id} failed to be renewed: ${providerOutput.errorReason}")
        }

        nginxService.reload()
        logger.info("Certificate auto-renew cycle completed")
    }

    private fun removeSensitiveParameters(certificate: Certificate): Certificate {
        val provider = providers.first { it.id == certificate.providerId }
        return certificate
            .copy(parameters = DynamicFields.removeSensitiveParameters(provider.dynamicFields, certificate.parameters))
    }
}
