package br.com.dillmann.nginxignition.core.certificate

import br.com.dillmann.nginxignition.core.certificate.command.*
import br.com.dillmann.nginxignition.core.certificate.model.AvailableCertificateProvider
import br.com.dillmann.nginxignition.core.certificate.provider.CertificateProvider
import br.com.dillmann.nginxignition.core.certificate.provider.CertificateRequest
import br.com.dillmann.nginxignition.core.common.log.logger
import br.com.dillmann.nginxignition.core.common.pagination.Page
import br.com.dillmann.nginxignition.core.nginx.NginxService
import java.util.*

internal class CertificateService(
    private val repository: CertificateRepository,
    private val validator: CertificateValidator,
    private val providers: List<CertificateProvider>,
    private val nginxService: NginxService,
) : DeleteCertificateCommand, GetCertificateCommand, IssueCertificateCommand,
    ListCertificateCommand, RenewCertificateCommand, GetAvailableProvidersCommand {
    private val logger = logger<CertificateService>()

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

        val providerOutput = renew(certificate)
        return RenewCertificateCommand.Output(providerOutput.success, providerOutput.errorReason)
    }

    private suspend fun renew(certificate: Certificate): CertificateProvider.Output {
        val provider = providers.first { it.uniqueId == certificate.providerId }
        val providerOutput = provider.renew(certificate)
        if (providerOutput.success)
            repository.save(providerOutput.certificate!!)

        return providerOutput
    }

    override suspend fun getAvailableProviders(): List<AvailableCertificateProvider> =
        providers.map {
            AvailableCertificateProvider(
                name = it.name,
                uniqueId = it.uniqueId,
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
}
