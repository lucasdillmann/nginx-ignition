package br.com.dillmann.nginxignition.application.controller.user.model

import br.com.dillmann.nginxignition.application.common.pagination.PageResponse
import br.com.dillmann.nginxignition.core.common.pagination.Page
import br.com.dillmann.nginxignition.core.user.User
import br.com.dillmann.nginxignition.core.user.model.SaveUserRequest
import org.mapstruct.Mapper
import org.mapstruct.Mapping
import org.mapstruct.ReportingPolicy

@Mapper(unmappedTargetPolicy = ReportingPolicy.IGNORE)
interface UserConverter {
    @Mapping(target = "id", expression = "java(java.util.UUID.randomUUID())")
    fun toDomainModel(input: UserRequest): SaveUserRequest

    fun toResponse(input: User): UserResponse

    fun toResponse(page: Page<User>): PageResponse<UserResponse>
}
