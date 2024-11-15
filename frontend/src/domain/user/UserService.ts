import UserGateway from "./UserGateway";
import UserResponse from "./model/UserResponse";

export default class UserService {
    private gateway: UserGateway

    constructor() {
        this.gateway = new UserGateway()
    }

    async current(): Promise<UserResponse | undefined> {
        return this.gateway
            .getCurrent()
            .then(response => response.body)
            .catch(() => undefined)
    }
}
