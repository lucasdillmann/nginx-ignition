package br.com.dillmann.nginxignition.api.integration

import br.com.dillmann.nginxignition.api.common.jsonobject.toJsonObject
import br.com.dillmann.nginxignition.api.common.pagination.PageResponse
import br.com.dillmann.nginxignition.api.integration.model.IntegrationConfigurationResponse
import br.com.dillmann.nginxignition.api.integration.model.IntegrationOptionResponse
import br.com.dillmann.nginxignition.api.integration.model.IntegrationResponse
import br.com.dillmann.nginxignition.core.common.pagination.Page
import br.com.dillmann.nginxignition.core.integration.command.GetIntegrationByIdCommand
import br.com.dillmann.nginxignition.core.integration.command.ListIntegrationsCommand
import br.com.dillmann.nginxignition.core.integration.model.IntegrationOption
import kotlinx.serialization.json.JsonObject
import org.mapstruct.Mapper
import org.mapstruct.Mapping
import org.mapstruct.Named
import org.mapstruct.ReportingPolicy

@Mapper(unmappedTargetPolicy = ReportingPolicy.IGNORE)
internal abstract class IntegrationConverter {
    abstract fun toResponse(input: ListIntegrationsCommand.Output): IntegrationResponse

    abstract fun toResponse(input: Page<IntegrationOption>): PageResponse<IntegrationOptionResponse>

    abstract fun toResponse(input: IntegrationOption): IntegrationOptionResponse

    @Mapping(source = "parameters", target = "parameters", qualifiedByName = ["toResponseParameters"])
    abstract fun toResponse(input: GetIntegrationByIdCommand.Output): IntegrationConfigurationResponse

    @Named("toResponseParameters")
    protected fun toResponseParameters(input: Map<String, Any?>?): JsonObject =
        input?.toJsonObject() ?: JsonObject(emptyMap())

}
