export interface Version {
    current?: string
    latest?: string
}

export default interface Configuration {
    version: Version
}
