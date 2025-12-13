import DynamicField from "../../../../core/dynamicfield/DynamicField"

export default interface AvailableProviderResponse {
    id: string
    name: string
    priority: number
    dynamicFields: DynamicField[]
}
