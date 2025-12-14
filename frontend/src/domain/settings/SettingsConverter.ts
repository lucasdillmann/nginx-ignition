import ServerCertificateService from "../certificate/server/ServerCertificateService"
import SettingsDto from "./model/SettingsDto"
import SettingsFormValues from "./model/SettingsFormValues"
import { HostBinding } from "../host/model/HostRequest"
import { HostFormBinding } from "../host/model/HostFormValues"

class SettingsConverter {
    private readonly serverCertificateService: ServerCertificateService

    constructor() {
        this.serverCertificateService = new ServerCertificateService()
    }

    private notNull(value?: any) {
        return value !== undefined && value !== null
    }

    private async bindingToFormValues(binding: HostBinding): Promise<HostFormBinding> {
        const serverCertificate = this.notNull(binding.serverCertificateId)
            ? await this.serverCertificateService.getById(binding.serverCertificateId!!)
            : undefined

        return {
            ...binding,
            serverCertificate,
        }
    }

    private formValuesToBinding(binding: HostFormBinding): HostBinding {
        const output = {
            ...binding,
            serverCertificateId: binding.serverCertificate?.id,
        }

        delete output.serverCertificate
        return output
    }

    async settingsToFormValues(settings: SettingsDto): Promise<SettingsFormValues> {
        const { nginx, certificateAutoRenew, logRotation } = settings

        const globalBindings = await Promise.all(
            settings.globalBindings.map(binding => this.bindingToFormValues(binding)),
        )

        return {
            nginx,
            certificateAutoRenew,
            logRotation,
            globalBindings,
        }
    }

    formValuesToSettings(formValues: SettingsFormValues): SettingsDto {
        const { nginx, certificateAutoRenew, logRotation } = formValues

        const globalBindings = formValues.globalBindings.map(binding => this.formValuesToBinding(binding))

        return {
            nginx,
            certificateAutoRenew,
            logRotation,
            globalBindings,
        }
    }
}

export default new SettingsConverter()
