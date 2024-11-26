export default interface AvailableProviderResponse {
    id: string
    name: string
    dynamicFields: DynamicField[]
}

export interface DynamicField {
    id: string
    priority: number
    description: string
    required: boolean
    sensitive: boolean
    type: DynamicFieldType
    enumOptions: DynamicFieldEnumOption[]
    condition?: DynamicFieldCondition
    helpText?: string
}

export enum DynamicFieldType {
    SINGLE_LINE_TEXT = "SINGLE_LINE_TEXT",
    MULTI_LINE_TEXT = "MULTI_LINE_TEXT",
    EMAIL = "EMAIL",
    BOOLEAN = "BOOLEAN",
    ENUM = "ENUM",
    FILE = "FILE",
}

export interface DynamicFieldEnumOption {
    id: string
    description: string
}

export interface DynamicFieldCondition {
    parentField: string
    value: string
}
