package br.com.dillmann.nginxignition.api.host

import br.com.dillmann.nginxignition.api.common.pagination.PageResponse
import br.com.dillmann.nginxignition.api.host.model.HostRequest
import br.com.dillmann.nginxignition.api.host.model.HostResponse
import br.com.dillmann.nginxignition.core.common.pagination.Page
import br.com.dillmann.nginxignition.core.host.Host
import org.mapstruct.Mapper
import org.mapstruct.Mapping
import org.mapstruct.ReportingPolicy

@Mapper(unmappedTargetPolicy = ReportingPolicy.IGNORE)
internal interface HostConverter {

    fun toResponse(input: Host): HostResponse

    fun toResponse(page: Page<Host>): PageResponse<HostResponse>

    @Mapping(target = "id", expression = "java(java.util.UUID.randomUUID())")
    fun toDomainModel(input: HostRequest): Host

    @Mapping(target = "id", expression = "java(java.util.UUID.randomUUID())")
    fun toDomainModel(input: HostRequest.Route): Host.Route

    @Mapping(target = "id", expression = "java(java.util.UUID.randomUUID())")
    fun toDomainModel(input: HostRequest.Binding): Host.Binding
}
