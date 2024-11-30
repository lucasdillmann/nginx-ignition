import ApiClientEventListener from "./ApiClientEventListener"
import ApiResponse from "../ApiResponse"

class ApiClientEventDispatcher {
    private readonly listeners: ApiClientEventListener[]

    constructor() {
        this.listeners = []
    }

    register(listener: ApiClientEventListener) {
        this.listeners.push(listener)
    }

    notifyRequest(request: RequestInit) {
        for (const listener of this.listeners) {
            try {
                listener.handleRequest(request)
            } catch (ex) {}
        }
    }

    notifyResponse(request: RequestInit, response: ApiResponse<any>) {
        for (const listener of this.listeners) {
            try {
                listener.handleResponse(request, response)
            } catch (ex) {}
        }
    }
}

// eslint-disable-next-line import/no-anonymous-default-export
export default new ApiClientEventDispatcher()
