package br.com.dillmann.nginxignition.application.controller.certificate.model

import br.com.dillmann.nginxignition.core.certificate.command.IssueCertificateCommand
import br.com.dillmann.nginxignition.core.certificate.model.AvailableCertificateProvider
import br.com.dillmann.nginxignition.core.certificate.provider.CertificateRequest
import kotlinx.serialization.json.*
import org.mapstruct.Mapper
import org.mapstruct.Mapping
import org.mapstruct.Named
import org.mapstruct.ReportingPolicy

@Mapper(unmappedTargetPolicy = ReportingPolicy.IGNORE)
abstract class CertificateConverter {
    abstract fun toResponse(input: AvailableCertificateProvider): AvailableProviderResponse

    abstract fun toResponse(input: IssueCertificateCommand.Output): IssueCertificateResponse

    @Mapping(source = "parameters", target = "parameters", qualifiedByName = ["toDomainModelParameters"])
    abstract fun toDomainModel(input: IssueCertificateRequest): CertificateRequest

    @Named("toDomainModelParameters")
    protected fun toDomainModelParameters(input: JsonObject): Map<String, Any?> =
        input.entries
            .map { (key, rawValue) ->
                val value =
                    if (rawValue is JsonObject) toDomainModelParameters(rawValue)
                    else rawValue.jsonPrimitive.unwrap()

                key to value
            }
            .toMap()

    private fun JsonPrimitive.unwrap(): Any? =
        booleanOrNull
            ?: longOrNull?.toBigInteger()
            ?: doubleOrNull?.toBigDecimal()
            ?: contentOrNull
}