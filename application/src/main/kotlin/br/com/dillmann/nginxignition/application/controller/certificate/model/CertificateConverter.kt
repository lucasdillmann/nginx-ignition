package br.com.dillmann.nginxignition.application.controller.certificate.model

import br.com.dillmann.nginxignition.core.certificate.model.AvailableCertificateProvider
import org.mapstruct.Mapper
import org.mapstruct.ReportingPolicy

@Mapper(unmappedTargetPolicy = ReportingPolicy.IGNORE)
interface CertificateConverter {
    fun toResponse(input: AvailableCertificateProvider): AvailableProviderResponse
}
