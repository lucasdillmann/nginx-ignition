package br.com.dillmann.nginxignition.api.certificate

import br.com.dillmann.nginxignition.api.certificate.model.AvailableProviderResponse
import br.com.dillmann.nginxignition.api.certificate.model.CertificateResponse
import br.com.dillmann.nginxignition.api.certificate.model.IssueCertificateRequest
import br.com.dillmann.nginxignition.api.certificate.model.IssueCertificateResponse
import br.com.dillmann.nginxignition.api.certificate.model.RenewCertificateResponse
import br.com.dillmann.nginxignition.api.common.jsonobject.toJsonObject
import br.com.dillmann.nginxignition.api.common.jsonobject.toUnwrappedMap
import br.com.dillmann.nginxignition.api.common.pagination.PageResponse
import br.com.dillmann.nginxignition.core.certificate.Certificate
import br.com.dillmann.nginxignition.core.certificate.command.IssueCertificateCommand
import br.com.dillmann.nginxignition.core.certificate.command.RenewCertificateCommand
import br.com.dillmann.nginxignition.core.certificate.model.AvailableCertificateProvider
import br.com.dillmann.nginxignition.core.certificate.provider.CertificateRequest
import br.com.dillmann.nginxignition.core.common.pagination.Page
import kotlinx.serialization.json.*
import org.mapstruct.Mapper
import org.mapstruct.Mapping
import org.mapstruct.Named
import org.mapstruct.ReportingPolicy

@Mapper(unmappedTargetPolicy = ReportingPolicy.IGNORE)
internal abstract class CertificateConverter {
    abstract fun toResponse(input: AvailableCertificateProvider): AvailableProviderResponse

    abstract fun toResponse(input: IssueCertificateCommand.Output): IssueCertificateResponse

    abstract fun toResponse(input: RenewCertificateCommand.Output): RenewCertificateResponse

    abstract fun toResponse(input: Page<Certificate>): PageResponse<CertificateResponse>

    @Mapping(source = "parameters", target = "parameters", qualifiedByName = ["toResponseParameters"])
    abstract fun toResponse(input: Certificate): CertificateResponse

    @Mapping(source = "parameters", target = "parameters", qualifiedByName = ["toDomainModelParameters"])
    abstract fun toDomainModel(input: IssueCertificateRequest): CertificateRequest

    @Named("toDomainModelParameters")
    protected fun toDomainModelParameters(input: JsonObject?): Map<String, Any?> = input?.toUnwrappedMap() ?: emptyMap()

    @Named("toResponseParameters")
    protected fun toResponseParameters(input: Map<String, Any>?): JsonObject =
        input?.toJsonObject() ?: JsonObject(emptyMap())
}
