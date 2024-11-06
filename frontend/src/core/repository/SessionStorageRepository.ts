import AbstractStorageRepository from "./AbstractStorageRepository"

export default class SessionStorageRepository<T> extends AbstractStorageRepository<T> {
    constructor(key: string) {
        super(sessionStorage, key)
    }
}
