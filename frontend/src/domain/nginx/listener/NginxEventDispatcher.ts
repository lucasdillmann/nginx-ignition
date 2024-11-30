import { NginxEventListener, NginxOperation } from "./NginxEventListener"

class NginxEventDispatcher {
    private listeners: NginxEventListener[]

    constructor() {
        this.listeners = []
    }

    register(listener: NginxEventListener) {
        this.listeners.push(listener)
    }

    remove(listener: NginxEventListener) {
        this.listeners = this.listeners.filter(element => element !== listener)
    }

    notify(operation: NginxOperation) {
        this.listeners.forEach(listener => listener(operation))
    }
}

// eslint-disable-next-line import/no-anonymous-default-export
export default new NginxEventDispatcher()
