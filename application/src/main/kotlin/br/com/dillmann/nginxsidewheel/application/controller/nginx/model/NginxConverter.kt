package br.com.dillmann.nginxsidewheel.application.controller.nginx.model

import br.com.dillmann.nginxsidewheel.core.nginx.exception.NginxCommandException
import org.mapstruct.Mapper
import org.mapstruct.ReportingPolicy

@Mapper(unmappedTargetPolicy = ReportingPolicy.IGNORE)
interface NginxConverter {
    fun toResponse(exception: NginxCommandException): NginxActionErrorResponse
}
