import LocalStorageRepository from "../repository/LocalStorageRepository"

class AuthenticationService {
    private readonly repository: LocalStorageRepository<string>

    constructor() {
        this.repository = new LocalStorageRepository("nginxIgnition.authentication.token")
    }

    getCurrentToken(): string | null {
        return this.repository.get()
    }

    setToken(token: string) {
        this.repository.set(token)
    }

    deleteToken() {
        this.repository.clear()
    }
}

// eslint-disable-next-line import/no-anonymous-default-export
export default new AuthenticationService()
