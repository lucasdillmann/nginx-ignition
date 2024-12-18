package br.com.dillmann.nginxignition.api.accesslist

import br.com.dillmann.nginxignition.api.accesslist.model.AccessListRequest
import br.com.dillmann.nginxignition.api.accesslist.model.AccessListResponse
import br.com.dillmann.nginxignition.api.common.pagination.PageResponse
import br.com.dillmann.nginxignition.api.host.HostConverter
import br.com.dillmann.nginxignition.core.accesslist.AccessList
import br.com.dillmann.nginxignition.core.common.pagination.Page
import org.mapstruct.Mapper
import org.mapstruct.Mapping
import org.mapstruct.ReportingPolicy

@Mapper(unmappedTargetPolicy = ReportingPolicy.IGNORE, uses = [HostConverter::class])
internal interface AccessListConverter {
    @Mapping(target = "id", expression = "java(java.util.UUID.randomUUID())")
    fun toDomain(input: AccessListRequest): AccessList
    fun toResponse(input: AccessList): AccessListResponse
    fun toResponse(input: Page<AccessList>): PageResponse<AccessListResponse>
}
