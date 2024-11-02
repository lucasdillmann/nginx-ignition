package br.com.dillmann.nginxsidewheel.application.controller.host.model

import br.com.dillmann.nginxsidewheel.application.common.pagination.PageResponse
import br.com.dillmann.nginxsidewheel.core.common.pagination.Page
import br.com.dillmann.nginxsidewheel.core.host.Host
import org.mapstruct.Mapper
import org.mapstruct.Mapping
import org.mapstruct.ReportingPolicy

@Mapper(unmappedTargetPolicy = ReportingPolicy.IGNORE)
interface HostConverter {

    fun toResponse(input: Host): HostResponse

    fun toResponse(page: Page<Host>): PageResponse<HostResponse>

    @Mapping(target = "id", expression = "java(java.util.UUID.randomUUID())")
    fun toDomainModel(input: HostRequest): Host

    @Mapping(target = "id", expression = "java(java.util.UUID.randomUUID())")
    fun toDomainModel(input: HostRequest.Route): Host.Route

    @Mapping(target = "id", expression = "java(java.util.UUID.randomUUID())")
    fun toDomainModel(input: HostRequest.Binding): Host.Binding
}
