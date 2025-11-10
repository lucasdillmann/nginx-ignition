import DynamicField from "../../../core/dynamicfield/DynamicField"

export default interface AvailableDriverResponse {
    id: string
    name: string
    importantInstructions: string[]
    configurationFields: DynamicField[]
}
