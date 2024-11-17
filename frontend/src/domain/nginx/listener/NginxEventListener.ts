export enum NginxOperation {
    START,
    STOP,
    RELOAD,
}

export type NginxEventListener = (operation: NginxOperation) => void
