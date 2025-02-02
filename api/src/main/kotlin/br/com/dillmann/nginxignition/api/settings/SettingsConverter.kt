package br.com.dillmann.nginxignition.api.settings

import br.com.dillmann.nginxignition.api.host.HostConverter
import br.com.dillmann.nginxignition.api.settings.model.SettingsDto
import br.com.dillmann.nginxignition.core.settings.Settings
import org.mapstruct.Mapper
import org.mapstruct.ReportingPolicy

@Mapper(unmappedTargetPolicy = ReportingPolicy.IGNORE, uses = [HostConverter::class])
internal interface SettingsConverter {
    fun toResponse(input: Settings): SettingsDto
    fun toDomain(input: SettingsDto): Settings
}
