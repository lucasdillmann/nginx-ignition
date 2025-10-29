import DynamicField from "../../../core/dynamicfield/DynamicField"

export default interface AvailableDriverResponse {
    id: string
    name: string
    description: string
    configurationFields: DynamicField[]
}
