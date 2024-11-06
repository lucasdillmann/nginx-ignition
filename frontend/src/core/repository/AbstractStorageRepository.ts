export default abstract class AbstractStorageRepository<T> {
    private readonly key: string
    private readonly delegate: Storage

    protected constructor(delegate: Storage, key: string) {
        this.delegate = delegate
        this.key = key
    }

    store(value: T) {
        const json = JSON.stringify(value)
        this.delegate.setItem(this.key, json)
    }

    clear() {
        this.delegate.removeItem(this.key)
    }

    get(): T | null {
        const json = this.delegate.getItem(this.key)
        if (json == null) {
            return null
        }

        return JSON.parse(json)
    }

    getOrDefault(defaultValue: T): T {
        const json = this.delegate.getItem(this.key)
        if (json == null) {
            return defaultValue
        }

        return JSON.parse(json)
    }

    update(defaultValue: T, valueTransformer: (value: T) => T) {
        const currentValue = this.getOrDefault(defaultValue)
        const newValue = valueTransformer(currentValue)
        this.store(newValue)
    }
}
