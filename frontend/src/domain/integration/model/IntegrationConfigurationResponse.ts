import { IntegrationResponse } from "./IntegrationResponse"
import DynamicField from "../../../core/dynamicfield/DynamicField"

export interface IntegrationConfigurationResponse extends IntegrationResponse {
    configurationFields: DynamicField[]
    parameters: Record<string, any>
}
