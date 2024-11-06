import ApiClientEventListener from "./ApiClientEventListener";
import ApiResponse from "../ApiResponse";

class Dispatcher {
    private readonly listeners: ApiClientEventListener[]

    constructor() {
        this.listeners = []
    }

    register(listener: ApiClientEventListener) {
        this.listeners.push(listener)
    }

    notifyRequest(request: RequestInit) {
        for(const listener of this.listeners) {
            try {
                listener.handleRequest(request)
            } catch (ex) {
            }
        }
    }

    notifyResponse(request: RequestInit, response: ApiResponse<any>) {
        for(const listener of this.listeners) {
            try {
                listener.handleResponse(request, response)
            } catch (ex) {
            }
        }
    }
}

const ApiClientEventDispatcher = new Dispatcher()
export default ApiClientEventDispatcher
