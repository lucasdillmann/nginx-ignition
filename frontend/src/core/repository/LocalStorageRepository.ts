import AbstractStorageRepository from "./AbstractStorageRepository"

export default class LocalStorageRepository<T> extends AbstractStorageRepository<T> {
    constructor(key: string) {
        super(localStorage, key)
    }
}
