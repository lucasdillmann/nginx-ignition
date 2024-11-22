export default interface AvailableProviderResponse {
    uniqueId: string
    name: string
    dynamicFields: DynamicField[]
}

export interface DynamicField {
    id: string
    description: string
    required: boolean
    type: DynamicFieldType
    enumOptions: DynamicFieldEnumOption[]
    condition?: DynamicFieldCondition
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