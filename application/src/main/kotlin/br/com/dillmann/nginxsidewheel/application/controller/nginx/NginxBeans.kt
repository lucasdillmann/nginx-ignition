package br.com.dillmann.nginxsidewheel.application.controller.nginx

import br.com.dillmann.nginxsidewheel.application.controller.host.handler.NginxReloadHandler
import br.com.dillmann.nginxsidewheel.application.controller.host.handler.NginxStartHandler
import br.com.dillmann.nginxsidewheel.application.controller.host.handler.NginxStopHandler
import br.com.dillmann.nginxsidewheel.application.controller.nginx.model.NginxConverter
import org.koin.core.module.Module
import org.mapstruct.factory.Mappers

internal fun Module.nginxBeans() {
    single { Mappers.getMapper(NginxConverter::class.java) }
    single { NginxStartHandler(get(), get()) }
    single { NginxStopHandler(get(), get()) }
    single { NginxReloadHandler(get(), get()) }
}
