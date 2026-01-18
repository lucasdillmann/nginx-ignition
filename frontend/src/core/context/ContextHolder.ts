export type ContextListener<T> = (newContext: T) => void

export default class ContextHolder<T> {
    private readonly listeners: ContextListener<T>[]
    private current: T

    constructor(initialValue: T) {
        this.current = initialValue
        this.listeners = []
    }

    register(listener: ContextListener<T>) {
        const index = this.listeners.indexOf(listener)
        if (index !== -1) {
            return
        }

        this.listeners.push(listener)
    }

    deregister(listener: ContextListener<T>) {
        const index = this.listeners.indexOf(listener)
        if (index !== -1) {
            this.listeners.splice(index, 1)
        }
    }

    replace(newContext: T) {
        this.current = newContext

        for (const listener of this.listeners) {
            listener(newContext)
        }
    }

    get(): T {
        return this.current
    }
}
