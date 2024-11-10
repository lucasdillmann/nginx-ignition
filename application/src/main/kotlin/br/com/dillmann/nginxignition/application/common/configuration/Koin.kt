package br.com.dillmann.nginxignition.application.common.configuration

import br.com.dillmann.nginxignition.application.ApplicationModule
import br.com.dillmann.nginxignition.certificate.custom.CustomCertificateModule
import br.com.dillmann.nginxignition.certificate.selfsigned.SelfSignedCertificateModule
import br.com.dillmann.nginxignition.core.CoreModule
import br.com.dillmann.nginxignition.database.DatabaseModule
import br.com.dillmann.nginxignition.certificate.acme.AcmeCertificateModule
import io.ktor.server.application.*
import org.koin.ktor.plugin.Koin

fun Application.configureKoin() {
    install(Koin) {
        modules(
            CoreModule.initialize(),
            DatabaseModule.initialize(),
            AcmeCertificateModule.initialize(),
            CustomCertificateModule.initialize(),
            SelfSignedCertificateModule.initialize(),
            ApplicationModule.initialize(),
        )
    }
}
