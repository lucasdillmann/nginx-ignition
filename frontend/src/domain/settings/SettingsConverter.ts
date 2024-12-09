import CertificateService from "../certificate/CertificateService"
import SettingsDto from "./model/SettingsDto"
import SettingsFormValues from "./model/SettingsFormValues"
import { HostBinding } from "../host/model/HostRequest"
import { HostFormBinding } from "../host/model/HostFormValues"

class SettingsConverter {
    private readonly certificateService: CertificateService

    constructor() {
        this.certificateService = new CertificateService()
    }

    private notNull(value?: any) {
        return value !== undefined && value !== null
    }

    private async bindingToFormValues(binding: HostBinding): Promise<HostFormBinding> {
        const certificate = this.notNull(binding.certificateId)
            ? await this.certificateService.getById(binding.certificateId!!)
            : undefined

        return {
            ...binding,
            certificate,
        }
    }

    private formValuesToBinding(binding: HostFormBinding): HostBinding {
        const output = {
            ...binding,
            certificateId: binding.certificate?.id,
        }

        delete output.certificate
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

// eslint-disable-next-line import/no-anonymous-default-export
export default new SettingsConverter()
