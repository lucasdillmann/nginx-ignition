import DynamicField from "../../../core/dynamicfield/DynamicField"

export default interface AvailableProviderResponse {
    id: string
    name: string
    dynamicFields: DynamicField[]
}
