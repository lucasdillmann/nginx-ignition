package br.com.dillmann.nginxignition.application.controller.nginx.model

import br.com.dillmann.nginxignition.core.nginx.exception.NginxCommandException
import org.mapstruct.Mapper
import org.mapstruct.ReportingPolicy

@Mapper(unmappedTargetPolicy = ReportingPolicy.IGNORE)
interface NginxConverter {
    fun toResponse(exception: NginxCommandException): NginxActionErrorResponse
}
