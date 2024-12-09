import { CertificateAutoRenewSettingsDto, LogRotationSettingsDto, NginxSettingsDto } from "./SettingsDto"
import { HostFormBinding } from "../../host/model/HostFormValues"

export default interface SettingsFormValues {
    nginx: NginxSettingsDto
    logRotation: LogRotationSettingsDto
    certificateAutoRenew: CertificateAutoRenewSettingsDto
    globalBindings: HostFormBinding[]
}
