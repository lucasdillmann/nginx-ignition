export default class ContextHolder<T> {
    private current: T

    constructor(initialValue: T) {
        this.current = initialValue
    }

    replace(newContext: T) {
        this.current = newContext
    }

    get(): T {
        return this.current
    }
}
