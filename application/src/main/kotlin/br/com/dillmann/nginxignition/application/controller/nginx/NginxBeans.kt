package br.com.dillmann.nginxignition.application.controller.nginx

import br.com.dillmann.nginxignition.application.controller.nginx.handler.NginxReloadHandler
import br.com.dillmann.nginxignition.application.controller.nginx.handler.NginxStartHandler
import br.com.dillmann.nginxignition.application.controller.nginx.handler.NginxStatusHandler
import br.com.dillmann.nginxignition.application.controller.nginx.handler.NginxStopHandler
import br.com.dillmann.nginxignition.application.controller.nginx.model.NginxConverter
import org.koin.core.module.Module
import org.mapstruct.factory.Mappers

internal fun Module.nginxBeans() {
    single { Mappers.getMapper(NginxConverter::class.java) }
    single { NginxStartHandler(get(), get()) }
    single { NginxStopHandler(get(), get()) }
    single { NginxReloadHandler(get(), get()) }
    single { NginxStatusHandler(get()) }
}
