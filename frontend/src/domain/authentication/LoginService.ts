import UserGateway from "../user/UserGateway";
import UserLoginRequest from "../user/model/UserLoginRequest";
import {requireSuccessResponse} from "../../core/apiclient/ApiResponse";
import AuthenticationService from "../../core/authentication/AuthenticationService";

export default class LoginService {
    private readonly gateway: UserGateway

    constructor() {
        this.gateway = new UserGateway()
    }

    login(username: string, password: string): Promise<any> {
        const request: UserLoginRequest = { username, password }
        return this.gateway
            .login(request)
            .then(requireSuccessResponse)
            .then((response) => {
                AuthenticationService.setToken(response.token)
            })
    }
}
